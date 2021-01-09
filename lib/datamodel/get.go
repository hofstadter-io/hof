package datamodel

import (
	// "fmt"
	"os"

	"cuelang.org/go/cue"
	"github.com/olekukonko/tablewriter"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	// "github.com/hofstadter-io/hof/lib/cuetils"
)

func RunGetFromArgs(args []string, cmdpflags flags.DatamodelPflagpole) error {
	// fmt.Println("lib/datamodel.Get", args, cmdpflags)

	// Loadup our Cue files
	val, err := LoadDatamodel(args)
	if err != nil {
		return err
	}

	// TODO: find values from flags / attributes, condition printing

	err = printTable(val)
	if err != nil {
		return err
	}

	//syn, err := cuetils.PrintCueValue(val)
	//if err != nil {
		//return err
	//}

	//fmt.Println(syn)

	return nil
}

func printTable(val cue.Value) error {

	models := val.LookupPath(cue.ParsePath("Models.MigrateOrder"))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Model", "Field", "Type", "Status"})
	// table.SetRowLine(true)

	iter, err := models.List()
	if err != nil {
		return err
	}

	for iter.Next() {
		model := iter.Value()
		err = printModel(model, table)
		if err != nil {
			return err
		}
	}


	table.Render()

	return nil
}

func printModel(model cue.Value, table *tablewriter.Table) error {
	// Note, we can probably avoid checking errors because Cue has already validated the data

	// Model name
	name, _ := model.LookupPath(cue.ParsePath("Name")).String()
	table.Append([]string{ name, "", "", "" })

	// Model fields
	fields := model.LookupPath(cue.ParsePath("Fields"))
	st, _ := fields.Struct()
	iter := st.Fields()

	for iter.Next() {
		f := iter.Value()

		fn, _ := f.LookupPath(cue.ParsePath("Name")).String()
		ft, _ := f.LookupPath(cue.ParsePath("type")).String()

		table.Append([]string{ "", fn, ft, "" })
	}

	return nil
}

