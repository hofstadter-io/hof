package config

import (
	// hof "github.com/hofstadter-io/hof/schema"

	"github.com/hofstadter-io/hofmod-cuefig/schema"
)

// Local context
#HofContext: schema._Config & {
	Name:         "context"
	Entrypoint:   ".hofctx.cue"
	ConfigSchema: #ContextSchema
}

// Local config
#HofConfig: schema._Config & {
	Name:         "config"
	Entrypoint:   ".hofcfg.cue"
	ConfigSchema: #WorkspaceSchema
}

// Local secret
#HofSecret: schema._Config & {
	Sensative:    true
	Name:         "secret"
	Entrypoint:   ".hofshh.cue"
	ConfigSchema: #SecretSchema
}

// (user/app config dir) context
#HofUserContext: schema._Config & {
	Name:         "hofctx"
	Entrypoint:   ".hofctx.cue"
	Workpath:     "hof"
	Location:     "user"
	ConfigSchema: #ContextSchema
}

// (user/app config dir) config
#HofUserConfig: schema._Config & {
	Name:         "hofcfg"
	Entrypoint:   ".hofcfg.cue"
	Workpath:     "hof"
	Location:     "user"
	ConfigSchema: #ConfigSchema
}

// (user/app config dir) secret
#HofUserSecret: schema._Config & {
	Sensative:    true
	Name:         "hofshh"
	Entrypoint:   ".hofshh.cue"
	Workpath:     "hof"
	Location:     "user"
	ConfigSchema: #SecretSchema
}

#ContextSchema: {
	Current?: #ContextItemSchema
	Contexts?: [ContextName=string]: #ContextItemSchema & {name: ContextName}
}

#ContextItemSchema: {
	Name:         string
	Credentials?: string
	Environment?: string
	Account?:     string
	Billing?:     string
	Project?:     string
	Package?:     string
	...
}

// Secrets at both local / global level
#SecretSchema: {
	[Group=string]: {
		[Cred=string]: {
			[Key=string]: string
		}
	}
}

// Hof tool configuration
#ConfigSchema: {
	// This should only be used from the global context, local ought to be determined from walking up to find a .hofcfg.cue file
	// Unless... we want to subdivide workspaces, monorepo style (probably do want ot do this)
	// We can also associate developer setups with this
	// ... rethinking having multiple workspaces per repo, doesn't fit with the latest UX (in particular workspace/workflow integration)
	Workspaces?: [WorkspaceName=string]:     #WorkspaceSchema & {name:   WorkspaceName}

	...
}

// Workspace specfic config
#WorkspaceSchema: {
	Name: string | *""
	Dir:  string | *""

	ModelsDir: string | *"models"
	ResourcesDir: string | *"resources"
}
