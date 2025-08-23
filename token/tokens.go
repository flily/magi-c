package token

import (
	"github.com/flily/magi-c/context"
)

type TokenType int

const (
	Invalid TokenType = iota
	literalBegin
	Null
	False
	True
	Integer
	Float
	String
	literalEnd

	keywordBegin
	Auto
	Var
	Const
	Global
	Function
	Structure
	Type
	If
	Elif
	Else
	For
	While
	Do
	Foreach
	Break
	Continue
	And
	Or
	Not
	New
	Delete
	Ref
	Return
	Call
	Export
	Import
	Module
	Sizeof
	keywordEnd

	operatorBegin
	Plus               // +
	Sub                // -
	Asterisk           // *
	Slash              // /
	Backslash          // \
	Percent            // %
	Equal              // ==
	NotEqual           // !=
	InstanceEqual      // ===
	InstanceNotEqual   // !==
	LessThan           // <
	LessThanOrEqual    // <=
	GreaterThan        // >
	GreaterThanOrEqual // >=
	Ampersand          // &
	VerticalBar        // |
	Tilde              // ~
	Caret              // ^
	ShiftLeft          // <<
	ShiftRight         // >>
	PointerAdd         // +>>
	PointerSub         // -<<
	// #

	punctuationBegin
	Assign          // =
	InferenceAssign // :=
	LeftParen       // (
	RightParen      // )
	LeftBracket     // [
	RightBracket    // ]
	LeftBrace       // {
	RightBrace      // }
	Comma           // ,
	Period          // .
	Colon           // :
	Semicolon       // ;
	DualColon       // ::
	QuestionMark    // ?
	Bang            // !
	Hash            // #
	At              // @
	CommentStart    // //

	punctuationEnd
	operatorEnd
	LastToken

	SIllegal            = "Illegal"
	SEOF                = "EOF"
	SAuto               = "auto"
	SVar                = "var"
	SConst              = "const"
	SGlobal             = "global"
	SFunction           = "func"
	SStructure          = "struct"
	SType               = "type"
	SIf                 = "if"
	SElif               = "elif"
	SElse               = "else"
	SFor                = "for"
	SWhile              = "while"
	SDo                 = "do"
	SForeach            = "foreach"
	SBreak              = "break"
	SContinue           = "continue"
	SAnd                = "and"
	SOr                 = "or"
	SNot                = "not"
	SNew                = "new"
	SDelete             = "delete"
	SRef                = "ref"
	SReturn             = "return"
	SCall               = "call"
	SExport             = "export"
	SImport             = "import"
	SModule             = "module"
	SSizeof             = "sizeof"
	SInclude            = "include"
	SPlus               = "+"
	SSub                = "-"
	SAsterisk           = "*"
	SSlash              = "/"
	SBackslash          = "\\"
	SPercent            = "%"
	SEqual              = "=="
	SNotEqual           = "!="
	SInstanceEqual      = "==="
	SInstanceNotEqual   = "!=="
	SLessThan           = "<"
	SLessThanOrEqual    = "<="
	SGreaterThan        = ">"
	SGreaterThanOrEqual = ">="
	SAmpersand          = "&"
	SVerticalBar        = "|"
	STilde              = "~"
	SCaret              = "^"
	SShiftLeft          = "<<"
	SShiftRight         = ">>"
	SPointerAdd         = "+>>"
	SPointerSub         = "-<<"
	SAssign             = "="
	SInferenceAssign    = ":="
	SLeftParen          = "("
	SRightParen         = ")"
	SLeftBracket        = "["
	SRightBracket       = "]"
	SLeftBrace          = "{"
	SRightBrace         = "}"
	SComma              = ","
	SPeriod             = "."
	SColon              = ":"
	SSemicolon          = ";"
	SDualColon          = "::"
	SQuestionMark       = "?"
	SBang               = "!"
	SHashTag            = "#"
	SAt                 = "@"
	SCommentStart       = "//"
)

func (t TokenType) IsOperator() bool {
	return t > operatorBegin && t < operatorEnd
}

type Token struct {
	Type    TokenType
	Context *context.Context
}

func NewToken(typ TokenType, ctx *context.Context) *Token {
	t := &Token{
		Type:    typ,
		Context: ctx,
	}

	return t
}
