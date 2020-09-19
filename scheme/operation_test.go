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

func TestPlus(t *testing.T) {
	testCode := [][]string{
		{"(+ 1 2)", "3"},
		{"(+ 3 1.5)", "4.5"},
		{"(+ 1 1.5 1.25)", "3.75"},

		{"(+ 10.2)", "E1007"},
		{"(+ )", "E1007"},
		{"(+ #t 10.2)", "E1003"},
	}
	executeTest(testCode, "plus", t)
}
func TestMinus(t *testing.T) {
	testCode := [][]string{
		{"(- 6 1)", "5"},
		{"(- 3 1.5 0.25)", "1.25"},
		{"(- 5.75 1.5)", "4.25"},

		{"(- 10.2 #f)", "E1003"},
		{"(- 10.2)", "E1007"},
		{"(-)", "E1007"},
	}
	executeTest(testCode, "minus", t)
}
func TestMulti(t *testing.T) {
	testCode := [][]string{
		{"(* 3 6)", "18"},
		{"(* 0.5 5.75)", "2.875"},
		{"(* 2 0.5 1.25)", "1.25"},

		{"(* 10.2 #f)", "E1003"},
		{"(* 10.2)", "E1007"},
		{"(*)", "E1007"},
	}
	executeTest(testCode, "multi", t)
}
func TestDiv(t *testing.T) {
	testCode := [][]string{
		{"(/ 4 3)", "1"},
		{"(/ 0.75 0.25)", "3"},
		{"(/ 9.5 5)", "1.9"},
		{"(/ 3 0.5 2)", "3"},

		{"(/ 10.2 #f)", "E1003"},
		{"(/ 10.2)", "E1007"},
		{"(/)", "E1007"},
		{"(/ 10 0)", "E1013"},
		{"(/ 10 2 0 3)", "E1013"},
	}
	executeTest(testCode, "div", t)
}
func TestModulo(t *testing.T) {
	testCode := [][]string{
		{"(modulo 18 12)", "6"},
		{"(modulo 11 (+ 1 2))", "2"},
		{"(modulo  3 5)", "3"},

		{"(modulo 1)", "E1007"},
		{"(modulo 10 3 2)", "E1007"},
		{"(modulo 10 2.5)", "E1002"},
		{"(modulo 10 0)", "E1013"},
	}
	executeTest(testCode, "modulo", t)
}
func TestQuotient(t *testing.T) {
	testCode := [][]string{
		{"(quotient 18 12)", "1"},
		{"(quotient 11 (+ 1 2))", "3"},
		{"(quotient 3 5)", "0"},

		{"(quotient 1)", "E1007"},
		{"(quotient 10 3 2)", "E1007"},
		{"(quotient 10 2.5)", "E1002"},
		{"(quotient 10 0)", "E1013"},
	}
	executeTest(testCode, "quotient", t)
}
func TestEq(t *testing.T) {
	testCode := [][]string{
		{"(= 5 5)", "#t"},
		{"(= 5.5 5.5)", "#t"},
		{"(= 5 5.0)", "#t"},
		{"(= 5.0 5)", "#t"},
		{"(= 5 6)", "#f"},
		{"(= 5.5 6.6)", "#f"},
		{"(= 5 6.6)", "#f"},
		{"(= 5.0 6)", "#f"},
		{"(= (+ 1 1)(+ 0 2))", "#t"},

		{"(= 10 3 2)", "E1007"},
		{"(= 10.2 #f)", "E1003"},
		{"(= 10.2)", "E1007"},
		{"(=)", "E1007"},
	}
	executeTest(testCode, "eq", t)
}
func TestThan(t *testing.T) {
	testCode := [][]string{
		{"(> 6 5)", "#t"},
		{"(> 6.5 5.5)", "#t"},
		{"(> 6.1 6)", "#t"},
		{"(> 6 5.9)", "#t"},
		{"(> 6 6)", "#f"},
		{"(> 4.5 5.5)", "#f"},
		{"(> 4 5.5)", "#f"},
		{"(> 4.5 5)", "#f"},
		{"(> (+ 3 3) 5)", "#t"},

		{"(> 10 3 2)", "E1007"},
		{"(> 10.2 #f)", "E1003"},
		{"(> 10.2)", "E1007"},
		{"(>)", "E1007"},
	}
	executeTest(testCode, "than", t)
}
func TestLess(t *testing.T) {
	testCode := [][]string{
		{"(< 5 6)", "#t"},
		{"(< 5.6 6.5)", "#t"},
		{"(< 5 6.1)", "#t"},
		{"(< 5 6.5)", "#t"},
		{"(> 6 6)", "#f"},
		{"(> 6.5 6.6)", "#f"},
		{"(> 6 6.0)", "#f"},
		{"(> 5.9 6)", "#f"},
		{"(< 5 (+ 3 3))", "#t"},

		{"(< 10.2 #f)", "E1003"},
		{"(< 10 3 2)", "E1007"},
		{"(< 10.2)", "E1007"},
		{"(<)", "E1007"},
	}
	executeTest(testCode, "less", t)
}
func TestThanEq(t *testing.T) {
	testCode := [][]string{
		{"(>= 6 6)", "#t"},
		{"(>= 6 5)", "#t"},
		{"(>= 6.1 5)", "#t"},
		{"(>= 7.6 7.6)", "#t"},
		{"(>= 6.3 5.2)", "#t"},
		{"(>= 6 5.1)", "#t"},
		{"(>= 5 6)", "#f"},
		{"(>= 5.1 6.2)", "#f"},
		{"(>= 5.9 6)", "#f"},
		{"(>= 5 6.1)", "#f"},
		{"(>= (+ 2 3 1) 6)", "#t"},

		{"(>= 10 3 2)", "E1007"},
		{"(>= 10.2 #f)", "E1003"},
		{"(>= 10.2)", "E1007"},
		{"(>=)", "E1007"},
	}
	executeTest(testCode, "than_eq", t)
}
func TestLessEq(t *testing.T) {
	testCode := [][]string{
		{"(<= 6 6)", "#t"},
		{"(<= 6 5)", "#f"},
		{"(<= 6.1 5)", "#f"},
		{"(<= 7.6 7.6)", "#t"},
		{"(<= 6.3 5.2)", "#f"},
		{"(<= 6 5.1)", "#f"},
		{"(<= 5 6)", "#t"},
		{"(<= 5.1 6.2)", "#t"},
		{"(<= 5.9 6)", "#t"},
		{"(<= 5 6.1)", "#t"},
		{"(<= (+ 3 3) 6)", "#t"},

		{"(<= 10.2 #f)", "E1003"},
		{"(<= 10 3 2)", "E1007"},
		{"(<= 10.2)", "E1007"},
		{"(<=)", "E1007"},
	}
	executeTest(testCode, "less_eq", t)
}
