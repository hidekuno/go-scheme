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

func TestStringAppend(t *testing.T) {
	testCode := [][]string{
		{"(string-append \"a\" \"b\")", "\"ab\""},
		{"(string-append \"a\" \"b\"  \"c\")", "\"abc\""},
		{"(string-append)", "E1007"},
		{"(string-append 10)", "E1007"},
		{"(string-append a b)", "E1008"},
		{"(string-append \"a\" 10)", "E1015"},
	}
	executeTest(testCode, "string-append", t)
}
func TestFormat(t *testing.T) {
	testCode := [][]string{
		{"(format \"~D\" 10)", "\"10\""},
		{"(format \"~d\" 10)", "\"10\""},
		{"(format \"~X\" 10)", "\"A\""},
		{"(format \"~x\" 10)", "\"a\""},
		{"(format \"~O\" 10)", "\"12\""},
		{"(format \"~o\" 10)", "\"12\""},
		{"(format \"~B\" 10)", "\"1010\""},
		{"(format \"~b\" 10)", "\"1010\""},
		{"(define a \"~D\")", "a"},
		{"(define b 100)", "b"},
		{"(format a b)", "\"100\""},

		{"(format)", "E1007"},
		{"(format \"~B\")", "E1007"},
		{"(format \"~B\" 10 12)", "E1007"},
		{"(format 10 12)", "E1015"},
		{"(format \"~A\" #f)", "E1002"},
		{"(format \"~A\" 10)", "E1018"},
	}
	executeTest(testCode, "format", t)
}
