package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

DaggerooCommand: schema.Command & {
	Name:  "daggeroo"
	Usage: "daggeroo [args]"
	Short: "dagger run helper for the extension server"
	Long:  "dagger run helper for the extension server"

	Args: [{
		Name:     "id"
		Type:     "string"
		Required: true
		Help:     "Dagger ID for a directory to mount"
	}]

	Flags: [...schema.Flag] & [ {
		Name:    "Workdir"
		Type:    "string"
		Default: "\"/work\""
		Help:    "work directory to start in"
		Long:    "workdir"
		// Short:   "W"
	}, {
		Name:    "Image"
		Type:    "string"
		Default: "\"qmcgaw/godevcontainer:debian\""
		Help:    "image uri to launch into"
		Long:    "image"
		// Short:   "I"
	}]
}

