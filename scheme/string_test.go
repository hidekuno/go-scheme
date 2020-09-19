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
