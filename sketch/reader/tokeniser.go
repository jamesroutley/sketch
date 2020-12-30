package reader

import (
	"regexp"
)

type Token struct {
	Value string
	// TODO: these can't be negative - could make uint
	Line int
	Row  int
	// TODO: this would store the character number - is this useful?
	char int
}

func Tokenize(s string) []string {
	re := regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" + `~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" + `,;)]*)`)
	return re.FindAllString(s, -1)
}

// func TokenizeNew(s string) []*Token {
// 	// Iterate over the string line by line. Sketch doesn't have significant
// 	// whitespace, but comments comment out the rest of the line
// 	lines := strings.Split(s, "\n")
// 	for _, line := range lines {

// 	}
// }
