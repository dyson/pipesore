package pipeline

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestPipeline(t *testing.T) {
	t.Parallel()

	input := "apple\nbanana\ncherry\n"
	filter := CountLines

	want := "3\n"
	got := &bytes.Buffer{}

	p := NewPipeline(strings.NewReader(input))
	p.Filter(filter())

	if _, err := p.Output(got); err != nil {
		t.Fatalf("error executing pipeline: %v", err)
	}

	if want != got.String() {
		log.Fatalf("wanted: %s, got: %s", want, got.String())
	}
}
