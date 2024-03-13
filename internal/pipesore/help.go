package pipesore

import (
	"fmt"
	"strings"

	"github.com/dyson/pipesore/pkg/pipeline"
)

func printHelp() {
	sb := &strings.Builder{}

	w := func(s string) {
		wrap(sb, s)
	}

	w("pipesore - command-line text processor")
	w("")
	w("Usage:")
	w("  pipesore '<filter>[ | <filter>]...'")
	w("  pipesore [option]")
	w("")
	w("Example:")
	w("  $ echo \"cat cat cat dog bird bird bird bird\" | \\")
	w("  pipesore 'Replace(\" \", \"\\n\") | Frequency() | First(1)'")
	w("  4 bird")
	w("")
	w("Filters:")
	w("  All filters can be '|' (piped) together in any order, although not all ordering is logical.")
	w("")
	w("  All filter arguments are required. There are no assumptions about default values.")
	w("")
	w("  A filter prefixed with an \"!\" will return the opposite result of the non prefixed filter of the same name. For example `First(1)` would return only the first line of the input and `!First(1)` (read as not first) would skip the first line of the input and return all other lines.")
	w("")
	w("  ---")
	w("")
	for _, name := range pipeline.Filters.GetOrderedNames() {
		filter := pipeline.Filters[name]
		w("  " + filter.Definition)
		w("    " + filter.Description)
		w("")
	}
	w("Options:")
	w("  -h, --help     show this help message")
	w("  -v, --version  show pipesore version")

	fmt.Printf(sb.String())
}

func wrap(sb *strings.Builder, s string) {
	width := 80

	prefix := ""
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' {
			prefix = s[:i]
			s = s[i:]
			break
		}
	}

	width = width - len(prefix)

	lastSpace := 0
	lastBreak := 0

	for i := 0; i < len(s); {
		if i == lastBreak+width {
			sb.WriteString(prefix)
			sb.WriteString(s[lastBreak:lastSpace])
			sb.WriteString("\n")
			lastBreak = lastSpace + 1
			i = lastBreak
		}

		if s[i] == ' ' {
			lastSpace = i
		}

		i++
	}

	sb.WriteString(prefix)
	sb.WriteString(s[lastBreak:])
	sb.WriteString("\n")
}
