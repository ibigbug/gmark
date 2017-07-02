package gmark

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	newline    = regexp.MustCompile("^\n+")
	blockcode  = regexp.MustCompile("^( {4}[^\n]+\n*)+")
	fences     = regexp.MustCompile("^ *(`{3,}|~{3,}) *(\\S+)? *\n([\\s\\S]+?)\\s*(`{3,}|~{3,}) *(?:\n+|$)")
	hrule      = regexp.MustCompile("^ {0,3}[-*_](?: *[-*_]){2,} *(?:\n+|$)")
	heading    = regexp.MustCompile("^ *(#{1,6}) *([^\n]+?) *#* *(?:\n+|$)")
	lheading   = regexp.MustCompile("^([^\n]+)\n *(=|-)+ *(?:\n+|$)")
	blockquote = regexp.MustCompile("^( *>[^\n]+(\n[^\n]+)*\n*)+")
	paragraph  = regexp.MustCompile("^((?:[^\n]+\n?)+)\n*")
	listblock  = regexp.MustCompile("^( *)([*+-]|\\d\\.) [\\s\\S]+?\n{2,}" +
		"(?:(?: *)(?:[*+-]|\\d\\.) [\\s\\S]+?\n{2,})*",
	)
	listitem   = regexp.MustCompile("(?m)^(( *)(?:[*+-]|\\d\\.) [^\n]*)(\n+)+")
	listbullet = regexp.MustCompile("^ *(?:[*+-]|\\d+\\.) +")

	paragraphLookAhead = regexp.MustCompile(fmt.Sprintf("(?:%s|%s)",
		strings.TrimPrefix(heading.String(), "^"),
		strings.TrimPrefix(lheading.String(), "^")))
	listblockLookAhead = regexp.MustCompile(strings.TrimPrefix(hrule.String(), "^"))
)

var DefaultBlockGrammer = map[string]*regexp.Regexp{
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
