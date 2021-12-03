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
func TestAppendEffect(t *testing.T) {
	testCode := [][]string{
		{"(append! (list 1)(list 2))", "(1 2)"},
		{"(append! (list 1)(list 2)(list 3))", "(1 2 3)"},
		{"(append! (list (list 10))(list 2)(list 3))", "((10) 2 3)"},
		{"(append! (iota 5) (list 100))", "(0 1 2 3 4 100)"},
		{"(define a (iota 5))", "a"},
		{"(define b a)", "b"},
		{"(append! a (iota 5 5))", "(0 1 2 3 4 5 6 7 8 9)"},
		{"a", "(0 1 2 3 4 5 6 7 8 9)"},
		{"b", "(0 1 2 3 4 5 6 7 8 9)"},

		{"(append!)", "E1007"},
		{"(append! 10)", "E1005"},
		{"(append! (list 1) 105)", "E1005"},
		{"(append! (list 1) c)", "E1008"},
	}
	executeTest(testCode, "append!", t)
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
		{"(map list (list (list 1 2)(list 3 4)))", "(((1 2)) ((3 4)))"},
		{"(map list (list (cons 1 2)(cons 3 4)))", "(((1 . 2)) ((3 . 4)))"},

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
		{"(filter list? (list (list 1 2)(list 3 4)))", "((1 2) (3 4))"},
		{"(filter list? (list (cons 1 2)(cons 3 4)))", "()"},

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
		{"(for-each (lambda (n) (set! cnt (+ cnt n)))(list 1 2 3 4))", "#<nil>"},
		{"cnt", "10"},
		{"(for-each display (iota 10))", "#<nil>"},
		{"(for-each list? (list (list 1 2)(list 3 4)))", "#<nil>"},
		{"(for-each list? (list (cons 1 2)(cons 3 4)))", "#<nil>"},

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
func TestDeleteEffect(t *testing.T) {
	testCode := [][]string{

		{"(define a (list 10 10.5 3/5 \"ABC\" #\\a #t))", "a"},
		{"(delete! 10 a)", "(10.5 3/5 \"ABC\" #\\a #t)"},
		{"a", "(10.5 3/5 \"ABC\" #\\a #t)"},
		{"(delete! 10.5 a)", "(3/5 \"ABC\" #\\a #t)"},
		{"(delete! 3/5 a)", "(\"ABC\" #\\a #t)"},
		{"(delete! \"ABC\" a)", "(#\\a #t)"},
		{"a", "(#\\a #t)"},
		{"(delete! #\\a a)", "(#t)"},
		{"a", "(#t)"},
		{"(delete! #f a)", "(#t)"},
		{"a", "(#t)"},
		{"(delete! #t a)", "()"},
		{"a", "()"},

		{"(delete!)", "E1007"},
		{"(delete! 10)", "E1007"},
		{"(delete! 10 (list 10 20) 3)", "E1007"},
		{"(delete! 10 20)", "E1005"},
		{"(delete! 10 c)", "E1008"},
	}
	executeTest(testCode, "delete!", t)
}
func TestListRef(t *testing.T) {
	testCode := [][]string{
		{"(list-ref (iota 10) 0)", "0"},
		{"(list-ref (iota 10) 1)", "1"},
		{"(list-ref (iota 10) 8)", "8"},
		{"(list-ref (iota 10) 9)", "9"},
		{"(list-ref (list #\\a #\\b #\\c) 1)", "#\\b"},
		{"(list-ref (list (list 0 1) 1 2 3) 0)", "(0 1)"},

		{"(list-ref)", "E1007"},
		{"(list-ref (iota 10))", "E1007"},
		{"(list-ref (iota 10) 1 2)", "E1007"},
		{"(list-ref 10 -1)", "E1005"},
		{"(list-ref (iota 10) #t)", "E1002"},
		{"(list-ref a #t)", "E1008"},
		{"(list-ref (iota 10) a)", "E1008"},
		{"(list-ref (iota 10) -1)", "E1011"},
		{"(list-ref (iota 10) 10)", "E1011"},
	}
	executeTest(testCode, "list-ref", t)
}
func TestListSet(t *testing.T) {
	testCode := [][]string{
		{"(define a (list 1 2 3 4 5))", "a"},
		{"(define b a)", "b"},
		{"(list-set! a 0 100)", "#<nil>"},
		{"a", "(100 2 3 4 5)"},
		{"b", "(100 2 3 4 5)"},

		{"(list-set!)", "E1007"},
		{"(list-set! (iota 10))", "E1007"},
		{"(list-set! (iota 10) 1 2 3)", "E1007"},
		{"(list-set! 10 0 -1)", "E1005"},
		{"(list-set! (iota 10) #t 0)", "E1002"},
		{"(list-set! c 0 #t)", "E1008"},
		{"(list-set! (iota 10) 0 d)", "E1008"},
		{"(list-set! (iota 10) -1 0)", "E1011"},
		{"(list-set! (iota 10) 10 0)", "E1011"},
	}
	executeTest(testCode, "list-ref", t)
}
func TestSetCar(t *testing.T) {
	testCode := [][]string{
		{"(define a (list 1 2 3 4 5))", "a"},
		{"(set-car! a 100)", "#<nil>"},
		{"a", "(100 2 3 4 5)"},
		{"(set-car! a (list 10 20))", "#<nil>"},
		{"a", "((10 20) 2 3 4 5)"},

		{"(set-car!)", "E1007"},
		{"(set-car! (list 1))", "E1007"},
		{"(set-car! c a)", "E1008"},
		{"(set-car! 10 20)", "E1005"},
		{"(set-car! () 20)", "E1011"},
	}
	executeTest(testCode, "set-car!", t)
}
func TestSetCdr(t *testing.T) {
	testCode := [][]string{
		{"(define a (list 1 2 3 4 5))", "a"},
		{"(set-cdr! a 100)", "#<nil>"},
		{"a", "(1 100)"},
		{"(set-cdr! a (list 10 20))", "#<nil>"},
		{"a", "(1 10 20)"},

		{"(set-cdr!)", "E1007"},
		{"(set-cdr! (list 1))", "E1007"},
		{"(set-cdr! c a)", "E1008"},
		{"(set-cdr! 100 200)", "E1005"},
		{"(set-cdr! () 20)", "E1011"},
	}
	executeTest(testCode, "set-cdr!", t)
}
func TestSort(t *testing.T) {
	testCode := [][]string{
		{"(sort (list 10 1 9 5 3 4 7 6 5))", "(1 3 4 5 5 6 7 9 10)"},
		{"(sort (list 10 1.5 9 5 2/3 4 7 6 5))", "(2/3 1.5 4 5 5 6 7 9 10)"},
		{"(sort (list \"z\" \"a\" \"b\" \"m\" \"l\" \"d\" \"A\" \"c\" \"0\"))",
			"(\"0\" \"A\" \"a\" \"b\" \"c\" \"d\" \"l\" \"m\" \"z\")"},
		{"(sort (list \"AZ\" \"AA\" \"AB\" \"CA\" \"BB\"))", "(\"AA\" \"AB\" \"AZ\" \"BB\" \"CA\")"},
		{"(sort (list #\\z #\\a #\\b #\\m #\\l #\\d #\\A #\\c #\\0))",
			"(#\\0 #\\A #\\a #\\b #\\c #\\d #\\l #\\m #\\z)"},
		{"(sort (list 10 1 9 5 3 4 7 6 5) <)", "(1 3 4 5 5 6 7 9 10)"},
		{"(sort (list 10 1 9 5 3 4 7 6 5) >)", "(10 9 7 6 5 5 4 3 1)"},
		{"(sort (list 10 1 9 5 3 4 7 6 5) <=)", "(1 3 4 5 5 6 7 9 10)"},
		{"(sort (list 10 1 9 5 3 4 7 6 5) >=)", "(10 9 7 6 5 5 4 3 1)"},
		{"(sort (list 10 1.5 9 5 2/3 4 7 6 5) >)", "(10 9 7 6 5 5 4 1.5 2/3)"},
		{"(sort (list \"z\" \"a\" \"b\" \"m\" \"l\" \"d\" \"A\" \"c\" \"0\") string>?)",
			"(\"z\" \"m\" \"l\" \"d\" \"c\" \"b\" \"a\" \"A\" \"0\")"},
		{"(sort (list \"AZ\" \"AA\" \"AB\" \"CA\" \"BB\") string>?)", "(\"CA\" \"BB\" \"AZ\" \"AB\" \"AA\")"},
		{"(sort (list \"AZ\" \"AA\" \"AB\" \"CA\" \"BB\") string>=?)", "(\"CA\" \"BB\" \"AZ\" \"AB\" \"AA\")"},
		{"(sort (list \"AZ\" \"AA\" \"AB\" \"CA\" \"BB\") string<?)", "(\"AA\" \"AB\" \"AZ\" \"BB\" \"CA\")"},
		{"(sort (list \"AZ\" \"AA\" \"AB\" \"CA\" \"BB\") string<=?)", "(\"AA\" \"AB\" \"AZ\" \"BB\" \"CA\")"},
		{"(sort (list #\\z #\\a #\\b #\\m #\\l #\\d #\\A #\\c #\\0) char>?)",
			"(#\\z #\\m #\\l #\\d #\\c #\\b #\\a #\\A #\\0)"},
		{"(sort (list #\\z #\\a #\\b #\\m #\\l #\\d #\\A #\\c #\\0) char>=?)",
			"(#\\z #\\m #\\l #\\d #\\c #\\b #\\a #\\A #\\0)"},
		{"(sort (list #\\z #\\a #\\b #\\m #\\l #\\d #\\A #\\c #\\0) char<?)",
			"(#\\0 #\\A #\\a #\\b #\\c #\\d #\\l #\\m #\\z)"},
		{"(sort (list #\\z #\\a #\\b #\\m #\\l #\\d #\\A #\\c #\\0) char<=?)",
			"(#\\0 #\\A #\\a #\\b #\\c #\\d #\\l #\\m #\\z)"},
		{"(sort (list 10 1.5 9 5 2/3 4 7 6 5) (lambda (a b) (> a b)))", "(10 9 7 6 5 5 4 1.5 2/3)"},

		{"(sort)", "E1007"},
		{"(sort (list 1) + +)", "E1007"},
		{"(sort +)", "E1005"},
		{"(sort (list 1) 10)", "E1006"},
	}
	executeTest(testCode, "sort", t)
}
func TestSortEffect(t *testing.T) {
	testCode := [][]string{
		{"(define a (list 10 1 9 5 3 4 7 6 5))", "a"},
		{"(sort! a)", "(1 3 4 5 5 6 7 9 10)"},
		{"a", "(1 3 4 5 5 6 7 9 10)"},

		{"(define a (list 10 1.5 9 5 2/3 4 7 6 5))", "a"},
		{"(sort! a)", "(2/3 1.5 4 5 5 6 7 9 10)"},
		{"a", "(2/3 1.5 4 5 5 6 7 9 10)"},

		{"(define a (list \"z\" \"a\" \"b\" \"m\" \"l\" \"d\" \"A\" \"c\" \"0\"))", "a"},
		{"(sort! a)", "(\"0\" \"A\" \"a\" \"b\" \"c\" \"d\" \"l\" \"m\" \"z\")"},
		{"a", "(\"0\" \"A\" \"a\" \"b\" \"c\" \"d\" \"l\" \"m\" \"z\")"},

		{"(define a (list \"AZ\" \"AA\" \"AB\" \"CA\" \"BB\"))", "a"},
		{"(sort! a)", "(\"AA\" \"AB\" \"AZ\" \"BB\" \"CA\")"},
		{"a", "(\"AA\" \"AB\" \"AZ\" \"BB\" \"CA\")"},
		{"(define a (list #\\z #\\a #\\b #\\m #\\l #\\d #\\A #\\c #\\0))", "a"},
		{"(sort! a)", "(#\\0 #\\A #\\a #\\b #\\c #\\d #\\l #\\m #\\z)"},
		{"a", "(#\\0 #\\A #\\a #\\b #\\c #\\d #\\l #\\m #\\z)"},

		{"(define a (list 10 1 9 5 3 4 7 6 5))", "a"},
		{"(sort! a >)", "(10 9 7 6 5 5 4 3 1)"},
		{"a", "(10 9 7 6 5 5 4 3 1)"},

		{"(define a (list 10 1.5 9 5 2/3 4 7 6 5))", "a"},
		{"(sort! a >)", "(10 9 7 6 5 5 4 1.5 2/3)"},
		{"a", "(10 9 7 6 5 5 4 1.5 2/3)"},

		{"(define a (list \"z\" \"a\" \"b\" \"m\" \"l\" \"d\" \"A\" \"c\" \"0\") )", "a"},
		{"(sort! a string>?)", "(\"z\" \"m\" \"l\" \"d\" \"c\" \"b\" \"a\" \"A\" \"0\")"},
		{"a", "(\"z\" \"m\" \"l\" \"d\" \"c\" \"b\" \"a\" \"A\" \"0\")"},

		{"(define a (list \"AZ\" \"AA\" \"AB\" \"CA\" \"BB\"))", "a"},
		{"(sort! a string>?)", "(\"CA\" \"BB\" \"AZ\" \"AB\" \"AA\")"},
		{"a", "(\"CA\" \"BB\" \"AZ\" \"AB\" \"AA\")"},

		{"(define a (list #\\z #\\a #\\b #\\m #\\l #\\d #\\A #\\c #\\0) )", "a"},
		{"(sort! a char>?)", "(#\\z #\\m #\\l #\\d #\\c #\\b #\\a #\\A #\\0)"},
		{"a", "(#\\z #\\m #\\l #\\d #\\c #\\b #\\a #\\A #\\0)"},

		{"(define a (list 10 1.5 9 5 2/3 4 7 6 5))", "a"},
		{"(sort! a (lambda (a b)(> a b)))", "(10 9 7 6 5 5 4 1.5 2/3)"},
		{"a", "(10 9 7 6 5 5 4 1.5 2/3)"},

		{"(sort!)", "E1007"},
		{"(sort! (list 1) + +)", "E1007"},
		{"(sort! +)", "E1005"},
		{"(sort! (list 1) 10)", "E1006"},
	}
	executeTest(testCode, "sort!", t)
}
func TestStableSort(t *testing.T) {
	testCode := [][]string{
		{"(stable-sort (list 10 1 9 5 3 4 7 6 5))", "(1 3 4 5 5 6 7 9 10)"},
		{"(stable-sort (list 10 1.5 9 5 2/3 4 7 6 5))", "(2/3 1.5 4 5 5 6 7 9 10)"},
		{"(stable-sort (list \"z\" \"a\" \"b\" \"m\" \"l\" \"d\" \"A\" \"c\" \"0\"))",
			"(\"0\" \"A\" \"a\" \"b\" \"c\" \"d\" \"l\" \"m\" \"z\")"},
		{"(stable-sort (list \"AZ\" \"AA\" \"AB\" \"CA\" \"BB\"))", "(\"AA\" \"AB\" \"AZ\" \"BB\" \"CA\")"},
		{"(stable-sort (list #\\z #\\a #\\b #\\m #\\l #\\d #\\A #\\c #\\0))",
			"(#\\0 #\\A #\\a #\\b #\\c #\\d #\\l #\\m #\\z)"},
		{"(stable-sort (list 10 1 9 5 3 4 7 6 5) <)", "(1 3 4 5 5 6 7 9 10)"},
		{"(stable-sort (list 10 1 9 5 3 4 7 6 5) >)", "(10 9 7 6 5 5 4 3 1)"},
		{"(stable-sort (list 10 1 9 5 3 4 7 6 5) <=)", "(1 3 4 5 5 6 7 9 10)"},
		{"(stable-sort (list 10 1 9 5 3 4 7 6 5) >=)", "(10 9 7 6 5 5 4 3 1)"},
		{"(stable-sort (list 10 1.5 9 5 2/3 4 7 6 5) >)", "(10 9 7 6 5 5 4 1.5 2/3)"},
		{"(stable-sort (list \"z\" \"a\" \"b\" \"m\" \"l\" \"d\" \"A\" \"c\" \"0\") string>?)",
			"(\"z\" \"m\" \"l\" \"d\" \"c\" \"b\" \"a\" \"A\" \"0\")"},
		{"(stable-sort (list \"AZ\" \"AA\" \"AB\" \"CA\" \"BB\") string>?)",
			"(\"CA\" \"BB\" \"AZ\" \"AB\" \"AA\")"},
		{"(stable-sort (list \"AZ\" \"AA\" \"AB\" \"CA\" \"BB\") string>=?)",
			"(\"CA\" \"BB\" \"AZ\" \"AB\" \"AA\")"},
		{"(stable-sort (list \"AZ\" \"AA\" \"AB\" \"CA\" \"BB\") string<?)",
			"(\"AA\" \"AB\" \"AZ\" \"BB\" \"CA\")"},
		{"(stable-sort (list \"AZ\" \"AA\" \"AB\" \"CA\" \"BB\") string<=?)",
			"(\"AA\" \"AB\" \"AZ\" \"BB\" \"CA\")"},
		{"(stable-sort (list #\\z #\\a #\\b #\\m #\\l #\\d #\\A #\\c #\\0) char>?)",
			"(#\\z #\\m #\\l #\\d #\\c #\\b #\\a #\\A #\\0)"},
		{"(stable-sort (list #\\z #\\a #\\b #\\m #\\l #\\d #\\A #\\c #\\0) char>=?)",
			"(#\\z #\\m #\\l #\\d #\\c #\\b #\\a #\\A #\\0)"},
		{"(stable-sort (list #\\z #\\a #\\b #\\m #\\l #\\d #\\A #\\c #\\0) char<?)",
			"(#\\0 #\\A #\\a #\\b #\\c #\\d #\\l #\\m #\\z)"},
		{"(stable-sort (list #\\z #\\a #\\b #\\m #\\l #\\d #\\A #\\c #\\0) char<=?)",
			"(#\\0 #\\A #\\a #\\b #\\c #\\d #\\l #\\m #\\z)"},
		{"(stable-sort (list 10 1.5 9 5 2/3 4 7 6 5) (lambda (a b) (> a b)))", "(10 9 7 6 5 5 4 1.5 2/3)"},

		{"(stable-sort)", "E1007"},
		{"(stable-sort (list 1) + +)", "E1007"},
		{"(stable-sort +)", "E1005"},
		{"(stable-sort (list 1) 10)", "E1006"},
	}
	executeTest(testCode, "stable-sort", t)
}
func TestStableSortEffect(t *testing.T) {
	testCode := [][]string{
		{"(define a (list 10 1 9 5 3 4 7 6 5))", "a"},
		{"(stable-sort! a)", "(1 3 4 5 5 6 7 9 10)"},
		{"a", "(1 3 4 5 5 6 7 9 10)"},

		{"(define a (list 10 1.5 9 5 2/3 4 7 6 5))", "a"},
		{"(stable-sort! a)", "(2/3 1.5 4 5 5 6 7 9 10)"},
		{"a", "(2/3 1.5 4 5 5 6 7 9 10)"},

		{"(define a (list \"z\" \"a\" \"b\" \"m\" \"l\" \"d\" \"A\" \"c\" \"0\"))", "a"},
		{"(stable-sort! a)", "(\"0\" \"A\" \"a\" \"b\" \"c\" \"d\" \"l\" \"m\" \"z\")"},
		{"a", "(\"0\" \"A\" \"a\" \"b\" \"c\" \"d\" \"l\" \"m\" \"z\")"},

		{"(define a (list \"AZ\" \"AA\" \"AB\" \"CA\" \"BB\"))", "a"},
		{"(stable-sort! a)", "(\"AA\" \"AB\" \"AZ\" \"BB\" \"CA\")"},
		{"a", "(\"AA\" \"AB\" \"AZ\" \"BB\" \"CA\")"},
		{"(define a (list #\\z #\\a #\\b #\\m #\\l #\\d #\\A #\\c #\\0))", "a"},
		{"(stable-sort! a)", "(#\\0 #\\A #\\a #\\b #\\c #\\d #\\l #\\m #\\z)"},
		{"a", "(#\\0 #\\A #\\a #\\b #\\c #\\d #\\l #\\m #\\z)"},

		{"(define a (list 10 1 9 5 3 4 7 6 5))", "a"},
		{"(stable-sort! a >)", "(10 9 7 6 5 5 4 3 1)"},
		{"a", "(10 9 7 6 5 5 4 3 1)"},

		{"(define a (list 10 1.5 9 5 2/3 4 7 6 5))", "a"},
		{"(stable-sort! a >)", "(10 9 7 6 5 5 4 1.5 2/3)"},
		{"a", "(10 9 7 6 5 5 4 1.5 2/3)"},

		{"(define a (list \"z\" \"a\" \"b\" \"m\" \"l\" \"d\" \"A\" \"c\" \"0\") )", "a"},
		{"(stable-sort! a string>?)", "(\"z\" \"m\" \"l\" \"d\" \"c\" \"b\" \"a\" \"A\" \"0\")"},
		{"a", "(\"z\" \"m\" \"l\" \"d\" \"c\" \"b\" \"a\" \"A\" \"0\")"},

		{"(define a (list \"AZ\" \"AA\" \"AB\" \"CA\" \"BB\"))", "a"},
		{"(stable-sort! a string>?)", "(\"CA\" \"BB\" \"AZ\" \"AB\" \"AA\")"},
		{"a", "(\"CA\" \"BB\" \"AZ\" \"AB\" \"AA\")"},

		{"(define a (list #\\z #\\a #\\b #\\m #\\l #\\d #\\A #\\c #\\0) )", "a"},
		{"(stable-sort! a char>?)", "(#\\z #\\m #\\l #\\d #\\c #\\b #\\a #\\A #\\0)"},
		{"a", "(#\\z #\\m #\\l #\\d #\\c #\\b #\\a #\\A #\\0)"},

		{"(define a (list 10 1.5 9 5 2/3 4 7 6 5))", "a"},
		{"(stable-sort! a (lambda (a b)(> a b)))", "(10 9 7 6 5 5 4 1.5 2/3)"},
		{"a", "(10 9 7 6 5 5 4 1.5 2/3)"},

		{"(stable-sort!)", "E1007"},
		{"(stable-sort! (list 1) + +)", "E1007"},
		{"(stable-sort! +)", "E1005"},
		{"(stable-sort! (list 1) 10)", "E1006"},
	}
	executeTest(testCode, "stable-sort!", t)
}
func TestMerge(t *testing.T) {
	testCode := [][]string{

		{"(merge (iota 10 1 2) (iota 10 2 2))", "(1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20)"},
		{"(merge (list 1/3 1 2) (list 1/2 3/2 1.75))", "(1/3 1/2 1 3/2 1.75 2)"},
		{"(merge (list \"a\" \"c\" \"e\" \"g\")(list \"b\" \"d\" \"f\" \"h\")  string<=?)",
			"(\"a\" \"b\" \"c\" \"d\" \"e\" \"f\" \"g\" \"h\")"},
		{"(merge (list #\\a #\\c #\\e #\\g)(list #\\b #\\d #\\f #\\h) char<=?)",
			"(#\\a #\\b #\\c #\\d #\\e #\\f #\\g #\\h)"},
		{"(merge (reverse (iota 10 1 2)) (reverse (iota 10 2 2)) >)",
			"(20 19 18 17 16 15 14 13 12 11 10 9 8 7 6 5 4 3 2 1)"},
		{"(merge (list \"g\" \"e\" \"c\" \"a\")(list \"h\" \"f\" \"d\" \"b\") string>?)",
			"(\"h\" \"g\" \"f\" \"e\" \"d\" \"c\" \"b\" \"a\")"},
		{"(merge (list #\\g #\\e #\\c #\\a)(list #\\h #\\f #\\d #\\b) char>?)",
			"(#\\h #\\g #\\f #\\e #\\d #\\c #\\b #\\a)"},
		{"(merge)", "E1007"},
		{"(merge (list 1))", "E1007"},
		{"(merge (list 1)(list 2) + +)", "E1007"},
		{"(merge + (list 2) +)", "E1005"},
		{"(merge (list 1) + +)", "E1005"},
		{"(merge (list 1)(list 2)(list 3))", "E1006"},
		{"(merge (list 1)(list 2) +)", "E1001"},
	}
	executeTest(testCode, "merge", t)
}
func TestIsSorted(t *testing.T) {
	testCode := [][]string{
		{"(sorted? (list 1 2 3))", "#t"},
		{"(sorted? (list 1 2 3) <)", "#t"},
		{"(sorted? (list 3 2 1) >)", "#t"},
		{"(sorted? (list 1 2 3) >)", "#f"},
		{"(sorted? (list 3 2 1) <)", "#f"},
		{"(sorted? (list 1 2 3) (lambda (a b)(< a b)))", "#t"},
		{"(sorted? (list 1 2 3) (lambda (a b)(> a b)))", "#f"},
		{"(sorted? (list \"a\" \"b\" \"c\") string<?)", "#t"},
		{"(sorted? (list \"c\" \"b\" \"a\") string>?)", "#t"},
		{"(sorted? (list #\\a #\\b #\\c) char<?)", "#t"},
		{"(sorted? (list #\\c #\\b #\\a) char>?)", "#t"},
		{"(sorted? (list #\\a #\\b #\\c) char-ci<?)", "#t"},
		{"(sorted? (list #\\c #\\b #\\a) char-ci>?)", "#t"},

		{"(sorted?)", "E1007"},
		{"(sorted? (list 1) + +)", "E1007"},
		{"(sorted? +)", "E1005"},
		{"(sorted? (list 1) 10)", "E1006"},
	}
	executeTest(testCode, "sorted?", t)
}
func TestVector(t *testing.T) {
	testCode := [][]string{
		{"(vector 1 2)", "#(1 2)"},
		{"(vector 0.5 1)", "#(0.5 1)"},
		{"(vector #t #f)", "#(#t #f)"},
		{"(vector (list 1)(list 2))", "#((1) (2))"},
		{"(vector (vector (vector 1))(vector 2)(vector 3))",
			"#(#(#(1)) #(2) #(3))"},

		{"(vector c 10)", "E1008"},
		{"#", "E1008"},
	}
	executeTest(testCode, "vector", t)
}
func TestMakeVector(t *testing.T) {
	testCode := [][]string{
		{"(make-vector 10 0)", "#(0 0 0 0 0 0 0 0 0 0)"},
		{"(make-vector 4 (list 1 2 3))", "#((1 2 3) (1 2 3) (1 2 3) (1 2 3))"},
		{"(make-vector 8 'a)", "#(a a a a a a a a)"},
		{"(make-vector 0 'a)", "#()"},

		{"(make-vector)", "E1007"},
		{"(make-vector 10)", "E1007"},
		{"(make-vector 10 0 1)", "E1007"},
		{"(make-vector #t 0)", "E1002"},
		{"(make-vector -1 0)", "E1011"},
		{"(make-vector 10 c)", "E1008"},
	}
	executeTest(testCode, "make-vector", t)
}
func TestVectorLength(t *testing.T) {
	testCode := [][]string{
		{"(vector-length (vector))", "0"},
		{"(vector-length (vector 3))", "1"},
		{"(vector-length (list->vector (iota 10)))", "10"},

		{"(vector-length)", "E1007"},
		{"(vector-length (vector 1)(vector 2))", "E1007"},
		{"(vector-length (list 1 2))", "E1022"},
		{"(vector-length (cons 1 2))", "E1022"},
		{"(vector-length a)", "E1008"},
	}
	executeTest(testCode, "vector-length", t)
}
func TestVectorList(t *testing.T) {
	testCode := [][]string{
		{"(vector->list #(1 2 3))", "(1 2 3)"},
		{"(vector->list (vector 1 2 3))", "(1 2 3)"},
		{"(vector->list #())", "()"},
		{"(vector->list)", "E1007"},
		{"(vector->list #(1 2 3) #(1 2 3))", "E1007"},
		{"(vector->list '(1 2 3))", "E1022"},
	}
	executeTest(testCode, "vector->list", t)
}
func TestListVector(t *testing.T) {
	testCode := [][]string{
		{"(list->vector '(1 2 3))", "#(1 2 3)"},
		{"(list->vector (list 1 2 3))", "#(1 2 3)"},
		{"(list->vector '())", "#()"},
		{"(list->vector)", "E1007"},
		{"(list->vector (list 1 2 3) '(1 2 3))", "E1007"},
		{"(list->vector #(1 2 3))", "E1005"},
	}
	executeTest(testCode, "list->vector", t)
}
