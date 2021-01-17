// Package reader implements Sketch's reader - the component which converts raw
// source code into an AST.
package reader

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jamesroutley/sketch/sketch/types"
)

type Reader struct {
	Tokens   []string
	Position int
}

func NewReader(tokens []string) *Reader {
	return &Reader{
		Tokens:   tokens,
		Position: 0,
	}
}

func (r *Reader) Peek() (string, error) {
	if r.Position == len(r.Tokens) {
		return "", fmt.Errorf("EOF")
	}
	return strings.Trim(r.Tokens[r.Position], " ,\n\t"), nil
}

func (r *Reader) Next() (string, error) {
	if r.Position == len(r.Tokens) {
		return "", fmt.Errorf("EOF")
	}
	current := strings.Trim(r.Tokens[r.Position], " ,\n\t")
	r.Position++
	return current, nil
}

func Read(s string) (types.SketchType, error) {
	ast, err := ReadWithoutReaderMacros(s)
	if err != nil {
		return nil, err
	}

	if _, ok := ast.(*types.SketchComment); ok {
		// If the read line is a comment, return this magic error which the
		// REPL catches to not print anything.
		return nil, fmt.Errorf("read comment")

	}

	// Reader macros
	ast = stripComments2(ast)
	ast = expandModuleLookup(ast)
	return ast, nil
}

func ReadWithoutReaderMacros(s string) (types.SketchType, error) {
	tokens := Tokenize(s)
	reader := NewReader(tokens)
	return ReadForm(reader)
}

func ReadForm(reader *Reader) (types.SketchType, error) {
	token, err := reader.Peek()
	if err != nil {
		return nil, err
	}

	switch token {
	case "(":
		// Increment the position pointer
		_, err = reader.Next()
		if err != nil {
			return nil, err
		}
		return ReadList(reader)
	default:
		return ReadAtom(reader)
	}
}

func ReadList(reader *Reader) (types.SketchType, error) {
	var items []types.SketchType
	for {
		// TODO: error case when we hit file without closing bracket
		tok, err := reader.Peek()
		if err != nil {
			return nil, err
		}
		if tok == ")" {
			// Increment the position pointer
			_, err := reader.Next()
			if err != nil {
				return nil, err
			}
			return &types.SketchList{
				Items: items,
			}, nil
		}
		item, err := ReadForm(reader)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
}

func ReadAtom(reader *Reader) (types.SketchType, error) {
	token, err := reader.Next()
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(token, ";") {
		comment := strings.TrimLeft(token, "; ")
		return &types.SketchComment{
			Value: comment,
		}, nil
	}

	if num, err := strconv.Atoi(token); err == nil {
		return &types.SketchInt{
			Value: num,
		}, nil
	}

	if token == "true" {
		return &types.SketchBoolean{Value: true}, nil
	}
	if token == "false" {
		return &types.SketchBoolean{Value: false}, nil
	}

	if token == "nil" {
		return &types.SketchNil{}, nil
	}

	if strings.HasPrefix(token, `"`) {
		if !strings.HasSuffix(token, `"`) {
			return nil, fmt.Errorf("unclosed string")
		}
		token = strings.Trim(token, `"`)
		// Replace the literal characters 'slash n' with a newline
		token = strings.Replace(token, `\n`, "\n", -1)

		return &types.SketchString{
			Value: strings.Trim(token, `"`),
		}, nil
	}

	return &types.SketchSymbol{
		Value: token,
	}, nil
}
