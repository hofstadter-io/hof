package flags

import (
	"bytes"
	"fmt"
	"strings"
)

func PrintSubject(title, prefix, subject string, subjects map[string]string) bool {
	// skip if null or empty
	if subjects == nil || len(subjects) == 0 {
		return false
	}

	// print keys for list, so don't have a subject that name
	if subject == "list" {
		var S []string
		for k, _ := range subjects {
			S = append(S, k)
		}
		var b bytes.Buffer
		fmt.Fprintln(&b, title)
		for _, s := range S {
			fmt.Fprintln(&b, prefix+s)
		}

		fmt.Println(b.String())
		return true
	}

	// print pubject, indenting all lines
	S, ok := subjects[subject]
	if !ok {
		return false
	}
	S = strings.Replace(S, "¡", "`", -1)

	var b bytes.Buffer
	fmt.Fprintln(&b, title)
	for _, s := range strings.Split(S, "\n") {
		fmt.Fprintln(&b, prefix+s)
	}

	fmt.Println(b.String())
	return true
}

type RootPflagpole struct {
	Labels             []string
	Config             string
	Secret             string
	ContextFile        string
	Context            string
	Global             bool
	Local              bool
	Input              []string
	InputFormat        string
	Output             []string
	OutputFormat       string
	Error              []string
	ErrorFormat        string
	Account            string
	Billing            string
	Project            string
	Workspace          string
	DatamodelDir       string
	ResourcesDir       string
	RuntimesDir        string
	Package            string
	Errors             bool
	Ignore             bool
	Simplify           bool
	Trace              bool
	Strict             bool
	Verbose            string
	Quiet              bool
	ImpersonateAccount string
	TraceToken         string
	LogHTTP            string
	NoColor            bool
	Topic              string
	Example            string
	Tutorial           string
}

var RootPflags RootPflagpole
