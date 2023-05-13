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

func SendCommandPath(cmd string) {
	if cid == "disabled" {
		return
	}

	cs := strings.Fields(cmd)
	c := strings.Join(cs[1:], "/")
	// cmd runner type
	l := "user"
	if isCI {
		l = "ci"
	}

	vals := url.Values{}
	vals.Add("measurement_id", "G-6CYEVMZL4R")

	evt := map[string]any{
		"name": "cmd_run",
		"params": map[string]any{
			"cmd": c,
			"version": verinfo.Version,
			"commit": verinfo.Commit,
			"arch": verinfo.BuildArch,
			"os": verinfo.BuildOS,
			"runtype": l,
			"engagement_time_msec" : 250,
		},
	}

	//pv := map[string]any{
	//  "name": "page_view",
	//  "params": map[string]any{
	//    "page_location": c,
	//    "page_title": cmd,
	//    "engagement_time_msec" : 250,
	//  },
	//}

	obj := map[string]any{
		"client_id": cid,
		"user_id": cid,
		"events": []map[string]any { evt },
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

	gaURL := "https://docs.hofstadter.io/mp/collect?"
	//gaURL = "https://next.hofstadter.io/mp/collect?"
	//gaURL = "https://next.hofstadter.io/debug/mp/collect?"
	//gaURL = "http://localhost:8080/mp/collect?"
	//gaURL = "http://localhost:8080/mp/collect?"
	url := gaURL +	vals.Encode()
	if debug {
		fmt.Println(url)
	}

	resp, err := http.Post(url, "application/json", reqBody)

	if debug {
		fmt.Println(resp, err)
	}

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
