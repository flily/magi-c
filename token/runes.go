package token

func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func IsValidIdentifierRune(r rune) bool {
	if 'a' <= r && r <= 'z' {
		return true
	}

	if 'A' <= r && r <= 'Z' {
		return true
	}

	if '0' <= r && r <= '9' {
		return true
	}

	if r == '_' {
		return true
	}

	return false
}

func IsValidIdentifierInitialRune(r rune) bool {
	if 'a' <= r && r <= 'z' {
		return true
	}

	if 'A' <= r && r <= 'Z' {
		return true
	}

	if r == '_' {
		return true
	}

	return false
}

var validSymbolInitialRunes = []bool{
	'!':  true,
	'"':  true,
	'#':  true,
	'$':  true,
	'%':  true,
	'&':  true,
	'\'': true,
	'(':  true,
	')':  true,
	'*':  true,
	'+':  true,
	',':  true,
	'-':  true,
	'.':  true,
	'/':  true,
	'[':  true,
	'\\': true,
	']':  true,
	'^':  true,
	'~':  true,
}

func IsValidSymbolRune(r rune) bool {
	if int(r) < len(validSymbolInitialRunes) {
		return validSymbolInitialRunes[int(r)]
	}

	return false
}
