package api

import (
	"fmt"
	"io"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"github.com/parnurzeal/gorequest"

	hofcontext "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

/*  TODO
    - catch / retry on failed connection
*/

var debug = false

type call struct {
}

type Call struct{}

func NewCall(val cue.Value) (hofcontext.Runner, error) {
	return &Call{}, nil
}

func (T *Call) Run(ctx *hofcontext.Context) (interface{}, error) {
	val := ctx.Value
	init_schemas(val.Context())
	// unify with schema
	val = val.Unify(task_call)
	if val.Err() != nil {
		return nil, cuetils.ExpandCueError(val.Err())
	}

	var R *gorequest.SuperAgent
	var err error

	func() error {
		ctx.CUELock.Lock()
		defer func() {
			ctx.CUELock.Unlock()
		}()

		req := val.LookupPath(cue.ParsePath("req"))

		R, err = buildRequest(req)
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		return nil, err
	}

	resp, err := makeRequest(R)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var bodyVal interface{}

	func() {
		ctx.CUELock.Lock()
		defer func() {
			ctx.CUELock.Unlock()
		}()

		// TODO, build resp cue.Value from http.Response

		var isString, isBytes bool
		r := val.LookupPath(cue.ParsePath("resp.body"))
		if r.Exists() {
			if r.IncompleteKind() == cue.StringKind {
				isString = true
			}
			if r.IncompleteKind() == cue.BytesKind {
				isBytes = true
			}
		}

		// TODO, make response object more interesting
		// such as status, headers, body vs json
		if isString {
			bodyVal = string(body)
		} else if isBytes {
			bodyVal = body
		} else {
			bodyVal = val.Context().CompileBytes(body, cue.Filename("body"))
		}
	}()

	return map[string]interface{}{
		"resp": map[string]interface{}{
			"status":     resp.Status,
			"statusCode": resp.StatusCode,
			"body":       bodyVal,
			"header":     resp.Header,
			"trailer":    resp.Trailer,
		},
	}, nil
}

/********* old *********/

const HTTP2_GOAWAY_CHECK = "http2: server sent GOAWAY and closed the connection"

func buildRequest(val cue.Value) (R *gorequest.SuperAgent, err error) {
	req := val.Eval()
	R = gorequest.New()

	// this should be the default after unifying, and this should always exist
	// R.Method = "GET"
	method := req.LookupPath(cue.ParsePath("method"))
	// so we may not need this is not needed, but it is defensive
	if method.Exists() {
		R.Method, err = method.String()
		if err != nil {
			return R, err
		}
	}

	host := req.LookupPath(cue.ParsePath("host"))
	hostStr, err := host.String()
	if err != nil {
		return
	}

	path := req.LookupPath(cue.ParsePath("path"))
	pathStr, err := path.String()
	if err != nil {
		return
	}
	R.Url = hostStr + pathStr

	headers := req.LookupPath(cue.ParsePath("headers"))
	if headers.Exists() {
		H, err := headers.Struct()
		if err != nil {
			return R, err
		}
		hIter := H.Fields()
		for hIter.Next() {
			label := hIter.Label()
			value, err := hIter.Value().String()
			if err != nil {
				return R, err
			}
			R.Header.Add(label, value)
		}
	}

	query := req.LookupPath(cue.ParsePath("query"))
	if query.Exists() {
		Q, err := query.Struct()
		if err != nil {
			return R, err
		}
		qIter := Q.Fields()
		for qIter.Next() {
			label := qIter.Label()
			value, err := qIter.Value().String()
			if err != nil {
				return R, err
			}
			R.QueryData.Add(label, value)
		}
	}

	form := req.LookupPath(cue.ParsePath("form"))
	if form.Exists() {
		F, err := form.Struct()
		if err != nil {
			return R, err
		}
		fIter := F.Fields()
		for fIter.Next() {
			label := fIter.Label()
			value, err := fIter.Value().String()
			if err != nil {
				return R, err
			}
			R.FormData.Add(label, value)
		}
	}

	data := req.LookupPath(cue.ParsePath("data"))
	if data.Exists() {
		d := map[string]any{}
		err := data.Decode(&d)
		if err != nil {
			return R, err
		}
		R = R.Send(d)
	}

	timeout := req.LookupPath(cue.ParsePath("timeout"))
	if timeout.Exists() {
		to, err := timeout.String()
		if err != nil {
			return R, err
		}
		d, err := time.ParseDuration(to)
		if err != nil {
			return R, err
		}
		R.Timeout(d)
	}

	retry := req.LookupPath(cue.ParsePath("retry"))
	if retry.Exists() {
		C := 3
		count := retry.LookupPath(cue.ParsePath("count"))
		if count.Exists() {
			c, err := count.Int64()
			if err != nil {
				return R, err
			}
			C = int(c)
		}

		D := time.Second * 6
		timer := retry.LookupPath(cue.ParsePath("timer"))
		if timer.Exists() {
			t, err := timer.String()
			if err != nil {
				return R, err
			}
			d, err := time.ParseDuration(t)
			if err != nil {
				return R, err
			}
			D = d
		}

		CS := []int{}
		codes := retry.LookupPath(cue.ParsePath("codes"))
		if codes.Exists() {
			L, err := codes.List()
			if err != nil {
				return R, err
			}
			for L.Next() {
				v, err := L.Value().Int64()
				if err != nil {
					return R, err
				}
				CS = append(CS, int(v))
			}
		}

		R.Retry(C, D, CS...)
	}

	// Todo, add 'print curl' option
	// check if set, then fill on return
	//curl, err := R.AsCurlCommand()
	//if err != nil {
	//return nil, err
	//}
	//fmt.Println(curl)

	return
}

func makeRequest(R *gorequest.SuperAgent) (gorequest.Response, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered in HTTP: %v %v\n", R, r)
		}
	}()

	if debug {
		s, err := R.Clone().AsCurlCommand()
		if err != nil {
			return nil, err
		}
		fmt.Println("CURL:", s)
	}

	resp, body, errs := R.End()

	if len(errs) != 0 && resp == nil {
		return resp, fmt.Errorf("%v", errs)
	}

	if len(errs) != 0 && !strings.Contains(errs[0].Error(), HTTP2_GOAWAY_CHECK) {
		return resp, fmt.Errorf("Internal Weirdr Error:\b%v\n%s\n", errs, body)
	}
	if len(errs) != 0 {
		return resp, fmt.Errorf("Internal Error:\n%v\n%s\n", errs, body)
	}

	return resp, nil
}
