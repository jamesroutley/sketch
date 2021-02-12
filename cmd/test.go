package cmd

import (
	"os"

	"github.com/jamesroutley/sketch/sketch"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Tests a sketch file",
	Long: `Sketch currently only supports 'docstring tests', where a function's
tests are defined in its docstring. For example:

(defn
    plus
    "Adds two numbers

    Examples:
    > (plus 1 1)
    -> 2

    > (plus -1 1)
    -> 0"
    (a b)
    (+ a b))

This function's docstring contains two tests. Each test mirrors what it would
be like to run some code at the REPL. Any line starting with a '>' is evaluated,
and its output is compared to the next line starting with a '->'

To avoid accidentally parsing lines which happen to start with '>' or '->' as a
test, tests must be defined below a line which says 'Examples:'.

In order to test a file, Sketch has to evaluate it, so any side effects will
will happen.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := sketch.TestFile(args[0]); err != nil {
			printError(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
