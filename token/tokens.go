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
	Plus
	Sub
	Multiply
	Divide
	Modulo
	Equal
	NotEqual
	InstanceEqual
	InstanceNotEqual
	LessThan
	LessThanOrEqual 
	GreaterThan
	GreaterThanOrEqual
	BitwiseAnd
	BitwiseOr
	BitwiseNot
	BitwiseXor
	ShiftLeft
	ShiftRight
	AddressOf
	PointerAdd
	PointerSub
	operatorEnd

	punctuationBegin
	punctuationEnd
)
