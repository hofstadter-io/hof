package config

import (
	hof "github.com/hofstadter-io/hof/schema"

	"github.com/hofstadter-io/hofmod-cuefig/schema"
)

// Local config
#HofConfig: schema.#Config & {
	Name:       "config"
	Entrypoint: ".hofcfg.cue"

	ConfigSchema: #ConfigSchema
}

// Local secret
#HofSecret: schema.#Config & {
	Sensative:    true
	Name:         "secret"
	Entrypoint:   ".hofshh.cue"
	ConfigSchema: #SecretSchema
}

// (user/app config dir) config
#HofUserConfig: schema.#Config & {
	Name:         "hofcfg"
	Entrypoint:   ".hofcfg.cue"
	Workpath:     "hof"
	Location:     "user"
	ConfigSchema: #ConfigSchema
}

// (user/app config dir) secret
#HofUserSecret: schema.#Config & {
	Sensative:    true
	Name:         "hofshh"
	Entrypoint:   ".hofshh.cue"
	Workpath:     "hof"
	Location:     "user"
	ConfigSchema: #SecretSchema
}

#SecretSchema: {
	[Group=string]: {
		[Cred=string]: {
			[Key=string]: string
		}
	}
}

#ConfigSchema: {
	Current?: #ContextSchema
	Contexts?: [ContextName=string]: #ContextSchema & {name: ContextName}
	// This should only be used from the global context, local ought to be determined from walking up to find a .hofcfg.cue file
	// Unless... we want to subdivide workspaces, monorepo style (probably do want ot do this)
	// We can also associate developer setups with this
	Workspaces?: [WorkspaceName=string]:     #WorkspaceSchema & {name:   WorkspaceName}
	Environments?: [EnvironmentName=string]: #EnvironmentSchema & {name: EnvironmentName}

	// TODO add these to the config like above?
	Modelsets:  hof.#Modelsets
	Datastores: hof.#Datastores
}

#ContextSchema: {
	Name:         string
	Credentials?: string
	Workspace?:   string
	Environment?: string
	Account?:     string
	Billing?:     string
	Project?:     string
	Package?:     string
}

#WorkspaceSchema: {
	Name: string
	Dir:  string
	// TODO, reference these from the local schema
	Modelsets:  hof.#Modelsets
	Datastores: hof.#Datastores
}

#EnvironmentSchema: {
	Name:     string
	Provider: "local" | "docker" | *"local-kind" | "remote-kind" | "gke"
}
