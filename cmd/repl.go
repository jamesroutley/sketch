package cmd

import (
	"log"

	"github.com/jamesroutley/sketch/sketch"
	"github.com/spf13/cobra"
)

// replCmd represents the repl command
var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Launches the Sketch REPL",
	// Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := sketch.Repl(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(replCmd)
}
