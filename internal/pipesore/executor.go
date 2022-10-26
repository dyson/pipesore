package pipesore

import (
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"

	"github.com/dyson/pipesore/pkg/pipeline"
)

func Execute(input string, in io.Reader, out io.Writer) error {
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

	for _, inFunction := range e.tree.functions {
		name := strings.ToLower(inFunction.name)

		filter, ok := pipeline.Filters[name]
		if !ok {
			return fmt.Errorf("unknown function %s()", inFunction.name)
		}

		filterType := filter.Type()

		args, err := e.convertArguments(inFunction, filterType)
		if err != nil {
			return err
		}

		p.Filter(filter.Call(args)[0].Interface().(func(io.Reader, io.Writer) error))
	}

	if _, err := p.Output(e.writer); err != nil {
		return fmt.Errorf("error filtering pipeline: %w", err)
	}

	return nil
}

func (e executor) convertArguments(inFunction function, filterType reflect.Type) ([]reflect.Value, error) {
	if len(inFunction.arguments) != filterType.NumIn() {
		return nil, fmt.Errorf("wrong number of arguments in call to %s(): expected %d, got %d", inFunction.name, filterType.NumIn(), len(inFunction.arguments))
	}

	args := []reflect.Value{}

	for i := 0; i < len(inFunction.arguments); i++ {
		inArg := inFunction.arguments[i]
		filterArgType := filterType.In(i)

		switch filterArgType.String() {
		case "string":
			if reflect.TypeOf(inArg).String() != "string" {
				return nil, fmt.Errorf("expected argument %d in call to %s() to be string, got: %v (%T)", i+1, inFunction.name, inArg, inArg)
			}

			args = append(args, reflect.ValueOf(inArg))

		case "int":
			if reflect.TypeOf(inArg).String() != "int" {
				return nil, fmt.Errorf("expected argument %d in call to %s() to be int, got: %v (%T)", i+1, inFunction.name, inArg, inArg)
			}

			args = append(args, reflect.ValueOf(inArg))

		case "*regexp.Regexp":
			if reflect.TypeOf(inArg).String() != "string" {
				return nil, fmt.Errorf("expected argument %d in call to %s() to be valid regex.Regexp string, got: %v (%T)", i+1, inFunction.name, inArg, inArg)
			}

			re, err := regexp.Compile(inArg.(string))
			if err != nil {
				return nil, fmt.Errorf("expected argument %d in call to %s() to be valid regex.Regexp string, got: %v (%T), err: %v", i+1, inFunction.name, inArg, inArg, err)
			}

			args = append(args, reflect.ValueOf(re))
		}
	}

	return args, nil
}
