package pipesore

import "fmt"

type syntaxError struct {
	err error
	position
}

func newSyntaxError(err error, position position) *syntaxError {
	return &syntaxError{
		err:      err,
		position: position,
	}
}

func (pe *syntaxError) Error() string {
	return pe.err.Error()
}

type filterNameError struct {
	err error
	position
	name       string
	suggestion string
}

func newFilterNameError(err error, position position, name, suggestion string) *filterNameError {
	return &filterNameError{
		err:        err,
		position:   position,
		name:       name,
		suggestion: suggestion,
	}
}

func (fne *filterNameError) Error() string {
	return fne.err.Error()
}

type filterArgumentError struct {
	err  error
	name string
	position
}

func newFilterArgumentError(err error, position position, name string) *filterArgumentError {
	return &filterArgumentError{
		err:      err,
		position: position,
		name:     name,
	}
}

func (fne *filterArgumentError) Error() string {
	return fne.err.Error()
}

func newFormattedError(err error, input string, position position, help string) error {
	red := "\x1b[31m"
	undercurl := "\x1b[4:3m"
	reset := "\x1b[0m"

	var inputBefore, inputAfter string

	start := position.start
	end := position.end

	// handle EOF
	if len(input) == start {
		input += " "
	}

	if start > 0 {
		inputBefore = input[:start]
	}

	inputError := input[start:end]

	if len(input) > end {
		inputAfter = input[end:]
	}

	if help != "" {
		help = "\n" + help
	}

	return fmt.Errorf(
		"%w:\n\t%s%s%s%s%s%s%s.",
		err,
		inputBefore,
		red,
		undercurl,
		inputError,
		reset,
		inputAfter,
		help,
	)
}
