package gmark

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"
)

type ParseFunc func(*regexp2.Match) []*Token

var EmptyParseFunc = func(m *regexp2.Match) (t []*Token) {
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
	"text":      parseText,
}

func printM(m *regexp2.Match) {
	fmt.Println("Match")
	fmt.Println(m.String())
	for i, g := range m.Groups() {
		fmt.Println("Group", i)
		for ii, c := range g.Captures {
			fmt.Println("Cap", ii, c.String())
		}
	}
}

func parseText(m *regexp2.Match) []*Token {
	return []*Token{
		&Token{
			Type: TypeText,
			Text: m.String(),
		},
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
	matches := []*regexp2.Match{m}
	for {
		m, _ = listitem.FindNextMatch(m)
		if m == nil {
			break
		}
		matches = append(matches, m)
	}
	prev := false
	for idx, mm := range matches {
		item := mm.String()
		space := len(item)
		item, _ = listbullet.Replace(item, "", -1, -1)
		var loose bool
		if idx == len(matches)-1 {
			loose = prev == true
			if strings.HasSuffix(item, "\n") {
				item = strings.TrimRight(item, "\n")
			}
		} else {
			loose = mm.GroupByNumber(3).String() == "\n"
		}
		prev = loose

		if strings.Contains(item, "\n") {
			space = space - len(item)
			p := regexp2.MustCompile(fmt.Sprintf("^ {1,%d}", space), regexp2.Multiline)
			item, _ = p.Replace(item, "", -1, -1)
		}
		tokens = append(tokens,
			&Token{Type: TypeListItemStart, Loose: loose},
			&Token{Type: TypeListItem, Text: item},
			&Token{Type: TypeListItemEnd, Loose: loose},
		)
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
