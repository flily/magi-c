package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/flily/magi-c/coder"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) <= 0 {
		flag.Usage()
		return
	}

	filename := args[0]
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("error reading file\n%s\n", err)
		return
	}

	coder, err := coder.NewCoderFromBinary(content, filename)
	if err != nil {
		fmt.Printf("parse file '%s' error\n%s\n", filename, err)
		return
	}

	fmt.Printf("parsed success %+v\n", coder)
}
