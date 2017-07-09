package gmark

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"
)

type ParseFunc func(*regexp2.Match) []*Token

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
	"listblock": parseListBlock,
}

func printM(m *regexp2.Match) {
	for _, g := range m.Groups() {
		for _, c := range g.Captures() {
			fmt.Println(c)
		}
	}
}
func parseListBlock(m *regexp2.Match) []*Token {
	printM(m)
	bullet := m[0][2]
	ordered := strings.Contains(bullet, ".")
	tokens := []*Token{
		&Token{
			Type:    TypeListStart,
			Ordered: ordered,
		},
	}

	listitems := listitem.FindAllStringSubmatch(m[0][0], -1)
	tokens = append(tokens, parseListItem(listitems)...)
	tokens = append(tokens, &Token{Type: TypeListEnd, Ordered: ordered})
	return tokens
}

func parseListItem(m *regexp2.Match) []*Token {
	tokens := make([]*Token, 0)
	prev := false
	for idx, mm := range m {
		item := mm[1]
		item = listbullet.ReplaceAllString(item, "")
		var loose bool
		if idx == len(m)-1 {
			loose = prev == true
		} else {
			loose = mm[3] == "\n\n"
		}
		prev = loose
		tokens = append(tokens, &Token{Type: TypeListItem, Text: item, Loose: loose})
	}
	return tokens
}

func parseNewline(m *regexp2.Match) []*Token {
	return []*Token{
		&Token{
			Type: TypeNewline,
		},
	}
}

func parseHrule(m *regexp2.Match) []*Token {
	return []*Token{
		&Token{
			Type: TypeHrule,
		},
	}
}

func parseHeading(m *regexp2.Match) []*Token {
	return []*Token{
		&Token{
			Type:  TypeHeading,
			Level: len(m[0][1]),
			Text:  m[0][2],
		},
	}
}

func parseLheading(m *regexp2.Match) []*Token {
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

func parseParagraph(m *regexp2.Match) []*Token {
	return []*Token{
		&Token{
			Type: TypeParagraph,
			Text: strings.TrimSuffix(m[0][1], "\n"),
		},
	}
}
