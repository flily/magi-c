package ast

import (
	"github.com/flily/magi-c/context"
)

type PreprocessorDirectiveType int

const (
	PreprocessorDirectiveUnknown PreprocessorDirectiveType = iota
	PreprocessorDirectiveInclude
	PreprocessorDirectiveInline
)

type PreprocessorDirectiveInfo struct {
	Command string
	Type    PreprocessorDirectiveType
}

var preprocessorDirectives = []*PreprocessorDirectiveInfo{
	{"include", PreprocessorDirectiveInclude},
	{"inline", PreprocessorDirectiveInline},
}

func GetPreprocessorDirectiveInfo(command string) *PreprocessorDirectiveInfo {
	for _, info := range preprocessorDirectives {
		if info.Command == command {
			return info
		}
	}

	return nil
}

type PreprocessorCommon struct {
	NonTerminalNode
	Hash    *context.Context
	Command *context.Context
}

type PreprocessorInclude struct {
	PreprocessorCommon
	LBracket *context.Context
	Content  *context.Context
	RBracket *context.Context
}

func NewPreprocessorInclude(hash *context.Context, command *context.Context, lbracket *context.Context, content *context.Context, rbracket *context.Context) *PreprocessorInclude {
	p := &PreprocessorInclude{
		PreprocessorCommon: PreprocessorCommon{
			Hash:    hash,
			Command: command,
		},
		LBracket: lbracket,
		Content:  content,
		RBracket: rbracket,
	}

	p.Init(p)
	return p
}

func (p *PreprocessorInclude) Type() NodeType {
	return NodePreprocessorInclude
}

func (p *PreprocessorInclude) Context() *context.Context {
	return context.Join(p.Hash, p.Command, p.LBracket, p.Content, p.RBracket)
}

type PreprocessorInline struct {
	PreprocessorCommon
	CodeType    *context.Context
	Content     *context.Context
	HashEnd     *context.Context
	CommandEnd  *context.Context
	CodeTypeEnd *context.Context
}

func NewPreprocessorInline(hash *context.Context, command *context.Context, codeType *context.Context, content *context.Context, hashEnd *context.Context, commandEnd *context.Context, codeTypeEnd *context.Context) *PreprocessorInline {
	p := &PreprocessorInline{
		PreprocessorCommon: PreprocessorCommon{
			Hash:    hash,
			Command: command,
		},
		CodeType:    codeType,
		Content:     content,
		HashEnd:     hashEnd,
		CommandEnd:  commandEnd,
		CodeTypeEnd: codeTypeEnd,
	}

	p.Init(p)
	return p
}

func (p *PreprocessorInline) Type() NodeType {
	return NodePreprocessorInline
}

func (p *PreprocessorInline) Context() *context.Context {
	return context.Join(p.Hash, p.Command, p.CodeType, p.Content, p.HashEnd, p.CommandEnd, p.CodeTypeEnd)
}

func (p *PreprocessorInline) Empty() bool {
	return p.Content == nil
}
