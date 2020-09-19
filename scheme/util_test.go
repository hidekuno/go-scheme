/*
   Go lang 4th study program.
   This is test program for mini scheme subset.

   ex.) go test -v lisp.go lisp_test.go

   hidekuno@gmail.com
*/
package scheme

import (
	"os"
	"testing"
)

func TestIdentity(t *testing.T) {
	testCode := [][]string{
		{"(identity 100)", "100"},
		{"(identity \"ABC\")", "\"ABC\""},

		{"(identity 100 200)", "E1007"},
		{"(identity)", "E1007"},
	}
	executeTest(testCode, "identity", t)
}
func TestTime(t *testing.T) {
	testCode := [][]string{
		{"(time)", "E1007"},
		{"(time #\\abc)", "E0004"},
	}
	executeTest(testCode, "time", t)
}
func TestGetEnvironment(t *testing.T) {
	testCode := [][]string{
		{"(get-environment-variable)", "E1007"},
		{"(get-environment-variable 10)", "E1015"},
		{"(get-environment-variable a)", "E1008"},
		{"(get-environment-variable \"HOME\")", "\"" + os.Getenv("HOME") + "\""},
	}
	executeTest(testCode, "get-environment-variable", t)
}
