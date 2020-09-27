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

func TestNot(t *testing.T) {
	testCode := [][]string{
		{"(not (= 0.75 0.75))", "#f"},
		{"(not (= 0.15 0.75))", "#t"},

		{"(not 10)", "E1001"},
		{"(not #t #f)", "E1007"},
		{"(not)", "E1007"},
	}
	executeTest(testCode, "not", t)
}
func TestAnd(t *testing.T) {
	testCode := [][]string{
		{"(and (= 1 1)(= 2 2))", "#t"},
		{"(and (= 1 1)(= 2 3))", "#f"},
		{"(and (= 2 1)(= 2 2))", "#f"},
		{"(and (= 0 1)(= 2 3))", "#f"},
		{"(let ((a 10)(b 20)(c 30)) (and (< a b)(< a c)(< a c)))", "#t"},
		{"(let ((a 10)(b 20)(c 30)) (and (< a b)(< a c)(< c a)))", "#f"},

		{"(and #t)", "E1007"},
		{"(and 10.2 0 1)", "E1001"},
	}
	executeTest(testCode, "and", t)
}
func TestOr(t *testing.T) {
	testCode := [][]string{
		{"(and (= 1 1)(= 2 2))", "#t"},
		{"(and (= 1 1)(= 2 3))", "#f"},
		{"(and (= 2 1)(= 2 2))", "#f"},
		{"(and (= 0 1)(= 2 3))", "#f"},
		{"(let ((a 10)(b 20)(c 30)) (or (= a b)(< b c)))", "#t"},
		{"(let ((a 10)(b 20)(c 30)) (or (= c a)(< c b)))", "#f"},

		{"(or #t)", "E1007"},
		{"(or 10.2 0 1)", "E1001"},
	}
	executeTest(testCode, "or", t)
}
func TestIf(t *testing.T) {
	testCode := [][]string{
		{"(if (= 10 10) #\\a)", "#\\a"},
		{"(if (= 10 11) #\\a)", "nil"},
		{"(if (<= 1 6) #\\a #\\b)", "#\\a"},
		{"(if (<= 9 6) #\\a #\\b)", "#\\b"},
		{"(let ((a 10)(b 20))(if (= a b) #t))", "nil"},

		{"(if 10 1 2)", "E1001"},
		{"(if (= 10 10))", "E1007"},
	}
	executeTest(testCode, "if", t)
}
func TestCond(t *testing.T) {
	testCode := [][]string{
		{"(let ((a 10)(b 10))(cond ((= a b) \"ok\")(else \"ng\")))", "\"ok\""},
		{"(let ((a 10)(b 20))(cond ((= a b) \"ok\")(else \"ng\")))", "\"ng\""},
		{"(let ((a 10)(b 20))(cond ((= a b) \"ok\")((= b 20) \"sankaku\")(else \"ng\")))", "\"sankaku\""},
		{"(let ((a 10)(b 20))(cond ((= a b) #t)))", "nil"},
		{"(define a 10)", "a"},
		{"(cond ((= a 10) 10 11)(else 20 30))", "11"},
		{"(define a 100)", "a"},
		{"(cond ((= a 10) 10 11)(else 20 30))", "30"},

		{"(cond)", "E1007"},
		{"(cond 10)", "E1005"},
		{"(cond (10))", "E1007"},
		{"(let ((a 10)(b 20))(cond ((= a b) #t)(lse #f)))", "E1012"},
		{"(cond (10 10))", "E1012"},
		{"(cond ((+ 10 20) 10 11)(else 20 30))", "E1001"},
	}
	executeTest(testCode, "cond", t)
}
func TestDelayForce(t *testing.T) {
	testCode := [][]string{
		{"(force (+ 1 1))", "2"},
		{"(force ((lambda (a) (delay (* 10 a))) 3))", "30"},
		{"(force (delay (+ 1 2)))", "3"},

		{"(delay)", "E1007"},
		{"(delay 1 2)", "E1007"},
		{"(force)", "E1007"},
		{"(force 1 2)", "E1007"},
	}
	executeTest(testCode, "delay_force", t)
}

func TestDefine(t *testing.T) {
	testCode := [][]string{
		{"(define foo (lambda () (define hoge (lambda (a) (+ 1 a))) (hoge 10)))", "foo"},
		{"(foo)", "11"},
		{"(define a 100)", "a"},
		{"a", "100"},
		{"(define\ta\t200)", "a"},
		{"a", "200"},
		{"(define\na\n300)", "a"},
		{"a", "300"},
		{"(define\r\na\r\n400)", "a"},
		{"a", "400"},

		{"(define 10 10)", "E1004"},
		{"(define a)", "E1007"},
		{"(define a 10 11)", "E1007"},
		{"(define (a))", "E1007"},
		{"(define 10 11)", "E1004"},
		{"(define (fuga 1 b) (+ a b))", "E1004"},
		{"hoge", "E1008"},
	}
	executeTest(testCode, "define", t)
}
func TestQuote(t *testing.T) {
	testCode := [][]string{
		{"(quote a)", "a"},
		{"(quote (a b c))", "(a b c)"},

		{"(quote)", "E1007"},
		{"(quote 1 2)", "E1007"},
	}
	executeTest(testCode, "quote", t)
}
func TestLet(t *testing.T) {
	testCode := [][]string{
		{"(let ((a 10)(b 20))(+ a b)(* a b))", "200"},
		{"(let loop ((i 0)(j 0)) (if (<= 10 i) (+ i j) (loop (+ i 1)(+ j 2))))", "30"},

		{"(let loop ((i 0)(j 10)(k 10)) (if (<= 1000000 i) i (if (= j k) (loop (+ 100 i)(+ 1 i)))))", "E1007"},
		{"(let ((a 10)))", "E1007"},
		{"(let 10 ((a 10)))", "E1004"},
		{"(let loop ((a 10)))", "E1007"},
		{"(let loop 10 19)", "E1005"},
		{"(let ((a))(+ 1 1))", "E1007"},
		{"(let ((a 10)) b a)", "E1008"},
		{"(let ((1 1)(a 0)) a)", "E1004"},
	}
	executeTest(testCode, "let", t)
}
func TestLambda(t *testing.T) {
	testCode := [][]string{
		{"((lambda (a b)(+ a b)) 1 2)", "3"},
		{"(define hoge (lambda (a b) (+ a b)))", "hoge"},
		{"(hoge 6 8)", "14"},
		{"(define hoge (lambda (a b) b))", "hoge"},
		{"(hoge 6 8)", "8"},

		{"(lambda a (+ a b))", "E1005"},
		{"(lambda (+ n m))", "E1007"},
		{"(lambda 10 11)", "E1005"},
		{"((lambda (n m) (+ n m)) 1 2 3)", "E1007"},
		{"((lambda (n m) (+ n m)) 1)", "E1007"},
		{"(lambda (1 1) (+ 1 2))", "E1004"},
	}
	executeTest(testCode, "lambda", t)
}
func TestSet(t *testing.T) {
	testCode := [][]string{
		{"(define a 10)", "a"},
		{"(set! a 20)", "20"},
		{"a", "20"},

		{"(set! 10 10)", "E1004"},
		{"(set! a)", "E1007"},
		{"(set! a 10 11)", "E1007"},
		{"(set! hoge 10)", "E1008"},
	}
	executeTest(testCode, "set!", t)
}
func TestBegin(t *testing.T) {
	testCode := [][]string{
		{"(begin 1 2)", "2"},

		{"(begin)", "E1007"},
	}
	executeTest(testCode, "load-file", t)
}