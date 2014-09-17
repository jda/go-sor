package sor

import (
	"bufio"
	"fmt"
)

type SOR struct {
	Version SORVersion
	Blocks  []Block
}

func Parse(r *bufio.Reader) (SOR, error) {
	// init any maps in SOR struct here
	var s SOR

	// get version
	ver, err := Identify(r)
	if err != nil {
		return s, fmt.Errorf("error identifying file", err)
	}
	if ver != SORv1 && ver != SORv2 {
		return s, fmt.Errorf("error identifying version")
	}
	s.Version = ver

	err = parseBlocks(r, &s)
	if err != nil {
		return s, err
	}

	fmt.Printf("SOR: %+v\n", s)

	return s, nil
}
