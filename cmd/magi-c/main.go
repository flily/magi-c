package main

import (
	"flag"
	"fmt"

	"github.com/flily/magi-c/ast"
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
		switch e := err.(type) {
		case *ast.Error:
			fmt.Printf("Syntax error:\n%s\n", e)

		default:
			fmt.Printf("Error:\n%s\n", err)
		}

		return
	}

	fmt.Printf("parsed success %+v\n", coder)
}
