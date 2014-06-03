package lisp

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
)

var (
	NoBeginParen = errors.New("No beginning paren found")
	leftParen    = []byte{'('}
	rightParen   = []byte{')'}
)

// func Parse(source []byte) []interface{} {
// 	exps := splitExps(source)
// }

func removeComments(source []byte) []byte {
	re := regexp.MustCompile(`;.*\n`)
	return re.ReplaceAllLiteral(source, []byte("\n"))
}

// findMatchingParen takes a []byte, the index of an opening paren, and returns
// the index of the matching closing paren.
func findMatchingParen(source []byte, start int) (int, error) {
	if source[start] != '(' {
		return 0, NoBeginParen
	}
	pos := start
	openBrackets := 1
	for openBrackets > 0 {
		pos++
		if len(source) == pos {
			return 0, fmt.Errorf("Incomplete expression %s", source[start:])
		}
		if source[pos] == '(' {
			openBrackets++
		}
		if source[pos] == ')' {
			openBrackets--
		}
	}
	return pos, nil
}

// firstExpression splits a []byte into (exp, rest), where exp is the first
// expression in the []byte and rest is the rest of the []byte after this
// expression.
func firstExpression(source []byte) ([]byte, []byte, error) {
	source = bytes.TrimSpace(source)
	if source[0] == '\'' {
		exp, rest, err := firstExpression(source[1:])
		if err != nil {
			return nil, nil, err
		}
		exp = append([]byte(source), exp...)
		return exp, rest, nil
	} else if source[0] == '(' {
		last, err := findMatchingParen(source, 0)
		if err != nil {
			return nil, nil, err
		}
		return source[:last+1], source[last+1:], nil
	} else {
		re := regexp.MustCompile(`^[^\s)']+`)
		index := re.FindIndex(source)
		if index == nil {
			panic(fmt.Sprintf("No match for %s", source))
		}
		end := index[1]
		atom := source[:end]
		return atom, source[end:], nil
	}
}

// splitExps splits a source string into sub expressions that can be parsed
// individually.
func SplitExps(source []byte) ([][]byte, error) {
	rest := bytes.TrimSpace(source)
	exps := make([][]byte, 0)
	for len(rest) > 0 {
		var exp []byte
		var err error
		exp, rest, err = firstExpression(rest)
		if err != nil {
			return nil, err
		}
		exps = append(exps, exp)
	}
	return exps, nil
}

func Tokenize(source []byte) [][]byte {
	source = bytes.Replace(source, leftParen, []byte("( "), -1)
	source = bytes.Replace(source, rightParen, []byte(" )"), -1)
	return bytes.Split(source, []byte(" "))
}
