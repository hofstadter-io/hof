package ga

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	if os.Getenv("HOF_TELEMETRY_DISABLED") != "" {
		if debug {
			fmt.Println("telemetry disabled in env")
		}
		cid = "disabled"
		return
	}

	// check if in CI
	vendor := cinful.Info()
	if vendor != nil {
		if debug {
			fmt.Println("in CI")
		}
		// generate an ID
		id, _ := uuid.NewUUID()
		cid = id.String()

		isCI = true
		cid = "CI-" + cid
		return
	}

	// setup dir info
	ucd, err := os.UserConfigDir()
	if err != nil {
		cid = "disabled"
		return
	}
	dir = filepath.Join(ucd, "hof")
	fn = filepath.Join(dir, ".uuid")

	// try reading
	cid, err = readGaId()
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

	// create the ID for the first time
	// prompting user for approval

	// if not found, ask and write
	approve := askGaId()
	if !approve {
		err = writeGaId("disabled")
	} else {
		id, _ := uuid.NewUUID()
		cid = id.String()
		err = writeGaId(cid)
	}

	if err != nil {
		fmt.Println("Error writing telemetry config, please let the devs know")
		return
	}
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
	SendGaEvent(c, l, 0)
}

func SendGaEvent(action, label string, value int) {
	if os.Getenv("HOF_TELEMETRY_DISABLED") != "" {
		return
	}

	if cid == "disabled" {
		return
	}

	ua := fmt.Sprintf(
		"%s/%s %s/%s",
		"hof", verinfo.Version,
		verinfo.BuildOS, verinfo.BuildArch,
	)

	cfg := yagu.GaConfig{
		TID: "UA-103579574-5",
		CID: cid,
		UA:  ua,
		CN:  "hof",
		CS:  "hof/" + verinfo.Version,
		CM:  verinfo.Version,
	}

	evt := yagu.GaEvent{
		Source:   cfg.UA,
		Category: "hof",
		Action:   action,
		Label:    label,
	}

	if value >= 0 {
		evt.Value = value
	}

	if debug {
		fmt.Printf("sending:\n%#v\n%#v\n", cfg, evt)
	}

	yagu.SendGaEvent(cfg, evt)
}

func readGaId() (string, error) {
	// ucd := yagu.UserHomeDir()

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
