package idrac

import (
  "fmt"
  "encoding/json"
  "io/ioutil"
  "os"
)

type UrlMappings struct{
    BaseURL string
    ManagerBaseURL  string
    SystemBaseURL string
    Mappings  Mapping
}

type Mapping struct {
  Manager map[string]string
  System  map[string]string
}

func New(url string) *UrlMappings{
  jsonFile, err := os.Open("urls.json")
  if err != nil {
    fmt.Println(err)
  }
  defer jsonFile.Close()
  byteValue, _ := ioutil.ReadAll(jsonFile)
  var maps Mapping
  json.Unmarshal([]byte(byteValue), &maps)
  return &UrlMappings{
    BaseURL: url,
    ManagerBaseURL: "/Managers/iDRAC.Embedded.1",
    SystemBaseURL: "/Systems/System.Embedded.1",
    Mappings: maps,
  }
}

// func (urlmappings UrlMappings) UrlMappings(type, url_key string) string {
//     return mappings[type]["base_url"] + mappings[type][url_key]
// }

func (urlmappings UrlMappings) SystemURL(parts string) string {

  if parts == ""{
      return urlmappings.BaseURL + urlmappings.SystemBaseURL
  } else {

	return urlmappings.BaseURL + urlmappings.SystemBaseURL + urlmappings.Mappings.Manager[parts]
  }
}


func (urlmappings UrlMappings) ManagerURL(parts string) string {

  if parts == ""{
      return urlmappings.BaseURL + urlmappings.ManagerBaseURL
  } else {
	return urlmappings.BaseURL + urlmappings.ManagerBaseURL + urlmappings.Mappings.System[parts]
  }
}
