package ast

import (
	"github.com/flily/magi-c/context"
)

type InlineBlock struct {
	TerminalNode
	Begin     *context.Context
	BlockType *context.Context
	Content   string
	End       *context.Context
}

func NewInlineBlock(begin *context.Context, blockType *context.Context, end *context.Context, content string) *InlineBlock {
	block := &InlineBlock{
		TerminalNode: NewTerminalNode(context.Join(begin, blockType)),
		Begin:        begin,
		BlockType:    blockType,
		End:          end,
		Content:      content,
	}

	return block
}
