package cmd

import (
	"fmt"

	"github.com/jamesroutley/sketch/sketch/errors"
)

func printError(err error) {
	xerr, ok := err.(*errors.Error)
	if !ok {
		fmt.Println(err)
		return
	}

	fmt.Println(xerr)
	fmt.Printf("\nCall stack:\n")
	for i := len(xerr.Stack) - 1; i >= 0; i-- {
		fmt.Println("  ", xerr.Stack[i])
	}
}
