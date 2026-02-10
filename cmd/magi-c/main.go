package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/flily/magi-c/coder"
	"github.com/flily/magi-c/context"
)

func translateFile(c *coder.Coder, filename string) error {
	err := c.ParseFile(filename)
	if err != nil {
		switch e := err.(type) {
		case *context.Diagnostic:
			fmt.Printf("Syntax error:\n%s\n", e)

		default:
			fmt.Printf("Error:\n%s\n", err)
		}

		return err
	}

	err = c.Check(filename)
	if err != nil {
		fmt.Printf("Error:\n%s\n", err)
		return err
	}

	outputFilename := path.Base(coder.OutputFilename(filename))
	fmt.Printf("%s -> %s", filename, outputFilename)
	err = c.OutputToFile(filename, outputFilename)
	if err != nil {
		fmt.Printf("    failed\n")
		fmt.Printf("Error:\n%s\n", err)
		return err
	}

	fmt.Printf("    ok\n")
	return nil
}

func doTranslate(args []string) error {
	set := flag.NewFlagSet("translate", flag.ExitOnError)
	source := set.String("source", ".", "source base directory")
	output := set.String("output", "output", "output base directory")
	_ = set.Parse(args)

	c := coder.NewCoder(*source, *output)

	base := "."

	if set.NArg() > 0 {
		base = set.Arg(0)
	}

	stat, err := os.Stat(base)
	if err != nil {
		return err
	}

	if stat.IsDir() {

	} else {
		err = translateFile(c, base)
	}

	return err
}

func doBuild(args []string) error {
	return nil
}

var commandMap = map[string]func([]string) error{
	"translate": doTranslate,
	"build":     doBuild,
}

func showCommands() {
	fmt.Println("Available commands:")
	for cmd := range commandMap {
		fmt.Printf(" - %s\n", cmd)
	}
}

func main() {
	flag.Parse()

	if flag.NArg() <= 0 {
		flag.Usage()
		showCommands()
		return
	}

	command := flag.Arg(0)
	entry, found := commandMap[command]
	if !found {
		fmt.Printf("Unknown command: %s\n", command)
		showCommands()
		return
	}

	args := flag.Args()[1:]
	err := entry(args)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
