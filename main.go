package main

import (
	"fmt"

	"github.com/ecnahc515/diy-lisp/lisp"
)

func main() {
	program := "(* 100 5)"
	source := []byte(program)
	// exprs, _ := lisp.SplitExps(source)
	// for _, v := range exprs {
	// 	tokens := lisp.Tokenize(v)
	// 	lisp.PrintRepr(tokens)
	// }
	_, err := lisp.Parse(source)
	if err != nil {
		fmt.Println(err)
	}

}
