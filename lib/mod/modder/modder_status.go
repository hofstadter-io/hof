package modder

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/yagu"
)

func (mdr *Modder) Status() error {

	// Status Command Override
	if len(mdr.CommandStatus) > 0 {
		for _, cmd := range mdr.CommandStatus {
			out, err := yagu.Exec(cmd)
			fmt.Println(out)
			if err != nil {
				return err
			}
		}
	} else {
		// Otherwise, MVS venodiring
		err := mdr.StatusMVS()
		if err != nil {
			mdr.PrintErrors()
			return err
		}
	}

	return nil
}

// The entrypoint to the MVS internal verify process
func (mdr *Modder) StatusMVS() error {
	var err error

	// Load minimal root module
	err = mdr.LoadMetaFromFS(".")
	if err != nil {
		return err
	}

	fmt.Println("==================")

	mod := mdr.module
	sf := mod.SumFile

	err = mod.PrintSelfDeps()
	if err != nil {
		return err
	}

	fmt.Println("==================")

	if sf != nil {
		out, err := sf.Write()
		if err != nil {
			return err
		}
		fmt.Println(out)

	} else {
		fmt.Printf("No sum file %q found for lang %q\n", mdr.SumFile, mdr.Name)
	}
	fmt.Println("==================")

	return nil
}
