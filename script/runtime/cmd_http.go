package runtime

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

const HTTP2_GOAWAY_CHECK = "http2: server sent GOAWAY and closed the connection"

// http	makes an http call.
func (ts *Script) CmdHttp(neg int, args []string) {
	if len(args) < 1 {
		ts.Fatalf("usage: http function [args...]")
	}

	var err error
	ts.stdout, ts.stderr, ts.status, err = ts.http(args)
	if ts.stdout != "" {
		fmt.Fprintf(&ts.log, "[stdout]\n%s", ts.stdout)
	}
	if ts.stderr != "" {
		fmt.Fprintf(&ts.log, "[stderr]\n%s", ts.stderr)
	}
	if err == nil && neg > 0 {
		ts.Fatalf("unexpected http success")
	}

	if err != nil {
		fmt.Fprintf(&ts.log, "[%v]\n", err)
		if ts.ctxt.Err() != nil {
			ts.Fatalf("test timed out while making http request")
		} else if neg > 0 {
			ts.Fatalf("unexpected http failure")
		}
	}
}

// call runs the given function and then returns collected standard output and standard error.
func (ts *Script) http(args []string) (string, string, int, error) {
	// TODO, turn this into a log line
	if args[0] == "client" {
		err := ts.manageHttpClient(args[1:])
		ts.Check(err)
		return "", "", 0, nil
	}

	fmt.Println("args", args)

	req, err := ts.reqFromArgs(args)
	ts.Check(err)

	defer func() {
        if r := recover(); r != nil {
			bs := debug.Stack()
			ts.Fatalf("Recovered in HTTP: %v\n%s\n", r, string(bs))
        }
    }()

	resp, body, errs := req.End()

	if len(errs) != 0 && resp == nil {
		ts.Fatalf("Fatal error in HTTP: %v", errs)
	}

	if len(errs) != 0 && !strings.Contains(errs[0].Error(), HTTP2_GOAWAY_CHECK) {
		return "", body, resp.StatusCode, fmt.Errorf("Internal Weirdr Error:\b%v\n%s\n", errs, body)
	}
	if len(errs) != 0 {
		return "", body, resp.StatusCode, fmt.Errorf("Internal Error:\n%v\n%s\n", errs, body)
	}

	body += "\n"
	return body, "", resp.StatusCode, nil
}

func (ts *Script) manageHttpClient(args []string) error {
	L := len(args)
	if L < 1 {
		ts.Fatalf("usage: http client [new,del] <name> http-args...")
	}

	key, name := args[0], "default"
	if len(args) == 1 {
		args = args[1:]
	} else {
		name = args[1]
		args = args[2:]
	}

	if ts.httpClients == nil {
		ts.httpClients = make(map[string]*gorequest.SuperAgent)
	}

	switch key {
	case "new":
		req, err := ts.newReqFromArgs(args)
		ts.Check(err)
		ts.httpClients[name] = req

	case "mod":
		req, ok := ts.httpClients[name]
		if !ok {
			ts.Fatalf("unknown http client %q", name)
		}
		req, err := ts.applyArgsToReq(req, args)
		ts.Check(err)
		ts.httpClients[name] = req

	case "del":
		_, ok := ts.httpClients[name]
		if !ok {
			ts.Fatalf("unknown http client %q", name)
		}
		delete(ts.httpClients, name)

	default:
		ts.Fatalf("usage: http client <op> args...")
	}

	return nil
}

func (ts *Script) reqFromArgs(args []string) (*gorequest.SuperAgent, error) {
	// first arg is a known client
	if req, ok := ts.httpClients[args[0]]; ok {
		R := req.Clone()
		return ts.applyArgsToReq(R, args[1:])
	}
	return ts.newReqFromArgs(args)
}

func (ts *Script) newReqFromArgs(args []string) (*gorequest.SuperAgent, error) {
	// otherwise create a one-time req obj
	req := gorequest.New()
	req = ts.applyDefaultsToReq(req)
	return ts.applyArgsToReq(req, args)
}

func (ts *Script) applyDefaultsToReq(req *gorequest.SuperAgent) *gorequest.SuperAgent {

	req.Method = "GET"

	return req
}

func (ts *Script) applyArgsToReq(req *gorequest.SuperAgent, args []string) (*gorequest.SuperAgent, error) {
	var err error
	for _, arg := range args {
		req, err = ts.applyArgToReq(req, arg)
		if err != nil {
			return nil, err
		}
	}

	return req, nil
}

func (ts *Script) applyArgToReq(req *gorequest.SuperAgent, arg string) (*gorequest.SuperAgent, error) {
	// fmt.Printf("  APPLY: %q\n", flds)

	flds := strings.SplitN(arg, "=", 2)
	key := flds[0]
	val := ""
	if len(flds) == 2 {
		val = flds[1]
	}

	K := strings.ToUpper(key)

	switch K {
	case "U", "URL":
		req.Url = val

	case "T", "TYPE":
		req = req.Type(val)

	case "Q", "QUERY":
		if strings.HasPrefix(val, "@") {
			val = ts.ReadFile(val[1:])
		}
		req = req.Query(val)

	case "R", "RETRY":
		flds = strings.Fields(val)
		if len(flds) < 3 {
			ts.Fatalf("http retry usage: RETRY:'<count> <timer> [codes...]'")
		}
		cnt, tmr, codes := flds[0], flds[1], flds[2:]

		c, err := strconv.Atoi(cnt)
		ts.Check(err)

		t, err := time.ParseDuration(tmr)
		ts.Check(err)

		cs := []int{}
		for _, code := range codes {
			i, err := strconv.Atoi(code)
			ts.Check(err)
			cs = append(cs, i)
		}

		req = req.Retry(c, t, cs...)

	case "D", "DATA", "S", "SEND":
		if strings.HasPrefix(val, "@") {
			val = ts.ReadFile(val[1:])
		}
		req = req.Send(val)

	case "F", "FILE":
		flds := strings.Split(val, ":")
		filename, fieldname := strings.TrimSpace(flds[0]), ""
		if len(flds) > 1 {
			fieldname = strings.TrimSpace(flds[1])
		}
		content := ts.ReadFile(filename)
		req = req.SendFile([]byte(content), filename, fieldname)

	case "A", "AUTH":
		flds := strings.Split(val, ":")
		k, v := strings.TrimSpace(flds[0]), strings.TrimSpace(flds[1])
		req = req.SetBasicAuth(k, v)

	case "H", "HEADER":
		flds := strings.Split(val, ":")
		k, v := strings.TrimSpace(flds[0]), strings.TrimSpace(flds[1])
		req = req.Set(k, v)

	case "M", "METHOD":
		req.Method = K
	// Specially recognized key only args
	case "GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS":
		req.Method = K

	default:

		// check some special prefixes
		if strings.HasPrefix(key, "http") {
			req.Url = key
			return req, nil
		}

		return nil, fmt.Errorf("unknown http arg/key: %q / %q", arg, key)
	}

	return req, nil
}
