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
		filters: []filter{
			{name: "Replace", arguments: []any{" ", "\n"}, position: position{start: 0, end: 7}},
			{name: "Freq", arguments: nil, position: position{start: 21, end: 25}},
			{name: "First", arguments: []any{1}, position: position{start: 30, end: 35}},
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
