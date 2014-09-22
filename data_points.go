package sor

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

type DataPointsBlock struct {
	ID               string
	TotalDataPoints  int32
	ScaleFactorsUsed int16
	ScaleFactors     []ScaleFactorBlock
}

type ScaleFactorBlock struct {
	TotalDataPoints int32
	ScaleFactor     int16
	DataPoints      []uint16
}

func parseDataPoints(r *bufio.Reader, s *SOR) error {
	bk := DataPointsBlock{ID: "DataPts"}

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

	binary.Read(bkBuf, binary.LittleEndian, &bk.TotalDataPoints)
	binary.Read(bkBuf, binary.LittleEndian, &bk.ScaleFactorsUsed)

	for i := int16(0); i < bk.ScaleFactorsUsed; i++ {
		d := ScaleFactorBlock{}

		binary.Read(bkBuf, binary.LittleEndian, &d.TotalDataPoints)
		binary.Read(bkBuf, binary.LittleEndian, &d.ScaleFactor)

		for j := int32(0); j < d.TotalDataPoints; j++ {
			var dp uint16
			binary.Read(bkBuf, binary.LittleEndian, &dp)
			d.DataPoints = append(d.DataPoints, dp)
		}

		bk.ScaleFactors = append(bk.ScaleFactors, d)
	}

	s.DataPoints = bk

	return nil
}
