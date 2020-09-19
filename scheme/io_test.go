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

func TestDisplay(t *testing.T) {
	testCode := [][]string{
		{"(display)", "E1007"},
	}
	executeTest(testCode, "display", t)
}
func TestNewLine(t *testing.T) {
	testCode := [][]string{
		{"(newline 10)", "E1007"},
	}
	executeTest(testCode, "newline", t)
}
func TestLoadFile(t *testing.T) {
	testCode := [][]string{
		{"(load-file)", "E1007"},
		{"(load-file 10)", "E1015"},
		{"(load-file a)", "E1008"},
		{"(load-file \"example/no.scm\")", "E1014"},
		{"(load-file \"/tmp\")", "E1016"},
		{"(load-file \"/etc/sudoers\")", "E9999"},
	}
	executeTest(testCode, "load-file", t)
}
