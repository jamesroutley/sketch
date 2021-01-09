// Package sketch implements Sketch's interpreter
package sketch

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/chzyer/readline"
	"github.com/jamesroutley/sketch/sketch/environment"
	"github.com/jamesroutley/sketch/sketch/evaluator"
	"github.com/jamesroutley/sketch/sketch/printer"
	"github.com/jamesroutley/sketch/sketch/reader"
)

func RunFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	ast, err := reader.ReadStr(fmt.Sprintf("(do %s)", data))
	if err != nil {
		return err
	}

	_, err = evaluator.Evaluate(ast)
	return err
}

func Repl() error {
	env, err := evaluator.RootEnvironment()
	if err != nil {
		return err
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "user> ",
		HistoryFile: "/Users/jamesroutley/.sketchhistory",
	})
	if err != nil {
		return err
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}
		line = strings.TrimSuffix(line, "\n")
		output, err := Rep(line, env)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(output)
	}
	return nil
}

// Rep - read, evaluate, print
func Rep(s string, env *environment.Env) (string, error) {
	ast, err := reader.ReadStr(s)
	if err != nil {
		return "", err
	}
	evaluated, err := evaluator.Eval(ast, env)
	if err != nil {
		if err.Error() == "read comment" {
			return "", nil
		}
		return "", err
	}
	return printer.PrStr(evaluated), nil
}
