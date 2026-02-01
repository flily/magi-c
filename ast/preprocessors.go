package ast

import (
	"strings"

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

func (p *PreprocessorCommon) declarationNode() {}

func (p *PreprocessorCommon) statementNode() {}

type PreprocessorInclude struct {
	PreprocessorCommon
	LBracketCtx *context.Context
	ContentCtx  *context.Context
	RBracketCtx *context.Context
	Content     string
	LBracket    string
	RBracket    string
}

func NewPreprocessorInclude(hash *context.Context, command *context.Context, lbracket *context.Context, content *context.Context, rbracket *context.Context) *PreprocessorInclude {
	p := &PreprocessorInclude{
		PreprocessorCommon: PreprocessorCommon{
			Hash:    hash,
			Command: command,
		},
		LBracketCtx: lbracket,
		ContentCtx:  content,
		RBracketCtx: rbracket,
		Content:     content.Content(),
		LBracket:    lbracket.Content(),
		RBracket:    rbracket.Content(),
	}

	p.Init(p)
	return p
}

func ASTBuildIncludeAngle(name string) *PreprocessorInclude {
	p := &PreprocessorInclude{}
	p.Content = name
	p.LBracket = "<"
	p.RBracket = ">"

	return p
}

func ASTBuildIncludeQuote(name string) *PreprocessorInclude {
	p := &PreprocessorInclude{}
	p.Content = name
	p.LBracket = `"`
	p.RBracket = `"`

	return p
}

func (p *PreprocessorInclude) Type() TokenType {
	return NodePreprocessorInclude
}

func (p *PreprocessorInclude) EqualTo(_ context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(p, other)
	if err != nil {
		return err
	}

	if p.Content != o.Content {
		return p.ContentCtx.Error("wrong include content, expect '%s', got '%s'", o.Content, p.Content).With(o.Content)
	}

	if p.LBracket != o.LBracket {
		ctx := context.Join(p.LBracketCtx, p.RBracketCtx)
		_, _, col1 := p.LBracketCtx.Position()
		_, _, col2 := p.RBracketCtx.Position()
		spaces := strings.Repeat(" ", col2-col1-1)

		return ctx.Error("wrong include bracket, expect '%s' and '%s', got '%s' and '%s'",
			o.LBracket, o.RBracket, p.LBracket, p.RBracket).With("%s%s%s", o.LBracket, spaces, o.RBracket)
	}

	return nil
}

func (p *PreprocessorInclude) Context() *context.Context {
	return context.Join(p.Hash, p.Command, p.LBracketCtx, p.ContentCtx, p.RBracketCtx)
}

type PreprocessorInline struct {
	PreprocessorCommon
	CodeTypeCtx *context.Context
	ContentCtx  *context.Context
	HashEnd     *context.Context
	CommandEnd  *context.Context
	CodeTypeEnd *context.Context
	CodeType    string
	Content     string
}

func NewPreprocessorInline(hash *context.Context, command *context.Context, codeType string, codeTypeCtx *context.Context, content string, contentCtx *context.Context, hashEnd *context.Context, commandEnd *context.Context, codeTypeEnd *context.Context) *PreprocessorInline {
	p := &PreprocessorInline{
		PreprocessorCommon: PreprocessorCommon{
			Hash:    hash,
			Command: command,
		},
		CodeTypeCtx: codeTypeCtx,
		ContentCtx:  contentCtx,
		CodeType:    codeType,
		Content:     content,
		HashEnd:     hashEnd,
		CommandEnd:  commandEnd,
		CodeTypeEnd: codeTypeEnd,
	}

	p.Init(p)
	return p
}

func ASTBuildInline(codeType string, content string) *PreprocessorInline {
	p := &PreprocessorInline{
		CodeType: codeType,
		Content:  content,
	}
	p.Init(p)

	return p
}

func (p *PreprocessorInline) statementNode() {}

func (p *PreprocessorInline) Type() TokenType {
	return NodePreprocessorInline
}

func (p *PreprocessorInline) EqualTo(_ context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(p, other)
	if err != nil {
		return err
	}

	if p.CodeType != o.CodeType {
		return p.CodeTypeCtx.Error("wrong inline code type, expect '%s', got '%s'", o.CodeType, p.CodeType).With("%s", o.CodeType)
	}

	if p.Content != o.Content {
		return p.ContentCtx.Error("wrong inline content").With(o.Content)
	}

	return nil
}

func (p *PreprocessorInline) Context() *context.Context {
	return context.Join(p.Hash, p.Command, p.CodeTypeCtx, p.ContentCtx, p.HashEnd, p.CommandEnd, p.CodeTypeEnd)
}

func (p *PreprocessorInline) Empty() bool {
	return p.ContentCtx == nil
}
