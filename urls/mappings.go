package urls


type UrlMappings interface {

    ManagerURL(url_key string) string

    SystemURL(parts string)

    // UrlMappings(parts string)
}
