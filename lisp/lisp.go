package lisp

import "fmt"

func PrintRepr(s [][]byte) {
	var expList []string
	l := len(s)
	for index, exp := range s {
		atoms := "\"" + string(exp) + "\""
		expList = append(expList, atoms)
		if index+1 < l {
			expList = append(expList, ",")
		}

	}
	fmt.Println(expList)
}
