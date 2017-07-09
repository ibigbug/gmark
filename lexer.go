package gmark

import (
	"strings"

	"github.com/dlclark/regexp2"
)

var (
	replaceNewline = regexp2.MustCompile("\r\n|\n", regexp2.None)
	replaceTabs    = regexp2.MustCompile("\t", regexp2.None)
	trimEmptyLine  = regexp2.MustCompile("^ +$", regexp2.Multiline)
)

type Lexer interface {
	Process(string) []*Token
}

var DefaultBlockLexer = &BlockLexer{
	Rules:          DefaultBlockGrammer,
	ParseFuncs:     DefaultParseFuncs,
	SupportedRules: DefaultSupportedRules,
}

type BlockLexer struct {
	Rules          map[string]*regexp2.Regexp
	ParseFuncs     map[string]ParseFunc
	SupportedRules []string
}

func (l *BlockLexer) preprocess(text string) (cleaned string) {
	cleaned, _ = replaceNewline.Replace(text, "\n", -1, -1)
	cleaned, _ = replaceTabs.Replace(cleaned, "    ", -1, -1)
	cleaned = strings.Replace(cleaned, "\u2424", "\n", -1)
	cleaned, _ = trimEmptyLine.Replace(cleaned, "", -1, -1)
	return
}

func (l *BlockLexer) Process(text string) (t []*Token) {
	text = l.preprocess(text)

	var manipulate = func(text string) (rv [][]string, matched bool) {
		for _, ruleName := range l.SupportedRules {
			re := l.Rules[ruleName]
			m, _ := re.FindStringMatch(text, -1)
			if m == nil {
				continue
			}
			if pf, canParse := l.ParseFuncs[ruleName]; canParse {
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
			var l = len(m[0][0])
			text = text[l:]
			continue
		}
		if text != "" {
			panic("failed")
		}
	}
	return
}
