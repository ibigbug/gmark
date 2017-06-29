package gmark

type Token struct {
	Type  TokenType
	Level int // heading level
	Text  string
}

type TokenType string

const (
	TypeHeading   TokenType = "heading"
	TypeLheading            = "lheading"
	TypeParagraph           = "paragraph"
	TypeHrule               = "hrule"
)
