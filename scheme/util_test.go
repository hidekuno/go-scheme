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

func TestEqv(t *testing.T) {
	testCode := [][]string{
		{"(eqv? 1.1 1.1)", "#t"},
		{"(eq? 1.1 1.1)", "#t"},
		{"(eqv? 1.1 1.2)", "#f"},
		{"(eqv? 10 (+ 2 8))", "#t"},
		{"(eqv? 1 2)", "#f"},
		{"(eqv? 1 1.0)", "#f"},
		{"(eqv? 1.0 1)", "#f"},
		{"(eq? (quote a) (quote a))", "#t"},
		{"(eq? (quote a) (quote b))", "#f"},
		{"(eq? (quote a) 10)", "#f"},
		{"(eq? #f #f)", "#t"},
		{"(eq? #t #f)", "#f"},
		{"(eq? #t 10)", "#f"},
		{"(eq? #\\a #\\a)", "#t"},
		{"(eq? #\\a #\\b)", "#f"},
		{"(eq? #\\space #\\space)", "#t"},

		{"(eqv?)", "E1007"},
		{"(eqv? 10 10 10)", "E1007"},
		{"(eq? 10 10 10)", "E1007"},
		{"(eq? 10 a)", "E1008"},
		{"(eq? a 10)", "E1008"},
	}
	executeTest(testCode, "eqv", t)
}
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
