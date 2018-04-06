package util

import (
	"bufio"
	"encoding/csv"
	"io"
)

type reader struct {
	r *bufio.Reader
}

const (
	rByte byte = 13 // the byte that corresponds to the '\r' rune.
	nByte byte = 10 // the byte that corresponds to the '\n' rune.
)

// Read replaces CR line endings in the source reader with LF line endings if the CR is not followed by a LF.
func (r reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	bn, err := r.r.Peek(1)
	for i, b := range p {
		// if the current byte is a CR and the next byte is NOT a LF then replace the current byte with a LF
		if j := i + 1; b == rByte && ((j < len(p) && p[j] != nByte) || (len(bn) > 0 && bn[0] != nByte)) {
			p[i] = nByte
		}
	}
	return
}

func NewCSVReader(r io.Reader) *csv.Reader {
	bufReader := bufio.NewReader(r)
	return csv.NewReader(reader{bufReader})
}
