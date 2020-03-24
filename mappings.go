package redfish


type UrlMappings interface {

    ManagerURL(type, urk_key string) string

    SystemURL(parts string)

    UrlMappings(parts string)
}
