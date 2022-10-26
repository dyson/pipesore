package pipesore

import (
	"log"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	t.Parallel()

	filters := `Replace(" ", "\n") | Freq() | First(1)`

	want := &ast{
		functions: []function{
			{name: "Replace", arguments: []any{" ", "\n"}},
			{name: "Freq", arguments: nil},
			{name: "First", arguments: []any{1}},
		},
	}

	t.Run("ast", func(t *testing.T) {
		got, err := newParser(newLexer(filters)).parse()
		if err != nil {
			log.Fatal(err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Fatalf("\nwanted:\n\n%#v\n\ngot:\n\n%#v\n\n", want, got)
		}
	})
}
