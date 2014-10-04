package main

import (
	"../.."
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatalf("Usage: sor2json FILE.sor")
	}
	fname := args[0]

	f, err := os.Open(fname)
	if err != nil {
		log.Fatalf("Error opening file %s because %s", fname, err)
	}

	rd := bufio.NewReader(f)
	data, err := sor.Parse(rd)
	if err != nil {
		log.Fatalf("Error parsing file: %s", err)
	}

	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Fatalf("Could not format SOR: %s", err)
	}

	os.Stdout.Write(b)
}
