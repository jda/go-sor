package sor

import (
	"bufio"
	"bytes"
	"encoding/binary"
	//"fmt"
	"strconv"
	"time"
)

type Fixed struct {
	ID                       string
	Timestamp                time.Time
	Units                    string
	ActualWavelength         int16
	AcquisitionOffset        int32
	AcquisitionDistance      int32
	TotalPulses              int16
	PulsesUsed               []int16
	DataSpacing              []int32
	PointsPerPulse           []int32
	GroupIndex               float32
	Backscatter              float32
	NumberOfAverages         int32
	AveragingTime            uint16
	AcquisitionRange         int32
	AcquisitionRangeDistance int32
	FrontPanelOffset         int32
	NoiseFloor               float32
	NoiseFloorScale          int
	PowerOffset              uint16
	LossThreshold            uint16
	ReflectanceThreshold     uint16
	EndOfFiberThreshold      uint16
	TraceType                string
	WindowCoordinates        []int32
}

func parseFixed(r *bufio.Reader, s *SOR) error {
	bk := Fixed{ID: "FxdParams"}

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

	// if vs, need to get header out of the way
	if s.Version == SORv2 {
		_, _ = bkBuf.ReadBytes('\x00')
	}

	var ts uint32
	binary.Read(bkBuf, binary.LittleEndian, &ts)
	bk.Timestamp = time.Unix(int64(ts), 0)

	ud := make([]byte, 2)
	_, _ = bkBuf.Read(ud)
	bk.Units = string(ud)

	binary.Read(bkBuf, binary.LittleEndian, &bk.ActualWavelength)

	if s.Version == SORv2 {
		binary.Read(bkBuf, binary.LittleEndian, &bk.AcquisitionOffset)
	}

	binary.Read(bkBuf, binary.LittleEndian, &bk.AcquisitionDistance)

	binary.Read(bkBuf, binary.LittleEndian, &bk.TotalPulses)

	num_pulses := int(bk.TotalPulses)
	var pulses []int16
	for i := 0; i < num_pulses; i++ {
		var pwu int16
		binary.Read(bkBuf, binary.LittleEndian, &pwu)
		pulses = append(pulses, pwu)
	}
	bk.PulsesUsed = pulses

	var spacing []int32
	for i := 0; i < num_pulses; i++ {
		var ds int32
		binary.Read(bkBuf, binary.LittleEndian, &ds)
		spacing = append(spacing, ds)
	}
	bk.DataSpacing = spacing

	var datapoints []int32
	for i := 0; i < num_pulses; i++ {
		var nppw int32
		binary.Read(bkBuf, binary.LittleEndian, &nppw)
		datapoints = append(datapoints, nppw)
	}
	bk.PointsPerPulse = datapoints

	var gi int32
	binary.Read(bkBuf, binary.LittleEndian, &gi)
	gi_str := strconv.Itoa(int(gi))
	gi_end := string(gi_str[0]) + "." + gi_str[1:]
	gi_float, _ := strconv.ParseFloat(gi_end, 32)
	bk.GroupIndex = float32(gi_float)

	var bs int16
	binary.Read(bkBuf, binary.LittleEndian, &bs)
	bs_str := strconv.Itoa(int(bs))
	bs_end := bs_str[0:2] + "." + string(bs_str[2])
	bs_float, _ := strconv.ParseFloat(bs_end, 32)
	bk.Backscatter = float32(bs_float)

	binary.Read(bkBuf, binary.LittleEndian, &bk.NumberOfAverages)

	if s.Version == SORv2 {
		binary.Read(bkBuf, binary.LittleEndian, &bk.AveragingTime)
	}

	binary.Read(bkBuf, binary.LittleEndian, &bk.AcquisitionRange)

	if s.Version == SORv2 {
		binary.Read(bkBuf, binary.LittleEndian, &bk.AcquisitionDistance)
	}

	binary.Read(bkBuf, binary.LittleEndian, &bk.FrontPanelOffset)

	var nf uint32
	binary.Read(bkBuf, binary.LittleEndian, &nf)
	nf_str := strconv.Itoa(int(nf))
	nf_end := "-" + nf_str[0:2] + "." + string(nf_str[2:])
	nf_float, _ := strconv.ParseFloat(nf_end, 32)
	bk.NoiseFloor = float32(nf_float)

	//binary.Read(bkBuf, binary.LittleEndian, &bk.NoiseFloor)

	binary.Read(bkBuf, binary.LittleEndian, &bk.NoiseFloorScale)

	binary.Read(bkBuf, binary.LittleEndian, &bk.PowerOffset)

	binary.Read(bkBuf, binary.LittleEndian, &bk.LossThreshold)

	binary.Read(bkBuf, binary.LittleEndian, &bk.ReflectanceThreshold)

	binary.Read(bkBuf, binary.LittleEndian, &bk.EndOfFiberThreshold)

	if s.Version == SORv2 {
		tt, _ := bkBuf.ReadBytes('\x00')
		bk.TraceType = cleanString(tt)
	}

	// window coordinates
	for i := 1; i <= 4; i++ {
		var wc int32
		binary.Read(bkBuf, binary.LittleEndian, &wc)
		bk.WindowCoordinates = append(bk.WindowCoordinates, wc)
	}

	s.FixedParams = bk

	return nil
}
