package pipesore

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type lexer struct {
	input    string
	position int
}

func newLexer(input string) *lexer {
	l := &lexer{input: input}
	return l
}

func (l *lexer) getToken() token {
	var tt tokenType
	var tl string

	ch := l.getSignificantChar()
	start := l.position
	if isFilter(ch) {
		return token{
			ttype:   FILTER,
			literal: l.getString(isFilter),
			position: position{
				start: start,
				end:   l.position,
			},
		}
	} else if isDigit(ch) {
		return token{
			ttype:   INT,
			literal: l.getString(isDigit),
			position: position{
				start: start,
				end:   l.position,
			},
		}
	} else if isQuote(ch) {
		str, err := l.getQuotedString()
		if err == nil {
			tt = STRING
			tl = str
		} else {
			tt = ILLEGAL
			tl = fmt.Sprintf("%v '%s'", err, str)
		}
		return token{
			ttype:   tt,
			literal: tl,
			position: position{
				start: start,
				end:   l.position,
			},
		}
	}

	switch ch {
	case '(':
		tt = LPAREN
	case ')':
		tt = RPAREN
	case ',':
		tt = COMMA
	case '|':
		tt = PIPE
	case '\000':
		tt = EOF
	}

	if unicode.IsPrint(rune(ch)) {
		tl = string(ch)
	}

	l.position++

	return token{
		ttype:   tt,
		literal: tl,
		position: position{
			start: start,
			end:   l.position,
		},
	}
}

func (l *lexer) getSignificantChar() byte {
	ch := l.getChar(0)

	for isWhitespace(ch) {
		l.position++
		ch = l.getChar(0)
	}

	return ch
}

func (l *lexer) getString(fn func(byte) bool) string {
	startPosition := l.position
	l.position++

	for fn(l.getChar(0)) {
		l.position++
	}

	return l.input[startPosition:l.position]
}

func (l *lexer) getQuotedString() (string, error) {
	startPosition := l.position

	for {
		l.position++
		ch := l.getChar(0)
		if ch == '"' {
			l.position++
			break
		}
		if ch == '\000' {
			return l.input[startPosition:l.position], errors.New("unterminated string")
		}

		peekCh := l.getChar(1)
		if ch == '\\' && (peekCh == '\\' || peekCh == '"') {
			l.position++
		}
	}

	return strconv.Unquote(l.input[startPosition:l.position])
}

func (l *lexer) getChar(offset int) byte {
	position := l.position + offset
	if position >= len(l.input) {
		return '\000'
	}
	return l.input[position]
}

func isWhitespace(ch byte) bool {
	return ch == ' '
}

func isFilter(ch byte) bool {
	return ch == '!' || 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isQuote(ch byte) bool {
	return ch == '"'
}
