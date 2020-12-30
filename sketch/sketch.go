// Package sketch implements Sketch's interpreter
package sketch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
	"github.com/jamesroutley/sketch/sketch/core"
	"github.com/jamesroutley/sketch/sketch/environment"
	"github.com/jamesroutley/sketch/sketch/printer"
	"github.com/jamesroutley/sketch/sketch/reader"
	"github.com/jamesroutley/sketch/sketch/types"
)

func RunFile(filename string) error {
	env, err := rootEnvironment()
	if err != nil {
		return err
	}

	_, err = Rep(fmt.Sprintf(`(load-file "%s")`, filename), env)
	return err
}

func Repl() error {
	env, err := rootEnvironment()
	if err != nil {
		return err
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "user> ",
		HistoryFile: "/Users/jamesroutley/.malhistory",
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
	t, err := reader.ReadStr(s)
	if err != nil {
		return "", err
	}
	t, err = Eval(t, env)
	if err != nil {
		return "", err
	}
	return printer.PrStr(t), nil
}

func rootEnvironment() (*environment.Env, error) {
	env := environment.NewEnv()
	for _, item := range core.Namespace {
		env.Set(item.Symbol.Value, item.Func)
	}

	for _, item := range core.Namespace {
		env.Set(item.Symbol.Value, item.Func)
	}

	// Eval function. Needs to be here, because it closes over `env`
	env.Set("eval", &types.MalFunction{
		Func: func(args ...types.MalType) (types.MalType, error) {
			return Eval(args[0], env)
		},
	})

	// Builtin functions defined in lisp
	if _, err := Rep("(def not (fn (a) (if a false true)))", env); err != nil {
		return nil, err
	}

	if _, err := Rep("(def load-file (fn (f) (eval (read-string (+ \"(do \" (slurp f) \"\nnil)\")))))", env); err != nil {
		return nil, err
	}

	// Load stdlib
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return nil, fmt.Errorf("$GOPATH not set")
	}
	stdlibDir := filepath.Join(gopath, "src", "github.com", "jamesroutley", "sketch", "sketch", "stdlib")
	stdlibFiles, err := filepath.Glob(filepath.Join(stdlibDir, "*.skt"))
	if err != nil {
		return nil, err
	}

	for _, filename := range stdlibFiles {
		if _, err := Rep(fmt.Sprintf(`(load-file "%s")`, filename), env); err != nil {
			return nil, err
		}
	}

	return env, nil
}
