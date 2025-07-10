package context

type Cursor struct {
	Line   int
	Column int
	File   *FileContext
}
