package ga

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hofstadter-io/cinful"

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
	if cid == "disabled" {
		return
	}
	if debug {
		fmt.Println("try sending:", cmd)
	}

	vals := url.Values{}

	// always
	vals.Add("v", "2")
	if debug {
		// vals.Add("_dbg", "1")
	}

	// ids
	vals.Add("measurement_id", "G-6CYEVMZL4R")
	// vals.Add("api_secret", os.Getenv("GA_MP_APIKEY")) 
	vals.Add("cid", cid)
	// vals.Add("_p", fmt.Sprint(rand.Intn(10000000000-1)))

	// system info
	vals.Add("uaa", verinfo.BuildArch)
	vals.Add("uap", verinfo.BuildOS)
	vals.Add("dh", "cli.hofstadter.io")

	// event info
	cs := strings.Fields(cmd)
	c := strings.Join(cs[1:], "/")
	l := "user"
	if isCI {
		l = "ci"
	}

	vals.Add("en", "pageview")
	vals.Add("dt", c)
	vals.Add("dl", "http://cli.hofstadter.io/" + c)
	vals.Add("cs", l)
	vals.Add("cm", verinfo.Version)

	if debug {
		// fmt.Printf("vals:%v\n", vals)
	}

	gaURL := "https://next.hofstadter.io/mp/collect?"
	//gaURL = "https://www.google-analytics.com/debug/mp/collect?"
	//gaURL = "http://localhost:8080/mp/collect?"
	url := gaURL +	vals.Encode()
	if debug {
		fmt.Println(url)
	}

	evt := map[string]any{
		"name": "cmd_run",
		"params": map[string]any{
			"cmd": c,
			"version": verinfo.Version,
			"arch": verinfo.BuildArch,
			"os": verinfo.BuildOS,
			"runtype": l,
		},
	}

	obj := map[string]any{
		"client_id": cid,
		"events": []map[string]any { evt },
		//"user_properties": map[string]any{
		//  "id": cid,
		//  "version": verinfo.Version,
		//  "arch": verinfo.BuildArch,
		//  "os": verinfo.BuildOS,
		//  "runtype": l,
		//},
	}

	postBody, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		if debug {
			fmt.Println("Marshal Error: ", err)
		}
	}
	if debug {
		fmt.Println("body: ", string(postBody))
	}


	reqBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(url, "application/json", reqBody)

	if debug {
		fmt.Println(resp, err)
	}

	// fmt.Println(resp, body, errs)

	// hacky check for golang issue(?)

	//if len(errs) != 0 && !strings.Contains(errs[0].Error(), "http2: server sent GOAWAY and closed the connection") {
	//  return body, errs[0]
	//}

	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if debug {
			fmt.Println("Read Error: ", resp.StatusCode, string(body), err)
		}
	}

	if resp.StatusCode >= 500 {
		if debug {
			fmt.Println("Internal Error: ", resp.StatusCode, string(body), err)
		}
		// return body, errors.New("Internal Error: " + body)
	}
	if resp.StatusCode >= 400 {
		if debug {
			fmt.Println("Bad Request: ", resp.StatusCode, string(body))
		}
		// return body, errors.New("Bad Request: " + body)
	}

	if debug {
		fmt.Println(string(body))
	}
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
