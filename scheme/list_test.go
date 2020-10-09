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

func TestList(t *testing.T) {
	testCode := [][]string{
		{"(list 1 2 3)", "(1 2 3)"},

		{"((list 1 12) 10)", "E1006"},
	}
	executeTest(testCode, "list", t)
}
func TestNull(t *testing.T) {
	testCode := [][]string{
		{"(null? (list 1 2 3))", "#f"},
		{"(null? (list))", "#t"},
		{"(null? (null? 10))", "#f"},

		{"(null? (list 1)(list 2))", "E1007"},
		{"(null?)", "E1007"},
	}
	executeTest(testCode, "null?", t)
}
func TestLength(t *testing.T) {
	testCode := [][]string{
		{"(length (list 1 2 3 4))", "4"},

		{"(length 10)", "E1005"},
		{"(length (list 1)(list 2))", "E1007"},
		{"(length)", "E1007"},
	}
	executeTest(testCode, "length", t)
}
func TestCar(t *testing.T) {
	testCode := [][]string{
		{"(car (list 10 20 30 40))", "10"},

		{"(car 10)", "E1005"},
		{"(car (list 1)(list 2))", "E1007"},
		{"(car)", "E1007"},
		{"(car (list))", "E1011"},
	}
	executeTest(testCode, "car", t)
}

func TestCdr(t *testing.T) {
	testCode := [][]string{
		{"(cdr (cons 10 20))", "20"},

		{"(cdr (list 1)(list 2))", "E1007"},
		{"(cdr)", "E1007"},
		{"(cdr 10)", "E1005"},
		{"(cdr (list))", "E1011"},
	}
	executeTest(testCode, "cdr", t)
}

func TestCadr(t *testing.T) {
	testCode := [][]string{
		{"(cadr (list 1 2 3 4))", "2"},

		{"(cadr (list 1)(list 2))", "E1007"},
		{"(cadr)", "E1007"},
		{"(cadr (list 1))", "E1011"},
	}
	executeTest(testCode, "cadr", t)
}

func TestCons(t *testing.T) {
	testCode := [][]string{
		{"(car (cons 100 200))", "100"},
		{"(cdr (cons 100 200))", "200"},

		{"(cons 1 (list 1)(list 2))", "E1007"},
		{"(cons 1)", "E1007"},
	}
	executeTest(testCode, "cons", t)
}

func TestAppend(t *testing.T) {
	testCode := [][]string{
		{"(append (list 1 2)(list 3 4))", "(1 2 3 4)"},

		{"(append 10 10)", "E1005"},
		{"(append (list 1))", "E1007"},
	}
	executeTest(testCode, "append", t)
}

func TestLast(t *testing.T) {
	testCode := [][]string{
		{"(last (list 1 2 3))", "3"},

		{"(last 10)", "E1005"},
		{"(last (list 1)(list 2))", "E1007"},
		{"(last)", "E1007"},
		{"(last (list))", "E1011"},
	}
	executeTest(testCode, "last", t)
}

func TestReverse(t *testing.T) {
	testCode := [][]string{
		{"(reverse (list 1 2 3))", "(3 2 1)"},

		{"(reverse 10)", "E1005"},
		{"(reverse (list 1)(list 2))", "E1007"},
		{"(reverse)", "E1007"},
	}
	executeTest(testCode, "reverse", t)
}

func TestIota(t *testing.T) {
	testCode := [][]string{
		{"(iota 10)", "(0 1 2 3 4 5 6 7 8 9)"},
		{"(iota 5 2)", "(2 3 4 5 6)"},
		{"(iota -10 0 1)", "()"},
		{"(iota 1 10)", "(10)"},
		{"(iota 10 1 -1)", "(1 0 -1 -2 -3 -4 -5 -6 -7 -8)"},

		{"(iota 10.2 1 0)", "E1002"},
		{"(iota 1 10.2)", "E1002"},
		{"(iota 1 10 2.5)", "E1002"},
		{"(iota)", "E1007"},
		{"(iota 100 0 1 2)", "E1007"},
	}
	executeTest(testCode, "iota", t)
}

func TestMap(t *testing.T) {
	testCode := [][]string{
		{"(map (lambda (n) (* n 10))(list 1 2 3))", "(10 20 30)"},
		{"(map (lambda (n) (* n 10))(list))", "()"},

		{"(map (lambda (n) (* n 10)) 20)", "E1005"},
		{"(map (list 1 12) (list 10))", "E1006"},
		{"(map (lambda (n) (* n 10)))", "E1007"},
		{"(map (lambda (n) (* n 10))(list 1)(list 1))", "E1007"},
	}
	executeTest(testCode, "map", t)
}

func TestFilter(t *testing.T) {
	testCode := [][]string{
		{"(filter (lambda (n) (= n 1))(list 1 2 3))", "(1)"},
		{"(filter (lambda (n) (= n 1))(list))", "()"},

		{"(filter (lambda (n) (* n 10)) 20)", "E1005"},
		{"(filter (list 1 12) (list 10))", "E1006"},
		{"(filter (lambda (n) (* n 10)))", "E1007"},
		{"(filter (lambda (n) (* n 10))(list 1)(list 1))", "E1007"},
		{"(filter (lambda (n) 10.1) (list 1 2))", "E1001"},
	}
	executeTest(testCode, "filter", t)
}

func TestReduce(t *testing.T) {
	testCode := [][]string{
		{"(reduce (lambda (a b) (+ a b)) 0 (list 1 2 3))", "6"},
		{"(reduce (lambda (a b) (+ a b)) (* 10 10) (list))", "100"},
		{"(reduce (list 1 12) 0 (list 10))", "10"},

		{"(reduce (lambda (a b) (+ a b)) (+ 1 2) 20)", "E1005"},
		{"(reduce (list 1 12) 0 (list 1 2))", "E1006"},
		{"(reduce (lambda (a b) (+ a b)))", "E1007"},
		{"(reduce (lambda (a b) (+ a b)) (list 1 2))", "E1007"},
	}
	executeTest(testCode, "reduce", t)
}
func TestForEach(t *testing.T) {
	testCode := [][]string{
		{"(define cnt 0)", "cnt"},
		{"(for-each (lambda (n) (set! cnt (+ cnt n)))(list 1 1 1 1 1))", "nil"},
		{"cnt", "5"},

		{"(for-each (lambda (n) (* n 10)) 20)", "E1005"},
		{"(for-each (list 1 12) (list 10))", "E1006"},
		{"(for-each (lambda (n) (* n 10)))", "E1007"},
		{"(for-each (lambda (n) (* n 10))(list 1)(list 1))", "E1007"},
	}
	executeTest(testCode, "for-each", t)
}
