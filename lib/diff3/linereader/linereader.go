package linereader

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"strings"
)

type LineReader struct {
	reader io.Reader
	sample []byte
}

// ErrBinaryContent is returned when attempting to read lines
// out of a stream that does not look like text
var ErrBinaryContent = errors.New("Cannot read binary content")

func NewLineReader(reader io.Reader) *LineReader {
	return &LineReader{
		reader: reader,
	}
}

func (lr *LineReader) Read(p []byte) (int, error) {
	i, err := lr.reader.Read(p)
	r := 512 - len(lr.sample)
	if r > 0 {
		lr.sample = append(lr.sample, p[:min(r, i)]...)
		r = 512 - len(lr.sample)
		if r == 0 || err == io.EOF {
			if !isText(lr.sample) {
				return i, ErrBinaryContent
			}
		}
	}
	return i, err
}

func (lr *LineReader) GetLines() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(lr)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()

}

func isText(b []byte) bool {
	if strings.Contains(http.DetectContentType(b), "text") || len(b) == 0 {
		return true
	}
	return false
}

func GetLines(r io.Reader) ([]string, error) {
	return NewLineReader(r).GetLines()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
