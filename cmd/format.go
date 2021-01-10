package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jamesroutley/sketch/sketch/printer"
	"github.com/jamesroutley/sketch/sketch/reader"
	"github.com/spf13/cobra"
)

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:   "format [path]...",
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
--write/-w flag to get Sketchfmt to write the file itself.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, filename := range args {
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Fatal(err)
			}

			ast, err := reader.ReadStrWithComments(fmt.Sprintf("(do %s)", data))
			if err != nil {
				log.Fatal(err)
			}

			formatted := printer.PrettyPrintTopLevelDo(ast)

			write, err := cmd.Flags().GetBool("write")
			if err != nil {
				log.Fatal(err)
			}

			if write {
				if err := ioutil.WriteFile(filename, []byte(formatted+"\n"), 0777); err != nil {
					log.Fatal(err)
				}
				continue
			}

			fmt.Println(formatted)
		}
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)

	// TODO; move to other type of flag
	formatCmd.Flags().BoolP("write", "w", false, "Write to the file, instead of printing to stdout")
}
