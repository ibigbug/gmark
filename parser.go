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
	"listblock": parseListBlock,
}

func printM(m [][]string) {
	for i, mm := range m {
		fmt.Printf("M%d\n", i)
		for ii, mmm := range mm {
			fmt.Printf("Group%d\n", ii)
			fmt.Println(mmm)
		}
	}
}
func parseListBlock(m [][]string) []*Token {
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

func parseListItem(m [][]string) []*Token {
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
