package runtime

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"

	"github.com/hofstadter-io/hof/script/ast"
)

const HTTP2_GOAWAY_CHECK = "http2: server sent GOAWAY and closed the connection"

// Cmd_http	makes an http call.
func (RT *Runtime) Cmd_http(cmd *ast.Cmd, r *ast.Result) (err error) {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: http client [args...]")
	}

	args := cmd.Args

	if args[0] == "client" {
		return RT.manageHttpClient(args[1:])
	}

	req, err := RT.reqFromArgs(args)
	RT.Check(err)

	defer func() {
		if r := recover(); r != nil {
			bs := debug.Stack()
			RT.logger.Errorf("Recovered in HTTP: %v\n%s\n", r, string(bs))
		}
	}()

	resp, body, errs := req.End()

	r.Status = resp.StatusCode
	RT.status = r.Status

	if len(errs) != 0 && resp == nil {
		fmt.Fprintln(RT.stderr, body)
		return fmt.Errorf("Fatal error in HTTP: %v", errs)
	}
	if len(errs) != 0 && !strings.Contains(errs[0].Error(), HTTP2_GOAWAY_CHECK) {
		fmt.Fprintln(RT.stderr, body)
		return fmt.Errorf("Internal Weirdr Error:\b%v\n%s\n", errs, body)
	}
	if len(errs) != 0 {
		fmt.Fprintln(RT.stderr, body)
		return fmt.Errorf("Internal Error:\n%v\n%s\n", errs, body)
	}

	fmt.Fprintln(RT.stdout, body)

	return nil
}

func (RT *Runtime) manageHttpClient(args []string) error {
	L := len(args)
	if L < 1 {
		RT.Fatalf("usage: http client [new,del] <name> http-args...")
	}

	key, name := args[0], "default"
	if len(args) == 1 {
		args = args[1:]
	} else {
		name = args[1]
		args = args[2:]
	}

	if RT.httpClients == nil {
		RT.httpClients = make(map[string]*gorequest.SuperAgent)
	}

	switch key {
	case "new":
		req, err := RT.newReqFromArgs(args)
		RT.Check(err)
		RT.httpClients[name] = req

	case "mod":
		req, ok := RT.httpClients[name]
		if !ok {
			RT.Fatalf("unknown http client %q", name)
		}
		req, err := RT.applyArgsToReq(req, args)
		RT.Check(err)
		RT.httpClients[name] = req

	case "del":
		_, ok := RT.httpClients[name]
		if !ok {
			RT.Fatalf("unknown http client %q", name)
		}
		delete(RT.httpClients, name)

	default:
		RT.Fatalf("usage: http client <op> args...")
	}

	return nil
}

func (RT *Runtime) reqFromArgs(args []string) (*gorequest.SuperAgent, error) {
	// first arg is a known client
	if req, ok := RT.httpClients[args[0]]; ok {
		R := req.Clone()
		return RT.applyArgsToReq(R, args[1:])
	}
	return RT.newReqFromArgs(args)
}

func (RT *Runtime) newReqFromArgs(args []string) (*gorequest.SuperAgent, error) {
	// otherwise create a one-time req obj
	req := gorequest.New()
	req = RT.applyDefaultsToReq(req)
	return RT.applyArgsToReq(req, args)
}

func (RT *Runtime) applyDefaultsToReq(req *gorequest.SuperAgent) *gorequest.SuperAgent {

	req.Method = "GET"

	return req
}

func (RT *Runtime) applyArgsToReq(req *gorequest.SuperAgent, args []string) (*gorequest.SuperAgent, error) {
	var err error
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if i + 2 < len(args) && args[i+1] == "=" {
			k, v := arg, args[i+2]
			i += 2
			if v[0] == '@' {
				fname := v[1:] // for error messages
				if fname == "stdout" {
					v = RT.GetStdout()
				} else if fname == "stderr" {
					v = RT.GetStderr()
				} else {
					data := RT.ReadFile(fname)
					v = string(data)
				}
			} else {
				v = RT.expand(v)
			}

			arg = fmt.Sprintf("%s=%s", k, v)
		}

		req, err = RT.applyArgToReq(req, arg)
		if err != nil {
			return nil, err
		}
	}

	return req, nil
}

func (RT *Runtime) applyArgToReq(req *gorequest.SuperAgent, arg string) (*gorequest.SuperAgent, error) {
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
			val = RT.ReadFile(val[1:])
		}
		req = req.Query(val)

	case "R", "RETRY":
		flds = strings.Fields(val)
		if len(flds) < 3 {
			RT.Fatalf("http retry usage: RETRY:'<count> <timer> [codes...]'")
		}
		cnt, tmr, codes := flds[0], flds[1], flds[2:]

		c, err := strconv.Atoi(cnt)
		RT.Check(err)

		t, err := time.ParseDuration(tmr)
		RT.Check(err)

		cs := []int{}
		for _, code := range codes {
			i, err := strconv.Atoi(code)
			RT.Check(err)
			cs = append(cs, i)
		}

		req = req.Retry(c, t, cs...)

	case "D", "DATA", "S", "SEND":
		if strings.HasPrefix(val, "@") {
			val = RT.ReadFile(val[1:])
		}
		req = req.Send(val)

	case "F", "FILE":
		flds := strings.Split(val, ":")
		filename, fieldname := strings.TrimSpace(flds[0]), ""
		if len(flds) > 1 {
			fieldname = strings.TrimSpace(flds[1])
		}
		content := RT.ReadFile(filename)
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
		req.Method = val
	// Specially recognized key only args
	case "GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS":
		req.Method = K


	// Graphql Helpers
	case "G", "GQL", "GRAPHQL":
		if strings.HasPrefix(val, "@") {
			val = RT.ReadFile(val[1:])
		}
		req = req.Send(fmt.Sprintf(`{ "query": %q }`, val))

	case "V", "VARS", "VARIABLES":
		if strings.HasPrefix(val, "@") {
			val = RT.ReadFile(val[1:])
		}
		req = req.Send(fmt.Sprintf(`{ "variables": %s }`, val))

	case "O", "OP", "OPNAME", "OPERATION":
		req = req.Send(fmt.Sprintf(`{ "operationName": %q }`, val))


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

/*
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
*/
