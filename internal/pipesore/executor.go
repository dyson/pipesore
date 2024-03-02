package pipesore

import (
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"

	"github.com/dyson/pipesore/pkg/levenshtein"
	"github.com/dyson/pipesore/pkg/pipeline"
)

func execute(input string, in io.Reader, out io.Writer) error {
	tree, err := newParser(newLexer(input)).parse()
	if err != nil {
		return fmt.Errorf("error parsing pipeline: %w", err)
	}

	return newExecutor(tree, in, out).execute()
}

type executor struct {
	tree   *ast
	reader io.Reader
	writer io.Writer
}

func newExecutor(tree *ast, r io.Reader, w io.Writer) *executor {
	return &executor{tree: tree, reader: r, writer: w}
}

func (e executor) execute() error {
	p := pipeline.NewPipeline(e.reader)

	for _, inFilter := range e.tree.filters {
		name := strings.ToLower(inFilter.name)

		filter, ok := pipeline.Filters[name]
		if !ok {
			lowestScore := len(name)
			suggestion := ""
			for f := range pipeline.Filters {
				distance := levenshtein.Distance(name, f)
				if distance < lowestScore {
					lowestScore = distance
					suggestion = f
				}
			}

			return newFilterNameError(
				fmt.Errorf("error running pipeline: unknown filter '%s()'", inFilter.name),
				inFilter.position,
				inFilter.name,
				suggestion,
			)
		}

		filterType := filter.Value.Type()

		args, err := e.convertArguments(inFilter, filterType)
		if err != nil {
			return newFilterArgumentError(
				fmt.Errorf("error running pipeline: %w", err),
				inFilter.position,
				name,
			)
		}

		p.Filter(filter.Value.Call(args)[0].Interface().(func(io.Reader, io.Writer) error))
	}

	if _, err := p.Output(e.writer); err != nil {
		return fmt.Errorf("error filtering pipeline: %w", err)
	}

	return nil
}

func (e executor) convertArguments(inFilter filter, filterType reflect.Type) ([]reflect.Value, error) {
	if len(inFilter.arguments) != filterType.NumIn() {
		argument := "argument"
		if filterType.NumIn() > 1 {
			argument += "s"
		}

		return nil, fmt.Errorf("expected %d %s in call to '%s()', got %d", filterType.NumIn(), argument, inFilter.name, len(inFilter.arguments))
	}

	args := []reflect.Value{}

	for i := 0; i < len(inFilter.arguments); i++ {
		inArg := inFilter.arguments[i]
		filterArgType := filterType.In(i)

		switch filterArgType.String() {
		case "string":
			if reflect.TypeOf(inArg).String() != "string" {
				return nil, fmt.Errorf("expected argument %d in call to '%s()' to be a string, got %v (%T)", i+1, inFilter.name, inArg, inArg)
			}

			args = append(args, reflect.ValueOf(inArg))

		case "int":
			if reflect.TypeOf(inArg).String() != "int" {
				return nil, fmt.Errorf("expected argument %d in call to '%s()' to be an int, got %v (%T)", i+1, inFilter.name, inArg, inArg)
			}

			args = append(args, reflect.ValueOf(inArg))

		case "*regexp.Regexp":
			if reflect.TypeOf(inArg).String() != "string" {
				return nil, fmt.Errorf("expected argument %d in call to '%s()' to be a valid regex.Regexp string, got %v (%T)", i+1, inFilter.name, inArg, inArg)
			}

			re, err := regexp.Compile(inArg.(string))
			if err != nil {
				return nil, fmt.Errorf("expected argument %d in call to '%s()' to be a valid regex.Regexp string, got %v (%T), err %v", i+1, inFilter.name, inArg, inArg, err)
			}

			args = append(args, reflect.ValueOf(re))
		}
	}

	return args, nil
}
