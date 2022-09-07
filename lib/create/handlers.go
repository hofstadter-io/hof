package create

import (
	"fmt"

	"cuelang.org/go/cue"
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
		"subgroup": handleSubgroup,

		// custom handlers
	}
}

func runPrompt(genVal cue.Value) (result cue.Value, err error) {
	// run while there are unanswered questions
	// we run in an extra loop to fill back answers
	// and recalculate the prompt questions each iteration
	done := false

	for !done {
		done = true

		inputVal := genVal.LookupPath(cue.ParsePath("Create.Input"))
		if inputVal.Err() != nil {
			return genVal, inputVal.Err()
		}

		// fmt.Printf("outer loop input: %#v\n", inputVal)

		prompt := genVal.LookupPath(cue.ParsePath("Create.Prompt"))
		if prompt.Err() != nil {
			return genVal, prompt.Err()
		}
		if !prompt.IsConcrete() || !prompt.Exists() {
			// to have a promptless generator, set it to the empty list
			return genVal, fmt.Errorf("Generator is missing Create.Prompt, set to empty list for promptless")
		}

		// prompt should be an ordered list of questions
		iter, err := prompt.List()
		if err != nil {
			return genVal, err
		}

		// loop over prompt questions, recursing as needed
		for iter.Next() {
			// todo, get label and check if input[label] is concrete
			value := iter.Value()

			Q := map[string]any{}
			err := value.Decode(&Q)
			if err != nil {
				return genVal, err
			}

			// fmt.Printf("%#v\n", Q)

			name := Q["Name"].(string)
			namePath := cue.ParsePath(name)

			// check if done already by inspececting in input
			i := inputVal.LookupPath(namePath)
			if i.Err() != nil {
				if i.Exists() {
					return genVal, i.Err()
				}
			}
			if i.Exists() && i.IsConcrete() {
				// question answer already exists in input
				// fmt.Println("continuing: ", name)
				continue
			}

			// there is a question to answer
			done = false

			// fmt.Println("q:", Q)
			// todo, extract Name
			A, err := handleQuestion(Q)
			if err != nil {
				if err == terminal.InterruptErr {
					return genVal, fmt.Errorf("user interrupt")
				}
				return genVal, err
			}

			// update input val
			inputVal = inputVal.FillPath(namePath, A)
			genVal = genVal.FillPath(cue.ParsePath("Create.Input"), inputVal)

			// restart the prompt loop
			break
		}
	}

	return genVal, nil
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
			panic("question 'Type' is not set using a string format")
		}

		h, ok := handlers[S]
		if !ok {
			panic("unknown question type: " + S)
		}

		a, err := h(Q)

		if err != nil {
			if err == terminal.InterruptErr {
				return nil, err
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
	if err != nil {
		if err == terminal.InterruptErr {
			return nil, err
		}
	}
	A = a

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

func handleSubgroup(Q map[string]any) (A any, err error) {
	// get some strings
	name := Q["Name"].(string)
	msg := Q["Prompt"].(string)
	fmt.Println(msg)

	// get subgroup Questions
	QS, ok := Q["Questions"]
	if !ok {
		return nil, fmt.Errorf("subgroup prompt %q is missing 'Questions' field", name)
	}

	// gather nested answers
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

	// set as nested value in A
	a := map[string]any{}
	a[name] = A2

	return a, err
}

