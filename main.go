package main

import (
	"fmt"

	"github.com/ecnahc515/diy-lisp/lisp"
)

func main() {
	program := `
(define fact
   ;; Factorial function
   (lambda (n)
       (if (eq n 0)
           1 ; Factorial of 0 is 1, and we deny
             ; the existence of negative numbers
           (* n (fact (- n 1))))))

(fact 1 10)
`
	source := []byte(program)
	exps, err := lisp.SplitExps(source)
	if err != nil {
		fmt.Println(err)
	}
	var expList []string
	for _, v := range exps {
		atoms := string(v[:])
		expList = append(expList, atoms)
	}
	// Print the list of expressions with each expression as a string
	fmt.Println(expList)
}
