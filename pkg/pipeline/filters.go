package pipeline

import (
	"bufio"
	"container/ring"
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

type filters map[string]filter

func (f filters) GetOrderedNames() []string {
	names := []string{}
	for name := range f {
		names = append(names, name)
	}

	sort.SliceStable(names, func(i, j int) bool {
		iName := names[i]
		iNot := false

		jName := names[j]

		if iName[0] == '!' {
			iName = iName[1:]
			iNot = true
		}

		if jName[0] == '!' {
			jName = jName[1:]
		}

		if iName != jName {
			return iName < jName
		}

		if iNot == true {
			return false
		}

		return true
	})

	return names
}

type filter struct {
	Value       reflect.Value
	Definition  string
	Description string
}

var (
	Filters = filters{
		"columns": {
			reflect.ValueOf(Columns),
			"Columns(delimiter string, columns string)",
			"Returns the selected `columns` in order where `columns is a 1-indexed comma separated list of column positions. Columns are defined by splitting with the `delimiter`.",
		},
		"columnscsv": {
			reflect.ValueOf(ColumnsCSV),
			"ColumnsCSV(delimiter string, columns string)",
			"Returns the selected `columns` in order where `columns` is a 1-indexed comma separated list of column positions. Parsing is CSV aware so quoted columns containing the `delimiter` when splitting are preserved.",
		},
		"countlines": {
			reflect.ValueOf(CountLines),
			"CountLines()",
			"Returns the line count. Lines are delimited by `\\r\\n`.",
		},
		"countrunes": {
			reflect.ValueOf(CountRunes),
			"CountRunes()",
			"Returns the rune (Unicode code points) count. Erroneous and short encodings are treated as single runes of width 1 byte.",
		},
		"countwords": {
			reflect.ValueOf(CountWords),
			"CountWords()",
			"Returns the word count. Words are delimited by `\\t|\\n|\\v|\\f|\\r|\u00A0|0x85|0xA0`.",
		},
		"first": {
			reflect.ValueOf(First),
			"First(n int)",
			"Returns first `n` lines where `n` is a positive integer. If the input has less than `n` lines, all lines are returned.",
		},
		"!first": {
			reflect.ValueOf(NotFirst),
			"!First(n int)",
			"Returns all but the the first `n` lines where `n` is a positive integer. If the input has less than `n` lines, no lines are returned.",
		},
		"frequency": {
			reflect.ValueOf(Frequency),
			"Frequency()",
			"Returns a descending list containing frequency and unique line. Lines with equal frequency are sorted alphabetically.",
		},
		"join": {
			reflect.ValueOf(Join),
			"Join(delimiter string)",
			"Joins all lines together seperated by `delimiter`.",
		},
		"last": {
			reflect.ValueOf(Last),
			"Last(n int)",
			"Returns last `n` lines where `n` is a positive integer. If the input has less than `n` lines, all lines are returned.",
		},
		"!last": {
			reflect.ValueOf(NotLast),
			"!Last(n int)",
			"Returns all but the last `n` lines where `n` is a positive integer. If the input has less than `n` lines, no lines are returned.",
		},
		"match": {
			reflect.ValueOf(Match),
			"Match(substring string)",
			"Returns all lines that contain `substring`.",
		},
		"!match": {
			reflect.ValueOf(NotMatch),
			"!Match(substring string)",
			"Returns all lines that don't contain `substring`.",
		},
		"matchregex": {
			reflect.ValueOf(MatchRegex),
			"MatchRegex(regex string)",
			"Returns all lines that match the compiled regular expression 'regex'. Regex is in the form of Re2 (https://github.com/google/re2/wiki/Syntax).",
		},
		"!matchregex": {
			reflect.ValueOf(NotMatchRegex),
			"!MatchRegex(regex string)",
			"Returns all lines that don't match the compiled regular expression 'regex'. Regex is in the form of Re2 (https://github.com/google/re2/wiki/Syntax).",
		},
		"replace": {
			reflect.ValueOf(Replace),
			"Replace(old string, replace string)",
			"Replaces all non-overlapping instances of `old` with `replace`.",
		},
		"replaceregex": {
			reflect.ValueOf(ReplaceRegex),
			"ReplaceRegex(regex string, replace string)",
			"Replaces all matches of the compiled regular expression `regex` with `replace`. Inside `replace`, `$` signs represent submatches. For example `$1` represents the text of the first submatch. Regex is in the form of Re2 (https://github.com/google/re2/wiki/Syntax).",
		},
	}
)

// Columns returns a filter that writes the selected 'columns' in the order
// provided where 'columns' is a 1-indexed comma separated list of column positions.
// Columns are defined by splitting with the 'delimiter'.
func Columns(delimiter string, columns string) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		order := []int{}
		for _, column := range strings.Split(columns, ",") {
			index, err := strconv.Atoi(strings.TrimSpace(column))
			if err != nil {
				return fmt.Errorf("list of columns must be comma separated list of ints, got: %v", columns)
			}

			order = append(order, index)
		}

		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			lineColumns := strings.Split(scanner.Text(), delimiter)

			output := []string{}
			for _, v := range order {
				if v-1 < len(lineColumns) {
					output = append(output, lineColumns[v-1])
				}
			}

			fmt.Fprintln(w, strings.Join(output, delimiter))
		}

		return scanner.Err()
	}
}

// ColumnsCSV returns a CSV aware filter that writes the selected 'columns' in
// the order provided where 'columns' is a 1-indexed comma separated list of
// column positions. Columns are defined by splitting with the 'delimiter'.
func ColumnsCSV(delimiter string, columns string) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		if utf8.RuneCount([]byte(delimiter)) > 1 {
			return fmt.Errorf("delimeter must be a single rune, got: %s", delimiter)
		}

		order := []int{}
		for _, column := range strings.Split(columns, ",") {
			index, err := strconv.Atoi(strings.TrimSpace(column))
			if err != nil {
				return fmt.Errorf("list of columns must be comma separated list of ints, got: %v", columns)
			}

			order = append(order, index)
		}

		reader := csv.NewReader(r)
		reader.Comma, _ = utf8.DecodeRuneInString(delimiter)
		// We really shouldn't be tolerant of malformed CSV input (and should
		// error) however we can set LazyQuotes to be less strict for commonly
		// incorrect quoting.
		//
		// Unfortunately how incorrect quoting should be interpreted is highly
		// dependent on how it was incorrectly implemented and so with LazyQuotes
		// enabled we will in some cases silently parse malformed CSV in a possibly
		// unexpected way to the user.
		//
		// On the other hand users don't always have control over the generation of
		// the CSV input and so it is hoped that the trade-off in using LazyQuotes
		// will allow for a better experience overall. If this is not that case we
		// can disable LazyQuotes and only parse valid rfc4180
		// (https://www.rfc-editor.org/rfc/rfc4180.html) csv.
		reader.LazyQuotes = true

		writer := csv.NewWriter(w)
		defer writer.Flush()

		for {
			lineColumns, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			output := []string{}
			for _, v := range order {
				if v-1 < len(lineColumns) {
					output = append(output, lineColumns[v-1])
				}
			}

			writer.Write(output)
		}

		return nil
	}
}

// CountLines returns a filter that writes the number of lines read.
func CountLines() func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		lines := 0

		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			lines++
		}

		fmt.Fprintln(w, lines)

		return scanner.Err()
	}
}

// CountRunes returns a filter that writes the number of runes (utf8
// characters) read.
func CountRunes() func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		runes := 0

		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			runes += utf8.RuneCountInString(scanner.Text())
		}

		fmt.Fprintln(w, runes)

		return scanner.Err()
	}
}

// CountWords returns a filter that writes the number of words (as defined by
// strings.Fields()) read.
func CountWords() func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		words := 0

		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			words += len(strings.Fields(scanner.Text()))
		}

		fmt.Fprintln(w, words)

		return scanner.Err()
	}
}

// First returns a filter that writes the first 'n' lines.
func First(n int) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		i := 0

		scanner := bufio.NewScanner(r)

		for i < n {
			if scanner.Scan() {
				fmt.Fprintln(w, scanner.Text())
			} else {
				return nil
			}

			i++
		}

		return scanner.Err()
	}
}

// NotFirst returns a filter that writes lines after the first 'n' lines.
func NotFirst(n int) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		i := 0

		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			if i >= n {
				fmt.Fprintln(w, scanner.Text())
			}

			i++
		}

		return scanner.Err()
	}
}

// Join returns a filter that writes all lines as a single string separated by
// 'delimiter'.
func Join(delimiter string) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		scanner := bufio.NewScanner(r)

		if scanner.Scan() {
			fmt.Fprint(w, scanner.Text())
		}

		for scanner.Scan() {
			fmt.Fprintf(w, "%s%s", delimiter, scanner.Text())
		}

		fmt.Fprintln(w)

		return scanner.Err()
	}
}

// Last returns a filter that writes the last 'n' lines.
func Last(n int) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		if n <= 0 {
			return nil
		}

		input := ring.New(n)

		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			input.Value = scanner.Text()
			input = input.Next()
		}

		input.Do(func(p any) {
			if p != nil {
				fmt.Fprintln(w, p)
			}
		})

		return scanner.Err()
	}
}

// NotLast returns a filter that writes up to the last 'n' lines.
func NotLast(n int) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		if n <= 0 {
			return nil
		}

		i := 0
		input := make([]string, n)

		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			j := i % n

			if i < n {
				input[j] = scanner.Text()
			} else {
				fmt.Fprintln(w, input[j])
				input[j] = scanner.Text()
			}

			i++
		}

		return scanner.Err()
	}
}

// Match returns a filter that writes lines containing 'substring'.
func Match(substring string) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		if substring == "" {
			return nil
		}

		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			if strings.Contains(scanner.Text(), substring) {
				fmt.Fprintln(w, scanner.Text())
			}
		}

		return scanner.Err()
	}
}

// NotMatch returns a filter that writes lines not containing 'substring'.
func NotMatch(substring string) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		if substring == "" {
			return nil
		}

		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			if !strings.Contains(scanner.Text(), substring) {
				fmt.Fprintln(w, scanner.Text())
			}
		}

		return scanner.Err()
	}
}

// MatchRegex returns a filter that writes lines matching the compiled regular
// expression 'regex'.
func MatchRegex(regex *regexp.Regexp) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		if regex.String() == "" {
			return nil
		}

		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			if regex.MatchString(scanner.Text()) {
				fmt.Fprintln(w, scanner.Text())
			}
		}

		return scanner.Err()
	}
}

// NotMatchRegex returns a filter that writes lines not matching the compiled
// regular expression 'regex'.
func NotMatchRegex(regex *regexp.Regexp) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		if regex.String() == "" {
			return nil
		}

		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			if !regex.MatchString(scanner.Text()) {
				fmt.Fprintln(w, scanner.Text())
			}
		}

		return scanner.Err()
	}
}

// Replace returns a filter that writes all lines replacing non-overlapping
// instances of 'old' with 'replace'.
func Replace(old, replace string) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			fmt.Fprintln(w, strings.ReplaceAll(scanner.Text(), old, replace))
		}

		return scanner.Err()
	}
}

// ReplaceRegex returns a filter that writes all lines replacing matches of the
// compiled regular expression 'regex' with 'replace'.
func ReplaceRegex(regex *regexp.Regexp, replace string) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			fmt.Fprintln(w, regex.ReplaceAllString(scanner.Text(), replace))
		}

		return scanner.Err()
	}
}

// MIT License

// Copyright (c) 2019 John Arundel, 2022 Dyson Simmons

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Frequency returns a filter that writes unique lines from the input, prefixed
// with a frequency count, in descending numerical order (most frequent lines
// first). Lines with equal frequency will be sorted alphabetically.
//
// This is a common pattern in shell scripts to find the most
// frequently-occurring lines in a file:
//
// sort testdata/freq.input.txt |uniq -c |sort -rn
//
// Frequency's behaviour is like the combination of Unix `sort`, `uniq -c`, and
// `sort -rn` used here.
//
// Like `uniq -c`, Freq left-pads its count values if necessary to make them
// easier to read:
//
//	10 apple
//	 4 banana
//	 2 orange
//	 1 kumquat
func Frequency() func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		freq := map[string]int{}

		type frequency struct {
			line  string
			count int
		}

		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			freq[scanner.Text()]++
		}

		freqs := make([]frequency, 0, len(freq))

		var maxCount int

		for line, count := range freq {
			freqs = append(freqs, frequency{line, count})

			if count > maxCount {
				maxCount = count
			}
		}

		sort.Slice(freqs, func(i, j int) bool {
			x, y := freqs[i].count, freqs[j].count

			if x == y {
				return freqs[i].line < freqs[j].line
			}

			return x > y
		})

		fieldWidth := len(strconv.Itoa(maxCount))

		for _, item := range freqs {
			fmt.Fprintf(w, "%*d %s\n", fieldWidth, item.count, item.line)
		}

		return nil
	}
}
