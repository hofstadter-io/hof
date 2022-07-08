package gen

import (
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	// Sep1 signifies the start of a conflict.
	Sep1 = "<<<<<<<"
	// Sep2 signifies the middle of a conflict.
	Sep2 = "======="
	// Sep3 signifies the end of a conflict.
	Sep3 = ">>>>>>>"
)

// DiffMatchPatch contains the diff algorithm settings.
var DiffMatchPatch = diffmatchpatch.New()

// Merge implements the diff3 algorithm to merge two texts into a common base.
func Merge(textO, textA, textB string) string {
	runesO, runesA, linesA := DiffMatchPatch.DiffLinesToRunes(textO, textA)
	_, runesB, linesB := DiffMatchPatch.DiffLinesToRunes(textO, textB)

	diffsA := DiffMatchPatch.DiffMainRunes(runesO, runesA, false)
	diffsB := DiffMatchPatch.DiffMainRunes(runesO, runesB, false)

	matchesA := matches(diffsA, runesA)
	matchesB := matches(diffsB, runesB)

	var result strings.Builder
	indexO, indexA, indexB := 0, 0, 0
	for {
		i := nextMismatch(indexO, indexA, indexB, runesA, runesB, matchesA, matchesB)

		o, a, b := 0, 0, 0
		if i == 1 {
			o, a, b = nextMatch(indexO, runesO, matchesA, matchesB)
		} else if i > 1 {
			o, a, b = indexO+i, indexA+i, indexB+i
		}

		if o == 0 || a == 0 || b == 0 {
			break
		}

		chunk(indexO, indexA, indexB, o-1, a-1, b-1, runesO, runesA, runesB, linesA, linesB, &result)
		indexO, indexA, indexB = o-1, a-1, b-1
	}

	chunk(indexO, indexA, indexB, len(runesO), len(runesA), len(runesB), runesO, runesA, runesB, linesA, linesB, &result)
	return result.String()
}

// matches returns a map of the non-crossing matches.
func matches(diffs []diffmatchpatch.Diff, runes []rune) map[int]int {
	matches := make(map[int]int)
	for _, d := range diffs {
		if d.Type != diffmatchpatch.DiffEqual {
			continue
		}

		for _, r := range d.Text {
			matches[int(r)] = indexOf(runes, r) + 1
		}
	}
	return matches
}

// nextMismatch searches for the next index where a or b is not equal to o.
func nextMismatch(indexO, indexA, indexB int, runesA, runesB []rune, matchesA, matchesB map[int]int) int {
	for i := 1; i <= len(runesA) && i <= len(runesB); i++ {
		a, okA := matchesA[indexO+i]
		b, okB := matchesB[indexO+i]

		if !okA || a != indexA+i || !okB || b != indexB+i {
			return i
		}
	}
	return 0
}

// nextMatch searches for the next index where a and b are equal to o.
func nextMatch(indexO int, runesO []rune, matchesA, matchesB map[int]int) (int, int, int) {
	for o := indexO + 1; o <= len(runesO); o++ {
		a, okA := matchesA[o]
		b, okB := matchesB[o]

		if okA && okB {
			return o, a, b
		}
	}
	return 0, 0, 0
}

// chunk merges the lines from o, a, and b into a single text.
func chunk(indexO, indexA, indexB, o, a, b int, runesO, runesA, runesB []rune, linesA, linesB []string, result *strings.Builder) {
	chunkO := buildChunk(linesA, runesO[indexO:o])
	chunkA := buildChunk(linesA, runesA[indexA:a])
	chunkB := buildChunk(linesB, runesB[indexB:b])

	switch {
	case chunkA == chunkB:
		fmt.Fprint(result, chunkO)
	case chunkO == chunkA:
		fmt.Fprint(result, chunkB)
	case chunkO == chunkB:
		fmt.Fprint(result, chunkA)
	default:
		fmt.Fprintf(result, "%s\n%s%s\n%s%s\n", Sep1, chunkA, Sep2, chunkB, Sep3)
	}
}

// indexOf returns the index of the first occurance of the given value.
func indexOf(runes []rune, value rune) int {
	for i, r := range runes {
		if r == value {
			return i
		}
	}
	return -1
}

// buildChunk assembles the lines of the chunk into a string.
func buildChunk(lines []string, runes []rune) string {
	var chunk strings.Builder
	for _, r := range runes {
		fmt.Fprint(&chunk, lines[int(r)])
	}
	return chunk.String()
}
