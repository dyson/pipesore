package pipeline

import (
	"bufio"
	"container/ring"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	Filters = map[string]reflect.Value{
		// "columns":      reflect.ValueOf(Columns),
		"countlines":   reflect.ValueOf(CountLines),
		"countrunes":   reflect.ValueOf(CountRunes),
		"countwords":   reflect.ValueOf(CountWords),
		"first":        reflect.ValueOf(First),
		"!first":       reflect.ValueOf(NotFirst),
		"frequency":    reflect.ValueOf(Frequency),
		"join":         reflect.ValueOf(Join),
		"last":         reflect.ValueOf(Last),
		"!last":        reflect.ValueOf(NotLast),
		"match":        reflect.ValueOf(Match),
		"!match":       reflect.ValueOf(NotMatch),
		"matchregex":   reflect.ValueOf(MatchRegex),
		"!matchregex":  reflect.ValueOf(NotMatchRegex),
		"replace":      reflect.ValueOf(Replace),
		"replaceregex": reflect.ValueOf(ReplaceRegex),
	}
)

// TODO
func Columns(delimiter int, columns string) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
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

// Join returns a filter that writes all lines as a single string separated by 'delimiter'.
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

// MIT License

// Frequency returns a filter that writes unique lines from the input, prefixed with a frequency
// count, in descending numerical order (most frequent lines first). Lines with
// equal frequency will be sorted alphabetically.
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
