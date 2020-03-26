package lib

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/kr/pretty"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"

	// "github.com/hofstadter-io/hof/lib/util"
)

func Gen(entrypoints, expressions []string, mode string) (string, error) {
	fmt.Println("Gen", entrypoints, expressions)

	var rt cue.Runtime

	out := make(map[string]interface{})

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

			ev := value
			// ev := value.Eval()

			/*
				err = ev.Validate()
				if err != nil {
					fmt.Println(err)
					continue
				}
			*/

			fmt.Printf(" - %v %b\n\n%#+v\n\n\n", label, ev.IsConcrete(), ev)

			/*
				vi := 0

				value.Walk(func(val cue.Value) bool {
					// l, _ := val.Label()
					// k := val.Kind()
					// fmt.Println(vi, l, k)
					vi += 1

					return true
				}, nil)

				fmt.Println("VI: ", vi)
			*/

			// Put anything starting with Gen into
			// our out map

			if strings.HasPrefix(label, "Gen") {
				var gen map[string]interface{}
				err = value.Decode(&gen)
				if err != nil {
					fmt.Println(err)
					continue
				}

				IN, ok := gen["In"].(map[string]interface{})

				fmt.Println("IN:===============")
				fmt.Printf("%# v\n", pretty.Formatter(IN))
				fmt.Println("IN:===============")


				OUT, ok := gen["Out"].([]interface{})
				if !ok {
					return "", fmt.Errorf("Generator: %q is missing 'Out' field.", label)
				}

				for _, O := range OUT {
					o := O.(map[string]interface{})
					renderFile(IN, o)
				}

				// obj := OUT
				// fmt.Printf("%# v\n\n", pretty.Formatter(obj))
				out[label] = OUT

				break
			}
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

	/*
		stdout, err := util.Exec([]string{"goimports", "-w", "-l", "."})
		if err != nil {
			return "", err
		}
	*/

	return "", nil
}

func renderFile(IN, file map[string]interface{}) error {
	// Look for input on the file
	fn := file["Filename"].(string)
	tp := file["Template"].(string)
	in, ok := file["In"].(map[string]interface{})

	fmt.Println(file["Filename"])

	// If not there, use the global IN
	if !ok {
		fmt.Println("missing In, replacing with global")
		in = IN

	} else {
		fmt.Println("checking In, filling in gaps with global")

		// Else, 'IN' has key and 'in' does not, add it
		for key, val := range IN {
			if _, ok := in[key]; !ok {
				fmt.Println("checking In, filling", key)
				in[key] = val
			}
		}
	}

	for key, _ := range in {
		fmt.Println(" -", key)
	}

	// alt := file["Alt"].(bool)

	t := template.Must(template.New(fn).Parse(tp))

	/*
		f, err := os.Create(fn)
		if err != nil {
			return err
		}
		defer f.Close()

		err = t.Execute(f, in)
		if err != nil {
			return err
		}
	*/

	var b bytes.Buffer
	var err error

	err = t.Execute(&b, in)
	if err != nil {
		return err
	}

	fmtd, err := format.Source(b.Bytes())
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(fn), 0755)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fn, fmtd, 0644)
	if err != nil {
		return err
	}

	return nil
}
