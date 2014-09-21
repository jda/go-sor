package sor

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

type General struct {
	ID                string
	Language          string
	CableID           string
	FiberID           string
	FiberType         int
	NominalWavelength int16
	OriginatingLoc    string
	TerminatingLoc    string
	CableCode         string
	Condition         string
	Offset            int32
	OffsetDistance    int32
	Operator          string
	Comment           string
}

func parseGeneral(r *bufio.Reader, s *SOR) error {
	bk := General{ID: "GenParams"}

	header, err := getBlock(bk.ID, s)
	if err != nil {
		return err
	}

	// capture entire block
	bkBytes, err := readNBytes(r, int(header.Bytes))
	if err != nil {
		return ErrIncompleteBlock
	}
	bkBuf := bytes.NewBuffer(bkBytes)

	// if v2, need to get header out of the way
	if s.Version == SORv2 {
		_, _ = bkBuf.ReadBytes('\x00')
	}

	lc := make([]byte, 2)
	_, _ = bkBuf.Read(lc)
	bk.Language = string(lc)

	cid, _ := bkBuf.ReadBytes('\x00')
	bk.CableID = cleanString(cid)

	fid, _ := bkBuf.ReadBytes('\x00')
	bk.FiberID = cleanString(fid)

	// fiber type not in v1?
	if s.Version == SORv2 {
		var ft uint16
		binary.Read(bkBuf, binary.LittleEndian, &ft)
		bk.FiberType = int(ft)
	}

	binary.Read(bkBuf, binary.LittleEndian, &bk.NominalWavelength)

	ol, _ := bkBuf.ReadBytes('\x00')
	bk.OriginatingLoc = cleanString(ol)

	tl, _ := bkBuf.ReadBytes('\x00')
	bk.TerminatingLoc = cleanString(tl)

	cc, _ := bkBuf.ReadBytes('\x00')
	bk.CableCode = cleanString(cc)

	cdf := make([]byte, 2)
	_, err = bkBuf.Read(cdf)
	bk.Condition = string(cdf)

	binary.Read(bkBuf, binary.LittleEndian, &bk.Offset)

	// offset distance not in v1?
	if s.Version == SORv2 {
		binary.Read(bkBuf, binary.LittleEndian, &bk.OffsetDistance)
	}

	op, _ := bkBuf.ReadBytes('\x00')
	bk.Operator = cleanString(op)

	cmt, _ := bkBuf.ReadBytes('\x00')
	bk.Comment = cleanString(cmt)

	s.GeneralParams = bk

	return nil
}
