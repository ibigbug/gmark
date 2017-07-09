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
	for i, g := range m.Groups() {
		fmt.Println("Group", i)
		for ii, c := range g.Captures {
			fmt.Println("Cap", ii, c.String())
		}
	}
}
func parseListBlock(m *regexp2.Match) []*Token {
	bullet := m.GroupByNumber(2).Capture.String()
	ordered := strings.Contains(bullet, ".")
	tokens := []*Token{
		&Token{
			Type:    TypeListStart,
			Ordered: ordered,
		},
	}

	listitems, _ := listitem.FindStringMatch(m.Group.String())
	tokens = append(tokens, parseListItem(listitems)...)
	tokens = append(tokens, &Token{Type: TypeListEnd, Ordered: ordered})
	return tokens
}

func parseListItem(m *regexp2.Match) []*Token {
	tokens := make([]*Token, 0)
	prev := false
	for idx, mm := range m.Groups() {
		item := mm.Capture.String()
		item, _ = listbullet.Replace(item, "", -1, -1)
		var loose bool
		if idx == m.GroupCount()-1 {
			loose = prev == true
		} else {
			loose = mm.Captures[3].String() == "\n\n"
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
			Level: len(m.GroupByNumber(1).String()),
			Text:  m.GroupByNumber(2).String(),
		},
	}
}

func parseLheading(m *regexp2.Match) []*Token {
	var level int
	if m.GroupByNumber(2).String() == "=" {
		level = 1
	} else {
		level = 2
	}
	return []*Token{
		&Token{
			Type:  TypeHeading,
			Level: level,
			Text:  m.GroupByNumber(1).String(),
		},
	}
}

func parseParagraph(m *regexp2.Match) []*Token {
	return []*Token{
		&Token{
			Type: TypeParagraph,
			Text: strings.TrimSuffix(m.GroupByNumber(1).String(), "\n"),
		},
	}
}
