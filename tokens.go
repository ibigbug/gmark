package gmark

import (
	"fmt"
)

type Token struct {
	Type TokenType
	Text string

	// Extra
	Level   int  // heading level
	Ordered bool // ordered or unordered list
	Loose   bool // if loose list item
}

func (t *Token) String() string {
	return fmt.Sprintf("<Type:%s, Text: %s>", t.Type, t.Text)
}

type TokenType string

const (
	TypeHeading       TokenType = "heading"
	TypeLheading                = "lheading"
	TypeParagraph               = "paragraph"
	TypeHrule                   = "hrule"
	TypeNewline                 = "newline"
	TypeListStart               = "liststart"
	TypeListItemStart           = "listitemstart"
	TypeListItem                = "listitem"
	TypeListItemEnd             = "listitemend"
	TypeListEnd                 = "listend"
	TypeText                    = "text"
)
