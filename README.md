# Sketch

Sketch is a simple lisp interpreter, written in Go.

## Install

1. [Install Go](https://golang.org/doc/install)
2. Run `go get -u github.com/jamesroutley/sketch`

You can now run the Sketch interpreter in a shell:

```sh
$ sketch
```

Running it without any arguments will drop you in a REPL (read, eval, print
loop). Here, you can use the language interactively. For example, you can:

**Perform maths calculations**

```
user> (+ 1 1)
2
```

Sketch (and most other Lisps) use prefix notation, sometimes called Polish
notation. Operators, such as `+` which adds two numbers, are written before the
arguments they operate on. This might seem unusual, but what it lacks in
familiarity it makes up in consistency.

Unlike other languages which use infix notation for maths (e.g. `1 + 1`),
there's no special syntax for maths functions. Addition just calls a function,
which takes two arguments, and returns the sum of them.

**Define variables**

```
user> (def! a 1)
1
user> a
1
```

**Define your own functions**

```
user> (defn add-1 (x) (+ x 1))
#<function>
user> (add-1 2)
3
```
