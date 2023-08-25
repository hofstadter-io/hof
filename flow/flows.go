package flow

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"github.com/bmatcuk/doublestar/v4"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/flow/flow"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/structural"
)

func hasFlowAttr(val cue.Value, args []string) (attr cue.Attribute, found, keep bool) {
	attrs := val.Attributes(cue.ValueAttr)

	for _, attr := range attrs {
		if attr.Name() == "flow" {
			// found a flow, stop recursion
			found = true
			// if it matches our args, create and append
			keep = matchFlow(attr, args)
			if keep {
				return attr, true, true
			}
		}
	}

	return cue.Attribute{}, found, false
}

func matchFlow(attr cue.Attribute, args []string) (keep bool) {
	// fmt.Println("matching 1:", attr, args, len(args), attr.NumArgs())
	// if no args, match flows without args
	if len(args) == 0 {
		if attr.NumArgs() == 0 {
			return true
		}
		// extra check for one arg which is empty
		if attr.NumArgs() == 1 {
			s, err := attr.String(0)
			if err != nil {
				fmt.Println("bad flow tag:", err)
				return false
			}
			return s == ""
		}

		return false
	}

	// for now, match any
	// upgrade logic for user later
	for _, tag := range args {
		for p := 0; p < attr.NumArgs(); p++ {
			s, err := attr.String(p)
			if err != nil {
				fmt.Println("bad flow tag:", err)
				return false
			}

			// exact match
			if s == tag {
				return true
			}

			// glob match
			if strings.Contains(tag, "*") {
				include, err := doublestar.PathMatch(tag, s)
				if err != nil {
					fmt.Println("warning:", err)
				}
				return include
			}
		}
	}

	return false
}

func flowList(val cue.Value, opts *flags.RootPflagpole, popts *flags.FlowFlagpole) []string {
	var names []string

	accum := func(v cue.Value) bool {
		attrs := v.Attributes(cue.ValueAttr)

		for _, attr := range attrs {
			if attr.Name() == "flow" {
				if attr.NumArgs() == 0 {
					names = append(names, "<unnamed>")
				} else {
					name, _ := attr.String(0)
					names = append(names, name)
				}
				return false
			}
		}

		return true
	}

	structural.Walk(val, accum, nil, walkOptions...)

	return names
}

func printFlows(val cue.Value, opts *flags.RootPflagpole, popts *flags.FlowFlagpole) error {
	args := popts.Flow
	foundAny := false

	printer := func(v cue.Value) bool {
		attrs := v.Attributes(cue.ValueAttr)

		for _, attr := range attrs {
			if attr.Name() == "flow" {
				if len(args) == 0 || matchFlow(attr, args) {
					foundAny = true
					//if popts.Docs {
					//  s := ""
					//  docs := v.Doc()
					//  for _, d := range docs {
					//    s += d.Text()
					//  }
					//  fmt.Print(s)
					//}
					if opts.Verbosity > 0 {
						s, _ := cuetils.FormatCue(v)
						fmt.Printf("%s: %s\n", v.Path(), s)
					} else {
						fmt.Println(attr)
					}
				}
				return false
			}
		}

		return true
	}

	structural.Walk(val, printer, nil, walkOptions...)

	if !foundAny {
		fmt.Println("no flows found")
	}

	return nil
}

// maybe this becomes recursive so we can find
// nested flows, but not recurse when we find one
// only recurse when we have something which is not a flow or task
func findFlows(ctx *context.Context, val cue.Value, opts *flags.RootPflagpole, popts *flags.FlowFlagpole) ([]*flow.Flow, error) {
	flows := []*flow.Flow{}

	// TODO, look for lists?
	s, err := val.Struct()
	if err != nil {
		// not a struct, so don't recurse
		// fmt.Println("not a struct", err)
		return nil, nil
	}

	args := popts.Flow

	// does our top-level (file-level) have @flow()
	_, found, keep := hasFlowAttr(val, args)
	if keep {
		// invoke TaskFactory
		p, err := flow.NewFlow(ctx, val)
		if err != nil {
			return flows, err
		}
		flows = append(flows, p)
	}

	if found {
		return flows, nil
	}

	iter := s.Fields(
		cue.Attributes(true),
		cue.Docs(true),
	)

	// loop over top-level fields in the cue instance
	for iter.Next() {
		v := iter.Value()

		_, found, keep := hasFlowAttr(v, args)
		if keep {
			p, err := flow.NewFlow(ctx, v)
			if err != nil {
				return flows, err
			}
			flows = append(flows, p)
		}

		// break recursion if flow found
		if found {
			continue
		}

		// recurse to search for more flows
		ps, err := findFlows(ctx, v, opts, popts)
		if err != nil {
			return flows, nil
		}
		flows = append(flows, ps...)
	}

	return flows, nil
}
