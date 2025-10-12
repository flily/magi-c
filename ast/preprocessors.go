package ast

import (
	"github.com/flily/magi-c/context"
)

type PreprocessorDirectiveType int

const (
	PreprocessorOneLine PreprocessorDirectiveType = iota
	PreprocessorBlock
)

type PreprocessorDirectiveInfo struct {
	Command string
	Type    PreprocessorDirectiveType
}

var preprocessorDirectives = []*PreprocessorDirectiveInfo{
	{"include", PreprocessorOneLine},
}

func GetPreprocessorDirectiveInfo(command string) *PreprocessorDirectiveInfo {
	for _, info := range preprocessorDirectives {
		if info.Command == command {
			return info
		}
	}

	return nil
}

type PreprocessorOneLineDirective struct {
	TerminalNode
	Hash      *context.Context
	Command   *context.Context
	Arguments []*context.Context
}

type PreprocessorInlineBlock struct {
	TerminalNode
	Hash      *context.Context
	Begin     *context.Context
	Arguments []*context.Context
	Content   string
	End       *context.Context
}
