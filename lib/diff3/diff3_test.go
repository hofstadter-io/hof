package diff3_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/hofstadter-io/hof/lib/diff3"
)

func rd(st string) io.Reader {
	return bytes.NewReader([]byte(st))
}

func compareReader(t *testing.T, a, b io.Reader) bool {
	abytes, err := ioutil.ReadAll(a)
	if err != nil {
		t.Fatal(err)
	}
	bbytes, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.Equal(abytes, bbytes)
}

func TestDiff3(t *testing.T) {

	const testDir = "./testdata"
	const generate = true // set to true to genereate the expected files *-m.txt and *-error.txt
	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		log.Fatal(err)
	}

	tests := make(map[string]bool)

	for _, f := range files {
		if len(f.Name()) > 3 {
			tests[f.Name()[:3]] = true
		}
	}

	for k, _ := range tests {
		func() {
			a, err := os.Open(fmt.Sprintf("%s/%s-a.txt", testDir, k))
			if err != nil {
				log.Fatal(err)
			}
			defer a.Close()
			b, err := os.Open(fmt.Sprintf("%s/%s-b.txt", testDir, k))
			if err != nil {
				log.Fatal(err)
			}
			defer b.Close()
			o, err := os.Open(fmt.Sprintf("%s/%s-o.txt", testDir, k))
			if err != nil {
				log.Fatal(err)
			}
			defer o.Close()
			var m io.ReadCloser
			if !generate {
				m, err = os.Open(fmt.Sprintf("%s/%s-m.txt", testDir, k))
				if err != nil {
					log.Fatal(err)
				}
				defer m.Close()
			}

			mr, mergeError := diff3.Merge(a, o, b, true, "A", "B")

			if generate {
				m, err := os.OpenFile(
					fmt.Sprintf("%s/%s-m.txt", testDir, k),
					os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
					0666,
				)
				if err != nil {
					t.Fatal(err)
				}
				defer m.Close()
				if mergeError == nil {
					_, err = io.Copy(m, mr.Result)
					if err != nil {
						t.Fatal(err)
					}
				}
				if mergeError != nil {
					err = ioutil.WriteFile(fmt.Sprintf("%s/%s-error.txt", testDir, k), []byte(mergeError.Error()), 0666)
					if err != nil {
						t.Fatal(err)
					}
				}
			} else {
				var expectError = false
				if _, err := os.Stat(fmt.Sprintf("%s/%s-error.txt", testDir, k)); !os.IsNotExist(err) {
					expectError = true
				}
				if mergeError != nil && !expectError {
					t.Fatalf("Did not expect merge error: %s", err)
				}
				if !expectError && !compareReader(t, mr.Result, m) {
					t.Fatalf("Test #%s does not match expected result", k)
				}
			}
		}()
	}
}
