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

func ReadStr(s string) (types.SketchType, error) {
	ast, err := ReadStrWithComments(s)
	if err != nil {
		return nil, err
	}
	switch castAST := ast.(type) {
	case *types.SketchComment:
		// If the read line is a comment, return empty string. This means
		// we don't print anything in the REPL
		return nil, fmt.Errorf("read comment")
	case *types.SketchList:
		ast = stripComments(castAST)
	default:
		// continue
	}
	return ast, nil
}

func ReadStrWithComments(s string) (types.SketchType, error) {
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

		return &types.SketchString{
			Value: strings.Trim(token, `"`),
		}, nil
	}

	return &types.SketchSymbol{
		Value: token,
	}, nil
}

func DebugType(m types.SketchType) {
	fmt.Println(debugType(m, 0))
}

func debugType(m types.SketchType, indent int) string {
	switch tok := m.(type) {
	case *types.SketchList:
		itemStrings := make([]string, len(tok.Items))
		for i, item := range tok.Items {
			itemStrings[i] = debugType(item, 0)
		}
		return fmt.Sprintf("(%s)", strings.Join(itemStrings, " "))
	case *types.SketchInt:
		return fmt.Sprintf("int:%d ", tok.Value)
	case *types.SketchSymbol:
		return fmt.Sprintf("symbol:`%s` ", tok.Value)
	case *types.SketchFunction:
		return "#<function>"
	case *types.SketchBoolean:
		if tok.Value {
			return "boolean:true"
		}
		return "boolean:false"
	case *types.SketchNil:
		return "nil"
	default:
		return tok.String()
	}
}

func DebugTokens(s string) {
	tokens := Tokenize(s)
	for _, tok := range tokens {
		fmt.Printf("`%s`\n", tok)
	}
}

func stripComments(list *types.SketchList) types.SketchType {
	var newItems []types.SketchType
	for _, item := range list.Items {
		switch item := item.(type) {
		case *types.SketchComment:
			// skip
		case *types.SketchList:
			newItem := stripComments(item)
			newItems = append(newItems, newItem)
		default:
			newItems = append(newItems, item)
		}
	}

	list.Items = newItems
	return list
}
