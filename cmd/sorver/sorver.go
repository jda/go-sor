package main

import (
	"../.."
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatalf("Usage: sorver FILE.sor")
	}
	fname := args[0]

	f, err := os.Open(fname)
	if err != nil {
		log.Fatalf("Error opening file %s because %s", fname, err)
	}

	rd := bufio.NewReader(f)
	ver, err := sor.Identify(rd)
	if err != nil {
		log.Fatalf("Error identifying file: %s", err)
	}
	fmt.Println(ver)
}
