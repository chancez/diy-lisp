package lisp

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var (
	NoBeginParen = errors.New("No beginning paren found")
	leftParen    = []byte{'('}
	rightParen   = []byte{')'}
)

func Parse(source []byte) ([]interface{}, error) {
	source = removeComments(source)
	exps, err := SplitExps(source)
	if err != nil {
		return nil, err
	}
	// For each expression, we need to tokenize them, and then convert
	// each item into the proper values
	for _, exp := range exps {
		tokens := Tokenize(exp)
		tks, _, err := ReadFrom(tokens)
		if err != nil {
			return nil, err
		}
		for k, v := range tks {
			switch vv := v.(type) {
			case string:
				fmt.Printf("%v: is string - %q\n", k, vv)
			case float64:
				fmt.Printf("%v: is float64 - %q\n", k, vv)
			case int64:
				fmt.Printf("%v: is int - %q\n", k, vv)
			default:
				fmt.Printf("%v: ", k)
			}
		}
	}
	return nil, nil
}

func ReadFrom(tokens [][]byte) ([]interface{}, [][]byte, error) {
	if len(tokens) == 0 {
		return nil, nil, fmt.Errorf("Unexpected EOF while reading")
	}
	var token []byte
	// Pop off a token
	token, tokens = tokens[0], tokens[1:]
	if bytes.Equal(token, leftParen) {
		var l []interface{}
		for !bytes.Equal(tokens[0], rightParen) {
			var atom []interface{}
			var err error
			atom, tokens, err = ReadFrom(tokens)
			if err != nil {
				return nil, nil, err
			}
			l = append(l, atom[0])
		}
		// pop off ')'
		tokens = tokens[1:]
		return l, tokens, nil
	} else if bytes.Equal(token, rightParen) {
		return nil, nil, fmt.Errorf("unexpected )")
	} else {
		atom := Atom(token)
		a := make([]interface{}, 1)
		a[0] = atom
		return a, tokens, nil
	}
	return nil, nil, fmt.Errorf("Unexpected error")
}

func Atom(token []byte) interface{} {
	token2 := string(token)
	i, err := strconv.ParseInt(token2, 10, 64)
	if err == nil {
		return i
	} else {
		fmt.Println("err", err)
	}
	f, err := strconv.ParseFloat(token2, 64)
	if err == nil {
		return f
	}
	b, err := strconv.ParseBool(token2)
	if err == nil {
		return b
	}
	return token2
}

func removeComments(source []byte) []byte {
	re := regexp.MustCompile(`;.*\n`)
	return re.ReplaceAllLiteral(source, []byte("\n"))
}

// FindMatchingParen takes a []byte, the index of an opening paren, and returns
// the index of the matching closing paren.
func FindMatchingParen(source []byte, start int) (int, error) {
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

// FirstExpression splits a []byte into (exp, rest), where exp is the first
// expression in the []byte and rest is the rest of the []byte after this
// expression.
func FirstExpression(source []byte) ([]byte, []byte, error) {
	source = bytes.TrimSpace(source)
	if source[0] == '\'' {
		exp, rest, err := FirstExpression(source[1:])
		if err != nil {
			return nil, nil, err
		}
		exp = append([]byte(source), exp...)
		return exp, rest, nil
	} else if source[0] == '(' {
		last, err := FindMatchingParen(source, 0)
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
		exp, rest, err = FirstExpression(rest)
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
