package test

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hofstadter-io/hof/lib/cuetils"
)

func printTests(suites []Suite, stats bool) {

	totalTests := 0
	totalStats := Stats{}

	for _, S := range suites {
		S.Stats.Time = 0
		for _, t := range S.Tests {
			S.Stats.add(t.Stats)
		}

		A := S.Value.Attribute("test")
		as := []string{}
		for k, v := range cuetils.AttrToMap(A) {
			if v != "" {
				as = append(as, fmt.Sprintf("%s=%s", k, v))
			} else {
				as = append(as, fmt.Sprintf("%s", k))
			}
		}
		a := strings.Join(as, " ")

		st := ""
		if stats {
			s := S.Stats
			st = fmt.Sprintf("  %d/%d/%d ~ %v", s.Pass,s.Fail,s.Skip,s.Time)
		}

		lt := fmt.Sprintf("[%d]", len(S.Tests))

		fmt.Printf( "[suite]     %-16s %6s %-32v%s\n", S.Name, lt, a, st)

		totalTests += len(S.Tests)

		for _, t := range S.Tests {
			A := t.Value.Attribute("test")
			as := []string{}
			for k, v := range cuetils.AttrToMap(A) {

				// if key is a knownTester, put it first
				known := false
				for _, kT := range knownTesters {
					if k == kT {
						as = append([]string{"["+k+"]"}, as...)
						known = true
						break
					}
				}
				if known {
					continue
				}

				// otherwise just add
				if v != "" {
					as = append(as, fmt.Sprintf("%s=%s", k, v))
				} else {
					as = append(as, fmt.Sprintf("%s", k))
				}
			}
			sort.Strings(as[1:])
			a := strings.Join(as, " ")
			st := ""
			if stats {
				totalStats.add(t.Stats)
				s := t.Stats
				st = fmt.Sprintf("  %d/%d/%d ~ %v", s.Pass,s.Fail,s.Skip,s.Time)
			}
			fmt.Printf( "[tester]      %-16s      %-32v%s\n", t.Name, a, st)
		}

		fmt.Println()

	}

	fmt.Println()
	fmt.Println("Total Suites: ", len(suites))
	fmt.Println("Total Tests:  ", totalTests)
	if stats {
		s := totalStats
		st := fmt.Sprintf("%d/%d/%d ~ %v", s.Pass,s.Fail,s.Skip,s.Time)
		fmt.Println("Total Stats:  ", st)
	}
}

