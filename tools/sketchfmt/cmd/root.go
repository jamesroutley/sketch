// Package cmd implements sketchfmt's CLI
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jamesroutley/sketch/sketch/printer"
	"github.com/jamesroutley/sketch/sketch/reader"
	"github.com/spf13/cobra"
)

var rewrite bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sketchfmt [path]",
	Short: "Automatically formats Sketch files",
	Long: `Sketchfmt is an autoformatter for the Sketch language.

Sketchfmt is ultimately pragmatic. It exists to simplify development and remove
discussions around formatting Sketch code. It might look different to how other
Lisp code is conventionally formatted, but the output of Sketchfmt is the
definition of idiomatic Sketch code.

It formats with the following rules:

1.  Indentation is done with tabs
2.  When formatting files, one newline is placed between each top level
	expression.
3.	Any top level comments are assumed to be about the expression which comes
	below it. The comment is placed on the line above the next expression.
4.  If a list doesn't contain a comment, it will be put on a single line, if
	that results in a line of less than 80 chars. If not, each item in the list
	will be placed on a separate line.
5.  If a list contains a comment, each element is always put on a separate
	line. The comment is assumed to be about the item which came before the
	comment, and the comment is placed on the same line as it.

By default, Sketchfmt prints the formatted file to stdout. You can supply the
--rewrite/-r flag to get Sketchfmt to rewrite the file itself.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}

		ast, err := reader.ReadStrWithComments(fmt.Sprintf("(do %s)", data))
		if err != nil {
			log.Fatal(err)
		}

		formatted := printer.PrettyPrintTopLevelDo(ast)

		if rewrite {
			if err := ioutil.WriteFile(filename, []byte(formatted), 0777); err != nil {
				log.Fatal(err)
			}
			return
		}

		fmt.Println(formatted)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&rewrite, "rewrite", "r", false, "Rewrite the file, instead of printing to stdout")
}
