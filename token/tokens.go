package token

type Token int

const (
	Invalid Token = iota
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
	Multiply           // *
	Divide             // /
	Modulo             // %
	Equal              // ==
	NotEqual           // !=
	InstanceEqual      // ===
	InstanceNotEqual   // !===
	LessThan           // <
	LessThanOrEqual    // <=
	GreaterThan        // >
	GreaterThanOrEqual // >=
	BitwiseAnd         // &
	BitwiseOr          // |
	BitwiseNot         // ~
	BitwiseXor         // ^
	ShiftLeft          // <<
	ShiftRight         // >>
	AddressOf          // &
	PointerAdd         // +>>
	PointerSub         // -<<

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
	Band            // !
	CommentStart    // //

	punctuationEnd
	operatorEnd
	LastToken
)

func (t Token) IsOperator() bool {
	return t > operatorBegin && t < operatorEnd
}
