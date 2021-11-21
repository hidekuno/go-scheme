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
		{"(if (= 10 11) #\\a)", "#<nil>"},
		{"(if (<= 1 6) #\\a #\\b)", "#\\a"},
		{"(if (<= 9 6) #\\a #\\b)", "#\\b"},
		{"(let ((a 10)(b 20))(if (= a b) #t))", "#<nil>"},

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
		{"(let ((a 10)(b 20))(cond ((= a b) #t)))", "#<nil>"},
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
func TestCase(t *testing.T) {
	testCode := [][]string{
		{"(case 10)", "#<nil>"},
		{"(case 10 ((1 2) \"A\"))", "#<nil>"},
		{"(case 1 ((1 2)))", "(1 2)"},
		{"(case 10 (else 20))", "20"},
		{"(case 10 (else))", "0"},

		{"(case 100 ((100 200) \"A\")(else \"B\"))", "\"A\""},
		{"(case 300 ((100 200) \"A\")(else \"B\"))", "\"B\""},
		{"(case 200 ((100 200) \"A\")(else \"B\"))", "\"A\""},
		{"(case 300 ((100 200) \"A\")((300 400) \"B\")(else \"C\"))", "\"B\""},
		{"(case 400 ((100 200) \"A\")((300 400) \"B\")(else \"C\"))", "\"B\""},
		{"(case 500 ((100 200) \"A\")((300 400) \"B\")(else \"C\"))", "\"C\""},

		{"(case)", "E1007"},
		{"(case 10 (hoge 20))", "E1017"},
		{"(case 10 10)", "E1005"},
		{"(case 10 (20))", "E1017"},
		{"(case a)", "E1008"},
		{"(case 10 ((10 20) a))", "E1008"},
		{"(case 10 ((20 30) 1)(else a))", "E1008"},
	}
	executeTest(testCode, "case", t)
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
		{"(define)", "E1007"},
		{"(define (a))", "E1007"},
		{"(define 10 11)", "E1004"},
		{"(define (fuga 1 b) (+ a b))", "E1004"},
		{"hoge", "E1008"},
	}
	executeTest(testCode, "define", t)
}
func TestApply(t *testing.T) {
	testCode := [][]string{
		{"(apply + (list 1 2 3))", "6"},
		{"(apply + (list (+ 1 1) 2 3))", "7"},
		{"(apply - (list 5 3 2))", "0"},
		{"(apply (lambda (a b) (+ a b)) (list 1 2))", "3"},
		{"(apply + (iota 10))", "45"},
		{"(define hoge (lambda (a b) (* a b)))", "hoge"},
		{"(apply hoge (list 3 4))", "12"},
		{"(apply append (list (list 1 2 3)(list 4 5 6)))", "(1 2 3 4 5 6)"},
		{"(apply (lambda (a) (map (lambda (n) (* n n)) a)) (list (list 1 2 3)))", "(1 4 9)"},

		{"(apply)", "E1007"},
		{"(apply -)", "E1007"},
		{"(apply + (list 1 2)(lis 3 4))", "E1007"},
		{"(apply + 10)", "E1005"},
		{"(apply fuga (list 1 2))", "E1008"},
	}

	executeTest(testCode, "apply", t)
}
func TestQuote(t *testing.T) {
	testCode := [][]string{
		{"(quote a)", "a"},
		{"'a", "a"},
		{"'a 'b", "b"},
		{"(quote (a b c))", "(a b c)"},
		{"'(a b c)", "(a b c)"},
		{"(append '(a b c)'(d e f))", "(a b c d e f)"},
		{"'10", "10"},
		{"'(1 2 (3 4 (5 6)))", "(1 2 (3 4 (5 6)))"},

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
		{"(let)", "E1007"},
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
	executeTest(testCode, "begin", t)
}
func TestDo(t *testing.T) {
	testCode := [][]string{
		{"(do ((i 0 (+ i 1)))((= i 10) i))", "10"},
		{"(do ((i 0 (+ i 1))(j 0 (+ i j)))((= i 10) j)(display j)(newline))", "45"},
		{"(do ((a '(0 1 2 3 4) (cdr a))(b 0 (+ b (car a))))((null? a) b)(display (car a))(newline))", "10"},
		{"(define x 100)", "x"},
		{"(do ((i 0 (+ i 1)))((= i 10) x)(set! x (+ i x)))", "145"},

		{"(do)", "E1007"},
		{"(do 1 2)", "E1005"},
		{"(do () 1)", "E1007"},
		{"(do (a) 1)", "E1005"},
		{"(do (()) 1)", "E1007"},
		{"(do ((10 1 1)) 1)", "E1004"},

		{"(do ((i 0 (+ 1))) 10)", "E1005"},
		{"(do ((i 0 (+ 1))) (10))", "E1007"},
		{"(do ((i 0 (+ 1))) (10 10))", "E1001"},
		{"(do ((i 0 (+ 1))) (#f 10) a)", "E1008"},
		{"(do ((i 0 (+ 1))) (#t a) 10)", "E1008"},
	}
	executeTest(testCode, "do", t)
}
func TestCallcc(t *testing.T) {
	testCode := [][]string{
		{"(+ 1 (* 2 (call/cc (lambda (cont) (cont 3)))))", "7"},
		{"(call/cc (lambda (throw)(+ 5 (* 10 (call/cc (lambda (escape) (* 100 (throw 3))))))))", "3"},
		{"(call/cc (lambda (throw)(+ 5 (* 10 (call/cc (lambda (escape) (* 100 (escape 3))))))))", "35"},
		{"(call/cc (lambda (hoge) (+ 3 (call/cc (lambda (throw)(+ 5 (* 10 (call/cc (lambda (escape)" +
			"(* 100 (throw 3)))))))))))", "6"},
		{"(call/cc (lambda (hoge) (+ 3 (call/cc (lambda (throw)(+ 5 (* 10 (call/cc (lambda (escape)" +
			"(* 100 (escape 3)))))))))))", "38"},
		{"(define (map-check fn chk ls)(call/cc (lambda (return)" +
			"(map (lambda (x) (if (chk x) (return '()) (fn x))) ls))))", "map-check"},
		{"(map-check (lambda (x) (* x x)) (lambda (x) (< x 0)) (list 1 2 3 4 5))", "(1 4 9 16 25)"},
		{"(map-check (lambda (x) (* x x)) (lambda (x) (< x 0)) (list 1 2 3 -1 5))", "()"},
	}
	executeTest(testCode, "call/cc", t)
}
