package consts

// Directory Settings
const VEG_SYS_DATA_DIR_VAR = "VEG_SYS_DATA_DIR"
const VEG_USER_DATA_DIR_VAR = "VEG_USER_DATA_DIR"
const VEG_SYS_DATA_DIR_DEFAULT = "/var/lib/veg/data"
const VEG_REPO_LOCAL_PATH = `.veg`

// Database settings
const VEG_DATABASE_TYPE_VAR = "VEG_DATABASE_TYPE"
const VEG_DATABASE_TYPE_DEFAULT = "sqlite"
const VEG_DATABASE_CONN_VAR = "VEG_DATABASE_CONN"
const VEG_DATABASE_CONN_DEFAULT = "veg.db"

// Registry settings
const VEG_REGISTRY_DEFAULT_CONTAINER_NAME = "veg-registry"
const VEG_REGISTRY_DEFAULT_IMAGE = "registry:3"
const VEG_REGISTRY_DEFAULT_HOST = "host.docker.internal:5000"
const VEG_REGISTRY_DATA_VAR = "VEG_REGISTRY_DATA"

// Dagger settings
const VEG_DAGGER_ENGINE_DEFAULT_CONTAINER_NAME = "veg-dagger-engine"
const VEG_DAGGER_ENGINE_DEFAULT_IMAGE = "registry.dagger.io/engine:v0.19.10"
const VEG_DAGGER_ENGINE_CONFIG_VAR = "VEG_DAGGER_ENGINE_CONFIG"
const VEG_DAGGER_ENGINE_CONFIG_DEFAULT = `
{
  "registries": {
    "host.docker.internal:5000": {
      "http": true
    }
  },
  "gc": {
    "enabled": true,
    "reservedSpace": "50GB",
    "maxUsedSpace": "100GB",
    "minFreeSpace": "5GB"
  }
}
`

// Formatter settings

// Agentic Settings
const VEG_DEFAULT_USER = "tony"
const VEG_USER_HEADER = "X-Veg-User"
