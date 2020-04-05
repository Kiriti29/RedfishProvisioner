package idrac_urls

type UrlMappings struct{
    BaseURL string
    ManagerBaseURL  string
    SystemBaseURL string
}


mappings := {
  "Manager": {
      "get_virtual_media": "/VirtualMedia/CD"
  },
  "System": {
      "get_virtual_disks": "/Storage/RAID.Slot.6-1/Volumes"
  }
}

func New(url string) *UrlMappings{
  return &UrlMappings{
    BaseURL: url,
    ManagerBaseURL: "/Managers/iDRAC.Embedded.1",
    SystemBaseURL: "/Systems/System.Embedded.1"
  }
}

// func (urlmappings UrlMappings) UrlMappings(type, url_key string) string {
//     return mappings[type]["base_url"] + mappings[type][url_key]
// }

func (urlmappings UrlMappings) SystemURL(parts string) string {

  if parts == ""{
      return urlmappings.BaseURL + urlmappings.SystemURL
  } else {

	return urlmappings.BaseURL + urlmappings.ManagerBaseURL + mappings["Manager"][parts]
  }
}


func (urlmappings UrlMappings) ManagerURL(parts string) string {

  if parts == ""{
      return urlmappings.BaseURL + urlmappings.ManagerURL
  } else {
	return urlmappings.BaseURL + urlmappings.SystemBaseURL + mappings["System"][parts]
  }
}
