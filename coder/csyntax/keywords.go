package csyntax

type Keyword = StringElement

const (
	KeywordIf           Keyword = "if"
	KeywordElse         Keyword = "else"
	KeywordFor          Keyword = "for"
	KeywordWhile        Keyword = "while"
	KeywordDo           Keyword = "do"
	KeywordSwitch       Keyword = "switch"
	KeywordCase         Keyword = "case"
	KeywordBreak        Keyword = "break"
	KeywordContinue     Keyword = "continue"
	KeywordReturn       Keyword = "return"
	KeywordStruct       Keyword = "struct"
	KeywordUnion        Keyword = "union"
	KeywordEnum         Keyword = "enum"
	KeywordTypedef      Keyword = "typedef"
	KeywordConst        Keyword = "const"
	KeywordVolatile     Keyword = "volatile"
	KeywordStatic       Keyword = "static"
	KeywordExtern       Keyword = "extern"
	KeywordAuto         Keyword = "auto"
	KeywordRegister     Keyword = "register"
	PreprocessorLine    Keyword = "#line"
	PreprocessorInclude Keyword = "#include"
	PreprocessorDefine  Keyword = "#define"
	PreprocessorIf      Keyword = "#if"
	PreprocessorElse    Keyword = "#else"
	PreprocessorElif    Keyword = "#elif"
	PreprocessorEndif   Keyword = "#endif"
)
