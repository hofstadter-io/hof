package yagu

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/parnurzeal/gorequest"
)

type GaConfig struct {
	TID string
	CID string
	IP  string
	UA  string

	CN string
	CS string
	CM string
}

// TODO... make this accurate, add a second function, make both funcs specialized
// XXX see if we can find the API spec and reverse it into gocode
type GaPageview struct {
	Type     string
	Source   string
	Category string
	Action   string
	Label    string
	Value    int
}

type GaEvent struct {
	Type     string
	Source   string
	Category string
	Action   string
	Label    string
	Value    int
}

func SendGaEvent(cfg GaConfig, evt GaEvent) (string, error) {

	gaURL := "https://www.google-analytics.com/collect"

	vals := url.Values{}

	vals.Add("tid", cfg.TID)
	vals.Add("cid", cfg.CID)
	vals.Add("cs", cfg.CS)
	vals.Add("cn", cfg.CN)
	vals.Add("cm", cfg.CM)
	if cfg.IP != "" {
		vals.Add("uip", cfg.IP)
	}
	if cfg.UA != "" {
		vals.Add("ua", cfg.UA)
	}

	if evt.Type != "" {
		vals.Add("t", "pageview")
	} else {
		vals.Add("t", evt.Type)
	}

	// TODO, move this parameter hackery for CLIs to hofmod-cli
	vals.Add("dh", cfg.UA)
	vals.Add("dt", cfg.CS)
	vals.Add("dp", evt.Action)
	vals.Add("v", "1")
	if evt.Source != "" {
		vals.Add("ds", evt.Source)
	}

	payload := vals.Encode()

	// fmt.Println("GA: ", payload)

	req := gorequest.New().Post(gaURL).Send(payload)

	resp, body, errs := req.End()

	// fmt.Println(resp, body, errs)

	if len(errs) != 0 && !strings.Contains(errs[0].Error(), HTTP2_GOAWAY_CHECK) {
		return body, errs[0]
	}

	if len(errs) != 0 || resp.StatusCode >= 500 {
		return body, errors.New("Internal Error: " + body)
	}
	if resp.StatusCode >= 400 {
		return body, errors.New("Bad Request: " + body)
	}

	return body, nil
}

func SendGaEvents(cfg GaConfig, evts []GaEvent) (string, error) {

	gaURL := "https://www.google-analytics.com/batch"

	payload := ""
	for _, evt := range evts {
		vals := url.Values{}

		vals.Add("tid", cfg.TID)
		vals.Add("cid", cfg.CID)
		vals.Add("cs", cfg.CS)
		vals.Add("cn", cfg.CN)
		if cfg.IP != "" {
			vals.Add("uip", cfg.IP)
		}
		if cfg.UA != "" {
			vals.Add("ua", cfg.UA)
		}

		if evt.Type != "" {
			vals.Add("t", "event")
		} else {
			vals.Add("t", evt.Type)
		}
		vals.Add("ec", evt.Category)
		vals.Add("ea", evt.Action)
		vals.Add("el", evt.Label)
		vals.Add("ev", fmt.Sprint(evt.Value))
		vals.Add("v", "1")
		if evt.Source != "" {
			vals.Add("ds", evt.Source)
		}

		payload += vals.Encode() + "\n"
	}

	req := gorequest.New().Post(gaURL).Send(payload)

	resp, body, errs := req.End()

	if len(errs) != 0 && !strings.Contains(errs[0].Error(), HTTP2_GOAWAY_CHECK) {
		return body, errs[0]
	}

	if len(errs) != 0 || resp.StatusCode >= 500 {
		return body, errors.New("Internal Error: " + body)
	}
	if resp.StatusCode >= 400 {
		return body, errors.New("Bad Request: " + body)
	}

	return body, nil
}
