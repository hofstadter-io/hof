package ga

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hofstadter-io/cinful"
	"github.com/hofstadter-io/hof/lib/yagu"

	"github.com/hofstadter-io/hof/cmd/hof/verinfo"
)

var dir, fn, cid string
var isCI bool

var debug = false

func init() {
	if debug {
		fmt.Println("init telemetry")
	}
	// short-circuit
	ev := os.Getenv("HOF_TELEMETRY_DISABLED")
	eb, _ := strconv.ParseBool(ev)
	if ev != "" {
		if eb {
			if debug {
				fmt.Println("telemetry disabled in env")
			}
			cid = "disabled"
			return
		}
	}

	// setup dir info
	ucd, err := os.UserConfigDir()
	if err != nil {
		cid = "disabled"
		return
	}
	// workaround for running in TestScript tool
	if strings.HasPrefix(ucd, "/no-home") {
		ucd = strings.TrimPrefix(ucd, "/")
	}
	dir = filepath.Join(ucd, "hof")
	fn = filepath.Join(dir, ".uuid")


	// try reading
	cid, err = readGaId()
	if debug {
		fmt.Println("realGaId:", cid, err)
	}
	if err != nil {
		cid = "missing"
	}
	if cid == "disabled" {
		if debug {
			fmt.Println("telemetry disabled in cfg")
		}
		return
	}

	if debug {
		fmt.Println("telemetry ok:", cid)
	}

	// does it exist already
	if cid != "missing" {
		return
	}

	// generate an ID
	id, _ := uuid.NewUUID()
	cid = id.String()

	// check if in CI, and add prefix
	vendor := cinful.Info()
	if vendor != nil {
		if debug {
			fmt.Println("in CI")
		}
		isCI = true
		cid = "CI-" + cid
	}

	err = writeGaId(cid)
	if err != nil {
		fmt.Println("Error writing telemetry config, please let the devs know")
		return
	}

	// create the ID for the first time
	// prompting user for approval

	// if not found, ask and write
	// if ev != "" {
	// 	if !eb {
	// 		approve := askGaId()
	// 	}
	// }
	// if !approve {
	// 	err = writeGaId("disabled")
	// } else {
	// 	id, _ := uuid.NewUUID()
	// 	cid = id.String()
	// }
}

func SendCommandPath(cmd string) {
	if debug {
		fmt.Println("try sending:", cmd)
	}
	cs := strings.Fields(cmd)
	c := strings.Join(cs[1:], "/")
	l := "user"
	if isCI {
		l = "ci"
	}

	if debug {
		fmt.Println("SendGaEvent:", c, l, cid)
	}
	if cid == "disabled" {
		return
	}

	ua := fmt.Sprintf(
		"%s %s (%s/%s)",
		"hof", verinfo.Version,
		verinfo.BuildOS, verinfo.BuildArch,
	)

	cfg := yagu.GaConfig{
		TID: "UA-103579574-5",
		CID: cid,
		UA:  ua,
		CS:  verinfo.Version,
		CM:  l,
	}

	evt := yagu.GaEvent{
		Action:   c, // path or cmd here
		Source:   fmt.Sprintf("%s/%s",verinfo.BuildOS, verinfo.BuildArch),
		Category: l,
	}

	if debug {
		fmt.Printf("sending:\n%#v\n%#v\n", cfg, evt)
	}

	yagu.SendGaEvent(cfg, evt)
}

func readGaId() (string, error) {
	_, err := os.Lstat(fn)
	if err != nil {
		return "missing", err
	}

	content, err := ioutil.ReadFile(fn)
	if err != nil {
		return "missing", err
	}

	if debug {
		fmt.Printf("read %q from %s\n", string(content), fn)
	}

	return string(content), nil
}

func writeGaId(value string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fn, []byte(value), 0644)
	if err != nil {
		return err
	}

	return nil
}

var askMsg = `We only send the command run, no args or input.
You can disable at any time by setting
  HOF_TELEMETRY_DISABLED=1

Would you like to help by sharing very basic usage stats?`

func askGaId() bool {
	prompt := "\n(y/n) >"
	fmt.Printf(askMsg + prompt)
	var ans string
	fmt.Scanln(&ans)
	a := strings.ToLower(ans)
	if a == "n" || a == "no" {
		return false
	}
	return true
}
