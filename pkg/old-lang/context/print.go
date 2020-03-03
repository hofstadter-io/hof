package context

import (
	"fmt"
	"github.com/kr/pretty"
)

func (ctx *Context) PrintErrors() bool {
	if len(ctx.Errors) == 0 {
		return false
	}

	fmt.Println("Errors:")
	for i, err := range ctx.Errors {
		fmt.Println(i, err, "\n")
	}

	return true
}

func (ctx *Context) Pretty() {
	fmt.Printf("%#v\n", pretty.Formatter(*ctx))
}

func (ctx *Context) Print() {

	for _, pkg := range ctx.Packages {
		fmt.Println("Package:", pkg.Path)
		for _, file := range pkg.Files {
			fmt.Println(" -", file.Path)
		}
	}
	/*
	hofData, err := hofFile.ToData()
	if err != nil {
		return err
	}

	var b bytes.Buffer
	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)

	err = encoder.Encode(hofData)
	if err != nil {
		return err
	}

	fmt.Printf("--- parseYaml:\n%s\n", b.String())
	*/

}

