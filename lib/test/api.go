package test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/encoding/json"
	"github.com/parnurzeal/gorequest"

	"github.com/hofstadter-io/hof/lib/cuetils"
)

const HTTP2_GOAWAY_CHECK = "http2: server sent GOAWAY and closed the connection"

func RunAPI(T *Tester, verbose int) (err error) {
	fmt.Println("api:", T.Name)

	// make sure we resolve references and unifications
	val := T.Value.Eval()

	vSyn, vErr := cuetils.ValueToSyntaxString(val)
	if vErr != nil {
		fmt.Println(vSyn)
		return vErr
	}
	if verbose > 0 {
		fmt.Println(vSyn)
	}

	return runCase(T, verbose, val)
}

func runCase(T *Tester, verbose int, val cue.Value) (err error) {

	req := val.LookupPath(cue.ParsePath("req"))
	expected := val.LookupPath(cue.ParsePath("resp"))

	R, err := buildRequest(T, verbose, req)
	if err != nil {
		return err
	}

	actual, err := makeRequest(T, verbose, R)
	if err != nil {
		return err
	}

	err = checkResponse(T, verbose, actual, expected)

	return err
}


func buildRequest(T *Tester, verbose int, val cue.Value) (R *gorequest.SuperAgent, err error) {
	req := val.Eval()
	R = gorequest.New()

	method := req.LookupPath(cue.ParsePath("method"))
	R.Method, err = method.String()
	if err != nil {
		return
	}

	host := req.LookupPath(cue.ParsePath("host"))
	path := req.LookupPath(cue.ParsePath("path"))
	hostStr, err := host.String()
	if err != nil {
		return
	}
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

	data := req.LookupPath(cue.ParsePath("data"))
	if data.Exists() {
		err := data.Decode(&R.Data)
		if err != nil {
			return R, err
		}
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

	return
}

func makeRequest(T *Tester, verbose int, R *gorequest.SuperAgent) (gorequest.Response, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in HTTP: %v %v\n", R, r)
		}
	}()

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

	if verbose > 0 {
		fmt.Println(body)
	}
	fmt.Println(body)

	return resp, nil
}

func checkResponse(T *Tester, verbose int, actual gorequest.Response, expect cue.Value) (err error) {
	expect = expect.Eval()

	S, err := expect.Struct()
	if err != nil {
		return err
	}
	iter := S.Fields()
	for iter.Next() {
		label := iter.Label()
		value := iter.Value()

		fmt.Println("checking:", label)

		switch label {
			case "status":
				status, err := value.Int64()
				if err != nil {
					return err
				}
				if int64(actual.StatusCode) != status {
					return fmt.Errorf("status code mismatch %v != %v", actual.StatusCode, status)
				}

			case "body":
				body, err := ioutil.ReadAll(actual.Body)
				if err != nil {
					return err
				}

				inst, err := json.Decode(T.CRT, "", body)
				if err != nil {
					return err
				}

				V := inst.Value()
				result := value.Unify(V)
				if result.Err() != nil {
					return result.Err()
				}
				fmt.Println("result: ", result)
				err = result.Validate()
				if err != nil {
					fmt.Println(value)
					fmt.Println(inst.Value())

					return err
				}


			default:
				return fmt.Errorf("Unknown field in expected response:", label)
		}
	}


	return nil
}
