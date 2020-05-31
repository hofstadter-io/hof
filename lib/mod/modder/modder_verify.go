package modder

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/yagu"
)

func (mdr *Modder) Verify() error {

	// Verify Command Override
	if len(mdr.CommandVerify) > 0 {
		for _, cmd := range mdr.CommandVerify {
			out, err := yagu.Exec(cmd)
			fmt.Println(out)
			if err != nil {
				return err
			}
		}
	} else {
		// Otherwise, MVS venodiring
		err := mdr.VerifyMVS()
		if err != nil {
			mdr.PrintErrors()
			return err
		}
	}

	return nil
}

// The entrypoint to the MVS internal verify process
func (mdr *Modder) VerifyMVS() error {

	valid := true

	// Load minimal root module
	err := mdr.LoadMetaFromFS(".")
	if err != nil {
		return err
	}

	// Load the root module's deps
	present, missing, local, err := mdr.PartitionSumEntries()
	if err != nil {
		return err
	}

	// Invalid if there are missing deps
	for _, m := range missing {
		valid = false
		R := mdr.module.SelfDeps[m]
		err := fmt.Errorf("Sumfile missing: %s@%s", R.NewPath, R.NewVersion)
		mdr.errors = append(mdr.errors, err)
	}

	for _, p := range present {
		R := mdr.module.SelfDeps[p]
		err := mdr.CompareSumEntryToVendor(R)
		// Something is wrong with the vendored copy or the hash
		if err != nil {
			valid = false
			mdr.errors = append(mdr.errors, err)
		}
	}

	for _, p := range local {
		R := mdr.module.SelfDeps[p]
		err := mdr.CompareLocalReplaceToVendor(R)
		// Something is wrong with the vendored copy or the hash
		if err != nil {
			valid = false
			mdr.errors = append(mdr.errors, err)
		}
	}

	if !valid {
		return fmt.Errorf("Vendoring is in an inconsistent state, please run 'mvs vendor %s' ", mdr.Name)
	}

	// We are OK!
	return nil
}
