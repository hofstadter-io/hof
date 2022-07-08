package gen

import (
	"testing"
)

const textO = `
celery
garlic
onions
salmon
tomatoes
wine
`

const textA = `
celery
salmon
tomatoes
garlic
onions
wine
`

const textB = `
celery
garlic
salmon
tomatoes
onions
wine
`

const expect = `
celery
salmon
tomatoes
garlic
<<<<<<<
onions
=======
salmon
tomatoes
onions
>>>>>>>
wine
`

func TestMerge(t *testing.T) {
	if Merge(textO, textA, textB) != expect {
		t.Errorf("unexpected merge result")
	}
}
