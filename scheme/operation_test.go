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
func TestAsh(t *testing.T) {
	testCode := [][]string{
		{"(ash 10 1)", "20"},
		{"(ash 10 -1)", "5"},
		{"(ash 10 0)", "10"},
		{"(ash)", "E1007"},
		{"(ash 10)", "E1007"},
		{"(ash 10 1 1)", "E1007"},
		{"(ash a 1)", "E1008"},
		{"(ash 10 a)", "E1008"},
		{"(ash 10.5 1)", "E1002"},
		{"(ash 10 1.5)", "E1002"},
	}

	executeTest(testCode, "ash", t)
}
func TestLogand(t *testing.T) {
	testCode := [][]string{
		{"(logand 10 2)", "2"},
		{"(logand 10 2 3)", "2"},
		{"(logand 10)", "10"},

		{"(logand)", "E1007"},
		{"(logand a 1)", "E1008"},
		{"(logand 10 a)", "E1008"},
		{"(logand 10.5 1)", "E1002"},
		{"(logand 10 1.5)", "E1002"},
	}

	executeTest(testCode, "logand", t)
}
func TestLogior(t *testing.T) {
	testCode := [][]string{
		{"(logior 10 2)", "10"},
		{"(logior 10 2 3)", "11"},
		{"(logior 10)", "10"},

		{"(logior)", "E1007"},
		{"(logior a 1)", "E1008"},
		{"(logior 10 a)", "E1008"},
		{"(logior 10.5 1)", "E1002"},
		{"(logior 10 1.5)", "E1002"},
	}
	executeTest(testCode, "logior", t)
}
func TestLogxor(t *testing.T) {
	testCode := [][]string{
		{"(logxor 10 2)", "8"},
		{"(logxor 10 2 2)", "10"},
		{"(logxor 10)", "10"},

		{"(logxor)", "E1007"},
		{"(logxor a 1)", "E1008"},
		{"(logxor 10 a)", "E1008"},
		{"(logxor 10.5 1)", "E1002"},
		{"(logxor 10 1.5)", "E1002"},
	}
	executeTest(testCode, "logxor", t)
}

func TestMax(t *testing.T) {
	testCode := [][]string{
		{"(max 10 12 11 1 2)", "12"},
		{"(max 10 12 11 1 12)", "12"},
		{"(max 10 12 13.5 1 1)", "13.5"},

		{"(max 10)", "E1007"},
		{"(max 9 a)", "E1008"},
		{"(max 1 3.4 #t)", "E1003"},
		{"(max 1)", "E1007"},
	}
	executeTest(testCode, "max", t)
}
func TestMin(t *testing.T) {
	testCode := [][]string{
		{"(min 10 12 11 3 9)", "3"},
		{"(min 3 12 11 3 12)", "3"},
		{"(min 10 12 0.5 1 1)", "0.5"},

		{"(min 10)", "E1007"},
		{"(min 9 a)", "E1008"},
		{"(min 1 3.4 #t)", "E1003"},
		{"(min 1)", "E1007"},
	}
	executeTest(testCode, "min", t)
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
