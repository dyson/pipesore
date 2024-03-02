package pipesore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dyson/pipesore/pkg/pipeline"
)

func Run(version, commit, date string) (int, error) {
	seeHelp := fmt.Sprintf("See '%s --help'", filepath.Base(os.Args[0]))

	if len(os.Args) != 2 {
		return 1, fmt.Errorf("error: define a single pipeline or option.\n%s.", seeHelp)
	}

	input := os.Args[1]

	switch input {
	case "-h", "--help":
		printHelp()
		return 0, nil
	case "-v", "--version":
		fmt.Printf("pipesore version %s, commit %s, date %s\n", version, commit, date)
		return 0, nil
	case "":
		return 1, fmt.Errorf("error: no pipeline defined.\n%s.", seeHelp)
	}

	err := execute(input, os.Stdin, os.Stdout)
	if err != nil {
		var syntaxError *syntaxError
		if errors.As(err, &syntaxError) {
			return 1, newFormattedError(err, input, syntaxError.position, seeHelp)
		}

		var filterNameError *filterNameError
		if errors.As(err, &filterNameError) {
			if filterNameError.suggestion != "" {
				definition := pipeline.Filters[filterNameError.suggestion].Definition
				seeHelp = fmt.Sprintf("Did you mean '%s'?\n%s", definition, seeHelp)
			}

			return 1, newFormattedError(err, input, filterNameError.position, seeHelp)
		}

		var filterArgumentError *filterArgumentError
		if errors.As(err, &filterArgumentError) {
			help := fmt.Sprintf("%s. %s", pipeline.Filters[filterArgumentError.name].Definition, seeHelp)

			return 1, newFormattedError(err, input, filterArgumentError.position, help)
		}

		return 1, fmt.Errorf("%w.\n%s.", err, seeHelp)
	}

	return 0, nil
}
