package pipesore

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetToken(t *testing.T) {
	t.Parallel()

	filters := `Replace(" ", "\n") | Freq() | First(1)`

	tests := []token{
		{ttype: FILTER, literal: "Replace", position: position{start: 0, end: 7}},
		{ttype: LPAREN, literal: "(", position: position{start: 7, end: 8}},
		{ttype: STRING, literal: " ", position: position{start: 8, end: 11}},
		{ttype: COMMA, literal: ",", position: position{start: 11, end: 12}},
		{ttype: STRING, literal: "\n", position: position{start: 13, end: 17}},
		{ttype: RPAREN, literal: ")", position: position{start: 17, end: 18}},

		{ttype: PIPE, literal: "|", position: position{start: 19, end: 20}},

		{ttype: FILTER, literal: "Freq", position: position{start: 21, end: 25}},
		{ttype: LPAREN, literal: "(", position: position{start: 25, end: 26}},
		{ttype: RPAREN, literal: ")", position: position{start: 26, end: 27}},

		{ttype: PIPE, literal: "|", position: position{start: 28, end: 29}},

		{ttype: FILTER, literal: "First", position: position{start: 30, end: 35}},
		{ttype: LPAREN, literal: "(", position: position{start: 35, end: 36}},
		{ttype: INT, literal: "1", position: position{start: 36, end: 37}},
		{ttype: RPAREN, literal: ")", position: position{start: 37, end: 38}},

		{ttype: EOF, literal: "", position: position{start: 38, end: 39}},
	}

	l := newLexer(filters)

	for k, tc := range tests {
		k := k
		tc := tc

		got := l.getToken()

		t.Run(fmt.Sprint(k), func(t *testing.T) {
			t.Parallel()

			if !reflect.DeepEqual(tc, got) {
				t.Fatalf("wanted: %#v, got: %#v", tc, got)
			}
		})
	}
}
