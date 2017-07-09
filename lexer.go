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
	Process(string, *[]string) []*Token
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

func (l *BlockLexer) Process(text string, rules *[]string) (t []*Token) {
	text = l.preprocess(text)

	if rules == nil {
		rules = &l.SupportedRules
	}
	var manipulate = func(text string) (rv *regexp2.Match, matched bool) {
		for _, ruleName := range *rules {
			re := l.Rules[ruleName]
			m, _ := re.FindStringMatch(text)
			if m == nil {
				continue
			}
			if pf, canParse := l.ParseFuncs[ruleName]; canParse {
				tokens := pf(m)
				if ruleName == "listblock" {
					newTokens := make([]*Token, 0)
					for _, t := range tokens {
						if t.Type == TypeListItem {
							recur := l.Process(t.Text, &listItemRules)
							newTokens = append(newTokens, recur...)
						} else {
							newTokens = append(newTokens, t)
						}
					}
					tokens = newTokens
				}
				t = append(t, tokens...)
				return m, true
			} else {
				panic("Unknow syntax: " + ruleName + m.Group.String())
			}
		}
		return nil, false
	}

	for text != "" {
		m, matched := manipulate(text)
		if matched {
			var l = len(m.Group.String())
			text = text[l:]
			continue
		}
		if text != "" {
			panic("failed" + text)
		}
	}
	return
}
