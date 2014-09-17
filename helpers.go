package sor

import (
	"bufio"
)

// read at most n bytes from reader.
func readNBytes(r *bufio.Reader, n int) (buf []byte, err error) {
	buf = make([]byte, n)

	for i := 0; i < n; i++ {
		buf[i], err = r.ReadByte()
		if err != nil {
			return buf, err
		}
	}

	return buf, nil
}
