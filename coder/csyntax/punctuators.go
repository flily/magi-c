package csyntax

type Punctuator int

const (
	PunctuatorInvalid Punctuator = iota
	OperatorLeftBracket
	OperatorRightBracket
	OperatorLeftParen
	OperatorRightParen
	OperatorLeftBrace
	OperatorRightBrace
	OperatorDot                 // expr . identifier        6.5.2.3
	OperatorArrow               // expr -> identifier       6.5.2.3
	OperatorIncrement           // expr ++, ++ expr         6.5.2.4 6.5.3.1
	OperatorDecrement           // expr --, -- expr         6.5.2.4 6.5.3.1
	OperatorSizeOf              // sizeof                   6.5.3.4
	OperatorAddressOf           // & expr                   6.5.3.2
	OperatorDereference         // * expr                   6.5.3.2
	OperatorPositive            // + expr                   6.5.3.3
	OperatorNegative            // - expr                   6.5.3.3
	OperatorBitwiseNot          // ~ expr                   6.5.3.3
	OperatorLogicalNot          // ! expr                   6.5.3.3
	OperatorMultiply            // expr * expr              6.5.5
	OperatorDivide              // expr / expr              6.5.5
	OperatorModulo              // expr % expr              6.5.5
	OperatorAdd                 // expr + expr              6.5.6
	OperatorSubtract            // expr - expr              6.5.6
	OperatorShiftLeft           // expr << expr             6.5.7
	OperatorShiftRight          // expr >> expr             6.5.7
	OperatorLessThan            // expr < expr              6.5.8
	OperatorGreaterThan         // expr > expr              6.5.8
	OperatorLessEqual           // expr <= expr             6.5.8
	OperatorGreaterEqual        // expr >= expr             6.5.8
	OperatorEqual               // expr == expr             6.5.9
	OperatorNotEqual            // expr != expr             6.5.9
	OperatorBitwiseAnd          // expr & expr              6.5.10
	OperatorBitwiseXor          // expr ^ expr              6.5.11
	OperatorBitwiseOr           // expr | expr              6.5.12
	OperatorLogicalAnd          // expr && expr             6.5.13
	OperatorLogicalOr           // expr || expr             6.5.14
	OperatorConditionalQuestion // expr ? expr : expr       6.5.15
	OperatorConditionalColon    // expr ? expr : expr.      6.5.15
	OperatorAssign              // expr = expr              6.5.16
	OperatorAssignAdd           // expr += expr             6.5.16
	OperatorAssignSubtract      // expr -= expr             6.5.16
	OperatorAssignMultiply      // expr *= expr             6.5.16
	OperatorAssignDivide        // expr /= expr             6.5.16
	OperatorAssignModulo        // expr %= expr             6.5.16
	OperatorAssignAnd           // expr &= expr             6.5.16
	OperatorAssignOr            // expr |= expr             6.5.16
	OperatorAssignXor           // expr ^= expr             6.5.16
	OperatorAssignShiftLeft     // expr <<= expr            6.5.16
	OperatorAssignShiftRight    // expr >>= expr            6.5.16
	PunctuatorComma             // ,
	PunctuatorColon             // :
	PunctuatorSemicolon         // ;
	PunctuatorEllipsis          // ...
	PunctuatorHash              // #
	PunctuatorDoubleHash        // ##
	PunctuatorLeftBracketAlt    // <:
	PunctuatorRightBracketAlt   // :>
	PunctuatorLeftBraceAlt      // <%
	PunctuatorRightBraceAlt     // %>
	PunctuatorHashAlt           // %:
	PunctuatorDoubleHashAlt     // %:%:
	PunctuatorSpace             // " "
)

var operatorString = map[Punctuator]string{
	OperatorLeftBracket:         "[",
	OperatorRightBracket:        "]",
	OperatorLeftParen:           "(",
	OperatorRightParen:          ")",
	OperatorLeftBrace:           "{",
	OperatorRightBrace:          "}",
	OperatorDot:                 ".",
	OperatorArrow:               "->",
	OperatorIncrement:           "++",
	OperatorDecrement:           "--",
	OperatorSizeOf:              "sizeof",
	OperatorAddressOf:           "&",
	OperatorDereference:         "*",
	OperatorPositive:            "+",
	OperatorNegative:            "-",
	OperatorBitwiseNot:          "~",
	OperatorLogicalNot:          "!",
	OperatorMultiply:            "*",
	OperatorDivide:              "/",
	OperatorModulo:              "%",
	OperatorAdd:                 "+",
	OperatorSubtract:            "-",
	OperatorShiftLeft:           "<<",
	OperatorShiftRight:          ">>",
	OperatorLessThan:            "<",
	OperatorGreaterThan:         ">",
	OperatorLessEqual:           "<=",
	OperatorGreaterEqual:        ">=",
	OperatorEqual:               "==",
	OperatorNotEqual:            "!=",
	OperatorBitwiseAnd:          "&",
	OperatorBitwiseXor:          "^",
	OperatorBitwiseOr:           "|",
	OperatorLogicalAnd:          "&&",
	OperatorLogicalOr:           "||",
	OperatorConditionalQuestion: "?",
	OperatorConditionalColon:    ":",
	OperatorAssign:              "=",
	OperatorAssignAdd:           "+=",
	OperatorAssignSubtract:      "-=",
	OperatorAssignMultiply:      "*=",
	OperatorAssignDivide:        "/=",
	OperatorAssignModulo:        "%=",
	OperatorAssignAnd:           "&=",
	OperatorAssignOr:            "|=",
	OperatorAssignXor:           "^=",
	OperatorAssignShiftLeft:     "<<=",
	OperatorAssignShiftRight:    ">>=",
	PunctuatorComma:             ",",
	PunctuatorColon:             ":",
	PunctuatorSemicolon:         ";",
	PunctuatorEllipsis:          "...",
	PunctuatorHash:              "#",
	PunctuatorDoubleHash:        "##",
	PunctuatorLeftBracketAlt:    "<:",
	PunctuatorRightBracketAlt:   ":>",
	PunctuatorLeftBraceAlt:      "<%",
	PunctuatorRightBraceAlt:     "%>",
	PunctuatorHashAlt:           "%:",
	PunctuatorDoubleHashAlt:     "%:%:",
	PunctuatorSpace:             " ",
}

func (p Punctuator) String() string {
	if s, ok := operatorString[p]; ok {
		return s
	}

	return "INVALID"
}

func (p Punctuator) codeElement() {}

func (p Punctuator) ItemString() string {
	return p.String()
}
