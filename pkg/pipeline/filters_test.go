package pipeline

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"testing"
)

func TestFilters(t *testing.T) {
	t.Parallel()

	input := "apple\nbanana\ncherry\n"

	tests := []struct {
		filter func(io.Reader, io.Writer) error
		input  string
		want   string
	}{
		{Columns(",", "3,2,1"), "one\t\tthree\n", "one\t\tthree\n"},
		{Columns("\t", "9"), "one\t\tthree\n", "\n"},
		{Columns("\t", "3,2,1"), "one\t\tthree\n", "three\t\tone\n"},
		{ColumnsCSV(",", "3,2,1"), "one\t\tthree\n", "one\t\tthree\n"},
		{ColumnsCSV("\t", "9"), "one\t\tthree\n", "\n"},
		{ColumnsCSV(",", "3,2,1"), "one,\"t,w,o\",\"th\"\"ree\"\n", "\"th\"\"ree\",\"t,w,o\",one\n"},
		{CountLines(), "", "0\n"},
		{CountLines(), input, "3\n"},
		{CountRunes(), "", "0\n"},
		{CountRunes(), input + " 🍎", "19\n"},
		{CountWords(), "", "0\n"},
		{CountWords(), input + " 🍎", "4\n"},
		{First(-1), input, ""},
		{First(2), input, "apple\nbanana\n"},
		{First(10), input, input},
		{NotFirst(-1), input, input},
		{NotFirst(2), input, "cherry\n"},
		{NotFirst(10), input, ""},
		{Frequency(), input + "apple\n", "2 apple\n1 banana\n1 cherry\n"},
		{Join(", "), input, "apple, banana, cherry\n"},
		{Last(-1), input, ""},
		{Last(2), input, "banana\ncherry\n"},
		{Last(10), input, input},
		{NotLast(-1), input, ""},
		{NotLast(2), input, "apple\n"},
		{NotLast(10), input, ""},
		{Match(""), input, ""},
		{Match("pl"), input, "apple\n"},
		{Match("z"), input, ""},
		{NotMatch(""), input, ""},
		{NotMatch("pl"), input, "banana\ncherry\n"},
		{NotMatch("z"), input, input},
		{MatchRegex(regexp.MustCompile("")), input, ""},
		{MatchRegex(regexp.MustCompile("a.+e")), input, "apple\n"},
		{MatchRegex(regexp.MustCompile("[0-9]")), input, ""},
		{NotMatchRegex(regexp.MustCompile("")), input, ""},
		{NotMatchRegex(regexp.MustCompile("a.+e")), input, "banana\ncherry\n"},
		{NotMatchRegex(regexp.MustCompile("[0-9]")), input, input},
		{Replace("", ""), input, input},
		{Replace("orange", ""), input, input},
		{Replace("ple", "ricot"), input, "apricot\nbanana\ncherry\n"},
		{ReplaceRegex(regexp.MustCompile(""), ""), input, input},
		{ReplaceRegex(regexp.MustCompile("[0-9]"), ""), input, input},
		{ReplaceRegex(regexp.MustCompile("(.+)y"), "${1}ies"), input, "apple\nbanana\ncherries\n"},
	}

	for k, tc := range tests {
		// TODO: remove shadowing once using go v1.22
		k := k
		tc := tc

		t.Run(fmt.Sprint(k), func(t *testing.T) {
			t.Parallel()

			got := &bytes.Buffer{}

			err := tc.filter(strings.NewReader(tc.input), got)
			if err != nil {
				t.Fatalf("(test: %d) error executing filter: %v: input: %v", k, err, tc.input)
			}

			if tc.want != got.String() {
				log.Fatalf("(test %d) wanted: %q, got: %q", k, tc.want, got.String())
			}
		})
	}
}
