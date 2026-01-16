package coder

import (
	"os"
	"path/filepath"
)

const (
	DefaultOutputBase = "output"
)

type Coder struct {
	SourceBase string
	OutputBase string
	Refs       *Cache
}

func NewCoder(sourceBase string, outputBase string) *Coder {
	c := &Coder{
		SourceBase: sourceBase,
		OutputBase: outputBase,
		Refs:       NewCache(),
	}

	return c
}

func (c *Coder) ParseAll() error {
	err := filepath.Walk(c.SourceBase, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		e := c.ParseFile(path)
		if e != nil {
			return e
		}

		return nil
	})

	return err
}

func (c *Coder) ParseFileContent(filename string, content []byte) error {
	doc, err := c.Refs.ParseContent(content, filename)
	if err != nil {
		return err
	}

	c.Refs.Add(filename, doc)

	return nil
}

func (c *Coder) ParseFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return c.ParseFileContent(filename, content)
}
