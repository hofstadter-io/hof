package pkg

import (
	"fmt"
	"os"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	cueyaml "cuelang.org/go/encoding/yaml"
	"gopkg.in/yaml.v2"
)

func Gen(entrypoints, expressions []string, mode string) (string, error) {
	fmt.Println("Gen", entrypoints, expressions)

	var rt cue.Runtime

	out := make(map[string]interface{})

	bis := load.Instances([]string{}, nil)
	for _, bi := range bis {
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

			// Put anything starting with Gen into
			// our out map
			bytes, err := cueyaml.Encode(value)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if strings.HasPrefix(label, "Gen") {
				var ym interface{}
				err = yaml.Unmarshal(bytes, &ym)
				if err != nil {
					fmt.Println(err)
					continue
				}
				out[label] = ym
			}
		}
	}

	bytes, err := yaml.Marshal(out)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	fmt.Println(string(bytes))

	// TODO see if we can parse and introspect *_tool.cue files

	return "", nil
}
