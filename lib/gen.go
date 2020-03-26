package lib

import (
	"fmt"
	"os"

	// "strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
)

func Gen(entrypoints, expressions []string, mode string) (string, error) {
	fmt.Println("Gen", entrypoints, expressions)

	var rt cue.Runtime

	// out := make(map[string]interface{})

	bis := load.Instances([]string{}, nil)
	for i, bi := range bis {
		fmt.Println("BI", i)
		if bi.Err != nil {
			fmt.Println(bi.Err)
			os.Exit(1)
		}
		i, err := rt.Build(bi)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Loop through all top level fields
		toplevel, err := i.Value().Struct()
		if err != nil {
			fmt.Println(err)
			continue
		}

		iter := toplevel.Fields()
		for iter.Next() {

			label := iter.Label()
			value := iter.Value()

			// 	ev := value
			ev := value.Eval()

			/*
				err = ev.Validate()
				if err != nil {
					fmt.Println(err)
					continue
				}
			*/

			fmt.Printf(" - %v %b\n\n%#+v\n\n\n", label, ev.IsConcrete(), ev)

			vi := 0

			value.Walk(func(val cue.Value) bool {
				// l, _ := val.Label()
				// k := val.Kind()
				// fmt.Println(vi, l, k)
				vi += 1

				return true
			}, nil)

			fmt.Println("VI: ", vi)

			// Put anything starting with Gen into
			// our out map

			/*
				if strings.HasPrefix(label, "Gen") {
					var obj map[string]interface{}
					err = value.Decode(&obj)
					if err != nil {
						fmt.Println(err)
						continue
					}
					out[label] = obj
				}
			*/
		}
	}

	/*
		GenCli := out["GenCli"].(map[string]interface{})
		All := GenCli["All"].([]interface{})
		Zero := All[0]

		what := Zero

		bytes, err := yaml.Marshal(what)
		if err != nil {
			fmt.Println(err)
			return "", nil
		}
		fmt.Println(string(bytes))

		// TODO see if we can parse and introspect *_tool.cue files
	*/

	return "", nil
}
