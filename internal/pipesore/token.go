package pipesore

import "fmt"

const (
	ILLEGAL = iota
	EOF

	FUNCTION // First
	INT      // 1234
	STRING   // hello, world!

	QUOTE // "

	LPAREN // (
	RPAREN // )

	COMMA // ,
	PIPE  // |
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	FUNCTION: "FUNCTION",
	INT:      "INT",
	STRING:   "STRING",

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

type token struct {
	ttype   tokenType
	literal string
}

func (t token) String() string {
	return fmt.Sprintf("{%s %v}", t.ttype, t.literal)
}
