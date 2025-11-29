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

func (p *PreprocessorCommon) declarationNode() {}

func (p *PreprocessorCommon) statementNode() {}

type PreprocessorInclude struct {
	PreprocessorCommon
	LBracketCtx *context.Context
	ContentCtx  *context.Context
	RBracketCtx *context.Context
	Content     string
	Bracket     string
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
		Bracket:     lbracket.Content(),
	}

	p.Init(p)
	return p
}

func ASTBuildIncludeAngle(name string) *PreprocessorInclude {
	p := &PreprocessorInclude{}
	p.Content = name
	p.Bracket = "<"

	return p
}

func ASTBuildIncludeQuote(name string) *PreprocessorInclude {
	p := &PreprocessorInclude{}
	p.Content = name
	p.Bracket = `"`

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
		return NewError(p.Context(), "wrong include content, expect %s, got %s", o.Content, p.Content)
	}

	if p.Bracket != o.Bracket {
		return NewError(p.Context(), "wrong include bracket, expect %s, got %s", o.Bracket, p.Bracket)
	}

	return nil
}

func (p *PreprocessorInclude) Context() *context.Context {
	return context.Join(p.Hash, p.Command, p.LBracketCtx, p.ContentCtx, p.RBracketCtx)
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

func (p *PreprocessorInline) statementNode() {}

func (p *PreprocessorInline) Type() TokenType {
	return NodePreprocessorInline
}

func (p *PreprocessorInline) EqualTo(_ context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(p, other)
	if err != nil {
		return err
	}

	if p.CodeType.Content() != o.CodeType.Content() {
		return NewError(p.Context(), "wrong inline code type, expect %s, got %s", o.CodeType.Content(), p.CodeType.Content())
	}

	if p.Content.Content() != o.Content.Content() {
		return NewError(p.Context(), "wrong inline content, expect %s, got %s", o.Content.Content(), p.Content.Content())
	}

	return nil
}

func (p *PreprocessorInline) Context() *context.Context {
	return context.Join(p.Hash, p.Command, p.CodeType, p.Content, p.HashEnd, p.CommandEnd, p.CodeTypeEnd)
}

func (p *PreprocessorInline) Empty() bool {
	return p.Content == nil
}
