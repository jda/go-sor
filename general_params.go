package sor

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

type General struct {
	ID             string
	Language       string
	CableID        string
	FiberID        string
	FiberType      int
	Wavelength     int
	OriginatingLoc string
	TerminatingLoc string
	CableCode      string
	Condition      string
	Offset         int32
	OffsetDistance int32
	Operator       string
	Comment        string
}

func parseGeneral(r *bufio.Reader, s *SOR) error {
	g := General{ID: "GenParams"}

	header, err := getBlock("GenParams", s)
	if err != nil {
		return err
	}

	// capture entire block
	gBytes, err := readNBytes(r, int(header.Bytes))
	if err != nil {
		return ErrIncompleteBlock
	}
	gBuf := bytes.NewBuffer(gBytes)

	// if v2, need to get header out of the way
	if s.Version == SORv2 {
		_, _ = gBuf.ReadBytes('\x00')
	}

	lc := make([]byte, 2)
	_, err = gBuf.Read(lc)
	g.Language = string(lc)

	cid, _ := gBuf.ReadBytes('\x00')
	g.CableID = string(bytes.TrimSpace(cid))

	fid, _ := gBuf.ReadBytes('\x00')
	g.FiberID = string(bytes.TrimSpace(fid))

	// fiber type not in v1?
	if s.Version == SORv2 {
		var ft uint16
		binary.Read(gBuf, binary.LittleEndian, &ft)
		g.FiberType = int(ft)
	}

	var nw uint16
	binary.Read(gBuf, binary.LittleEndian, &nw)
	g.Wavelength = int(nw)

	ol, _ := gBuf.ReadBytes('\x00')
	g.OriginatingLoc = string(bytes.TrimSpace(ol))

	tl, _ := gBuf.ReadBytes('\x00')
	g.TerminatingLoc = string(bytes.TrimSpace(tl))

	cc, _ := gBuf.ReadBytes('\x00')
	g.CableCode = string(bytes.TrimSpace(cc))

	cdf := make([]byte, 2)
	_, err = gBuf.Read(cdf)
	g.Condition = string(cdf)

	binary.Read(gBuf, binary.LittleEndian, &g.Offset)

	// offset distance not in v1?
	if s.Version == SORv2 {
		binary.Read(gBuf, binary.LittleEndian, &g.OffsetDistance)
	}

	op, _ := gBuf.ReadBytes('\x00')
	g.Operator = string(bytes.TrimSpace(op))

	cmt, _ := gBuf.ReadBytes('\x00')
	g.Comment = string(bytes.TrimSpace(cmt))

	s.GeneralParams = g

	return nil
}
