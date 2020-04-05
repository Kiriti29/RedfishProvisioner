package redfish

import (
  //"fmt"
  "strings"
	"net/http"
	"github.com/imroc/req"
	"crypto/tls"
	"encoding/base64"
)

type RedfishClient struct {
	Name 		string
	AuthType    string
	BaseURL		string
	Header 		req.Header
	HttpClient  *req.Req
}

func New(base_url, username, password string) (RedfishClient) {

	https_base_url := strings.Replace(base_url, "redfish", "https", 1)

	//Initialize Redfish Client
	redfish_client := RedfishClient {
		Name:       "RedfishClient",
		AuthType:   "Basic",
		BaseURL:    https_base_url,
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

type RedfishBase interface{

    // Inspects hardware according to the Hardware Profile
    InspectHardware(hp string)

    // Sets the hardware profile given by the user. This includes configuring
    // Raids, network devices
    HardwareProfile(hp string)

    // Reads the storage partitioning, Platform config which includes kernel,
    // grub, cpus parameters and creates an iso image by injecting all these
    // parameters into preseed file
    HostProfile(hp string)

    // Triggers the provisioning function which installs OS on baremetal node
    // using the redfish urls and the iso created
    Provision(userdata, hp string) bool

    // Removes the BareMetalHost and all its related resources(secrets,
    // configmaps etc)
    Deprovision(uuid string)

    // Power on a baremetal node
    PowerOn(uuid string)

    // Power off a baremetal node
    PowerOff(uuid string)
}
