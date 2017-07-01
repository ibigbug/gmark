package gmark

type Token struct {
	Type TokenType
	Text string

	// Extra
	Level   int  // heading level
	Ordered bool // ordered or unordered list
	Loose   bool // if loose list item
}

type TokenType string

const (
	TypeHeading   TokenType = "heading"
	TypeLheading            = "lheading"
	TypeParagraph           = "paragraph"
	TypeHrule               = "hrule"
	TypeNewline             = "newline"
	TypeListStart           = "liststart"
	TypeListItem            = "listitem"
	TypeListEnd             = "listend"
)
