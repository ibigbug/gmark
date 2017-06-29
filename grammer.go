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
	listitem   = regexp.MustCompile("^(( *)(?:[*+-]|\\d\\.) [^\n]*)")

	paragraphLookAhead = regexp.MustCompile(fmt.Sprintf("(?:%s|%s)",
		strings.TrimPrefix(heading.String(), "^"),
		strings.TrimPrefix(lheading.String(), "^")))
)

var DefaultBlockGrammer = map[string]*regexp.Regexp{
	"newline":   newline,
	"hrule":     hrule,
	"heading":   heading,
	"lheading":  lheading,
	"listitem":  listitem,
	"paragraph": paragraph,
}

var DefaultSupportedRules = []string{
	"newline", "hrule",
	"heading", "lheading", "listitem",
	"paragraph",
}
