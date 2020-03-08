package util

import (
	// "fmt"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	// "path/filepath"
	"strings"

	"github.com/aymerick/raymond"
)

// Used somewhere in here?
// https://blog.depado.eu/post/copy-files-and-directories-in-go [03-04-2-19]

func init() {
	raymond.RegisterHelper("pw", func(content, width string) string {
		return fmt.Sprintf("%-"+width+"s", content)
	})
}

func RenderString(template string, data interface{}) (string, error) {

	output, err := raymond.Render(template, data)
	if err != nil {
		return "", err
	}

	return output, nil
}

// File copies a single file from src to dst
func RenderFile(src, dst string, data interface{}) error {

	content, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	output, err := raymond.Render(string(content), data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, []byte(output), 0644)
	if err != nil {
		return err
	}

	return nil
}

// Dir copies a whole directory recursively
func RenderDir(src string, dst string, data interface{}) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	srcinfo, err = os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, srcinfo.Mode())
	if err != nil {
		return err
	}

	fds, err = ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			err = RenderDir(srcfp, dstfp, data)
			if err != nil {
				return err
			}
		} else {
			err = RenderFile(srcfp, dstfp, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func RenderFileNameSub(src, dst string, data interface{}) error {

	dst = subNames(dst, data)
	// fmt.Println(src, "->", filepath.Join(dst))

	c, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	content := string(c)

	swapped := SwapDelimits(
		content,
		"{{{", "}}}", "{{", "}}",
		"{%%", "%%}", "{%", "%}",
	)

	output, err := raymond.Render(swapped, data)
	if err != nil {
		return err
	}

	final := SwapDelimits(
		output,
		"{%%", "%%}", "{%", "%}",
		"{{{", "}}}", "{{", "}}",
	)

	err = ioutil.WriteFile(dst, []byte(final), 0644)
	if err != nil {
		return err
	}

	return nil
}

// Dir copies a whole directory recursively
func RenderDirNameSub(src string, dst string, data interface{}) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	srcinfo, err = os.Stat(src)
	if err != nil {
		return err
	}

	dst = subNames(dst, data)
	err = os.MkdirAll(dst, srcinfo.Mode())
	if err != nil {
		return err
	}

	fds, err = ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if fd.Name() == ".git" {
				continue
			}

			err = RenderDirNameSub(srcfp, dstfp, data)
			if err != nil {
				return err
			}
		} else {
			err = RenderFileNameSub(srcfp, dstfp, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

var names = []string{
	"AppName",
	"ModuleName",
	"TypeName",
	"PageName",
	"ComponentName",
	"FuncName",
}

func subNames(name string, data interface{}) string {
	ctx := data.(map[string]interface{})

	for _, N := range names {
		// what is the substitution?
		sub, ok := ctx[N]
		if ok {
			name = strings.Replace(name, N, sub.(string), -1)
		}
	}

	return name
}

func SwapDelimits(content, fromLL, fromLR, fromSL, fromSR, toLL, toLR, toSL, toSR string) string {
	tmpLL, tmpLR, tmpSL, tmpSR := "{@@", "@@}", "{@", "@}"

	// Change "to" -> "tmp"
	content = strings.Replace(content, toLL, tmpLL, -1)
	content = strings.Replace(content, toLR, tmpLR, -1)
	content = strings.Replace(content, toSL, tmpSL, -1)
	content = strings.Replace(content, toSR, tmpSR, -1)

	// Change "from" -> "to"
	content = strings.Replace(content, fromLL, toLL, -1)
	content = strings.Replace(content, fromLR, toLR, -1)
	content = strings.Replace(content, fromSL, toSL, -1)
	content = strings.Replace(content, fromSR, toSR, -1)

	// Change "to" (now tmp) -> "from"
	content = strings.Replace(content, tmpLL, fromLL, -1)
	content = strings.Replace(content, tmpLR, fromLR, -1)
	content = strings.Replace(content, tmpSL, fromSL, -1)
	content = strings.Replace(content, tmpSR, fromSR, -1)

	return content
}
