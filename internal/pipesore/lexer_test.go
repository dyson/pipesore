package pipesore

import (
	"fmt"
	"testing"
)

func TestGetToken(t *testing.T) {
	t.Parallel()

	filters := `Replace(" ", "\n") | Freq() | First(1)`

	tests := []struct {
		wantedType    tokenType
		wantedLiteral string
	}{
		{FUNCTION, "Replace"},
		{LPAREN, "("},
		{STRING, " "},
		{COMMA, ","},
		{STRING, "\n"},
		{RPAREN, ")"},

		{PIPE, "|"},

		{FUNCTION, "Freq"},
		{LPAREN, "("},
		{RPAREN, ")"},

		{PIPE, "|"},

		{FUNCTION, "First"},
		{LPAREN, "("},
		{INT, "1"},
		{RPAREN, ")"},

		{EOF, "\000"},
	}

	l := newLexer(filters)

	for k, tc := range tests {
		t.Run(fmt.Sprint(k), func(t *testing.T) {
			got := l.getToken()

			if got.ttype != tc.wantedType || got.literal != tc.wantedLiteral {
				t.Logf("tokentype wanted=%q, got=%q", tc.wantedType, got.ttype)
				t.Fatalf("literal wanted=%q, got=%q", tc.wantedLiteral, got.literal)
			}
		})
	}
}
