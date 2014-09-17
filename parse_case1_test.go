package sor

import (
	"bufio"
	"os"
	"testing"
)

func TestParseCase1(t *testing.T) {
	fname := "test-v1-0.sor"
	path := "test_data" + RuneToAscii(os.PathSeparator) + fname
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("could not open test data: %s because %s", path, err)
	}

	rd := bufio.NewReader(f)
	Parse(rd)

}
