package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
	"github.com/hofstadter-io/hof/cmd/hof/verinfo"
)

var UpdateLong = `Print the build version for hof`

var (
	UpdateCheckFlag   bool
	UpdateVersionFlag string

	UpdateStarted   bool
	UpdateErrored   bool
	UpdateChecked   bool
	UpdateAvailable *ProgramVersion
	UpdateData      []interface{}
)

func init() {
	UpdateCmd.Flags().BoolVarP(&UpdateCheckFlag, "check", "", false, "only check for an update")
	UpdateCmd.Flags().BoolVarP(&UpdateCheckFlag, "list", "", false, "show local installed version")
	UpdateCmd.Flags().StringVarP(&UpdateVersionFlag, "version", "V", "", "the version to update or install")
}

const updateMessage = `
Updates available. v%s -> %s

  run 'hof update' to get the latest.

`

// TODO, add a curl to the above? or os specific?

var UpdateCmd = &cobra.Command{

	Use: "update",

	Short: "update the hof tool",

	Long: UpdateLong,

	PreRun: func(cmd *cobra.Command, args []string) {
		ga.SendGaEvent("update", "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {

		latest, err := CheckUpdate(true)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		// Semver Check?
		cur := ProgramVersion{Version: "v" + verinfo.Version}
		if latest.Version == cur.Version || (UpdateVersionFlag == "" && cur.Version == "vLocal") {
			return
		} else {
			if UpdateCheckFlag {
				PrintUpdateAvailable()
				return
			}
		}

		err = InstallUpdate()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	},
}

func init() {
	go CheckUpdate(false)

	help := UpdateCmd.HelpFunc()
	usage := UpdateCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		if UpdateCmd.Name() == cmd.Name() {
			ga.SendGaEvent("update/help", "<omit>", 0)
		}
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		if UpdateCmd.Name() == cmd.Name() {
			ga.SendGaEvent("update/usage", "<omit>", 0)
		}
		return usage(cmd)
	}
	UpdateCmd.SetHelpFunc(thelp)
	UpdateCmd.SetUsageFunc(tusage)

}

type ProgramVersion struct {
	Version string
	URL     string
}

func CheckUpdate(manual bool) (ver ProgramVersion, err error) {
	if !manual && os.Getenv("HOF_UPDATES_DISABLED") != "" {
		return
	}
	UpdateStarted = true
	cur := ProgramVersion{Version: "v" + verinfo.Version}

	checkURL := "https://api.github.com/repos/hofstadter-io/hof/releases/latest"
	if UpdateVersionFlag != "" {
		checkURL = "https://api.github.com/repos/hofstadter-io/hof/releases/tags/" + UpdateVersionFlag
		manual = true
	}

	req := gorequest.New()
	if os.Getenv("GITHUB_TOKEN") != "" {
		req = req.SetBasicAuth("github-token", os.Getenv("GITHUB_TOKEN"))
	}
	resp, b, errs := req.Get(checkURL).EndBytes()
	UpdateErrored = true

	check := "http2: server sent GOAWAY and closed the connection"
	if len(errs) != 0 && !strings.Contains(errs[0].Error(), check) {
		// fmt.Println("errs:", errs)
		return ver, errs[0]
	}

	if len(errs) != 0 || resp.StatusCode >= 500 {
		return ver, fmt.Errorf("Internal Error: " + string(b))
	}
	if resp.StatusCode >= 400 {
		if resp.StatusCode == 404 {
			fmt.Println("404?!", checkURL)
			return ver, fmt.Errorf("No releases available :[")
		}
		return ver, fmt.Errorf("Bad Request: " + string(b))
	}

	UpdateErrored = false
	// fmt.Println(string(b))

	var gh map[string]interface{}
	err = json.Unmarshal(b, &gh)
	if err != nil {
		return ver, err
	}

	nameI, ok := gh["name"]
	if !ok {
		return ver, fmt.Errorf("Internal Error: could not find version in update check response")
	}
	name, ok := nameI.(string)
	if !ok {
		return ver, fmt.Errorf("Internal Error: version is not a string in update check response")
	}
	ver.Version = name

	if !manual {
		UpdateChecked = true

		// Semver Check?
		if ver.Version != cur.Version && cur.Version != "vLocal" {
			UpdateAvailable = &ver
		}

		return ver, nil
	}

	// This goes here and signals else where that we got the request back
	UpdateChecked = true

	// Semver Check?
	if ver.Version != cur.Version && (manual || cur.Version != "vLocal") {
		UpdateAvailable = &ver
		aI, ok := gh["assets"]
		if ok {
			a, aok := aI.([]interface{})
			if aok {
				UpdateData = a
			}
		}
	}

	return ver, nil
}

func WaitPrintUpdateAvailable() {
	if os.Getenv("HOF_UPDATES_DISABLED") != "" {
		return
	}
	for i := 0; i < 20 && !UpdateStarted && !UpdateChecked && !UpdateErrored; i++ {
		time.Sleep(50 * time.Millisecond)
	}
	PrintUpdateAvailable()
}

func PrintUpdateAvailable() {
	if UpdateChecked && UpdateAvailable != nil {
		fmt.Printf(updateMessage, verinfo.Version, UpdateAvailable.Version)
	}
}

func InstallUpdate() (err error) {
	fmt.Printf("Installing hof@%s\n", UpdateAvailable.Version)

	if UpdateData == nil {
		return fmt.Errorf("No update available")
	}
	/*
		vers, err := json.MarshalIndent(UpdateData, "", "  ")
		if err == nil {
			fmt.Println(string(vers))
		}
	*/

	fmt.Println("OS/Arch", verinfo.BuildOS, verinfo.BuildArch)

	url := ""
	for _, Asset := range UpdateData {
		asset := Asset.(map[string]interface{})
		U := asset["browser_download_url"].(string)
		u := strings.ToLower(U)

		osOk, archOk := false, false

		switch verinfo.BuildOS {
		case "linux":
			if strings.Contains(u, "linux") {
				osOk = true
			}

		case "darwin":
			if strings.Contains(u, "darwin") {
				osOk = true
			}

		case "windows":
			if strings.Contains(u, "windows") {
				osOk = true
			}
		}

		switch verinfo.BuildArch {
		case "amd64":
			if strings.Contains(u, "x86_64") {
				archOk = true
			}
		case "arm64":
			if strings.Contains(u, "arm64") {
				archOk = true
			}
		case "arm":
			if strings.Contains(u, "arm") && !strings.Contains(u, "arm64") {
				archOk = true
			}
		}

		if osOk && archOk {
			url = u
			break
		}
	}

	fmt.Println("Download URL: ", url, "\n")

	switch verinfo.BuildOS {
	case "linux":
		fallthrough
	case "darwin":

		return downloadAndInstall(url)

	case "windows":
		fmt.Println("Please downlaod and install manually from the link above.\n")
		return nil
	}

	return nil
}

func downloadAndInstall(url string) error {
	req := gorequest.New()
	if os.Getenv("GITHUB_TOKEN") != "" {
		req = req.SetBasicAuth("github-token", os.Getenv("GITHUB_TOKEN"))
	}

	resp, content, errs := req.Get(url).EndBytes()

	check := "http2: server sent GOAWAY and closed the connection"
	if len(errs) != 0 && !strings.Contains(errs[0].Error(), check) {
		fmt.Println("errs:", errs)
		fmt.Println("resp:", resp)
		return errs[0]
	}

	if len(errs) != 0 || resp.StatusCode >= 400 {
		return fmt.Errorf("Error %v - %s", resp.StatusCode, string(content))
	}

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		return err
	}

	// defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		return err
	}
	if err := tmpfile.Close(); err != nil {
		return err
	}

	ex, err := os.Executable()
	if err != nil {
		return err
	}

	real, err := filepath.EvalSymlinks(ex)
	if err != nil {
		return err
	}

	// Sudo copy the file
	cmd := exec.Command("/bin/sh", "-c",
		fmt.Sprintf("export OWNER=$(ls -l %s | awk '{ print $3 \":\" $4 }') && sudo mv %s %s-v%s && sudo cp %s %s && sudo chown $OWNER %s && sudo chmod 0755 %s",
			real,                        // get owner
			real, real, verinfo.Version, // backup
			tmpfile.Name(), real, // cp
			real, // chown
			real, // chmod
			real, // rm
		),
	)

	// prep stdin for password
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")
	}()

	stdoutStderr, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", stdoutStderr)
	if err != nil {
		return err
	}

	UpdateAvailable = nil
	UpdateData = nil
	return nil
}
