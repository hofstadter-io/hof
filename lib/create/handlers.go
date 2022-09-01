package create

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

type handler func (Q map[string]any) (A any, err error)

var handlers map[string]handler

func init() {
	handlers = map[string]handler{
		// builtin handlers
		"input": handleInput,
		"multiline": handleMultiline,
		"password": handlePassword,
		"confirm": handleConfirm,
		"select": handleSelect,
		"multiselect": handleMultiselect,

		// custom handlers
	}
}

func handleQuestion(Q map[string]any) (A any, err error) {
	// ask question until we get an answer or interrupt
	for {
		s, ok := Q["Type"]
		if !ok {
			panic("question type not set")
		}

		S, ok := s.(string)
		if !ok {
			panic("question type not a string")
		}

		h, ok := handlers[S]
		if !ok {
			panic("unknown question type: " + S)
		}

		a, err := h(Q)

		if err != nil {
			if err == terminal.InterruptErr {
				return nil, fmt.Errorf("user interrupt")
			}
			fmt.Println("error:", err)
			continue
		}

		A = a

		// we got an answer
		break
	}

	return A, nil
}

func handleInput(Q map[string]any) (A any, err error) {
	dval := ""
	if d, ok := Q["Default"]; ok {
		dval = d.(string)
	}
	prompt := &survey.Input {
		// todo, rename prompt to message in test and schema
		Message: Q["Prompt"].(string),
		Default: dval,
	}
	var a string
	err = survey.AskOne(prompt, &a)
	A = a

	return A, err
}

func handleMultiline(Q map[string]any) (A any, err error) {
	dval := ""
	if d, ok := Q["Default"]; ok {
		dval = d.(string)
	}
	prompt := &survey.Multiline {
		// todo, rename prompt to message in test and schema
		Message: Q["Prompt"].(string),
		Default: dval,
	}
	var a string
	err = survey.AskOne(prompt, &a)
	A = a

	return A, err
}

func handlePassword(Q map[string]any) (A any, err error) {
	prompt := &survey.Password {
		// todo, rename prompt to message in test and schema
		Message: Q["Prompt"].(string),
	}
	var a string
	err = survey.AskOne(prompt, &a)
	A = a

	return A, err
}

func handleConfirm(Q map[string]any) (A any, err error) {
	dval := false
	if d, ok := Q["Default"]; ok {
		dval = d.(bool)
	}
	prompt := &survey.Confirm {
		// todo, rename prompt to message in test and schema
		Message: Q["Prompt"].(string),
		Default: dval,
	}
	var a bool
	err = survey.AskOne(prompt, &a)
	if !a {
		return nil, nil
	} else {
		A = a
	}

	if err != nil {
		return A, err
	}

	// possibly recurse if has own Questions
	QS, ok := Q["Questions"]
	if a && ok {
		A2 := map[string]any{}
		for _, Q2 := range QS.([]any) {
			q2 := Q2.(map[string]any)
			a2, e2 := handleQuestion(q2)
			// todo, think about if/how to handle this error
			if e2 != nil {
				return nil, e2
			}
			A2[q2["Name"].(string)] = a2
		}
		A = A2
	}

	return A, err
}

func handleSelect(Q map[string]any) (A any, err error) {
	opts := []string{}
	for _, o := range Q["Options"].([]any) {
		opts = append(opts, o.(string))
	}
	prompt := &survey.Select {
		// todo, rename prompt to message in test and schema
		Message: Q["Prompt"].(string),
		// todo, probably need to error handle options
		Options: opts,
		Default: Q["Default"],
	}
	var a string
	err = survey.AskOne(prompt, &a)
	A = a

	return A, err
}
	
func handleMultiselect(Q map[string]any) (A any, err error) {
	opts := []string{}
	for _, o := range Q["Options"].([]any) {
		opts = append(opts, o.(string))
	}
	prompt := &survey.MultiSelect {
		// todo, rename prompt to message in test and schema
		Message: Q["Prompt"].(string),
		// todo, probably need to error handle options
		Options: opts,
		Default: Q["Default"],
	}
	var a []string
	err = survey.AskOne(prompt, &a)
	A = a

	return A, err
}
