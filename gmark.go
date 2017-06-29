package gmark

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	hrule     = regexp.MustCompile("^ {0,3}[-*_](?: *[-*_]){2,} *(?:\n+|$)")
	heading   = regexp.MustCompile("^ *(#{1,6}) *([^\n]+?) *#* *(?:\n+|$)")
	lheading  = regexp.MustCompile("^([^\n]+)\n *(=|-)+ *(?:\n+|$)")
	paragraph = regexp.MustCompile("^((?:[^\n]+\n?)+)\n*")

	paragraphLookAhead = regexp.MustCompile(fmt.Sprintf("(?:%s|%s)",
		strings.TrimPrefix(heading.String(), "^"),
		strings.TrimPrefix(lheading.String(), "^")))
)

type Token struct {
	Type  string
	Level int // heading level
	Text  string
}

type Lexer interface {
	Process(string) []*Token
}

type BlockLexer struct {
	Rules          map[string]*regexp.Regexp
	ParseFuncs     map[string]ParseFunc
	SupportedRules []string
}

func (l *BlockLexer) Process(text string) (t []*Token) {
	text = strings.TrimRight(text, "\n")

	var manipulate = func(text string) (rv [][]string, matched bool) {
		for _, ruleName := range l.SupportedRules {
			re := l.Rules[ruleName]
			m := re.FindAllStringSubmatch(text, -1)
			if len(m) == 0 {
				continue
			}
			if pf, canParse := l.ParseFuncs[ruleName]; canParse {
				if ruleName == "paragraph" {
					for idx, mm := range m {
						p := mm[0]
						unpeek := paragraphLookAhead.FindStringIndex(p)
						if len(unpeek) != 0 {
							mm[0] = mm[0][:unpeek[0]]
							mm[1] = mm[0]
							m[idx] = mm
						}
					}
				}
				tokens := pf(m)
				t = append(t, tokens...)
				return m, true
			} else {
				panic("Unknow syntax: " + ruleName + m[0][0])
			}
		}
		return nil, false
	}

	for text != "" {
		m, matched := manipulate(text)
		if matched {
			var l = 0
			for _, g := range m {
				l += len(g[0])
			}
			text = text[l:]
			continue
		}
		if text != "" {
			panic("failed")
		}
	}
	return
}

type Markdown struct {
	Inline   Lexer
	Block    Lexer
	Renderer Renderer

	Tokens []*Token
}

func (m *Markdown) Parse(text string) string {
	m.Tokens = m.Block.Process(text)
	return m.Output()
}

func (m *Markdown) Output() (out string) {
	for _, tok := range m.Tokens {
		out += m.Renderer.Render(tok)
	}
	fmt.Println(out)
	return
}

var DefaultBlockGrammer = map[string]*regexp.Regexp{
	"hrule":     hrule,
	"heading":   heading,
	"lheading":  lheading,
	"paragraph": paragraph,
}

var DefaultSupportedRules = []string{
	"hrule", "heading", "lheading", "paragraph",
}

var DefaultBlockLexer = &BlockLexer{
	Rules:          DefaultBlockGrammer,
	ParseFuncs:     DefaultParseFuncs,
	SupportedRules: DefaultSupportedRules,
}

func Convert(text string) string {
	m := Markdown{
		Block:    DefaultBlockLexer,
		Renderer: DefaultRenderer,
	}
	return m.Parse(text)
}
