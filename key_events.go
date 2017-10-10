package sor

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

type KeyEventsBlock struct {
	ID          string
	TotalEvents int16
	Events      []KeyEvent
	TotalLoss	int16
}

type KeyEvent struct {
	Number                   int16
	PropagationTime          int32
	AttenuationCoefficient   int16
	Loss                     int16
	Reflectance              int32
	Code                     string
	LossMeasurementTechnique string
	MarkerLocations          []int32
	Comment                  string
}

func parseKeyEvents(r *bufio.Reader, s *SOR) error {
	bk := KeyEventsBlock{ID: "KeyEvents"}

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

	binary.Read(bkBuf, binary.LittleEndian, &bk.TotalEvents)

	for i := int16(0); i < bk.TotalEvents; i++ {
		e := KeyEvent{}

		binary.Read(bkBuf, binary.LittleEndian, &e.Number)
		binary.Read(bkBuf, binary.LittleEndian, &e.PropagationTime)
		binary.Read(bkBuf, binary.LittleEndian, &e.AttenuationCoefficient)
		binary.Read(bkBuf, binary.LittleEndian, &e.Loss)
		binary.Read(bkBuf, binary.LittleEndian, &e.Reflectance)

		ec := make([]byte, 6)
		_, _ = bkBuf.Read(ec)
		e.Code = cleanString(ec)

		lmt := make([]byte, 2)
		_, _ = bkBuf.Read(lmt)
		e.LossMeasurementTechnique = cleanString(lmt)

		for j := 0; j < 5; j++ {
			var ml int32
			binary.Read(bkBuf, binary.LittleEndian, &ml)
			e.MarkerLocations = append(e.MarkerLocations, ml)
		}

		cmt, _ := bkBuf.ReadBytes('\x00')
		e.Comment = cleanString(cmt)

		bk.Events = append(bk.Events, e)
	}

	s.KeyEvents = bk
	binary.Read(bkBuf, binary.LittleEndian, &bk.TotalLoss)
	return nil
}
