package hardware

import (
    // hardware "github.com/Kiriti29/RedfishProvisioner/utils/hardware"
    rp "github.com/Kiriti29/RedfishProvisioner"
    url_mappings "github.com/Kiriti29/RedfishProvisioner/urls/idrac"
    "github.com/imroc/req"
  )

type hardwareProfile struct {
    // BaseURL string
    // hp  []hardware.HardwareProfile
    // HttpClient  *req.Req
    // Header 	req.Header
    // AuthType:   string
    RedfishClient *rp.RedfishClient
    UrlMappings: *url_mappings.UrlMappings
}

func New(hp, rp_client rp.RedfishClient) *hardwareProfile {
    hard_prof := hardware.GetHardwareProfile(hp)
    return &hardwareProfile{
        BaseURL: "url",
        hp: hard_prof,
        AuthType:   "Basic",
        HttpClient: http_client,
        Header: header,
        UrlMappings: url_mappings.New(baseurl)
    }
}

func (hp hardwareProfile) checkJobStatus(job_url string) bool {
  var result bool = false
  header := hp.Header
  for {
      r, err := hp.HttpClient.Get(job_url, header)
      resp := CheckErrorAndReturn(r,err)
      var data map[string]interface{}
      resp.ToJSON(&data)
      if data["JobState"] == "Completed" {
          result = true
          break
      }
  }
  return result
}

func (hp hardwareProfile) cleanVirtualDIskIfEExists() bool {
		// url := "https://32.67.151.80/redfish/v1/Systems/System.Embedded.1/Storage/Volumes"
		var result bool = false
		endpoint := hp.UrlMappings.SystemURL("Storage", "RAID.Slot.6-1", "Volumes")
    header := hp.Header
    r, err := hp.HttpClient.Get(endpoint, header)
    resp := CheckErrorAndReturn(r,err)
    var data map[string]interface{}
    resp.ToJSON(&data)
    var disks []map[string]string
    tmp, _ := json.Marshal(data["Members"])
    json.Unmarshal(tmp, &disks)
		if len(disks) == 0{
				log.Println("No existing RAID found. Creating Virtual disks")
				return true
		}
		log.Println("Found existing RAID config. Deleting existing RAID")
    tmp1 := strings.Split(hp.BaseURL, "/")
    disk_url := tmp1[0] + "//" + tmp1[2]
    for _,disk := range disks{
				msg := fmt.Sprintf("Deleting the virtual disk %s", disk["@odata.id"])
				log.Println(msg)
				endpoint = disk_url + disk["@odata.id"]
        r, _ = hp.HttpClient.Delete(endpoint, header)
        job_id := strings.Split(r.Response().Header["Location"][0], "/")
        job := job_id[len(job_id) - 1]
        job_url := hp.UrlMappings.ManagerURL() + "/Jobs/" + job
				log.Println("Waiting for the delete job to finish")
        result = hp.checkJobStatus(job_url)
    }
    return result
}

func (hp hardwareProfile) CreateVirtualDisks(hp []hardware.Disks) bool{
		if !hp.cleanVirtualDIskIfEExists() {
				return false
		}
		var result bool = false
		for _, d := range hp {
				var VolumeType string
				switch {
				case d.RaidType == "50":
						VolumeType = "SpannedStripesWithParity"
				case d.RaidType == "1":
						VolumeType = "Mirrored"
				case d.RaidType == "5":
						VolumeType = "StripedWithParity"
				case d.RaidType == "10":
						VolumeType = "SpannedMirrors"
				default:
						VolumeType = "NonRedundant"
				}
				var drives []string
        for _, disk := range d.Disk{
            drives = append(drives, fmt.Sprintf(`{"@odata.id": "/redfish/v1/Systems/System.Embedded.1/Storage/Drives/%s"}`, disk))
						log.Println(fmt.Sprintf("Creating RAID for disk %s", disk))
        }
        tmp := strings.Join(drives, ",")
        drive := `[` + tmp + `]`
				payload := fmt.Sprintf(`{
					"VolumeType": "%s",
					"Name": "%s",
					"Drives": %s
					}`, VolumeType, d.Name, drive)
				endpoint := hp.UrlMappings.SystemURL(get_virtual_disks)
				header := hp.Header
				r, err := hp.HttpClient.Post(endpoint, header, payload)
				resp := CheckErrorAndReturn(r,err)
				job_id := strings.Split(resp.Response().Header["Location"][0], "/")
				job := job_id[len(job_id) - 1]
				job_url := hp.UrlMappings.ManagerURL() + "/Jobs/" + job
				// body := `{}`
				log.Println(fmt.Sprintf("Waiting for job %s to complete", job))
				result = hp.checkJobStatus(job_url)
		}
		return result
}
