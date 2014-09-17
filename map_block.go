package sor

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

type Block struct {
	ID      string
	Version uint16
	Bytes   int32
}

// unsigned short & signed long
func parseBlocks(r *bufio.Reader, s *SOR) error {
	var mb Block
	mb.ID = "Map"

	// no map is 'Map' on v1 so skip past it
	if s.Version == SORv2 {
		_, _ = readNBytes(r, 4)
	}

	binary.Read(r, binary.LittleEndian, &mb.Version)
	binary.Read(r, binary.LittleEndian, &mb.Bytes)

	var blocks int16
	binary.Read(r, binary.LittleEndian, &blocks)

	s.Blocks = append(s.Blocks, mb)

	// parse remaining blocks
	for i := 1; i < int(blocks); i++ {
		var m Block
		block_name, _ := r.ReadBytes('\x00')
		m.ID = string(bytes.Trim(block_name, "\x00"))
		binary.Read(r, binary.LittleEndian, &m.Version)
		binary.Read(r, binary.LittleEndian, &m.Bytes)
		s.Blocks = append(s.Blocks, m)
	}

	return nil
}
