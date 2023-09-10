package playground

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/parnurzeal/gorequest"
)

const HTTP2_GOAWAY_CHECK = "http2: server sent GOAWAY and closed the connection"

func (C *Playground) PushToPlayground() (string, error) {
	src := C.edit.GetText()

	url := "https://cuelang.org/.netlify/functions/snippets"
	req := gorequest.New().Post(url)
	req.Set("Content-Type", "text/plain")
	req.Send(src)

	resp, body, errs := req.End()

	if len(errs) != 0 && !strings.Contains(errs[0].Error(), HTTP2_GOAWAY_CHECK) {
		fmt.Println("errs:", errs)
		fmt.Println("resp:", resp)
		fmt.Println("body:", body)
		return body, errs[0]
	}

	if len(errs) != 0 || resp.StatusCode >= 500 {
		return body, fmt.Errorf("Internal Error: " + body)
	}
	if resp.StatusCode >= 400 {
		return body, fmt.Errorf("Bad Request: " + body)
	}

	return body, nil
}

func (C *Playground) WriteEditToFile(filename string) (error) {
	src := C.edit.GetText()

	return os.WriteFile(filename, []byte(src), 0644)
}

func (C *Playground) ExportFinalToFile(filename string) (error) {
	ext := filepath.Ext(filename)
	ext = strings.TrimPrefix(ext, ".")
	src, err := C.final.viewer.GetValueText(ext)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(src), 0644)
}
