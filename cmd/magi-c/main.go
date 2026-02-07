package main

import (
	"flag"
	"fmt"
	"path"

	"github.com/flily/magi-c/coder"
	"github.com/flily/magi-c/context"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) <= 0 {
		flag.Usage()
		return
	}

	filename := args[0]

	c := coder.NewCoder(".", "output")
	err := c.ParseFile(filename)
	if err != nil {
		switch e := err.(type) {
		case *context.Diagnostic:
			fmt.Printf("Syntax error:\n%s\n", e)

		default:
			fmt.Printf("Error:\n%s\n", err)
		}

		return
	}

	err = c.Check(filename)
	if err != nil {
		fmt.Printf("Error:\n%s\n", err)
		return
	}

	outputFilename := path.Base(coder.OutputFilename(filename))
	fmt.Printf("%s -> %s", filename, outputFilename)
	err = c.OutputToFile(filename, outputFilename)
	if err != nil {
		fmt.Printf("    failed\n")
		fmt.Printf("Error:\n%s\n", err)
		return
	}

	fmt.Printf("    ok\n")
}
