package pipesore

import "fmt"

const (
	ILLEGAL = iota
	EOF

	FILTER // First
	INT    // 1234
	STRING // hello, world!

	QUOTE // "

	LPAREN // (
	RPAREN // )

	COMMA // ,
	PIPE  // |
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	FILTER: "FILTER",
	INT:    "INT",
	STRING: "STRING",

	QUOTE: "\"",

	LPAREN: "(",
	RPAREN: ")",

	COMMA: ",",
	PIPE:  "|",
}

type tokenType int

func (t tokenType) String() string {
	return tokens[t]
}

type position struct {
	start int
	end   int
}

type token struct {
	ttype   tokenType
	literal string
	position
}

func (t token) String() string {
	s := fmt.Sprintf("'%s'", t.ttype)
	if t.literal != "" && t.literal != t.ttype.String() {
		s += fmt.Sprintf(" (%s)", t.literal)
	}

	return s
}
