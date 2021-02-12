package cmd

import (
	"os"

	"github.com/jamesroutley/sketch/sketch"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "run <file.skt>",
	Short: "Runs a Sketch program",
	// Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := sketch.RunFile(args[0]); err != nil {
			printError(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
