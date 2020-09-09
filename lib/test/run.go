package test

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func RunSuites(suites []Suite, verbose int) (TS Stats, err error) {

	// set start time
	TS.Start = time.Now()

	// loop over tests
	for s, S := range suites {
		// Short circuit for empty suites
		if len(S.Tests) == 0 {
			continue
		}

		// set start time
		S.Stats.Start = time.Now()

		for t, T := range S.Tests {

			// run the test
			err = RunTest(&T, verbose)
			if err != nil {
				S.Errors = append(S.Errors, err)
			}

			// update total time
			S.Stats.Time += T.Stats.Time

			// set the element because of Go value copy in the loop header
			S.Tests[t] = T
		}

		// stats work
		S.Stats.End = time.Now()
		TS.Time += S.Stats.Time

		// set the element again because of Go copy
		suites[s] = S
	}

	// stats work
	TS.End = time.Now()

	return TS, nil
}

func RunTest(T *Tester, verbose int) (err error) {
	// stats work
	T.Stats.Start = time.Now()

	// TODO find and possibly skip
	skip := T.Value.Lookup("skip")
	if skip.Exists() {
		doskip, err := skip.Bool()
		if err != nil {
			return err
		}
		if doskip {
			T.Stats.Skip += 1
			return nil
		}
	}


	// switch type
	switch T.Type {

	case "bash", "shell":
		err = RunBash(T, verbose)

	case "exec", "custom":
		err = RunExec(T, verbose)

	case "hls", "script":
		err = RunScript(T, verbose)

	default:
		err = fmt.Errorf("unknown tester type %q", T.Type)
	}

	T.Stats.End = time.Now()
	T.Stats.Time = T.Stats.End.Sub(T.Stats.Start)

	if err != nil {
		T.Stats.Fail += 1
		T.Errors = append(T.Errors, err)
		err = fmt.Errorf("Test Failed: %s", T.Name)

		fmt.Println("!!!!!!!  BEG Failure: ", T.Name)
		fmt.Println(T.Output)
		fmt.Println(T.Errors)
		fmt.Println("!!!!!!   END Failure: ", T.Name)
	} else {
		T.Stats.Pass += 1
	}

	return err
}

type BaseTester struct {
	Dir    string
	Env    map[string]string
	Sysenv bool
}

type BashTester struct {
	BaseTester

	Script string
}

func RunBash(T *Tester, verbose int) (err error) {
	// Decode our BT
	var BT BashTester
	err = T.Value.Decode(&BT)

	// Check for errors and validate
	if err != nil {
		return err
	}
	if BT.Script == "" {
		return fmt.Errorf("Bash tester %q has empty script field", T.Name)
	}

	// Prep our command
	cmd := exec.Command("bash", "-p", "-c", BT.Script)
	cmd.Dir = BT.Dir

	// add env vars if needed
	if len(BT.Env) > 0 {
		for k,v := range BT.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	// Run and save output
	out, err := cmd.CombinedOutput()
	T.Output = string(out)

	return err
}

type ExecTester struct {
	BaseTester

	Command string
}

func RunExec(T *Tester, verbose int) (err error) {
	// Decode our ET
	var ET ExecTester
	err = T.Value.Decode(&ET)

	// Check for errors and validate
	if err != nil {
		return err
	}
	if ET.Command == "" {
		return fmt.Errorf("Bash tester %q has empty script field", T.Name)
	}

	args := strings.Fields(ET.Command)

	// Prep our command
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = ET.Dir

	// add env vars if needed
	if len(ET.Env) > 0 {
		for k,v := range ET.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	// Run and save output
	out, err := cmd.CombinedOutput()
	T.Output = string(out)

	return err
}

func RunScript(T *Tester, verbose int) (err error) {
	// fmt.Println("hls:", T.Name)

	return nil
}

