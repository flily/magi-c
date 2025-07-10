package context

type Highlight struct {
	Start int
	End   int
}

type LineContext struct {
	Content    *LineContent
	File       *FileContext
	Highlights []*Highlight
}

func (l *LineContext) Rune(n int) rune {
	return 0
}

type Context struct {
	File      *FileContext
	PrevLines []*LineContent
	NextLines []*LineContent
	Line      *LineContext
}
