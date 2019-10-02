/*
   Go lang 4th study program.
   This is test program for mini scheme subset.

   ex.) go test -v lisp.go lisp_test.go

   hidekuno@gmail.com
*/
package scheme

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
)

func executeTest(testCode [][]string, testName string, t *testing.T) {

	BuildFunc()
	rootEnv := NewSimpleEnv(nil, nil)
	for i, c := range testCode {
		exp, err := DoCoreLogic(c[0], rootEnv)
		if err != nil {
			if e, ok := err.(*SyntaxError); ok {
				if e.MsgCode != c[1] {
					t.Log(i)
					t.Fatal("failed "+testName+" test", e.MsgCode)
				}
			} else if e, ok := err.(*RuntimeError); ok {
				if e.MsgCode != c[1] {
					t.Log(i)
					t.Fatal("failed "+testName+" test", e.MsgCode)
				}
			} else {
				t.Log(i)
				t.Fatal("failed "+testName+" test", err.Error())
			}
		} else {
			if exp.String() != c[1] {
				t.Log(i)
				t.Fatal("failed "+testName+" test", exp)
			}
		}
	}
}

var (
	program = []string{
		"(define test-list (list 36 27 14 19 2 8 7 6 0 9 3))",
		"(define (counter) (let ((c 0)) (lambda () (set! c (+ 1 c)) c)))",
		"(define a (counter))",
		"(define b (counter))",
		"(define (step-counter s) (let ((c 0)) (lambda () (set! c (+ s c)) c)))",
		"(define x (step-counter 10))",
		"(define y (step-counter 100))",
		"(define (gcm n m) (let ((mod (modulo n m))) (if (= 0 mod) m (gcm m mod))))",
		"(define (lcm n m) (/(* n m)(gcm n m)))",
		"(define hanoi (lambda (from to work n) (if (>= 0 n)(list)(append (hanoi from work to (- n 1))(list (list (cons from to) n))(hanoi work to from (- n 1))))))",
		"(define prime (lambda (l) (if (> (car l)(sqrt (last l))) l (cons (car l)(prime (filter (lambda (n) (not (= 0 (modulo n (car l))))) (cdr l)))))))",
		"(define qsort (lambda (l pred) (if (null? l) l (append (qsort (filter (lambda (n) (pred n (car l))) (cdr l)) pred) (cons (car l) (qsort (filter (lambda (n) (not (pred n (car l))))(cdr l)) pred))))))",
		"(define comb (lambda (l n) (if (null? l) l (if (= n 1) (map (lambda (n) (list n)) l) (append (map (lambda (p) (cons (car l) p)) (comb (cdr l)(- n 1))) (comb (cdr l) n))))))",
		"(define delete (lambda (x l) (filter (lambda (n) (not (= x n))) l)))",
		"(define perm (lambda (l n)(if (>= 0 n) (list (list))(reduce (lambda (a b)(append a b)) (list) (map (lambda (x) (map (lambda (p) (cons x p)) (perm (delete x l)(- n 1)))) l)))))",
		"(define bubble-iter (lambda (x l)(if (or (null? l)(< x (car l)))(cons x l)(cons (car l)(bubble-iter x (cdr l))))))",
		"(define bsort (lambda (l)(if (null? l) l (bubble-iter (car l)(bsort (cdr l))))))",
		"(define merge (lambda (a b)(if (or (null? a)(null? b)) (append a b) (if (< (car a)(car b))(cons (car a)(merge (cdr a) b))(cons (car b) (merge a (cdr b)))))))",
		"(define take (lambda (l n)(if (>= 0 n) (list)(cons (car l)(take (cdr l)(- n 1))))))",
		"(define drop (lambda (l n)(if (>= 0 n) l (drop (cdr l)(- n 1)))))",
		"(define msort (lambda (l)(let ((n (length l)))(if (>= 1 n ) l (if (= n 2) (if (< (car l)(cadr l)) l (reverse l))(let ((mid (quotient n 2)))(merge (msort (take l mid))(msort (drop l mid)))))))))",
		"(define stream-car (lambda (l)(car l)))",
		"(define stream-cdr (lambda (l)(force (cdr l))))",
		"(define make-generator (lambda (generator inits)(cons (car inits)(delay (make-generator generator (generator inits))))))",
		"(define inf-list (lambda (generator inits limit)(let loop ((l (make-generator generator inits))(c limit)) (if (>= 0 c) (list)(cons (stream-car l)(loop (stream-cdr l)(- c 1)))))))",
		"(define fact/cps (lambda (n cont)(if (= n 0)(cont 1)(fact/cps (- n 1) (lambda (a) (cont (* n a)))))))",
		"(define fact (lambda (n) (fact/cps n identity)))",
		"(define (testf x) (lambda () (* x 10)))",
		"(define (create-testf  x) (testf (* 2 x)))",
	}
)

func TestAtom(t *testing.T) {
	testCode := [][]string{
		{"10", "10"},
		{"10.5", "10.5"},
		{"\"ABC\"", "\"ABC\""},
		{"\"AB\\\"C\"", "\"AB\\\"C\""},
		{"\"(A B C)\"", "\"(A B C)\""},
		{"\"(\"", "\"(\""},
		{"\"  a  \"", "\"  a  \""},
		{"#t", "#t"},
		{"#f", "#f"},
		{"#\\tab", "#\\tab"},
		{"#\\space", "#\\space"},
		{"#\\newline", "#\\newline"},
		{"#\\return", "#\\return"},
		{"#\\A", "#\\A"},

		{"(", "E0001"},
		{"(a (b", "E0002"},
		{")", "E0003"},
		{"(a))", "E0003"},
		{"1)", "E0003"},
		{"#\\abc", "E0004"},
	}
	executeTest(testCode, "atom", t)
}
func TestAtomUTF8(t *testing.T) {
	testCode := [][]string{
		{"#\\山", "#\\山"},
		{"\"山田太郎\"", "\"山田太郎\""},
		{"\"山田(太郎\"", "\"山田(太郎\""},
		{"(define 山 25)", "山"},
		{"山", "25"},
	}
	executeTest(testCode, "atom_utf8", t)
}
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
func TestIdentity(t *testing.T) {
	testCode := [][]string{
		{"(identity 100)", "100"},
		{"(identity \"ABC\")", "\"ABC\""},

		{"(identity 100 200)", "E1007"},
		{"(identity)", "E1007"},
	}
	executeTest(testCode, "identity", t)
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
	executeTest(testCode, "identity", t)
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
		{"(iota 5 2)", "(2 3 4 5 6)"},

		{"(iota 10.2 1)", "E1002"},
		{"(iota 1 10.2)", "E1002"},
		{"(iota)", "E1007"},
		{"(iota 1 2 3)", "E1007"},
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

		{"(reduce (lambda (a b) (+ a b)) (+ 1 2) 20)", "E1005"},
		{"(reduce (list 1 12) 0 (list 10))", "E1006"},
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
		{"(expt 2 3)", "8"},

		{"(expt 10)", "E1007"},
		{"(expt 10 10 10)", "E1007"},
		{"(expt 11.5 10)", "E1002"},
		{"(expt 11 12.5)", "E1002"},
	}
	executeTest(testCode, "expt", t)
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

func TestTime(t *testing.T) {
	testCode := [][]string{
		{"(time)", "E1007"},
		{"(time #\\abc)", "E0004"},
	}
	executeTest(testCode, "time", t)
}
func TestBegin(t *testing.T) {
	testCode := [][]string{
		{"(begin 1 2)", "2"},

		{"(begin)", "E1007"},
	}
	executeTest(testCode, "load-file", t)
}
func TestLoadFile(t *testing.T) {
	testCode := [][]string{
		{"(load-file)", "E1007"},
		{"(load-file 10)", "E1015"},
		{"(load-file \"example/no.scm\")", "E1014"},
		{"(load-file \"/tmp\")", "E1016"},
		{"(load-file \"/etc/sudoers\")", "E9999"},
	}
	executeTest(testCode, "load-file", t)
}

func TestLispSampleProgram(t *testing.T) {
	var (
		exp Expression
	)
	BuildFunc()
	rootEnv := NewSimpleEnv(nil, nil)
	for _, p := range program {
		exp, _ = DoCoreLogic(p, rootEnv)
	}

	exp, _ = DoCoreLogic("(let loop ((a 0)(r (list 1 2 3))) (if (null? r) a (loop (+ (car r) a)(cdr r))))", rootEnv)
	if exp.String() != "6" {
		t.Fatal("failed test: let loop")
	}
	exp, _ = DoCoreLogic("(a)", rootEnv)
	exp, _ = DoCoreLogic("(a)", rootEnv)
	exp, _ = DoCoreLogic("(a)", rootEnv)
	if exp.String() != "3" {
		t.Fatal("failed test: closure")
	}
	exp, _ = DoCoreLogic("(b)", rootEnv)
	exp, _ = DoCoreLogic("(b)", rootEnv)
	if exp.String() != "2" {
		t.Fatal("failed test: closure")
	}
	exp, _ = DoCoreLogic("(x)", rootEnv)
	exp, _ = DoCoreLogic("(x)", rootEnv)

	if exp.String() != "20" {
		t.Fatal("failed test: closure")
	}
	exp, _ = DoCoreLogic("(y)", rootEnv)
	exp, _ = DoCoreLogic("(y)", rootEnv)
	if exp.String() != "200" {
		t.Fatal("failed test: closure")
	}

	exp, _ = DoCoreLogic("(gcm 36 27)", rootEnv)
	if exp.String() != "9" {
		t.Fatal("failed test: gcm")
	}
	exp, _ = DoCoreLogic("(lcm 36 27)", rootEnv)
	if exp.String() != "108" {
		t.Fatal("failed test: lcm")
	}

	exp, _ = DoCoreLogic("(qsort test-list (lambda (a b)(< a b)))", rootEnv)
	if exp.String() != "(0 2 3 6 7 8 9 14 19 27 36)" {
		t.Fatal("failed test: qsort")
	}
	exp, _ = DoCoreLogic("(bsort test-list)", rootEnv)
	if exp.String() != "(0 2 3 6 7 8 9 14 19 27 36)" {
		t.Fatal("failed test: bsort")
	}

	exp, _ = DoCoreLogic("(prime (iota 30 2))", rootEnv)
	if exp.String() != "(2 3 5 7 11 13 17 19 23 29 31)" {
		t.Fatal("failed test: prime")
	}

	exp, _ = DoCoreLogic("(perm (list 1 2 3) 2)", rootEnv)
	if exp.String() != "((1 2) (1 3) (2 1) (2 3) (3 1) (3 2))" {
		t.Fatal("failed test: perm")
	}

	exp, _ = DoCoreLogic("(comb (list 1 2 3) 2)", rootEnv)
	if exp.String() != "((1 2) (1 3) (2 3))" {
		t.Fatal("failed test: comb")
	}
	exp, _ = DoCoreLogic("(hanoi (quote a)(quote b)(quote c) 3)", rootEnv)
	if exp.String() != "(((a . b) 1) ((a . c) 2) ((b . c) 1) ((a . b) 3) ((c . a) 1) ((c . b) 2) ((a . b) 1))" {
		t.Fatal("failed test: hanoi")
	}
	exp, _ = DoCoreLogic("(merge (list 1 3 5 7 9)(list 2 4 6 8 10))", rootEnv)
	if exp.String() != "(1 2 3 4 5 6 7 8 9 10)" {
		t.Fatal("failed test: merge")
	}
	exp, _ = DoCoreLogic("(take (list 2 4 6 8 10) 3)", rootEnv)
	if exp.String() != "(2 4 6)" {
		t.Fatal("failed test: take")
	}
	exp, _ = DoCoreLogic("(drop (list 2 4 6 8 10) 3)", rootEnv)
	if exp.String() != "(8 10)" {
		t.Fatal("failed test: drop")
	}
	exp, _ = DoCoreLogic("(msort test-list)", rootEnv)
	if exp.String() != "(0 2 3 6 7 8 9 14 19 27 36)" {
		t.Fatal("failed test: bsort")
	}
	exp, _ = DoCoreLogic("(inf-list (lambda (n) (list (cadr n)(+ (car n)(cadr n)))) (list 0 1) 10)", rootEnv)
	if exp.String() != "(0 1 1 2 3 5 8 13 21 34)" {
		t.Fatal("failed test: fibonacci")
	}

	exp, _ = DoCoreLogic("(fact/cps 5 (lambda (a) (+ 80 a)))", rootEnv)
	if exp.String() != "200" {
		t.Fatal("failed test: fact/cps")
	}
	exp, _ = DoCoreLogic("(fact 5)", rootEnv)
	if exp.String() != "120" {
		t.Fatal("failed test: fact")
	}
	exp, _ = DoCoreLogic("((create-testf 2))", rootEnv)
	if exp.String() != "40" {
		t.Fatal("failed test: create-testf")
	}
}
func TestErrCase(t *testing.T) {
	var (
		err error
	)
	BuildFunc()
	rootEnv := NewSimpleEnv(nil, nil)

	// Impossible absolute, But Program bug is except
	_, err = eval(NewFunction(rootEnv, NewList(nil), nil, "lambda"), rootEnv)
	if err != nil {
		if false == strings.Contains(err.Error(), "Not Enough Data Type") {
			t.Fatal("failed test: E1009")
		}
	} else {
		t.Fatal("failed test2: E1009")
	}
	// Error()
	_, err = DoCoreLogic(")", rootEnv)
	if err != nil {
		if false == strings.Contains(err.Error(), "Extra close") {
			t.Fatal("failed test: SyntaxError::Error()")
		}
	} else {
		t.Fatal("failed test2: SyntaxError::Error()")
	}
	_, err = DoCoreLogic("undef", rootEnv)
	if err != nil {
		if false == strings.Contains(err.Error(), "Undefine variable") {
			t.Fatal("failed test: RuntimeError::Error()")
		}
	} else {
		t.Fatal("failed test2: SyntaxError::Error()")
	}
	err = NewRuntimeError("E1008", "a", "b")
	if err != nil {
		if false == strings.Contains(err.Error(), "Undefine variable") {
			t.Fatal("failed test: RuntimeError::Error()")
		}
	} else {
		t.Fatal("failed test: RuntimeError::Error()")
	}
}
func TestInteractive(t *testing.T) {
	var iostub func(program string, ret string)

	iostub = func(program string, ret string) {
		inr, inw, _ := os.Pipe()
		outr, outw, _ := os.Pipe()
		errr, errw, _ := os.Pipe()
		orgStdin := os.Stdin
		orgStdout := os.Stdout
		orgStderr := os.Stderr
		inw.Write([]byte(program))
		inw.Close()
		os.Stdin = inr
		os.Stdout = outw
		os.Stderr = errw

		DoInteractive()

		os.Stdin = orgStdin
		os.Stdout = orgStdout
		os.Stderr = orgStderr
		outw.Close()
		outbuf, _ := ioutil.ReadAll(outr)
		errw.Close()
		errbuf, _ := ioutil.ReadAll(errr)

		s := string(outbuf)
		s = strings.Replace(s, "scheme.go>", "", -1)
		s = strings.Replace(s, "\n", "", -1)
		s = strings.Replace(s, "\t", "", -1)

		rep := regexp.MustCompile(`^ *`)
		s = rep.ReplaceAllString(s, "")
		rep = regexp.MustCompile(` *$`)
		s = rep.ReplaceAllString(s, "")
		rep = regexp.MustCompile(`: &.*$`)
		s = rep.ReplaceAllString(s, ":")
		if s != ret {
			t.Fatal(s)
			t.Fatal(string(errbuf))
		}
		t.Log(s)
	}
	iostub("(+ 1 2.5)", "3.5")
	iostub("((lambda \n(n m)(+ n m))\n 10 20)", "30")
	iostub("(define a 1)", "a")
	iostub("(= 10 10)", "#t")
	iostub("\"ABC\"", "\"ABC\"")
	iostub("(list 1 2 3 (list 4 5))", "(1 2 3 (4 5))")
	iostub("(cons 1 2)", "(1 . 2)")
	iostub("(lambda (n m) (+ n m))", "Function:")
	iostub("+", "Build In Function: +")
	iostub("if", "Build In Function: if")
	iostub("(delay 1)", "Promise:")
	iostub("#\\space", "#\\space")
	iostub("#\\newline", "#\\newline")
	iostub("#\\tab", "#\\tab")
	iostub("#\\return", "#\\return")
	iostub("#\\A", "#\\A")
}

//https://github.com/hidekuno/go-scheme/issues/46
func TestRecursive(t *testing.T) {
	var (
		exp Expression
	)
	BuildFunc()
	rootEnv := NewSimpleEnv(nil, nil)
	DoCoreLogic("(define (fact n result)(if (>= 1 n) result (fact (- n 1) (* result n))))", rootEnv)

	exp, _ = DoCoreLogic("(fact 5 1)", rootEnv)
	if exp.String() != "120" {
		t.Fatal("failed test: tail recursive")
	}

	exp, _ = DoCoreLogic("(let loop ((i 0)) (if (<= 1000000 i) i (loop (+ 1 i))))", rootEnv)
	if exp.String() != "1000000" {
		t.Fatal("failed test: tail recursive")
	}

	exp, _ = DoCoreLogic("(let loop ((i 0)(j 10)(k 10)) (if (<= 1000000 i) i (if (= j k) (loop (+ 50 i) j k)(loop (+ 1 i) j k))))", rootEnv)
	if exp.String() != "1000000" {
		t.Fatal("failed test: tail recursive")
	}
}

// go test -bench . -benchmem
func BenchmarkQsort(b *testing.B) {

	BuildFunc()
	rootEnv := NewSimpleEnv(nil, nil)

	DoCoreLogic("(define test-list (map (lambda (n) (rand-integer 10000))(iota 600)))", rootEnv)
	DoCoreLogic("(define qsort (lambda (l)(if (null? l) l (append (qsort (filter (lambda (n) (< n (car l)))(cdr l)))(cons (car l)(qsort (filter (lambda (n) (not (< n (car l))))(cdr l))))))))", rootEnv)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DoCoreLogic("(qsort test-list)", rootEnv)
	}
}
