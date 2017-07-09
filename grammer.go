package gmark

import (
	"fmt"

	"strings"

	"github.com/dlclark/regexp2"
)

var (
	newline    = regexp2.MustCompile(`^\n+`, regexp2.None)
	blockcode  = regexp2.MustCompile(`^( {4}[^\n]+\n*)+`, regexp2.None)
	fences     = regexp2.MustCompile("^ *(`{3,}|~{3,}) *(\\S+)? *\n([\\s\\S]+?)\\s*(`{3,}|~{3,}) *(?:\n+|$)", regexp2.None)
	hrule      = regexp2.MustCompile(`^ {0,3}[-*_](?: *[-*_]){2,} *(?:\n+|$)`, regexp2.None)
	heading    = regexp2.MustCompile(`^ *(#{1,6}) *([^\n]+?) *#* *(?:\n+|$)`, regexp2.None)
	lheading   = regexp2.MustCompile(`^([^\n]+)\n *(=|-)+ *(?:\n+|$)`, regexp2.None)
	blockquote = regexp2.MustCompile(`^( *>[^\n]+(\n[^\n]+)*\n*)+`, regexp2.None)
	paragraph  = regexp2.MustCompile(
		fmt.Sprintf(`^((?:[^\n]+\n?(?!%s|%s))+)\n*`, stringify(heading), stringify(lheading)),
		regexp2.None)
	listblock = regexp2.MustCompile(`^( *)([*+-]|\d+\.) [\s\S]+?`+
		`(?:`+
		`\n+(?=\1?(?:[-*_] *){3,}(?:\n+|$))`+ // hrule
		`|\n{2,}`+
		`(?! )`+
		`(?!\1(?:[*+-]|\d+\.) )\n*`+
		`|\s*$)`, regexp2.None)
	listitem = regexp2.MustCompile(`^(( *)(?:[*+-]|\d+\.) [^\n]*`+
		`(\n(?!\2(?:[*+-]|\d+\.) )[^\n]*)*)`, regexp2.Multiline)
	listbullet = regexp2.MustCompile(`^ *(?:[*+-]|\d+\.) +`, regexp2.None)
)

var DefaultBlockGrammer = map[string]*regexp2.Regexp{
	"newline":   newline,
	"hrule":     hrule,
	"heading":   heading,
	"lheading":  lheading,
	"listblock": listblock,
	"paragraph": paragraph,
}

var DefaultSupportedRules = []string{
	"newline", "hrule",
	"heading", "lheading", "listblock",
	"paragraph",
}

func stringify(p *regexp2.Regexp) (s string) {
	s = p.String()
	if strings.HasPrefix(s, "^") {
		s = s[1:]
	}
	return
}
