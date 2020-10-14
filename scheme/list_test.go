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
		{"(map list (list 1 2 3))", "((1) (2) (3))"},

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
		{"(filter null? (list () 10 20))", "(())"},

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
		{"(reduce + 0 (list 1 2 3))", "6"},

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
		{"(for-each (lambda (n) (set! cnt (+ cnt n)))(list 1 2 3 4))", "nil"},
		{"cnt", "10"},
		{"(for-each display (iota 10))", "nil"},

		{"(for-each (lambda (n) (* n 10)) 20)", "E1005"},
		{"(for-each (list 1 12) (list 10))", "E1006"},
		{"(for-each (lambda (n) (* n 10)))", "E1007"},
		{"(for-each (lambda (n) (* n 10))(list 1)(list 1))", "E1007"},
	}
	executeTest(testCode, "for-each", t)
}
func TestMakeList(t *testing.T) {
	testCode := [][]string{
		{"(make-list 10 0)", "(0 0 0 0 0 0 0 0 0 0)"},
		{"(make-list 4 (list 1 2 3))", "((1 2 3) (1 2 3) (1 2 3) (1 2 3))"},
		{"(make-list 8 (quote a))", "(a a a a a a a a)"},
		{"(make-list 8 #t)", "(#t #t #t #t #t #t #t #t)"},
		{"(make-list 10 1.0)", "(1 1 1 1 1 1 1 1 1 1)"},
		{"(make-list 0 1)", "()"},

		{"(make-list)", "E1007"},
		{"(make-list 10)", "E1007"},
		{"(make-list 10 0 1)", "E1007"},
		{"(make-list #t 0)", "E1002"},
		{"(make-list -1 0)", "E1011"},
		{"(make-list 10 c)", "E1008"},
	}
	executeTest(testCode, "for-each", t)
}
func TestTake(t *testing.T) {
	testCode := [][]string{
		{"(take (iota 10) 0)", "()"},
		{"(take (iota 10) 1)", "(0)"},
		{"(take (iota 10) 3)", "(0 1 2)"},
		{"(take (iota 10) 9)", "(0 1 2 3 4 5 6 7 8)"},
		{"(take (iota 10) 10)", "(0 1 2 3 4 5 6 7 8 9)"},

		{"(take)", "E1007"},
		{"(take (list 10 20))", "E1007"},
		{"(take (list 10 20) 1 2)", "E1007"},
		{"(take 1 (list 1 2))", "E1005"},
		{"(take (list 1 2) 10.5)", "E1002"},
		{"(take (list 1 2) 3)", "E1011"},
		{"(take (list 1 2) -1)", "E1011"},
		{"(take a 1)", "E1008"},
	}
	executeTest(testCode, "take", t)
}
func TestDrop(t *testing.T) {
	testCode := [][]string{
		{"(drop (iota 10) 0)", "(0 1 2 3 4 5 6 7 8 9)"},
		{"(drop (iota 10) 1)", "(1 2 3 4 5 6 7 8 9)"},
		{"(drop (iota 10) 3)", "(3 4 5 6 7 8 9)"},
		{"(drop (iota 10) 9)", "(9)"},
		{"(drop (iota 10) 10)", "()"},

		{"(drop)", "E1007"},
		{"(drop (list 10 20))", "E1007"},
		{"(drop (list 10 20) 1 2)", "E1007"},
		{"(drop 1 (list 1 2))", "E1005"},
		{"(drop (list 1 2) 10.5)", "E1002"},
		{"(drop (list 1 2) 3)", "E1011"},
		{"(drop (list 1 2) -1)", "E1011"},
		{"(drop a 1)", "E1008"},
	}
	executeTest(testCode, "drop", t)
}
func TestDelete(t *testing.T) {
	testCode := [][]string{
		{"(define a (list 10 10.5 \"ABC\" #\\a #t))", "a"},
		{"(delete 10 a)", "(10.5 \"ABC\" #\\a #t)"},
		{"(delete 10.5 a)", "(10 \"ABC\" #\\a #t)"},
		{"(delete \"ABC\" a)", "(10 10.5 #\\a #t)"},
		{"(delete #\\a a)", "(10 10.5 \"ABC\" #t)"},
		{"(delete #t a)", "(10 10.5 \"ABC\" #\\a)"},

		{"(delete)", "E1007"},
		{"(delete 10)", "E1007"},
		{"(delete 10 (list 10 20) 3)", "E1007"},
		{"(delete 10 20)", "E1005"},
		{"(delete 10 b)", "E1008"},
	}
	executeTest(testCode, "delete", t)
}
