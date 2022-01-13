package main

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/kr/pretty"
	"gopkg.in/yaml.v2"

	"github.com/hofstadter-io/hof/lib/dotpath"
)

func main() {
	fmt.Println("dotpath test\n----------------\n")

	//	dotpath.SetLogLevel("debug")

	test_array()

}

func test_json() error {

	return nil
}

func test_yaml() {
	fmt.Println("Testing yaml")

	data, err := read_yaml("data/test.yaml")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Printf("data:\n%# v\n\n", pretty.Formatter(data))

	paths := []string{
		"data.key",
		"data.array",
		"data.array.subarray",
		"data.array.[:].subarray",
		"data.array.elemA",
		"data.array.[elemA]",
		"data.array.[elemA,elemB]",
		"data.array.[name==elemB]",
		"data.array.[name==elemB,elemC]",
		"data.array.[value==foo]",
		"data.array.[value==foo,goo]",
		"data.object",
		"data.object.name",
		"data.object.myobject",
		"data.object.field1",
		"data.object.[name]",
		"data.object.[field1]",
		"data.object.[name,field1]",
		"data.object.[name,field1,field2]",
		"data.object.list",
		"data.object.list.u",
		"data.object.list.x",
		"data.object.list.[:2]",
		"data.object.list.[2]",
		"data.object.list.[2:]",
		"data.object.list.[:]",
		"data.object.list.[]", // should fail
		// "data.object.list.[value==0]",
		// "data.object.list.[value!=0]",
		// "data.object.list.[value>0]",
		// "data.object.list.[value>=0]",
		// "data.object.list.[value<0]",
		// "data.object.list.[value<=0]",
		"data.array.[:].subobject",
		"data.array.[:].subobject.array",
		"data.array.[:].subobject.array.[value==foo]",
		"data.array.[name==elemB,elemC].subobject.array.[elemA,elemC]",
	}

	for _, path := range paths {

		fmt.Printf("@%s:\n", path)
		d, err := dotpath.Get(path, data)
		if err != nil {
			fmt.Println("ERROR:", err, "\n\n")
			continue
		}
		fmt.Printf("%# v\n\n", pretty.Formatter(d))
	}

}

func test_array() {
	fmt.Println("Testing yaml array")

	data, err := read_yaml("data/array.yaml")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Printf("data:\n%# v\n\n", pretty.Formatter(data))
	M := data.(map[string]interface{})

	paths := []string{
		"elemA",
		// "array.elemA",
		// "array.[elemA]",
		// "array.[elemA,elemB]",
		// "array.[name==elemB]",
		// "array.[name==elemB,elemC]",
		// "array.[value==foo]",
		// "array.[value==foo,goo]",
	}

	for _, path := range paths {

		fmt.Printf("@%s:\n", path)
		d, err := dotpath.Get(path, M["array"])
		if err != nil {
			fmt.Println("ERROR:", err, "\n\n")
			continue
		}
		fmt.Printf("%# v\n\n", pretty.Formatter(d))
	}

}

func test_struct() error {

	return nil
}

func read_yaml(filename string) (interface{}, error) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "while reading yaml file: (readfile) %s\n", filename)
	}

	obj := map[string]interface{}{}
	err = yaml.Unmarshal(data, &obj)
	if err != nil {
		return nil, errors.Wrapf(err, "while reading yaml file: (unmarshal) %s\n", filename)
	}

	return obj, nil
}
