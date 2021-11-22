/*
   Go lang 4th study program.
   This is test program for mini scheme subset.

   ex.) go test -v lisp.go lisp_test.go

   hidekuno@gmail.com
*/
package scheme

import (
	"testing"
)

func TestNot(t *testing.T) {
	testCode := [][]string{
		{"(not (= 2 1))", "#t"},
		{"(not (= 1 1))", "#f"},
		{"(not 10)", "#t"},
		{"(not \"abc\")", "#t"},

		{"(not)", "E1007"},
		{"(not #t #f)", "E1007"},
		{"(not a)", "E1008"},
	}
	executeTest(testCode, "not", t)
}
func TestBoolean(t *testing.T) {
	testCode := [][]string{
		{"(boolean (= 1 1))", "#t"},
		{"(boolean (= 2 1))", "#f"},
		{"(boolean 10)", "#t"},
		{"(boolean \"abc\")", "#t"},

		{"(boolean)", "E1007"},
		{"(boolean 1 2)", "E1007"},
		{"(boolean a)", "E1008"},
	}
	executeTest(testCode, "boolean", t)
}
func TestBooleanEq(t *testing.T) {
	testCode := [][]string{
		{"(boolean=? #t #t)", "#t"},
		{"(boolean=? #f #f)", "#t"},
		{"(boolean=? #t #f)", "#f"},
		{"(boolean=? #t #t #t)", "#t"},
		{"(boolean=? #f #f #f)", "#t"},
		{"(boolean=? #f #f #t)", "#f"},

		{"(boolean=?)", "E1007"},
		{"(boolean=? #t)", "E1007"},
		{"(boolean=? 10 #f)", "E1001"},
		{"(boolean=? #t 10)", "E1001"},
	}
	executeTest(testCode, "boolean=?", t)
}
