package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/flily/magi-c/coder"
	"github.com/flily/magi-c/context"
)

func translateFile(c *coder.Coder, filename string) error {
	indexName, err := c.ParseFile(filename)
	if err != nil {
		switch e := err.(type) {
		case *context.Diagnostic:
			fmt.Printf("Syntax error:\n%s\n", e)

		default:
			fmt.Printf("Parse error:\n%s\n", err)
		}

		return err
	}

	err = c.Check(indexName)
	if err != nil {
		fmt.Printf("Check error:\n%s\n", err)
		return err
	}

	outputFilename := c.OutputFilename(indexName)
	fmt.Printf("%s -> [%s] %s", filename, indexName, outputFilename)
	err = c.Output(indexName)
	if err != nil {
		fmt.Printf("    failed\n")
		fmt.Printf("Output error:\n%s\n", err)
		return err
	}

	fmt.Printf("    ok\n")
	return nil
}

func translateDirectory(c *coder.Coder, base string) error {
	err := filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, coder.DefaultSourceSuffix) {
			return nil
		}

		return translateFile(c, path)
	})

	if err != nil {
		return err
	}

	return nil
}

func doTranslate(args []string) error {
	set := flag.NewFlagSet("translate", flag.ExitOnError)
	output := set.String("output", "output", "output base directory")
	_ = set.Parse(args)

	base := "."

	if set.NArg() > 0 {
		base = set.Arg(0)
	}

	stat, err := os.Stat(base)
	if err != nil {
		return err
	}

	c := coder.NewCoder(base, *output)

	if stat.IsDir() {
		err = translateDirectory(c, base)

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
