package structural_test

import (
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"

	"github.com/hofstadter-io/hof/lib/structural"
)

func TestPick(t *testing.T) {
	ctx := cuecontext.New()

	V, err := structural.LoadCueInputs([]string{"testdata/pick-cases.cue"}, ctx, nil)
	if err != nil {
    s := structural.FormatCueError(err)
		t.Fatalf("Error: %s", s)
	}

  cases := V.LookupPath(cue.ParsePath("pick_cases"))
  iter, _ := cases.List()

  for iter.Next() {
    curr   := iter.Value()
    value  := curr.LookupPath(cue.ParsePath("value"))
    pick   := curr.LookupPath(cue.ParsePath("pick"))
    // expect := curr.LookupPath(cue.ParsePath("expect"))

		result, err := structural.PickValue(pick, value, nil)
    if err != nil {
      t.Fatal(err)
    }
    if result.Err() != nil {
      t.Fatalf("in %q: error: %v\n", iter.Selector(), result.Err())
    }

    //same := expect.Unify(result)
    //if !same.Exists() || same.Err() != nil {
      //t.Fatalf("in %q: expected:\n%v\n\ngot:\n%v\n", iter.Selector(), expect, result)
    //}
		//o, err := cuetils.FormatOutput(r, "cue")
		//if err != nil {
			//t.Fatalf("failed to format output in TestPick %d", i)
		//}
	}
}

