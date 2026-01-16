package main

import (
	"flag"
	"fmt"

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

	coder := coder.NewCoder(".", "output")
	err := coder.ParseFile(filename)
	if err != nil {
		fmt.Printf("parse file '%s' error\n%s\n", filename, err)
		return
	}

	fmt.Printf("parsed success %+v\n", coder)
}
