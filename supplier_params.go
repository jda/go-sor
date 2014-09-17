package sor

import (
	"bufio"
	"bytes"
)

type Supplier struct {
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
	bk := Supplier{ID: "SupParams"}

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
	bk.Supplier = string(bytes.TrimSpace(sn))

	mfid, _ := bkBuf.ReadBytes('\x00')
	bk.MainframeID = string(bytes.TrimSpace(mfid))

	otdr, _ := bkBuf.ReadBytes('\x00')
	bk.MainframeSerial = string(bytes.TrimSpace(otdr))

	omid, _ := bkBuf.ReadBytes('\x00')
	bk.OpticalID = string(bytes.TrimSpace(omid))

	omsn, _ := bkBuf.ReadBytes('\x00')
	bk.OpticalSN = string(bytes.TrimSpace(omsn))

	sr, _ := bkBuf.ReadBytes('\x00')
	bk.SoftwareVer = string(bytes.TrimSpace(sr))

	ot, _ := bkBuf.ReadBytes('\x00')
	bk.Other = string(bytes.TrimSpace(ot))

	s.SupplerParams = bk

	return nil
}
