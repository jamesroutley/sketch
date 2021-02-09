// Package sketch implements Sketch's interpreter
package sketch

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/chzyer/readline"
	"github.com/jamesroutley/sketch/sketch/environment"
	"github.com/jamesroutley/sketch/sketch/errors"
	"github.com/jamesroutley/sketch/sketch/evaluator"
	"github.com/jamesroutley/sketch/sketch/printer"
	"github.com/jamesroutley/sketch/sketch/reader"
)

func RunFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	ast, err := reader.Read(fmt.Sprintf("(do %s)", data))
	if err != nil {
		return err
	}

	_, err = evaluator.Evaluate(ast)
	if err != nil {
		formatError(err)
	}
	return nil
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
	ast, err := reader.Read(s)
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

func formatError(err error) {
	if err == nil {
		return
	}

	fmt.Println(err)

	skerr, ok := err.(*errors.Error)
	if !ok {
		return
	}

	// If we've got it, print the stack trace pulled from the environment
	// list
	for env := skerr.Env; env != nil; env = env.Outer {
		// TODO
		// if !env.FunctionEnv {
		// 	continue
		// }
		fmt.Printf("- %+v\n", env.FunctionName)
		fmt.Printf("%+v\n", envKeys(env))
	}
}

func envKeys(env *environment.Env) (keys []string) {
	for key := range env.Data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}
