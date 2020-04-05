package hardware

import (
    // hardware "github.com/Kiriti29/RedfishProvisioner/utils/hardware"
    "encoding/json"
    "fmt"
    "log"
    "strings"
    rp "github.com/Kiriti29/RedfishProvisioner"
    url_mappings "github.com/Kiriti29/RedfishProvisioner/urls/idrac"
    // "github.com/imroc/req"
  )

type hardwareProfile struct {
    // BaseURL string
    HP  []HardwareProfile
    // HttpClient  *req.Req
    // Header 	req.Header
    // AuthType:   string
    RedfishClient rp.RedfishClient
    UrlMappings *url_mappings.UrlMappings
}

func init(){
		log.SetPrefix("INFO: ")
		log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func New(hp string, rp_client rp.RedfishClient) *hardwareProfile {
    hard_prof := GetHardwareProfile(hp)
    return &hardwareProfile{
        RedfishClient:  rp_client,
        UrlMappings: url_mappings.New(rp_client.BaseURL),
        HP: hard_prof,
    }
}

func (hp hardwareProfile) checkJobStatus(job_url string) bool {
  var result bool = false
  header := hp.RedfishClient.Header
  for {
      r, err := hp.RedfishClient.HttpClient.Get(job_url, header)
      resp := rp.CheckErrorAndReturn(r,err)
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
		endpoint := hp.UrlMappings.SystemURL("get_virtual_disks")
    header := hp.RedfishClient.Header
    r, err := hp.RedfishClient.HttpClient.Get(endpoint, header)
    resp := rp.CheckErrorAndReturn(r,err)
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
    tmp1 := strings.Split(hp.RedfishClient.BaseURL, "/")
    disk_url := tmp1[0] + "//" + tmp1[2]
    for _,disk := range disks{
				msg := fmt.Sprintf("Deleting the virtual disk %s", disk["@odata.id"])
				log.Println(msg)
				endpoint = disk_url + disk["@odata.id"]
        r, _ = hp.RedfishClient.HttpClient.Delete(endpoint, header)
        job_id := strings.Split(r.Response().Header["Location"][0], "/")
        job := job_id[len(job_id) - 1]
        job_url := hp.UrlMappings.ManagerURL("") + "/Jobs/" + job
				log.Println("Waiting for the delete job to finish")
        result = hp.checkJobStatus(job_url)
    }
    return result
}

func (h hardwareProfile) CreateVirtualDisks() bool{
		if !h.cleanVirtualDIskIfEExists() {
				return false
		}
		var result bool = false
    for _, tmp := range h.HP  {
  		for _, d := range tmp.Disk {
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
  				endpoint := h.UrlMappings.SystemURL("get_virtual_disks")
  				header := h.RedfishClient.Header
  				r, err := h.RedfishClient.HttpClient.Post(endpoint, header, payload)
  				resp := rp.CheckErrorAndReturn(r,err)
  				job_id := strings.Split(resp.Response().Header["Location"][0], "/")
  				job := job_id[len(job_id) - 1]
  				job_url := h.UrlMappings.ManagerURL("") + "/Jobs/" + job
  				// body := `{}`
  				log.Println(fmt.Sprintf("Waiting for job %s to complete", job))
  				result = h.checkJobStatus(job_url)
  		}
    }
		return result
}
