package sor

import (
	"encoding/json"
	"fmt"
	"os"
)

func prettyPrint(v interface{}) {
	out := os.Stdout

	js, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		fmt.Fprintf(out, "Could not pretty print: %s\n", err)
		return
	}

	fmt.Fprintf(out, "%s\n", js)
	return
}
