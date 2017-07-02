package gmark

import (
	"regexp"
	"strings"
)

var (
	replaceNewline = regexp.MustCompile("\r\n|\n")
	replaceTabs    = regexp.MustCompile("\t")
	trimEmptyLine  = regexp.MustCompile("(?m)^ +$")
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
	Rules          map[string]*regexp.Regexp
	ParseFuncs     map[string]ParseFunc
	SupportedRules []string
}

func (l *BlockLexer) preprocess(text string) (cleaned string) {
	cleaned = replaceNewline.ReplaceAllString(text, "\n")
	cleaned = replaceTabs.ReplaceAllString(cleaned, "    ")
	cleaned = strings.Replace(cleaned, "\u2424", "\n", -1)
	cleaned = trimEmptyLine.ReplaceAllString(cleaned, "")
	return
}

func (l *BlockLexer) Process(text string) (t []*Token) {
	text = l.preprocess(text)

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
				if ruleName == "listblock" {
					for idx, mm := range m {
						p := mm[0]
						unpeek := listblockLookAhead.FindStringIndex(p)
						if len(unpeek) != 0 {
							mm[0] = mm[0][:unpeek[0]]
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
