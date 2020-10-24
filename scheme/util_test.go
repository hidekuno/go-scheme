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
		{"(eq? #\\a 10)", "#f"},
		{"(eq? \"abc\" \"abc\")", "#t"},
		{"(eq? \"123\" \"456\")", "#f"},
		{"(eq? \"123\" 9)", "#f"},

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
func TestEvenEq(t *testing.T) {
	testCode := [][]string{
		{"(even? 2)", "#t"},
		{"(even? 4)", "#t"},
		{"(even? 0)", "#t"},
		{"(even? 1)", "#f"},
		{"(even? 5)", "#f"},
		{"(even?)", "E1007"},
		{"(even? 1 2)", "E1007"},
		{"(even? 10.5)", "E1002"},
		{"(even? a)", "E1008"},
	}
	executeTest(testCode, "even?", t)
}

func TestOddEq(t *testing.T) {
	testCode := [][]string{
		{"(odd? 2)", "#f"},
		{"(odd? 4)", "#f"},
		{"(odd? 0)", "#f"},
		{"(odd? 1)", "#t"},
		{"(odd? 5)", "#t"},
		{"(odd?)", "E1007"},
		{"(odd? 1 2)", "E1007"},
		{"(odd? 10.5)", "E1002"},
		{"(odd? a)", "E1008"},
	}
	executeTest(testCode, "odd?", t)
}
func TestListTypeEq(t *testing.T) {
	testCode := [][]string{
		{"(list? (list 1 2 3))", "#t"},
		{"(list? 90)", "#f"},
		{"(list?)", "E1007"},
		{"(list? (list 1)(list 2))", "E1007"},
		{"(list? a)", "E1008"},
	}
	executeTest(testCode, "list?", t)
}
func TestPairTypeEq(t *testing.T) {
	testCode := [][]string{
		{"(pair? (cons 1 2))", "#t"},
		{"(pair? 110)", "#f"},
		{"(pair?)", "E1007"},
		{"(pair? (cons 1 2)(cons 3 4))", "E1007"},
		{"(pair? a)", "E1008"},
	}
	executeTest(testCode, "pair?", t)
}
func TestCharTypeEq(t *testing.T) {
	testCode := [][]string{
		{"(char? #\\a)", "#t"},
		{"(char? 100)", "#f"},
		{"(char?)", "E1007"},
		{"(char? #\\a #\\b)", "E1007"},
		{"(char? a)", "E1008"},
	}
	executeTest(testCode, "char?", t)
}

func TestStringTypeEq(t *testing.T) {
	testCode := [][]string{
		{"(string? \"a\")", "#t"},
		{"(string? 100)", "#f"},
		{"(string?)", "E1007"},
		{"(string? \"a\" \"b\")", "E1007"},
		{"(string? a)", "E1008"},
	}
	executeTest(testCode, "string?", t)
}
func TestIntegerTypeEq(t *testing.T) {
	testCode := [][]string{
		{"(integer? 10)", "#t"},
		{"(integer? \"a\")", "#f"},
		{"(integer?)", "E1007"},
		{"(integer? 10 20)", "E1007"},
		{"(integer? a)", "E1008"},
	}
	executeTest(testCode, "integer?", t)
}
func TestNumberTypeEq(t *testing.T) {
	testCode := [][]string{
		{"(number? 10)", "#t"},
		{"(number? 10.5)", "#t"},
		{"(number? \"a\")", "#f"},
		{"(number?)", "E1007"},
		{"(number? 10 20)", "E1007"},
		{"(number? a)", "E1008"},
	}
	executeTest(testCode, "number?", t)
}
func TestProcedureTypeEq(t *testing.T) {
	testCode := [][]string{
		{"(procedure? (lambda (n)n))", "#t"},
		{"(procedure? +)", "#t"},
		{"(procedure? 10)", "#f"},
		{"(procedure?)", "E1007"},
		{"(procedure? (lambda (n) n)(lambda (n) n))", "E1007"},
		{"(procedure? a)", "E1008"},
	}
	executeTest(testCode, "procedure?", t)
}
func TestZeroEq(t *testing.T) {
	testCode := [][]string{
		{"(zero? 0)", "#t"},
		{"(zero? 0.0)", "#t"},
		{"(zero? 2)", "#f"},
		{"(zero? -3)", "#f"},
		{"(zero? 2.5)", "#f"},
		{"(zero?)", "E1007"},
		{"(zero? 1 2)", "E1007"},
		{"(zero? #f)", "E1003"},
		{"(zero? a)", "E1008"},
	}
	executeTest(testCode, "zero?", t)
}
func TestPositiveEq(t *testing.T) {
	testCode := [][]string{
		{"(positive? 0)", "#f"},
		{"(positive? 0.0)", "#f"},
		{"(positive? 2)", "#t"},
		{"(positive? -3)", "#f"},
		{"(positive? 2.5)", "#t"},
		{"(positive? -1.5)", "#f"},
		{"(positive?)", "E1007"},
		{"(positive? 1 2)", "E1007"},
		{"(positive? #f)", "E1003"},
		{"(positive? a)", "E1008"},
	}
	executeTest(testCode, "positive?", t)
}
func TestNegativeEq(t *testing.T) {
	testCode := [][]string{
		{"(negative? 0)", "#f"},
		{"(negative? 0.0)", "#f"},
		{"(negative? 2)", "#f"},
		{"(negative? -3)", "#t"},
		{"(negative? 2.5)", "#f"},
		{"(negative? -1.5)", "#t"},

		{"(negative?)", "E1007"},
		{"(negative? 1 2)", "E1007"},
		{"(negative? #f)", "E1003"},
		{"(negative? a)", "E1008"},
	}
	executeTest(testCode, "negative?", t)
}
