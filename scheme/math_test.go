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

func TestSqrt(t *testing.T) {
	testCode := [][]string{
		{"(sqrt 9)", "3"},

		{"(sqrt #t)", "E1003"},
		{"(sqrt 11 10)", "E1007"},
		{"(sqrt)", "E1007"},
	}
	executeTest(testCode, "sqrt", t)
}
func TestSin(t *testing.T) {

	testCode := [][]string{
		{"(sin (/ (* 30 (* (atan 1) 4)) 180))", "0.49999999999999994"},

		{"(sin #t)", "E1003"},
		{"(sin 11 10)", "E1007"},
		{"(sin)", "E1007"},
	}
	executeTest(testCode, "sin", t)
}

func TestCos(t *testing.T) {

	testCode := [][]string{
		{"(cos (/ (* 60 (* (atan 1) 4))180))", "0.5000000000000001"},

		{"(cos #t)", "E1003"},
		{"(cos 11 10)", "E1007"},
		{"(cos)", "E1007"},
	}
	executeTest(testCode, "cos", t)
}

func TestTan(t *testing.T) {

	testCode := [][]string{
		{"(tan (/ (* 45 (* (atan 1) 4)) 180))", "1"},

		{"(tan #t)", "E1003"},
		{"(tan 11 10)", "E1007"},
		{"(tan)", "E1007"},
	}
	executeTest(testCode, "tan", t)
}

func TestAtan(t *testing.T) {
	testCode := [][]string{
		{"(* 4 (atan 1))", "3.141592653589793"},

		{"(atan #t)", "E1003"},
		{"(atan 11 10)", "E1007"},
		{"(atan)", "E1007"},
	}
	executeTest(testCode, "atan", t)
}

func TestLog(t *testing.T) {
	testCode := [][]string{
		{"(/ (log 8)(log 2))", "3"},

		{"(log #t)", "E1003"},
		{"(log 11 10)", "E1007"},
		{"(log)", "E1007"},
	}
	executeTest(testCode, "log", t)
}

func TestExp(t *testing.T) {
	testCode := [][]string{
		{"(exp (/ (log 8) 3))", "2"},

		{"(exp #t)", "E1003"},
		{"(exp 11 10)", "E1007"},
		{"(exp)", "E1007"},
	}
	executeTest(testCode, "exp", t)
}

func TestRandInteger(t *testing.T) {
	testCode := [][]string{
		{"(rand-init 10.2)", "E1007"},
		{"(rand-integer 10.2)", "E1002"},
		{"(rand-integer)", "E1007"},
		{"(rand-integer 11 9)", "E1007"},
	}
	executeTest(testCode, "rand-integer", t)
}

func TestExpt(t *testing.T) {

	testCode := [][]string{
		{"(expt 2 0)", "1"},
		{"(expt 2 1)", "2"},
		{"(expt 2 (+ 1 2))", "8"},
		{"(define a 4)", "a"},
		{"(expt 2 a)", "16"},
		{"(expt 2.0 3.0)", "8"},
		{"(expt 2.0 3)", "8"},
		{"(expt 2 3.0)", "8"},

		{"(expt 10)", "E1007"},
		{"(expt 10 10 10)", "E1007"},
		{"(expt 11.5 #f)", "E1003"},
		{"(expt #t 12.5)", "E1003"},
	}
	executeTest(testCode, "expt", t)
}
func TestAsin(t *testing.T) {
	testCode := [][]string{
		{"(define pi (* 4 (atan 1)))", "pi"},
		{"(/ (asin 0.5)(/ pi 180))", "30.000000000000004"},

		{"(asin #t)", "E1003"},
		{"(asin 11 10)", "E1007"},
		{"(asin)", "E1007"},
	}
	executeTest(testCode, "asin", t)
}
func TestAcos(t *testing.T) {
	testCode := [][]string{
		{"(define pi (* 4 (atan 1)))", "pi"},
		{"(/ (acos 0.5)(/ pi 180))", "59.99999999999999"},

		{"(acos #t)", "E1003"},
		{"(acos 11 10)", "E1007"},
		{"(acos)", "E1007"},
	}
	executeTest(testCode, "acos", t)
}
func TestAbs(t *testing.T) {
	testCode := [][]string{
		{"(abs 10)", "10"},
		{"(abs -10)", "10"},
		{"(abs 10.5)", "10.5"},
		{"(abs -10.5)", "10.5"},

		{"(abs #t)", "E1003"},
		{"(abs 11 10)", "E1007"},
		{"(abs)", "E1007"},
	}
	executeTest(testCode, "abs", t)
}
