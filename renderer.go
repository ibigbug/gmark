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
	case "hrule":
		return "<hr>\n"
	case "heading":
		return fmt.Sprintf("<h%d>%s</h%d>\n", tok.Level, tok.Text, tok.Level)
	case "paragraph":
		return fmt.Sprintf("<p>%s</p>\n", tok.Text)
	}
	panic("Unknow Token type")
}
