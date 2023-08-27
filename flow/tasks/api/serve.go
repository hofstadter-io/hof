package api

import (
	gocontext "context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"cuelang.org/go/cue"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	hofcontext "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/flow/flow"
	"github.com/hofstadter-io/hof/flow/tasks/csp"
)

type Serve struct {
	sync.Mutex
}

func NewServe(val cue.Value) (hofcontext.Runner, error) {
	return &Serve{}, nil
}

func (T *Serve) Run(ctx *hofcontext.Context) (interface{}, error) {
	var err error

	// todo, check failure modes, fill, not return error?
	// (in all tasks)
	// do failed http handlings fail the client connection and server flow?

	val := ctx.Value

	var e *echo.Echo
	port := "2323"
	quit := ""

	err = func() error {
		ctx.CUELock.Lock()
		defer func() {
			ctx.CUELock.Unlock()
		}()

		logging := false
		l := val.LookupPath(cue.ParsePath("logging"))
		if l.Exists() {
			if l.Err() != nil {
				return l.Err()
			}
			logging, err = l.Bool()
			if err != nil {
				return err
			}
		}

		// get the port
		p := val.LookupPath(cue.ParsePath("port"))
		if p.Err() != nil {
			return p.Err()
		}
		port, err = p.String()
		if err != nil {
			return err
		}

		q := val.LookupPath(cue.ParsePath("quitMailbox"))
		if q.Exists() {
			if q.Err() != nil {
				return q.Err()
			}
			quit, err = q.String()
			if err != nil {
				return err
			}
		}

		// create server
		e = echo.New()
		e.HideBanner = true
		e.Use(middleware.Recover())
		if logging {
			e.Use(middleware.Logger())
		}

		// liveliness and metrics
		e.GET("/alive", func(c echo.Context) error {
			return c.NoContent(http.StatusNoContent)
		})

		// should we even have this route?
		if quit != "" {
			e.GET("/quit", func(c echo.Context) error {
				fmt.Println("quit handler!")
				fmt.Println("quitMailbox?:", quit)
				ci, loaded := ctx.Mailbox.Load(quit)
				if !loaded {
					return fmt.Errorf("channel %q not found", quit)
				}
				quitChan := ci.(chan csp.Msg)
				quitChan <- csp.Msg{Key: "quit"}
				return c.NoContent(http.StatusNoContent)
			})
		}

		prom := prometheus.NewPrometheus("echo", nil)
		prom.Use(e)

		//
		// Setup routes
		//
		routes := val.LookupPath(cue.ParsePath("routes"))
		iter, err := routes.Fields()
		if err != nil {
			return err
		}

		for iter.Next() {
			label := iter.Selector().String()
			route := iter.Value()

			// fmt.Println("route:", label)

			err := T.routeFromValue(label, route, e, ctx)
			if err != nil {
				fmt.Println("Error building route:", err)
				return err
			}
		}

		// put behind value field
		// print routes
		data, err := json.MarshalIndent(e.Routes(), "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(data))
		return nil
	}()

	// check return of our ad-hoc func scope
	if err != nil {
		return nil, err
	}

	if quit == "" {
		// run the server
		e.Logger.Fatal(e.Start(":" + port))
	} else {
		// load mailbox
		fmt.Println("quitMailbox?:", quit)
		qi, loaded := ctx.Mailbox.Load(quit)
		if !loaded {
			return nil, fmt.Errorf("channel %q not found", quit)
		}
		quitChan := qi.(chan csp.Msg)

		// Start server
		go func() {
			if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
				e.Logger.Fatal("shutting down the server")
			}
		}()

		// signal.Notify(quit, os.Interrupt)
		fmt.Println("waiting for quit...", loaded)
		<-quitChan
		fmt.Println("quit'n time!")

		ctx, cancel := gocontext.WithTimeout(gocontext.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}

	}

	fmt.Println("server exiting")

	return nil, err
}

func (T *Serve) routeFromValue(path string, route cue.Value, e *echo.Echo, ctx *hofcontext.Context) error {
	path = strings.Replace(path, "\"", "", -1)
	// fmt.Println(path + ":", route)

	// is this a flow handler?
	attrs := route.Attributes(cue.ValueAttr)
	isPipe := false
	for _, a := range attrs {
		if a.Name() == "flow" {
			isPipe = true
		}
	}

	local := route

	// fmt.Println("setting up route:", path, isPipe)

	// (1) can we read the flown once and reuse it
	// (2) or do we need to construct a new one on each call

	// setup handler, this will be invoked on all requests
	handler := func(c echo.Context) error {
		fmt.Println("start handling:", path)
		// pull apart c.request
		req, err := T.buildReqValue(c)
		if err != nil {
			fmt.Println("req build error: ", err)
			return err
		}
		// fmt.Println("reqVal", req)
		b := local.LookupPath(cue.ParsePath("req.body"))

		if b.Exists() {
			switch b.IncompleteKind() {
			case cue.BytesKind:
				// nothing, already bytes
			case cue.StringKind:
				req["body"] = string(req["body"].([]byte))

			case cue.StructKind:
				var body interface{}
				err = json.Unmarshal(req["body"].([]byte), &body)
				if err != nil {
					return err
				}
				req["body"] = body
			}
		}

		tmp := local.FillPath(cue.ParsePath("req"), req)
		if tmp.Err() != nil {
			fmt.Println("req fill error: ", tmp.Err())
			return tmp.Err()
		}
		// fmt.Println("tmp", tmp)

		if isPipe {
			p, err := flow.OldFlow(ctx, tmp)
			if err != nil {
				fmt.Println("handler pipe/new error:", err)
				return err
			}
			// p.Root = ctx.RootValue
			// p.Root = tmp

			err = p.Start()
			if err != nil {
				fmt.Println("handler pipe/run error:", err)
				return err
			}

			tmp = p.Final
		}

		resp := tmp.LookupPath(cue.ParsePath("resp"))
		if resp.Err() != nil {
			fmt.Println("handler resp error:", resp.Err())
			return resp.Err()
		}

		err = T.fillRespFromValue(resp, c)
		if err != nil {
			fmt.Println("handler fill error:", err)
			return err
		}

		fmt.Println("done handling:", path)
		return nil
	}

	// figure out route method(s): GET, POST, et al
	mv := route.LookupPath(cue.ParsePath("method"))
	methods := []string{}
	switch mv.IncompleteKind() {
	case cue.StringKind:
		m, err := mv.String()
		if err != nil {
			return err
		}
		m = strings.ToUpper(m)
		methods = append(methods, m)
	case cue.ListKind:
		iter, err := mv.List()
		if err != nil {
			return err
		}
		for iter.Next() {
			v := iter.Value()
			m, err := v.String()
			if err != nil {
				return err
			}
			m = strings.ToUpper(m)
			methods = append(methods, m)
		}

	case cue.BottomKind:
		methods = append(methods, "GET")

	default:
		return fmt.Errorf("unsupported type for method in %s %v", path, mv.IncompleteKind)
	}

	// fmt.Println("methods:", methods)
	e.Match(methods, path, handler)

	return nil
}

func (T *Serve) buildReqValue(c echo.Context) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	R := c.Request()

	req["method"] = R.Method
	req["header"] = R.Header
	req["url"] = R.URL
	req["query"] = c.QueryParams()

	b, err := io.ReadAll(R.Body)
	if err != nil {
		return nil, err
	}

	if len(b) > 0 {
		req["body"] = b
	}

	// form
	// path params
	return req, nil
}

func (T *Serve) fillRespFromValue(val cue.Value, c echo.Context) error {
	var ret map[string]interface{}

	{
		T.Lock()
		defer T.Unlock()

		err := val.Decode(&ret)
		if err != nil {
			return err
		}
	}

	// TODO, more http/response type things

	st, ok := ret["status"]
	if !ok {
		st = 200
	}
	status := st.(int)

	if ret["json"] != nil {
		return c.JSON(status, ret["json"])
	} else if ret["html"] != nil {
		// todo, better type casts
		return c.HTML(status, ret["html"].(string))
	} else if ret["body"] != nil {
		// todo, better type casts
		return c.String(status, ret["body"].(string))
	} else {
		return c.NoContent(status)
	}
}
