import "strings"

Server: {
	// ...

	// various casings of the resource Name
	serverName:  strings.ToCamel(Name)
	ServerName:  strings.ToTitle(Name)
	SERVER_NAME: strings.ToUpper(Name)
}
