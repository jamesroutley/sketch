package sketch

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jamesroutley/sketch/sketch/environment"
	"github.com/jamesroutley/sketch/sketch/evaluator"
	"github.com/jamesroutley/sketch/sketch/reader"
	"github.com/jamesroutley/sketch/sketch/types"
)

func TestFile(filename string) error {
	// We currently only support docstring tests, which is a test defined in a
	// function's docstring.

	// First, evaluate the file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	ast, err := reader.Read(fmt.Sprintf("(do %s)", data))
	if err != nil {
		return err
	}
	env, err := evaluator.RootEnvironment()
	if err != nil {
		return err
	}

	// Evaluate the file in a new child env this lets us just test items added
	// by the file
	child := env.ChildEnv()
	if _, err := evaluator.Eval(ast, child); err != nil {
		return err
	}

	// Then, iterate through any functions, and run docstring tests, if any are
	// specified
	for key, value := range child.Data {
		function, ok := value.(*types.SketchFunction)
		if !ok {
			continue
		}
		if err := runDocstringTest(key, child, function); err != nil {
			return err
		}
	}

	return nil
}

func runDocstringTest(name string, env *environment.Env, function *types.SketchFunction) error {
	docstring := function.Docs
	lines := strings.Split(docstring, "\n")

	examplesSection := false
	input := ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "Examples:" {
			fmt.Printf("Running docstring tests for %s\n", name)
			examplesSection = true
			continue
		}
		if !examplesSection {
			continue
		}

		if strings.HasPrefix(line, ">") {
			input = strings.TrimSpace(strings.TrimPrefix(line, ">"))
			continue
		}

		if !strings.HasPrefix(line, "->") {
			continue
		}
		expected := strings.TrimSpace(strings.TrimPrefix(line, "->"))

		ast, err := reader.Read(input)
		if err != nil {
			return err
		}
		actual, err := evaluator.Eval(ast, env)
		if err != nil {
			return err
		}

		if actual.String() != expected {
			return fmt.Errorf(
				"error running '%s' test: expected %s to eval to %s, got %s",
				name, input, expected, actual,
			)
		}

		input = ""
	}

	return nil
}
