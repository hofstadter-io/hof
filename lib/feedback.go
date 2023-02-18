package lib

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/cmd/hof/verinfo"
)

var info = `
Add more information here

---

Hof Metadata:

<pre>
Version:     v%s
Commit:      %s

BuildDate:   %s
GoVersion:   %s
OS / Arch:   %s %s
</pre>
`

func SendFeedback(args []string) error {
	title := url.QueryEscape(strings.Join(args, " "))

	body := fmt.Sprintf(
		info,
		verinfo.Version,
		verinfo.Commit,
		verinfo.BuildDate,
		verinfo.GoVersion,
		verinfo.BuildOS,
		verinfo.BuildArch,
	)
	body = url.QueryEscape(body)
	
	url := fmt.Sprintf("https://github.com/hofstadter-io/hof/issues/new?labels=feedback&title=%s&body=%s", title, body)
	yagu.OpenBrowserCmdSafe(url)

	return nil
}
