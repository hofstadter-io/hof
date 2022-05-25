package bbc

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/ktrysmt/go-bitbucket"
	"github.com/parnurzeal/gorequest"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/lib/yagu/repos/git"
)

func Fetch(FS billy.Filesystem, owner, repo, tag string, private bool) error {

	// If private, and no token auth, try git protocol
	// need to catch auth errors and suggest how to setup
	if private && os.Getenv(TokenEnv) == "" {
		fmt.Println("bitbucket git fallback")
		return git.Fetch(FS, "bitbucket.org", owner, repo, tag, private)
	}

	client, err := NewClient()
	if err != nil {
		return err
	}

	var zReader *zip.Reader

	if tag == "v0.0.0" {
		r, err := GetRepo(client, owner, repo)
		if err != nil {
			return err
		}

		zReader, err = FetchBranchZip(owner, repo, r.Mainbranch.Name)
		if err != nil {
			return fmt.Errorf("While fetching branch zipfile for %s/%s@%s\n%w\n", owner, repo, r.Mainbranch.Name, err)
		}

	} else {
		tags, err := GetTags(client, owner, repo)
		if err != nil {
			return err
		}

		// The tag we are looking for
		var T *bitbucket.RepositoryTag
		for _, t := range tags {
			if tag != "" && tag == t.Name {
				T = &t
				// fmt.Printf("FOUND  %v\n", *t.Name)
			}
		}

		if T == nil {
			return fmt.Errorf("Did not find tag %q for 'https://bitbucket.org/%s/%s' @%s", tag, owner, repo, tag)
		}

		zReader, err = FetchTagZip(owner, repo, T.Name)
		if err != nil {
			return fmt.Errorf("While fetching tag zipfile\n%w\n", err)
		}
	}

	if err != nil {
		return fmt.Errorf("While fetching from bitbucket\n%w\n", err)
	}

	if err := yagu.BillyLoadFromZip(zReader, FS, true); err != nil {
		return fmt.Errorf("While reading zipfile\n%w\n", err)
	}

	return nil
}

func FetchTagZip(owner, repo, ver string) (*zip.Reader, error) {

	// url := *tag.ZipballURL
	url := fmt.Sprintf("https://bitbucket.org/%s/%s/get/%s.zip", owner, repo, ver)

	req := gorequest.New().Get(url)

	// TODO, process auth logic here better, maybe find a way to DRY
	if token := os.Getenv(TokenEnv); token != "" {
		req.Header.Set("Authorization", "Bearer"+token)
	} else if password := os.Getenv(PasswordEnv); password != "" {
		username := os.Getenv(UsernameEnv)
		req.SetBasicAuth(username, password)
	}

	resp, data, errs := req.EndBytes()

	check := "http2: server sent GOAWAY and closed the connection"
	if len(errs) != 0 && !strings.Contains(errs[0].Error(), check) {
		fmt.Println("errs:", errs)
		fmt.Println("resp:", resp)
		fmt.Println("body:", len(data))
		return nil, errs[0]
	}

	if len(errs) != 0 || resp.StatusCode >= 500 {
		return nil, fmt.Errorf("Internal Error: " + string(resp.StatusCode))
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Bad Request: %v\n%v\n", resp.StatusCode, errs)
	}

	r := bytes.NewReader(data)

	zfile, err := zip.NewReader(r, int64(len(data)))

	return zfile, err
}

func FetchBranchZip(owner, repo, branch string) (*zip.Reader, error) {

	url := fmt.Sprintf("https://bitbucket.org/%s/%s/get/%s.zip", owner, repo, branch)

	req := gorequest.New().Get(url)

	// TODO, process auth logic here better, maybe find a way to DRY
	if token := os.Getenv(TokenEnv); token != "" {
		req.Header.Set("Authorization", "Bearer"+token)
	} else if password := os.Getenv(PasswordEnv); password != "" {
		username := os.Getenv(UsernameEnv)
		req.SetBasicAuth(username, password)
	}

	resp, data, errs := req.EndBytes()

	check := "http2: server sent GOAWAY and closed the connection"
	if len(errs) != 0 && !strings.Contains(errs[0].Error(), check) {
		fmt.Println("errs:", errs)
		fmt.Println("resp:", resp)
		fmt.Println("body:", len(data))
		return nil, errs[0]
	}

	if len(errs) != 0 || resp.StatusCode >= 500 {
		return nil, fmt.Errorf("Internal Error: " + string(resp.StatusCode))
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Bad Request: " + string(resp.StatusCode))
	}

	r := bytes.NewReader(data)

	zfile, err := zip.NewReader(r, int64(len(data)))

	return zfile, err
}
