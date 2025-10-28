package ast

import (
	"fmt"
)

type NodeType int

const (
	Invalid NodeType = iota
	notUsedToken
	literalBegin
	Null
	False
	True
	Integer
	Float
	String
	IdentifierName
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

	preprocessorBegin
	NodePreprocessorInclude
	NodePreprocessorInline
	preprocessorEnd

	LastToken

	SIllegal            = "Illegal"
	SEOF                = "EOF"
	SNull               = "null"
	SFalse              = "false"
	STrue               = "true"
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
	SHash               = "#"
	SAt                 = "@"
	SCommentStart       = "//"
)

var tokenStringMap = map[NodeType]string{
	Invalid:            SIllegal,
	Null:               SNull,
	False:              SFalse,
	True:               STrue,
	Auto:               SAuto,
	Var:                SVar,
	Const:              SConst,
	Global:             SGlobal,
	Function:           SFunction,
	Structure:          SStructure,
	Type:               SType,
	If:                 SIf,
	Elif:               SElif,
	Else:               SElse,
	For:                SFor,
	While:              SWhile,
	Do:                 SDo,
	Foreach:            SForeach,
	Break:              SBreak,
	Continue:           SContinue,
	And:                SAnd,
	Or:                 SOr,
	Not:                SNot,
	New:                SNew,
	Delete:             SDelete,
	Ref:                SRef,
	Return:             SReturn,
	Call:               SCall,
	Export:             SExport,
	Import:             SImport,
	Module:             SModule,
	Sizeof:             SSizeof,
	Plus:               SPlus,
	Sub:                SSub,
	Asterisk:           SAsterisk,
	Slash:              SSlash,
	Backslash:          SBackslash,
	Percent:            SPercent,
	Equal:              SEqual,
	NotEqual:           SNotEqual,
	InstanceEqual:      SInstanceEqual,
	InstanceNotEqual:   SInstanceNotEqual,
	LessThan:           SLessThan,
	LessThanOrEqual:    SLessThanOrEqual,
	GreaterThan:        SGreaterThan,
	GreaterThanOrEqual: SGreaterThanOrEqual,
	Ampersand:          SAmpersand,
	VerticalBar:        SVerticalBar,
	Tilde:              STilde,
	Caret:              SCaret,
	ShiftLeft:          SShiftLeft,
	ShiftRight:         SShiftRight,
	PointerAdd:         SPointerAdd,
	PointerSub:         SPointerSub,
	Assign:             SAssign,
	InferenceAssign:    SInferenceAssign,
	LeftParen:          SLeftParen,
	RightParen:         SRightParen,
	LeftBracket:        SLeftBracket,
	RightBracket:       SRightBracket,
	LeftBrace:          SLeftBrace,
	RightBrace:         SRightBrace,
	Comma:              SComma,
	Period:             SPeriod,
	Colon:              SColon,
	Semicolon:          SSemicolon,
	DualColon:          SDualColon,
	QuestionMark:       SQuestionMark,
	Bang:               SBang,
	Hash:               SHash,
	At:                 SAt,
	CommentStart:       SCommentStart,
}

func (t NodeType) IsOperator() bool {
	return t > operatorBegin && t < operatorEnd
}

func (t NodeType) String() string {
	s, ok := tokenStringMap[t]
	if ok {
		return s
	}

	return fmt.Sprintf("<Token %d>", t)
}

var keywordMap = map[string]NodeType{
	SNull:      Null,
	SFalse:     False,
	STrue:      True,
	SAuto:      Auto,
	SVar:       Var,
	SConst:     Const,
	SGlobal:    Global,
	SFunction:  Function,
	SStructure: Structure,
	SType:      Type,
	SIf:        If,
	SElif:      Elif,
	SElse:      Else,
	SFor:       For,
	SWhile:     While,
	SDo:        Do,
	SForeach:   Foreach,
	SBreak:     Break,
	SContinue:  Continue,
	SAnd:       And,
	SOr:        Or,
	SNot:       Not,
	SNew:       New,
	SDelete:    Delete,
	SRef:       Ref,
	SReturn:    Return,
	SCall:      Call,
	SExport:    Export,
	SImport:    Import,
	SModule:    Module,
	SSizeof:    Sizeof,
	SInclude:   Import,
}

func GetKeywordNodeType(s string) NodeType {
	if t, ok := keywordMap[s]; ok {
		return t
	}

	return Invalid
}

var operatorMap = map[string]NodeType{
	SPlus:               Plus,
	SSub:                Sub,
	SAsterisk:           Asterisk,
	SSlash:              Slash,
	SBackslash:          Backslash,
	SPercent:            Percent,
	SEqual:              Equal,
	SNotEqual:           NotEqual,
	SInstanceEqual:      InstanceEqual,
	SInstanceNotEqual:   InstanceNotEqual,
	SLessThan:           LessThan,
	SLessThanOrEqual:    LessThanOrEqual,
	SGreaterThan:        GreaterThan,
	SGreaterThanOrEqual: GreaterThanOrEqual,
	SAmpersand:          Ampersand,
	SVerticalBar:        VerticalBar,
	STilde:              Tilde,
	SCaret:              Caret,
	SShiftLeft:          ShiftLeft,
	SShiftRight:         ShiftRight,
	SPointerAdd:         PointerAdd,
	SPointerSub:         PointerSub,
	SAssign:             Assign,
	SInferenceAssign:    InferenceAssign,
	SLeftParen:          LeftParen,
	SRightParen:         RightParen,
	SLeftBracket:        LeftBracket,
	SRightBracket:       RightBracket,
	SLeftBrace:          LeftBrace,
	SRightBrace:         RightBrace,
	SComma:              Comma,
	SPeriod:             Period,
	SColon:              Colon,
	SSemicolon:          Semicolon,
	SDualColon:          DualColon,
	SQuestionMark:       QuestionMark,
	SBang:               Bang,
	SHash:               Hash,
	SAt:                 At,
	SCommentStart:       CommentStart,
}

func GetOperatorNodeType(s string) NodeType {
	if t, ok := operatorMap[s]; ok {
		return t
	}

	return Invalid
}
