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

func checkHanoi(exp Expression) bool {
	if l, ok := exp.(*List); ok {
		if len(l.Value) == 7 {
			return true
		}
	}
	return false
}
func checkLogicMatrix(exp Expression, items [][]int) bool {
	if l, ok := exp.(*List); ok {
		for i, e := range l.Value {
			if _, ok := e.(*List); !ok {
				return false
			}
			for j, r := range (e.(*List)).Value {
				if !checkLogicInt(r, items[i][j]) {
					return false
				}
			}
		}
		return true
	} else {
		return false
	}
}
func checkLogicList(exp Expression, items []int) bool {

	if l, ok := exp.(*List); ok {
		for i, e := range l.Value {
			if !checkLogicInt(e, items[i]) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func checkLogicInt(exp Expression, v int) bool {
	if i, ok := exp.(*Integer); ok {
		if i.Value != v {
			return false
		}
	} else {
		return false
	}
	return true
}

func checkErrorCode(err error, errCode string) bool {
	if e, ok := err.(*SyntaxError); ok {
		if e.MsgCode == errCode {
			return true
		}
	}
	if e, ok := err.(*RuntimeError); ok {
		if e.MsgCode == errCode {
			return true
		}
	}
	return false
}

var (
	program = []string{
		"(define test-list (list 36 27 14 19 2 8 7 6 0 9 3))",
		"(define (counter) (let ((c 0)) (lambda () (set! c (+ 1 c)) c)))",
		"(define a (counter))",
		"(define b (counter))",
		"(define (gcm n m) (let ((mod (modulo n m))) (if (= 0 mod) m (gcm m mod))))",
		"(define (lcm n m) (/(* n m)(gcm n m)))",
		"(define hanoi (lambda (from to work n) (if (>= 0 n)(list)(append (hanoi from work to (- n 1))(list (list (cons from to) n))(hanoi work to from (- n 1))))))",
		"(define prime (lambda (l) (if (> (car l)(sqrt (last l))) l (cons (car l)(prime (filter (lambda (n) (not (= 0 (modulo n (car l))))) (cdr l)))))))",
		"(define qsort (lambda (l pred) (if (null? l) l (append (qsort (filter (lambda (n) (pred n (car l))) (cdr l)) pred) (cons (car l) (qsort (filter (lambda (n) (not (pred n (car l))))(cdr l)) pred))))))",
		"(define comb (lambda (l n) (if (null? l) l (if (= n 1) (map (lambda (n) (list n)) l) (append (map (lambda (p) (cons (car l) p)) (comb (cdr l)(- n 1))) (comb (cdr l) n))))))",
		"(define delete (lambda (x l) (filter (lambda (n) (not (= x n))) l)))",
		"(define perm (lambda (l n)(if (>= 0 n) (list (list))(reduce (lambda (a b)(append a b))(map (lambda (x) (map (lambda (p) (cons x p)) (perm (delete x l)(- n 1)))) l)))))",
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
		"(define fact/cont (lambda (n)  (call/cc (lambda (cont)  (if (= n 3) (cont n) (if (>= 1 n) 1 (* n (fact/cont (- n 1)))))))))",
	}
)

func TestCheckFunclist(t *testing.T) {
	BuildFunc()
	for key, _ := range builtinFuncTbl {
		t.Log(key)
	}
	for key, _ := range specialFuncTbl {
		t.Log(key)
	}
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
	if !checkLogicInt(exp, 6) {
		t.Fatal("failed test: let loop")
	}
	exp, _ = DoCoreLogic("(a)", rootEnv)
	exp, _ = DoCoreLogic("(a)", rootEnv)
	exp, _ = DoCoreLogic("(a)", rootEnv)
	if !checkLogicInt(exp, 3) {
		t.Fatal("failed test: closure")
	}
	exp, _ = DoCoreLogic("(b)", rootEnv)
	exp, _ = DoCoreLogic("(b)", rootEnv)
	if !checkLogicInt(exp, 2) {
		t.Fatal("failed test: closure")
	}
	exp, _ = DoCoreLogic("(gcm 36 27)", rootEnv)
	if !checkLogicInt(exp, 9) {
		t.Fatal("failed test: gcm")
	}
	exp, _ = DoCoreLogic("(lcm 36 27)", rootEnv)
	if !checkLogicInt(exp, 108) {
		t.Fatal("failed test: lcm")
	}

	testSortData := []int{0, 2, 3, 6, 7, 8, 9, 14, 19, 27, 36}
	exp, _ = DoCoreLogic("(qsort test-list (lambda (a b)(< a b)))", rootEnv)
	if !checkLogicList(exp, testSortData) {
		t.Fatal("failed test: qsort")
	}
	exp, _ = DoCoreLogic("(bsort test-list)", rootEnv)
	if !checkLogicList(exp, testSortData) {
		t.Fatal("failed test: bsort")
	}

	prime_data := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31}
	exp, _ = DoCoreLogic("(prime (iota 30 2))", rootEnv)
	if !checkLogicList(exp, prime_data) {
		t.Fatal("failed test: prime")
	}

	perm_data := [][]int{{1, 2}, {1, 3}, {2, 1}, {2, 3}, {3, 1}, {3, 2}}
	exp, _ = DoCoreLogic("(perm (list 1 2 3) 2)", rootEnv)
	if !checkLogicMatrix(exp, perm_data) {
		t.Fatal("failed test: perm")
	}

	comb := [][]int{{1, 2}, {1, 3}, {2, 3}}
	exp, _ = DoCoreLogic("(comb (list 1 2 3) 2)", rootEnv)
	if !checkLogicMatrix(exp, comb) {
		t.Fatal("failed test: comb")
	}
	exp, _ = DoCoreLogic("(hanoi \"a\" \"b\" \"c\" 3)", rootEnv)
	if !checkHanoi(exp) {
		t.Fatal("failed test: hanoi")
	}
	exp, _ = DoCoreLogic("(merge (list 1 3 5 7 9)(list 2 4 6 8 10))", rootEnv)
	if !checkLogicList(exp, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Fatal("failed test: merge")
	}
	exp, _ = DoCoreLogic("(take (list 2 4 6 8 10) 3)", rootEnv)
	if !checkLogicList(exp, []int{2, 4, 6}) {
		t.Fatal("failed test: take")
	}
	exp, _ = DoCoreLogic("(drop (list 2 4 6 8 10) 3)", rootEnv)
	if !checkLogicList(exp, []int{8, 10}) {
		t.Fatal("failed test: drop")
	}
	exp, _ = DoCoreLogic("(msort test-list)", rootEnv)
	if !checkLogicList(exp, testSortData) {
		t.Fatal("failed test: bsort")
	}
	fibonacci := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}
	exp, _ = DoCoreLogic("(inf-list (lambda (n) (list (cadr n)(+ (car n)(cadr n)))) (list 0 1) 10)", rootEnv)
	if !checkLogicList(exp, fibonacci) {
		t.Fatal("failed test: fibonacci")
	}

	exp, _ = DoCoreLogic("(fact/cps 5 (lambda (a) (+ 80 a)))", rootEnv)
	if !checkLogicInt(exp, 200) {
		t.Fatal("failed test: fact/cps")
	}
	exp, _ = DoCoreLogic("(fact 5)", rootEnv)
	if !checkLogicInt(exp, 120) {
		t.Fatal("failed test: fact")
	}
	exp, _ = DoCoreLogic("(fact/cont 5)", rootEnv)
	if !checkLogicInt(exp, 60) {
		t.Fatal("failed test: fact/cont")
	}
}
func TestMathFunc(t *testing.T) {
	var (
		exp Expression
	)
	BuildFunc()
	rootEnv := NewSimpleEnv(nil, nil)

	exp, _ = DoCoreLogic("(sqrt 9)", rootEnv)
	if exp.(*Float).Value != 3.0 {
		t.Fatal("failed test: sqrt")
	}
	exp, _ = DoCoreLogic("(cos (/ (* 60 (* (atan 1) 4))180))", rootEnv)
	if exp.(*Float).Value != 0.5000000000000001 {
		t.Fatal("failed test: cos")
	}
	exp, _ = DoCoreLogic("(sin (/ (* 30 (* (atan 1) 4)) 180))", rootEnv)
	if exp.(*Float).Value != 0.49999999999999994 {
		t.Fatal("failed test: sin")
	}
	exp, _ = DoCoreLogic("(tan (/ (* 45 (* (atan 1) 4)) 180))", rootEnv)
	if exp.(*Float).Value != 1.0 {
		t.Fatal("failed test: tan")
	}
	exp, _ = DoCoreLogic("(/ (log 8)(log 2))", rootEnv)
	if exp.(*Float).Value != 3.0 {
		t.Fatal("failed test: log")
	}
	exp, _ = DoCoreLogic("(exp (/ (log 8) 3))", rootEnv)
	if exp.(*Float).Value != 2.0 {
		t.Fatal("failed test: exp")
	}
	exp, _ = DoCoreLogic("(expt 2 0)", rootEnv)
	if exp.(*Integer).Value != 1 {
		t.Fatal("failed test: expt")
	}
	exp, _ = DoCoreLogic("(expt 2 1)", rootEnv)
	if exp.(*Integer).Value != 2 {
		t.Fatal("failed test: expt")
	}
	exp, _ = DoCoreLogic("(expt 2 3)", rootEnv)
	if exp.(*Integer).Value != 8 {
		t.Fatal("failed test: expt")
	}
}
func TestListFunc(t *testing.T) {
	var (
		exp Expression
	)
	BuildFunc()
	rootEnv := NewSimpleEnv(nil, nil)

	exp, _ = DoCoreLogic("(list 1 2 3)", rootEnv)
	if !checkLogicList(exp, []int{1, 2, 3}) {
		t.Fatal("failed test: list")
	}
	exp, _ = DoCoreLogic("(null? (list 1 2 3))", rootEnv)
	if (exp.(*Boolean)).Value != false {
		t.Fatal("failed test: null?")
	}
	exp, _ = DoCoreLogic("(null? (list))", rootEnv)
	if (exp.(*Boolean)).Value != true {
		t.Fatal("failed test: null?")
	}
	exp, _ = DoCoreLogic("(length (list 1 2 3 4))", rootEnv)
	if !checkLogicInt(exp, 4) {
		t.Fatal("failed test: length")
	}
	exp, _ = DoCoreLogic("(car (list 10 20 30 40))", rootEnv)
	if !checkLogicInt(exp, 10) {
		t.Fatal("failed test: car")
	}
	exp, _ = DoCoreLogic("(cdr (cons 10 20))", rootEnv)
	if !checkLogicInt(exp, 20) {
		t.Fatal("failed test: cdr")
	}
	exp, _ = DoCoreLogic("(cadr (list 1 2 3 4))", rootEnv)
	if !checkLogicInt(exp, 2) {
		t.Fatal("failed test: cadr")
	}
	exp, _ = DoCoreLogic("(car (cons 100 200))", rootEnv)
	if !checkLogicInt(exp, 100) {
		t.Fatal("failed test: cons")
	}
	exp, _ = DoCoreLogic("(cdr (cons 100 200))", rootEnv)
	if !checkLogicInt(exp, 200) {
		t.Fatal("failed test: cons")
	}
	exp, _ = DoCoreLogic("(append (list 1 2)(list 3 4))", rootEnv)
	if !checkLogicList(exp, []int{1, 2, 3, 4}) {
		t.Fatal("failed test: append")
	}
	exp, _ = DoCoreLogic("(reverse (list 1 2 3))", rootEnv)
	if !checkLogicList(exp, []int{3, 2, 1}) {
		t.Fatal("failed test: list")
	}
	exp, _ = DoCoreLogic("(iota 5 2)", rootEnv)
	if !checkLogicList(exp, []int{2, 3, 4, 5, 6}) {
		t.Fatal("failed test: iota")
	}
	exp, _ = DoCoreLogic("(map (lambda (n) (* n 10))(list 1 2 3))", rootEnv)
	if !checkLogicList(exp, []int{10, 20, 30}) {
		t.Fatal("failed test: map")
	}
	exp, _ = DoCoreLogic("(filter (lambda (n) (= n 1))(list 1 2 3))", rootEnv)
	if !checkLogicList(exp, []int{1}) {
		t.Fatal("failed test: filter")
	}
	exp, _ = DoCoreLogic("(reduce (lambda (a b) (+ a b))(list 1 2 3))", rootEnv)
	if !checkLogicInt(exp, 6) {
		t.Fatal("failed test: reduce")
	}
	exp, _ = DoCoreLogic("()", rootEnv)
	if len((exp.(*List)).Value) != 0 {
		t.Fatal("failed test: ()")
	}
	DoCoreLogic("(define cnt 0)", rootEnv)
	DoCoreLogic("(for-each (lambda (n) (set! cnt (+ cnt n)))(list 1 1 1 1 1))", rootEnv)
	exp, _ = DoCoreLogic("cnt", rootEnv)
	if !checkLogicInt(exp, 5) {
		t.Fatal("failed test: for-each")
	}
}
func TestBasicOperation(t *testing.T) {
	var (
		exp Expression
	)
	BuildFunc()
	rootEnv := NewSimpleEnv(nil, nil)

	exp, _ = DoCoreLogic("10", rootEnv)
	if !checkLogicInt(exp, 10) {
		t.Fatal("failed test: atom")
	}
	exp, _ = DoCoreLogic("10.5", rootEnv)
	if (exp.(*Float)).Value != 10.5 {
		t.Fatal("failed test: atom")
	}
	exp, _ = DoCoreLogic("\"ABC\"", rootEnv)
	if (exp.(*String)).Value != "ABC" {
		t.Fatal("failed test: atom")
	}
	exp, _ = DoCoreLogic("\"AB\\\"C\"", rootEnv)
	if (exp.(*String)).Value != "AB\\\"C" {
		t.Fatal("failed test: atom")
	}
	exp, _ = DoCoreLogic("\"(A B C)\"", rootEnv)
	if (exp.(*String)).Value != "(A B C)" {
		t.Fatal("failed test: atom")
	}
	exp, _ = DoCoreLogic("\"(\"", rootEnv)
	if (exp.(*String)).Value != "(" {
		t.Fatal("failed test: atom")
	}
	exp, _ = DoCoreLogic("\"  a  \"", rootEnv)
	if (exp.(*String)).Value != "  a  " {
		t.Fatal("failed test: atom")
	}
	exp, _ = DoCoreLogic("#t", rootEnv)
	if (exp.(*Boolean)).Value != true {
		t.Fatal("failed test: atom")
	}
	exp, _ = DoCoreLogic("#f", rootEnv)
	if (exp.(*Boolean)).Value != false {
		t.Fatal("failed test: atom")
	}
	exp, _ = DoCoreLogic("(+ 1 1.5 1.25)", rootEnv)
	if (exp.(*Float)).Value != 3.75 {
		t.Fatal("failed test", (exp.(*Float)).Value)
	}
	exp, _ = DoCoreLogic("(- 3 1.5 0.25)", rootEnv)
	if (exp.(*Float)).Value != 1.25 {
		t.Fatal("failed test", (exp.(*Float)).Value)
	}
	exp, _ = DoCoreLogic("(* 2 0.5 1.25)", rootEnv)
	if (exp.(*Float)).Value != 1.25 {
		t.Fatal("failed test", (exp.(*Float)).Value)
	}
	exp, _ = DoCoreLogic("(/ 3 0.5 2)", rootEnv)
	if (exp.(*Float)).Value != 3 {
		t.Fatal("failed test", (exp.(*Float)).Value)
	}
	exp, _ = DoCoreLogic("(modulo 18 12)", rootEnv)
	if (exp.(*Integer)).Value != 6 {
		t.Fatal("failed test", (exp.(*Integer)).Value)
	}
	exp, _ = DoCoreLogic("(> 3 0.5)", rootEnv)
	if (exp.(*Boolean)).Value != true {
		t.Fatal("failed test: >")
	}
	exp, _ = DoCoreLogic("(>= 3 0.5)", rootEnv)
	if (exp.(*Boolean)).Value != true {
		t.Fatal("failed test: >=")
	}
	exp, _ = DoCoreLogic("(>= 0.5 0.5)", rootEnv)
	if (exp.(*Boolean)).Value != true {
		t.Fatal("failed test: >=")
	}
	exp, _ = DoCoreLogic("(< 0.25 0.5)", rootEnv)
	if (exp.(*Boolean)).Value != true {
		t.Fatal("failed test: <")
	}
	exp, _ = DoCoreLogic("(<= 0.25 0.5)", rootEnv)
	if (exp.(*Boolean)).Value != true {
		t.Fatal("failed test: <=")
	}
	exp, _ = DoCoreLogic("(<= 0.5 0.5)", rootEnv)
	if (exp.(*Boolean)).Value != true {
		t.Fatal("failed test: <=")
	}
	exp, _ = DoCoreLogic("(= 0.75 0.75)", rootEnv)
	if (exp.(*Boolean)).Value != true {
		t.Fatal("failed test: =")
	}
	exp, _ = DoCoreLogic("(not (= 0.75 0.75))", rootEnv)
	if (exp.(*Boolean)).Value != false {
		t.Fatal("failed test: not")
	}
	exp, _ = DoCoreLogic("(let ((a 10)(b 20)(c 30)) (and (< a b)(< a c)(< a c)))", rootEnv)
	if (exp.(*Boolean)).Value != true {
		t.Fatal("failed test: and")
	}
	exp, _ = DoCoreLogic("(let ((a 10)(b 20)(c 30)) (and (< a b)(< a c)(< c a)))", rootEnv)
	if (exp.(*Boolean)).Value != false {
		t.Fatal("failed test: and")
	}
	exp, _ = DoCoreLogic("(let ((a 10)(b 20)(c 30)) (or (= a b)(< b c)))", rootEnv)
	if (exp.(*Boolean)).Value != true {
		t.Fatal("failed test: or")
	}
	exp, _ = DoCoreLogic("(let ((a 10)(b 20)(c 30)) (or (= c a)(< c b)))", rootEnv)
	if (exp.(*Boolean)).Value != false {
		t.Fatal("failed test: or")
	}
	exp, _ = DoCoreLogic("(force ((lambda (a) (delay (* 10 a))) 3))", rootEnv)
	if !checkLogicInt(exp, 30) {
		t.Fatal("failed test: force, delay")
	}
	exp, _ = DoCoreLogic("(identity 100)", rootEnv)
	if !checkLogicInt(exp, 100) {
		t.Fatal("failed test: identity")
	}
	exp, _ = DoCoreLogic("(identity \"ABC\")", rootEnv)
	if (exp.(*String)).Value != "ABC" {
		t.Fatal("failed test: identity")
	}
	exp, _ = DoCoreLogic("(* 3 (call/cc (lambda (k)  (- 1 (k 2)))))", rootEnv)
	if !checkLogicInt(exp, 6) {
		t.Fatal("failed test: call/cc")
	}
	DoCoreLogic("(define hoge (lambda (a b) a))", rootEnv)
	exp, _ = DoCoreLogic("(* 3 (call/cc (lambda (k)  (hoge 1 (k 2)))))", rootEnv)
	if !checkLogicInt(exp, 6) {
		t.Fatal("failed test: call/cc")
	}
	exp, _ = DoCoreLogic("(* 3 (let ((n 3)) (call/cc (lambda (k) (+ 1 (k n))))))", rootEnv)
	if !checkLogicInt(exp, 9) {
		t.Fatal("failed test: call/cc")
	}
	exp, _ = DoCoreLogic("(call/cc (lambda (k) 10))", rootEnv)
	if !checkLogicInt(exp, 10) {
		t.Fatal("failed test: call/cc")
	}

	exp, _ = DoCoreLogic("(call/cc (lambda (k) (map (lambda (n) (map (lambda (m) (if (= m 6)(k m) (+ n m))) (iota 10)))(iota 10))))", rootEnv)
	if !checkLogicInt(exp, 6) {
		t.Fatal("failed test: force, call/cc")
	}

	exp, _ = DoCoreLogic("(call/cc (lambda (k) (reduce (lambda (a b) (if (= a 3)(k a)(+ a b))) (list 1 2 3 4 5))))", rootEnv)
	if !checkLogicInt(exp, 3) {
		t.Fatal("failed test: call/cc")
	}
	exp, _ = DoCoreLogic("(define foo (lambda () (define hoge (lambda (a) (+ 1 a))) (hoge 10)))", rootEnv)
	exp, _ = DoCoreLogic("(foo)", rootEnv)
	if !checkLogicInt(exp, 11) {
		t.Fatal("failed test: nested define")
	}
	exp, _ = DoCoreLogic("(define a 100)", rootEnv)
	exp, _ = DoCoreLogic("a", rootEnv)
	if !checkLogicInt(exp, 100) {
		t.Fatal("failed test: simple define")
	}
	exp, _ = DoCoreLogic("(define\ta\t200)", rootEnv)
	exp, _ = DoCoreLogic("a", rootEnv)
	if !checkLogicInt(exp, 200) {
		t.Fatal("failed test: tab define")
	}
	exp, _ = DoCoreLogic("(define\na\n300)", rootEnv)
	exp, _ = DoCoreLogic("a", rootEnv)
	if !checkLogicInt(exp, 300) {
		t.Fatal("failed test: newline define")
	}
	exp, _ = DoCoreLogic("(define\r\na\r\n400)", rootEnv)
	exp, _ = DoCoreLogic("a", rootEnv)
	if !checkLogicInt(exp, 400) {
		t.Fatal("failed test: newline define")
	}
	exp, _ = DoCoreLogic("(let ((a 10)(b 10))(cond ((= a b) \"ok\")(else \"ng\")))", rootEnv)
	if (exp.(*String)).Value != "ok" {
		t.Fatal("failed test: cond")
	}
	exp, _ = DoCoreLogic("(let ((a 10)(b 20))(cond ((= a b) \"ok\")(else \"ng\")))", rootEnv)
	if (exp.(*String)).Value != "ng" {
		t.Fatal("failed test: cond")
	}
	exp, _ = DoCoreLogic("(let ((a 10)(b 20))(cond ((= a b) \"ok\")((= b 20) \"sankaku\")(else \"ng\")))", rootEnv)
	if (exp.(*String)).Value != "sankaku" {
		t.Fatal("failed test: cond")
	}
	exp, _ = DoCoreLogic("(let ((a 10)(b 20))(cond ((= a b) #t)))", rootEnv)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: NilClass")
	}
	exp, _ = DoCoreLogic("(let ((a 10)(b 20))(if (= a b) #t))", rootEnv)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: NilClass")
	}
	exp, _ = DoCoreLogic("(quote a)", rootEnv)
	if _, ok := exp.(*Symbol); !ok {
		t.Fatal("failed test: quote")
	}
	exp, _ = DoCoreLogic("(quote (a b c))", rootEnv)
	if _, ok := exp.(*List); !ok {
		t.Fatal("failed test: quote")
	}
	exp, _ = DoCoreLogic("(let ((a 10)(b 20))(+ a b)(* a b))", rootEnv)
	if !checkLogicInt(exp, 200) {
		t.Fatal("failed test: call/cc")
	}
	exp, _ = DoCoreLogic("#\\tab", rootEnv)
	if (exp.(*Char)).Value != 0x09 {
		t.Fatal("failed test: char tab")
	}
	exp, _ = DoCoreLogic("#\\space", rootEnv)
	if (exp.(*Char)).Value != 0x20 {
		t.Fatal("failed test: char space")
	}
	exp, _ = DoCoreLogic("#\\newline", rootEnv)
	if (exp.(*Char)).Value != 0x0A {
		t.Fatal("failed test: char newline")
	}
	exp, _ = DoCoreLogic("#\\return", rootEnv)
	if (exp.(*Char)).Value != 0x0D {
		t.Fatal("failed test: char return")
	}
	exp, _ = DoCoreLogic("#\\A", rootEnv)
	if (exp.(*Char)).Value != 0x41 {
		t.Fatal("failed test: char A")
	}
}
func TestErrCase(t *testing.T) {
	var (
		err error
	)
	BuildFunc()
	rootEnv := NewSimpleEnv(nil, nil)

	testCode := [][]string{
		{"(", "E0001"},
		{"(a (b", "E0002"},
		{"(a))", "E0003"},
		{"#\\abc", "E0004"},

		{"(+ 10.2)", "E1007"},
		{"(+ )", "E1007"},
		{"(+ #t 10.2)", "E1003"},

		{"(- 10.2 #f)", "E1003"},
		{"(- 10.2)", "E1007"},
		{"(-)", "E1007"},

		{"(* 10.2 #f)", "E1003"},
		{"(* 10.2)", "E1007"},
		{"(*)", "E1007"},

		{"(/ 10.2 #f)", "E1003"},
		{"(/ 10.2)", "E1007"},
		{"(/)", "E1007"},
		{"(/ 10 0)", "E1013"},
		{"(/ 10 2 0 3)", "E1013"},

		{"(modulo 1)", "E1007"},
		{"(modulo 10 3 2)", "E1007"},
		{"(modulo 10 2.5)", "E1002"},
		{"(modulo 10 0)", "E1013"},

		{"(> 10 3 2)", "E1007"},
		{"(> 10.2 #f)", "E1003"},
		{"(> 10.2)", "E1007"},

		{"(< 10.2 #f)", "E1003"},
		{"(< 10 3 2)", "E1007"},
		{"(< 10.2)", "E1007"},

		{"(>= 10 3 2)", "E1007"},
		{"(>= 10.2 #f)", "E1003"},
		{"(>= 10.2)", "E1007"},

		{"(<= 10.2 #f)", "E1003"},
		{"(<= 10 3 2)", "E1007"},
		{"(<= 10.2)", "E1007"},

		{"(= 10.2 #f)", "E1003"},
		{"(= 10 3 2)", "E1007"},
		{"(= 10.2)", "E1007"},

		{"(not 10)", "E1001"},
		{"(not #t #f)", "E1007"},
		{"(not)", "E1007"},

		{"((list 1 12) 10)", "E1006"},
		{"(null? 10)", "E1005"},
		{"(null? (list 1)(list 2))", "E1007"},
		{"(null?)", "E1007"},

		{"(length 10)", "E1005"},
		{"(length (list 1)(list 2))", "E1007"},
		{"(length)", "E1007"},

		{"(car 10)", "E1005"},
		{"(car (list 1)(list 2))", "E1007"},
		{"(car)", "E1007"},
		{"(car (list))", "E1011"},

		{"(cdr (list 1)(list 2))", "E1007"},
		{"(cdr)", "E1007"},
		{"(cdr 10)", "E1005"},

		{"(cadr (list 1)(list 2))", "E1007"},
		{"(cadr)", "E1007"},
		{"(cadr (list 1))", "E1011"},

		{"(cons 1 (list 1)(list 2))", "E1007"},
		{"(cons 1)", "E1007"},

		{"(append 10 10)", "E1005"},
		{"(append (list 1))", "E1007"},

		{"(last 10)", "E1005"},
		{"(last (list 1)(list 2))", "E1007"},
		{"(last)", "E1007"},
		{"(last (list))", "E1011"},

		{"(reverse 10)", "E1005"},
		{"(reverse (list 1)(list 2))", "E1007"},
		{"(reverse)", "E1007"},

		{"(iota 10.2 1)", "E1002"},
		{"(iota 1 10.2)", "E1002"},
		{"(iota)", "E1007"},
		{"(iota 1 2 3)", "E1007"},

		{"(map (lambda (n) (* n 10)) 20)", "E1005"},
		{"(map (list 1 12) (list 10))", "E1006"},
		{"(map (lambda (n) (* n 10)))", "E1007"},
		{"(map (lambda (n) (* n 10))(list 1)(list 1))", "E1007"},

		{"(filter (lambda (n) (* n 10)) 20)", "E1005"},
		{"(filter (list 1 12) (list 10))", "E1006"},
		{"(filter (lambda (n) (* n 10)))", "E1007"},
		{"(filter (lambda (n) (* n 10))(list 1)(list 1))", "E1007"},
		{"(filter (lambda (n) 10.1) (list 1 2))", "E1001"},

		{"(reduce (lambda (a b) (+ a b)) 20)", "E1005"},
		{"(reduce (list 1 12) (list 10))", "E1006"},
		{"(reduce (lambda (a b) (+ a b)))", "E1007"},
		{"(reduce (lambda (a b) (+ a b))(list 1)(list 1))", "E1007"},

		{"(sqrt #t)", "E1003"},
		{"(sqrt 11 10)", "E1007"},
		{"(sqrt)", "E1007"},

		{"(sin #t)", "E1003"},
		{"(sin 11 10)", "E1007"},
		{"(sin)", "E1007"},

		{"(cos #t)", "E1003"},
		{"(cos 11 10)", "E1007"},
		{"(cos)", "E1007"},

		{"(tan #t)", "E1003"},
		{"(tan 11 10)", "E1007"},
		{"(tan)", "E1007"},

		{"(atan #t)", "E1003"},
		{"(atan 11 10)", "E1007"},
		{"(atan)", "E1007"},

		{"(log #t)", "E1003"},
		{"(log 11 10)", "E1007"},
		{"(log)", "E1007"},

		{"(exp #t)", "E1003"},
		{"(exp 11 10)", "E1007"},
		{"(exp)", "E1007"},

		{"(rand-init 10.2)", "E1007"},
		{"(rand-integer 10.2)", "E1002"},
		{"(rand-integer)", "E1007"},
		{"(rand-integer 11 9)", "E1007"},

		{"(expt 10)", "E1007"},
		{"(expt 10 10 10)", "E1007"},
		{"(expt 11.5 10)", "E1002"},
		{"(expt 11 12.5)", "E1002"},

		{"(if 10 1 2)", "E1001"},
		{"(if (= 10 10))", "E1007"},

		{"(define 10 10)", "E1004"},
		{"(define a)", "E1007"},
		{"(define a 10 11)", "E1007"},
		{"(define (a))", "E1007"},
		{"(define 10 11)", "E1004"},

		{"(set! 10 10)", "E1004"},
		{"(set! a)", "E1007"},
		{"(set! a 10 11)", "E1007"},
		{"(set! hoge 10)", "E1008"},
		{"hoge", "E1008"},

		{"(lambda a (+ a b))", "E1005"},
		{"(lambda (+ n m))", "E1007"},
		{"(lambda 10 11)", "E1005"},
		{"((lambda (n m) (+ n m)) 1 2 3)", "E1007"},
		{"((lambda (n m) (+ n m)) 1)", "E1007"},

		{"(let ((a 10)))", "E1007"},
		{"(let 10 ((a 10)))", "E1004"},
		{"(let loop ((a 10)))", "E1007"},
		{"(let loop 10 19)", "E1005"},
		{"(let ((a))(+ 1 1))", "E1007"},
		{"(let ((a 10)) b a)", "E1008"},

		{"(and #t)", "E1007"},
		{"(and 10.2 0 1)", "E1001"},
		{"(or #t)", "E1007"},
		{"(or 10.2 0 1)", "E1001"},

		{"(delay)", "E1007"},
		{"(delay 1 2)", "E1007"},

		{"(force)", "E1007"},
		{"(force 1 2)", "E1007"},
		{"(force 1)", "E1010"},

		{"(identity 100 200)", "E1007"},
		{"(identity)", "E1007"},

		{"(call/cc)", "E1007"},
		{"(call/cc (lambda () #t))", "E1007"},
		{"(call/cc (lambda (n) #t)(lambda (n) #t))", "E1007"},
		{"(call/cc 10)", "E1006"},

		{"(cond)", "E1007"},
		{"(cond 10)", "E1005"},
		{"(cond (10))", "E1007"},
		{"(let ((a 10)(b 20))(cond ((= a b) #t)(lse #f)))", "E1012"},
		{"(cond (10 10))", "E1012"},

		{"(quote)", "E1007"},
		{"(quote 1 2)", "E1007"},

		{"(time)", "E1007"},
		{"(time #\\abc)", "E0004"},
		{"(let loop ((i 0)(j 10)(k 10)) (if (<= 1000000 i) i (if (= j k) (loop (+ 100 i)(+ 1 i)))))", "E1007"},
		{"(load-file)", "E1007"},
		{"(load-file 10)", "E1015"},
		{"(load-file \"example/no.scm\")", "E1014"},
	}
	for _, e := range testCode {
		_, err = DoCoreLogic(e[0], rootEnv)
		if !checkErrorCode(err, e[1]) {
			t.Fatal("failed test: " + e[0])
		}
	}
	// Impossible absolute, But Program bug is except
	_, err = eval(NewFunction(rootEnv, NewList(nil), nil, "lambda"), rootEnv)
	if !checkErrorCode(err, "E1009") {
		t.Fatal("failed test: " + "E1009")
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
	iostub("+", "Operatotion or Builtin:")
	iostub("if", "Special Functon ex. if:")
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
	if !checkLogicInt(exp, 120) {
		t.Fatal("failed test: tail recursive")
	}

	exp, _ = DoCoreLogic("(let loop ((i 0)) (if (<= 1000000 i) i (loop (+ 1 i))))", rootEnv)
	if !checkLogicInt(exp, 1000000) {
		t.Fatal("failed test: tail recursive")
	}

	exp, _ = DoCoreLogic("(let loop ((i 0)(j 10)(k 10)) (if (<= 1000000 i) i (if (= j k) (loop (+ 50 i) j k)(loop (+ 1 i) j k))))", rootEnv)
	if !checkLogicInt(exp, 1000000) {
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
