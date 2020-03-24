import (
  "fmt"
	"net/http"
	"github.com/imroc/req"
	"crypto/tls"
	"crypto/md5"
	"encoding/hex"
	"encoding/base64"
	"encoding/json"
	"strings"
	"github.com/metal3-io/baremetal-operator/pkg/bmc"
  "github.com/redfishProvisioner/utils/hardware/disks"
  "github.com/redfishProvisioner/utils/preseed/iso"
	apiv1 "k8s.io/api/core/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  // batchv1 "k8s.io/api/batch/v1"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
  "log"

)

type RedfishClient struct {
	Name 		string
	AuthType    string
	BaseURL		string
	Header 		req.Header
	HttpClient  *req.Req
}

func init(){
		log.SetPrefix("INFO: ")
		log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func New(base_url, username, password string) (RedfishClient) {

	https_base_url := strings.Replace(base_url, "redfish", "https", 1)

	//Initialize Redfish Client
	redfish_client := RedfishClient {
		Name:       "RedfishClient",
		AuthType:   "Basic",
		BaseURL:    https_base_url
	}
	client := GetClient()

	//Set Proper Authorization and Other Headers
	bmc_username := strings.TrimSuffix(username, "\n")
	bmc_password := strings.TrimSuffix(password, "\n")
	b64_encoded_cred := redfish_client.EncodeString(bmc_username + ":" + bmc_password)
	auth_type := redfish_client.AuthType
	redfish_client.SetHeader("Authorization", auth_type + " " + b64_encoded_cred)
	redfish_client.SetHeader("Accept", "application/json")
	redfish_client.SetHeader("Content-Type","application/json")
	redfish_client.HttpClient = client
	return redfish_client
}

func GetClient() (*req.Req) {

	Req := req.New()
	trans, _ := Req.Client().Transport.(*http.Transport)
	trans.MaxIdleConns = 20
	trans.DisableKeepAlives = true
	trans.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return Req
}

func (redfishClient RedfishClient) EncodeString(data string) (string) {

	return base64.StdEncoding.EncodeToString([]byte(data))
}

func (redfishClient *RedfishClient) SetHeader(key string, value string) (req.Header) {

	if redfishClient.Header == nil {
		redfishClient.Header = req.Header{}
	}
	redfishClient.Header[key] = value
	return redfishClient.Header
}

// InspectHardware updates the HardwareDetails field of the host with
// details of devices discovered on the hardware. It may be called
// multiple times, and should return true for its dirty flag until the
// inspection is completed.
func (p *RedfishClient) InspectHardware() (result provisioner.Result, err error) {
	p.log.Info("inspecting hardware", "status", p.host.OperationalStatus())

	// The inspection is ongoing. We'll need to check the redfish
	// status for the server here until it is ready for us to get the
	// inspection details. Simulate that for now by creating the
	// hardware details struct as part of a second pass.
	if p.host.Status.HardwareDetails == nil {
		p.log.Info("continuing inspection by setting details")
		p.host.Status.HardwareDetails =
			&metalkubev1alpha1.HardwareDetails{
				RAMGiB: 128,
				NIC: []metalkubev1alpha1.NIC{
					metalkubev1alpha1.NIC{
						Name:      "nic-1",
						Model:     "virt-io",
						Network:   "Pod Networking",
						MAC:       "some:mac:address",
						IP:        "192.168.100.1",
						SpeedGbps: 1,
					},
					metalkubev1alpha1.NIC{
						Name:      "nic-2",
						Model:     "e1000",
						Network:   "Pod Networking",
						MAC:       "some:other:mac:address",
						IP:        "192.168.100.2",
						SpeedGbps: 1,
					},
				},
				Storage: []metalkubev1alpha1.Storage{
					metalkubev1alpha1.Storage{
						Name:    "disk-1 (boot)",
						Type:    "SSD",
						SizeGiB: 1024 * 93,
						Model:   "Dell CFJ61",
					},
					metalkubev1alpha1.Storage{
						Name:    "disk-2",
						Type:    "SSD",
						SizeGiB: 1024 * 93,
						Model:   "Dell CFJ61",
					},
				},
				CPUs: []metalkubev1alpha1.CPU{
					metalkubev1alpha1.CPU{
						Type:     "x86",
						SpeedGHz: 3,
					},
				},
			}
		p.publisher("InspectionComplete", "Hardware inspection completed")
		result.Dirty = true
		return result, nil
	}

	return result, nil
}

func (p *RedfishClient) HardwareProfile(hp srting) (result provisioner.Result, err error) {
    // var result bool = true
    log.Info("Setting RAID Levels")
    // base_url := p.BaseURL
    // client := redfish_client.New(base_url, p.bmcCreds)
    for _, i := range hp{
       _ = disks.CreateVirtualDisks(p.HttpClient, i.Disk)
    }
    result.Dirty = true
    return result, nil
}

func HostProfile(hp, iso_url, iso_checksum srting) (result provisioner.Result, err error) {

    iso.PrepareISO(iso_url, iso_checksum, hp)
}

func (p *RedfishClient) Provision(hp provisioner.ISOConfig, getUserData provisioner.UserDataSource) (result provisioner.Result, err error) {

	  p.log.Info("provisioning image to host", "state", p.host.Status.Provisioning.State)

		//result.Dirty = true
		p.log.Info("Testing Provisioner")

		base_url := p.host.Spec.BMC.Address

		// Step 0 Initialize Redfish Client
		client := redfish_client.New(base_url, p.bmcCreds)
		node_uuid := client.GetUniqueNodeId(p.host.Name)
		fmt.Printf("Node ID is %s\n", node_uuid)

		node,err  := rp.FindExistingNode(node_uuid)
		if err == nil &&  node.UUID == ""  {

			result.RequeueAfter = time.Second * 120
			result.Dirty = true
			fmt.Println("New Node, Creating new Node")

			node.Name = p.host.Name
			node.UUID = node_uuid
			node.DeployStatus = "NEW"
			_,err = rp.CreateNode(node)


			iso_url := p.host.Spec.Image.URL
			iso_checksum := p.host.Spec.Image.Checksum
			user_data, err1 := getUserData()
			fmt.Println(err1)

			isoconfig, err := hp()
			storageConfig, err := yaml.Marshal(&isoconfig.Storage)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(storageConfig))

			platformConfig, err := yaml.Marshal(&isoconfig.Platform)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(platformConfig))

			// go rp.PrepareISO(iso_url, iso_checksum, user_data, node, string(storageConfig), string(platformConfig))


		} else if err == nil &&  node.DeployStatus == "NEW"  && node.ImageURL != "" {

			fmt.Printf("\nImage URL is %s\n", node.ImageURL )

			result.RequeueAfter = time.Minute * 5
			result.Dirty = true


			fmt.Printf("\nStarting deployment of node %s\n" , node.UUID)

			//Step 1 Eject Existing ISO
			client.EjectISO()

			//Step 2 Insert Ubuntu ISO
			iso_inserted := client.InsertISO(node.ImageURL)
			if iso_inserted == false {
				fmt.Println("Inserting ISO Failed.")
				os.Exit(1)
			}


			//Step 3 Set Onetime boot to CD ROM
			client.SetOneTimeBoot()
			client.Reboot()

			node.DeployStatus = "OSINSTALLCALLBACKWAIT"
			_,err = rp.UpdateNode(node)

		} else if  err == nil &&  node.DeployStatus == "OSINSTALLCALLBACKWAIT"  {

			fmt.Printf("\nWaiting for OS installation to Finish \n")
			result.RequeueAfter = time.Minute * 2
			result.Dirty = true

		} else if  err == nil &&  node.DeployStatus == "OSINSTALLED"  {

			fmt.Printf(" \n deployment success for node %s\n" , node.UUID)

			result.Dirty = false
		} else {
			result.RequeueAfter = time.Minute * 1
			result.Dirty = true
		}
		if err != nil {
			fmt.Println(err)
			result.RequeueAfter = time.Minute * 5
			result.Dirty = true
		}

	return result, nil
}

// Deprovision prepares the host to be removed from the cluster. It
// may be called multiple times, and should return true for its dirty
// flag until the deprovisioning operation is completed.
func (p *RedfishClient) Deprovision(deleteIt bool) (result provisioner.Result, err error) {
	p.log.Info("ensuring host is removed")

	result.RequeueAfter = deprovisionRequeueDelay

	base_url := p.host.Spec.BMC.Address
	client := redfish_client.New(base_url, p.bmcCreds)
	node_uuid := client.GetUniqueNodeId(p.host.Name)
	node,err  := rp.FindExistingNode(node_uuid)

	if err == nil {
		fmt.Printf("\n\n Cleaning up node %s", node_uuid)
		_,err = rp.DeleteNode(node)
		if err != nil {
			fmt.Println(err)
		}
	}
	// NOTE(dhellmann): In order to simulate a multi-step process,
	// modify some of the status data structures. This is likely not
	// necessary once we really have redfish doing the deprovisioning
	// and we can monitor it's status.

	if p.host.Status.HardwareDetails != nil {
		p.publisher("DeprovisionStarted", "Image deprovisioning started")
		p.log.Info("clearing hardware details")
		p.host.Status.HardwareDetails = nil
		result.Dirty = true
		return result, nil
	}

	if p.host.Status.Provisioning.ID != "" {
		p.log.Info("clearing provisioning id")
		p.host.Status.Provisioning.ID = ""
		result.Dirty = true
		return result, nil
	}

  _ = rp.DeleteSecret(node_uuid + "-kubeconfig", p.host.Namespace)

  if p.host.Spec.SiteProfile.Name != "" {
       p.log.Info("Clearing site config")
       rp.DeleteSiteConfig(p.host.Spec.SiteProfile.Name, p.host.Spec.SiteProfile.Namespace)
  }

  if p.host.Spec.UserData != nil {
       p.log.Info("Clearing user data")
       rp.DeleteSecret(p.host.Spec.UserData.Name, p.host.Spec.UserData.Namespace)
      }

	p.publisher("DeprovisionComplete", "Image deprovisioning completed")
	return result, nil
}

// PowerOn ensures the server is powered on independently of any image
// provisioning operation.
func (p *RedfishClient) PowerOn() (result provisioner.Result, err error) {
	p.log.Info("ensuring host is powered on")

	if !p.host.Status.PoweredOn {
		p.publisher("PowerOn", "Host powered on")
		p.log.Info("changing status")
		p.host.Status.PoweredOn = true
		result.Dirty = true
		return result, nil
	}

	return result, nil
}

// PowerOff ensures the server is powered off independently of any image
// provisioning operation.
func (p *RedfishClient) PowerOff() (result provisioner.Result, err error) {
	p.log.Info("ensuring host is powered off")

	if p.host.Status.PoweredOn {
		p.publisher("PowerOff", "Host powered off")
		p.log.Info("changing status")
		p.host.Status.PoweredOn = false
		result.Dirty = true
		return result, nil
	}

	return result, nil
}

func (redfishClient RedfishClient) GetVirtualMediaStatus() (bool) {

	endpoint := redfishClient.ManagerURL("VirtualMedia","CD")
	header := redfishClient.Header
	r, err := redfishClient.HttpClient.Get(endpoint, header)

	res := CheckErrorAndReturn(r,err)
	var data map[string]interface{}
	res.ToJSON(&data)       // response => struct/map
	if data["ConnectedVia"] == "NotConnected" {
		return false
	}
	return true
}

func (redfishClient RedfishClient) InsertISO(image_url string) (bool) {
	fmt.Printf("Starting ISO attach\n")
	if redfishClient.GetVirtualMediaStatus() == true {
		fmt.Printf("Skipping Iso Insert. CD already Attached\n")
		return false
	} else {
		// image_url, err := redfishClient.PrepareISO(node_id)
        //         if err != nil {
		// 			fmt.Println(err)
		// 			return false
        //         }
		fmt.Printf("Attachig new ISO %s\n", image_url)
		endpoint := redfishClient.ManagerURL("VirtualMedia","CD", "Actions", "VirtualMedia.InsertMedia")
		header := redfishClient.Header
		body := `{"Image": "` + image_url +`"}`
		r, err := redfishClient.HttpClient.Post(endpoint, header, body)
		CheckErrorAndReturn(r,err)
		return true
	}

}


func(redfishClient RedfishClient) SetOneTimeBoot () (bool) {
		// Actions/Oem/EID_674_Manager.ImportSystemConfiguration
		fmt.Printf("Setting Onetime boot to VirtualMediaCDROM\n")
		endpoint := redfishClient.ManagerURL("Actions","Oem", "EID_674_Manager.ImportSystemConfiguration")
		header := redfishClient.Header
		body := `{
			"ShareParameters": {
				"Target": "ALL"
			},
			"ImportBuffer": "<SystemConfiguration><Component FQDD=\"iDRAC.Embedded.1\"><Attribute Name=\"ServerBoot.1#BootOnce\">Enabled</Attribute><Attribute Name=\"ServerBoot.1#FirstBootDevice\">VCD-DVD</Attribute></Component></SystemConfiguration>"
		}`
		//fmt.Printf("Body is %s \n", body)
		r, err := redfishClient.HttpClient.Post(endpoint, header, body)
		CheckErrorAndReturn(r,err)
		return true

}

func(redfishClient RedfishClient) Reboot () (bool) {
	///Systems/System.Embedded.1/Actions/ComputerSystem.Reset
	fmt.Printf("Starting OS installation. Rebooting the node\n")
	endpoint := redfishClient.SystemURL("Actions","ComputerSystem.Reset")
	header := redfishClient.Header
	body := `{"ResetType" : "ForceRestart" }`
	//fmt.Printf("Body is %s \n", body)
	r, err := redfishClient.HttpClient.Post(endpoint, header, body)
	CheckErrorAndReturn(r,err)
	return true
}


func (redfishClient RedfishClient) EjectISO() (bool) {
	fmt.Printf("Starting ISO Eject\n")
	if redfishClient.GetVirtualMediaStatus() == false {
		fmt.Printf("No CD to eject\n")
	} else {
		fmt.Printf("Ejecting existing CD\n")
		endpoint := redfishClient.ManagerURL("VirtualMedia","CD", "Actions", "VirtualMedia.EjectMedia")
		header := redfishClient.Header
		body := `{}`
		r, err := redfishClient.HttpClient.Post(endpoint, header, body)
		CheckErrorAndReturn(r,err)
	}
	return true
}

func (redfishClient RedfishClient)  GetUniqueNodeId(hostname string) (string) {

	h := md5.New()
    h.Write([]byte(strings.ToLower(hostname)))
    return hex.EncodeToString(h.Sum(nil))
}

func CheckErrorAndReturn(res *req.Resp, err error) (*req.Resp) {

	//fmt.Println(res)
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}

	return res
}
