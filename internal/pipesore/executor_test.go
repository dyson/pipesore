package pipesore

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	t.Parallel()

	input := "cat cat cat dog bird bird bird bird"
	filters := `Replace(" ", "\n") | Frequency() | First(1)`

	want := "4 bird\n"
	got := &bytes.Buffer{}

	err := execute(filters, strings.NewReader(input), got)
	if err != nil {
		t.Fatal(err)
	}

	if want != got.String() {
		log.Fatalf("wanted: %q, got: %q", want, got.String())
	}
}
