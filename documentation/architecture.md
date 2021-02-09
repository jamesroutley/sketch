# Architecture

This document covers the architecture of Sketch's interpreter. It's only
relevant if you want to understand or modify the interpreter itself.

## Stack traces

Stack traces are printed when an error occurs during evaluation. They show the
series of function calls that happened in the run up to an error.

They're implemented entirely within `evaluator.Eval`. When this function is
called, we initialise a slice of stack frames, which we append to whenever we:

1. explicitly call a non tail call optimised function
2. loop back to the top of `Eval` when calling a tail call optimised function.

If an error happens, we `errors.Wrap` it to add the stack frames to the error
itself. `Eval` is a recursive function, and so any errors which happen in sub
calls to `Eval` are also wrapped. In this way, we can build up the full list of
functions that were called which led to the error.
