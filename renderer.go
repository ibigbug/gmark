package gmark

import (
	"fmt"
)

type Renderer interface {
	Render(*Token) string
}

var DefaultRenderer = &renderer{}

type renderer struct{}

func (r *renderer) Render(tok *Token) string {
	switch tok.Type {
	case TypeHrule:
		return "<hr>\n"
	case TypeHeading:
		return fmt.Sprintf("<h%d>%s</h%d>\n", tok.Level, tok.Text, tok.Level)
	case TypeParagraph:
		return fmt.Sprintf("<p>%s</p>\n", tok.Text)
	case TypeNewline:
		return "<br>\n"
	case TypeListStart:
		if tok.Ordered {
			return "<ol>"
		}
		return "<ul>"
	case TypeListItem:
		if tok.Loose {
			return fmt.Sprintf("<li><p>%s</p></li>", tok.Text)
		}
		return fmt.Sprintf("<li>%s</li>", tok.Text)
	case TypeListEnd:
		if tok.Ordered {
			return "</ol>"
		}
		return "</ul>"
	}

	panic("Unknow Token type: " + string(tok.Type))
}
