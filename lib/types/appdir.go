package types

// More TBD
type Appdir struct {
	Accounts   map[string]interface{}
	Workspaces map[string]interface{}
	Contexts   map[string]interface{}

	Clouds    map[string]interface{}
	Environs  map[string]interface{}
	Resources map[string]interface{}
}
