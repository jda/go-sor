package sor

import (
	"bufio"
	"encoding/binary"
)

type SORVersion int

const (
	SORvUnknown SORVersion = 0
	SORv1       SORVersion = 1
	SORv2       SORVersion = 2
)

func Identify(rd *bufio.Reader) (ver SORVersion, e error) {
	pb, err := rd.Peek(4)
	if err != nil {
		return SORvUnknown, err
	}

	// starts with map, is v2
	checkMap := string(pb[0:3])
	if checkMap == "Map" {
		return SORv2, nil
	}

	// starts with 100, is v1
	checkNum, _ := binary.Uvarint(pb[0:4])
	if checkNum == 100 {
		return SORv1, nil
	}

	// unknown
	return SORvUnknown, nil
}
