package sor

import (
	"bufio"
	"bytes"
)

type SupplierParamBlock struct {
	ID              string
	Supplier        string
	MainframeID     string
	MainframeSerial string
	OpticalID       string
	OpticalSN       string
	SoftwareVer     string
	Other           string
}

func parseSupplier(r *bufio.Reader, s *SOR) error {
	bk := SupplierParamBlock{ID: "SupParams"}

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

	sn, _ := bkBuf.ReadBytes('\x00')
	bk.Supplier = cleanString(sn)

	mfid, _ := bkBuf.ReadBytes('\x00')
	bk.MainframeID = cleanString(mfid)

	otdr, _ := bkBuf.ReadBytes('\x00')
	bk.MainframeSerial = cleanString(otdr)

	omid, _ := bkBuf.ReadBytes('\x00')
	bk.OpticalID = cleanString(omid)

	omsn, _ := bkBuf.ReadBytes('\x00')
	bk.OpticalSN = cleanString(omsn)

	sr, _ := bkBuf.ReadBytes('\x00')
	bk.SoftwareVer = cleanString(sr)

	ot, _ := bkBuf.ReadBytes('\x00')
	bk.Other = cleanString(ot)

	s.SupplerParams = bk

	return nil
}
