package runtime

import "github.com/hofstadter-io/hof/lib/agent"

func extractMeta(a *agent.Agentic) (aname, akind, mname string) {
	akind = a.Hof.Agentic.Kind
	aname = a.Hof.Agentic.Name
	mname = a.Hof.Metadata.Name
	if aname == "" {
		aname = mname
	}
	return aname, akind, mname
}
