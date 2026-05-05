package config

var Veg *VegConfig

type VegConfig struct {
	SysDataDir  string
	UserDataDir string

	DatabaseType string
	DatabaseConn string

	RegistryData       string
	DaggerEngineConfig string
}
