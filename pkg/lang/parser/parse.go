package parser

import (
	"io/ioutil"
	"os"
)

func ParseHofFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	P := newParser(filename, b, opts...)
	P.cur.globalStore["filename"] = filename

	return P.parse(g)
}
