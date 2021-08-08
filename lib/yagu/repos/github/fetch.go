package github

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/google/go-github/v30/github"
	"github.com/parnurzeal/gorequest"

	"github.com/hofstadter-io/hof/lib/yagu"
)

func Fetch(FS billy.Filesystem, owner, repo, tag string) (error) {
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

		zReader, err = FetchBranchZip(client, owner, repo, *r.DefaultBranch)
		if err != nil {
			return fmt.Errorf("While fetching branch zipfile for %s/%s@%s\n%w\n", owner, repo, *r.DefaultBranch, err)
		}
	} else {
		tags, err := GetTags(client, owner, repo)
		if err != nil {
			return err
		}

		// The tag we are looking for
		var T *github.RepositoryTag
		for _, t := range tags {
			if tag != "" && tag == *t.Name {
				T = t
				// fmt.Printf("FOUND  %v\n", *t.Name)
			}
		}

		if T == nil {
			return fmt.Errorf("Did not find tag %q for 'https://github.com/%s/%s' @%s", tag, owner, repo, tag)
		}

		zReader, err = FetchTagZip(client, T)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return fmt.Errorf("While fetching tag zipfile\n%w\n", err)
		}
	}

	if err != nil {
		return fmt.Errorf("While fetching from github\n%w\n", err)
	}

	if err := yagu.BillyLoadFromZip(zReader, FS, true); err != nil {
		return fmt.Errorf("While reading zipfile\n%w\n", err)
	}

	return nil
}

func FetchTagZip(client *github.Client, tag *github.RepositoryTag) (*zip.Reader, error) {

	url := *tag.ZipballURL

	fmt.Println("url:", url)

	req := gorequest.New().Get(url)
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

func FetchBranchZip(client *github.Client, owner, repo, branch string) (*zip.Reader, error) {

	url := fmt.Sprintf("https://github.com/%s/%s/archive/refs/heads/%s.zip", owner, repo, branch)

	req := gorequest.New().Get(url)
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
