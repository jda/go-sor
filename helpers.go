package sor

import (
	"bufio"
)

// read at most n bytes from reader.
func readNBytes(r *bufio.Reader, n int) (buf []byte, err error) {
	buf = make([]byte, n)
	_, err = r.Read(buf)

	if err != nil {
		return buf, err
	}

	return buf, nil
}

// get named block from block array
func getBlock(name string, s *SOR) (b Block, err error) {
	for _, e := range s.Blocks {
		if e.ID == name {
			return e, nil
		}
	}
	var emptyBlock Block
	return emptyBlock, ErrNoBlock
}
