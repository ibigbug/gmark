package gmark

import (
	"fmt"
	"strings"
)

type ParseFunc func([][]string) []*Token

var EmptyParseFunc = func(m [][]string) (t []*Token) {
	fmt.Println(m)
	return
}

var DefaultParseFuncs = map[string]ParseFunc{
	"newline":   parseNewline,
	"hrule":     parseHrule,
	"heading":   parseHeading,
	"lheading":  parseLheading,
	"paragraph": parseParagraph,
	"listitem":  parseListItem,
}

func parseListItem(m [][]string) []*Token {
	return []*Token{
		&Token{
			Type: TypeListItem,
			Text: m[0][0],
		},
	}
}

func parseNewline(m [][]string) []*Token {
	return []*Token{
		&Token{
			Type: TypeNewline,
		},
	}
}

func parseHrule(m [][]string) []*Token {
	return []*Token{
		&Token{
			Type: TypeHrule,
		},
	}
}

func parseHeading(m [][]string) []*Token {
	return []*Token{
		&Token{
			Type:  TypeHeading,
			Level: len(m[0][1]),
			Text:  m[0][2],
		},
	}
}

func parseLheading(m [][]string) []*Token {
	var level int
	if m[0][2] == "=" {
		level = 1
	} else {
		level = 2
	}
	return []*Token{
		&Token{
			Type:  TypeHeading,
			Level: level,
			Text:  m[0][1],
		},
	}
}

func parseParagraph(m [][]string) []*Token {
	return []*Token{
		&Token{
			Type: TypeParagraph,
			Text: strings.TrimSuffix(m[0][1], "\n"),
		},
	}
}
