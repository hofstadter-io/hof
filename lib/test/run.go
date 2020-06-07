package test

import (
	"fmt"
	"strings"
	"time"
)

func RunSuites(suites []Suite, verbose int) (TS Stats, err error) {

	// set start time
	TS.Start = time.Now()

	// loop over tests
	for s, S := range suites {
		// set start time
		S.Stats.Start = time.Now()

		for t, T := range S.Tests {

			// run the test
			err = RunTest(&T, verbose)
			if err != nil {
				S.Errors = append(S.Errors, err)
			}

			// set the element again because of Go copy
			S.Stats.Time += T.Stats.Time
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

	// TODO< find and possibly skip
	// skip := T.Value.Lookup("skip")

	// switch type
	ts := strings.SplitN(T.Type, "=", 2)
	tkey, tval := ts[0], ""
	if len(ts) == 2 {
		tval = ts[1]
	}
	fmt.Printf("TYPE: '%v' '%v'\n", tkey, tval)
	switch tkey {

	case "exec":
		err = RunExec(T, verbose)

	case "script":
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
	} else {
		T.Stats.Pass += 1
	}

	fmt.Println()

	return err
}

func RunExec(T *Tester, verbose int) (err error) {
	fmt.Println("#!/", T.Name)
	ts := strings.SplitN(T.Type, "=", 2)
	tkey, tval := ts[0], ""
	if len(ts) == 2 {
		tval = ts[1]
	}
	fmt.Printf("  type: '%v' '%v'\n", tkey, tval)

	return nil
}

func RunScript(T *Tester, verbose int) (err error) {
	fmt.Println(">>>", T.Name)
	ts := strings.SplitN(T.Type, "=", 2)
	tkey, tval := ts[0], ""
	if len(ts) == 2 {
		tval = ts[1]
	}
	fmt.Printf("  type: '%v' '%v'\n", tkey, tval)

	return nil
}

