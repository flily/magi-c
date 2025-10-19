package ast

import (
	"github.com/flily/magi-c/context"
)

type PreprocessorDirectiveType int

const (
	PreprocessorDirectiveUnknown PreprocessorDirectiveType = iota
	PreprocessorDirectiveInclude
)

type PreprocessorDirectiveInfo struct {
	Command string
	Type    PreprocessorDirectiveType
}

var preprocessorDirectives = []*PreprocessorDirectiveInfo{
	{"include", PreprocessorDirectiveInclude},
}

func GetPreprocessorDirectiveInfo(command string) *PreprocessorDirectiveInfo {
	for _, info := range preprocessorDirectives {
		if info.Command == command {
			return info
		}
	}

	return nil
}

type PreprocessorInclude struct {
	TerminalNode
	Hash     *context.Context
	Command  *context.Context
	LBracket *context.Context
	Content  *context.Context
	RBracket *context.Context
}

func (p *PreprocessorInclude) Type() NodeType {
	return NodePreprocessorInclude
}
