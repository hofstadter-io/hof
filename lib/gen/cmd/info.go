package cmd

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func Info(args []string, rflags flags.RootPflagpole, gflags flags.GenFlagpole, iflags flags.Gen__InfoFlagpole) error {
	R, err := prepRuntime(args, rflags, gflags)
	if err != nil {
		return err
	}

	if len(iflags.Expression) == 0 {
		fmt.Println(R.Value)
		return nil
	}

	for _, ex := range iflags.Expression {
		val := R.Value.LookupPath(ex).CueValue()
		path := val.Path()
		fmt.Printf("%s: %v\n\n", path, val)
	}

	return nil


	//for _, G := range R.Generators {
		//if len(iflags.Expression) == 0 {
			//fmt.Printf("%s: %v\n\n", G.Hof.Metadata.Name, G.CueValue)
			//continue
		//}

		//for _, ex := range iflags.Expression {
			//val := G.CueValue.LookupPath(cue.ParsePath(ex))
			//path := G.Hof.Metadata.Name + "." + ex
			//fmt.Printf("%s: %v\n\n", path, val)
		//}
	//}

	// return nil
}

