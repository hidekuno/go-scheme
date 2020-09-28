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

		// ex. SEGV, BUS ERROR..
		defer func() {
			if err := recover(); err != nil {
				t.Log(i)
				t.Fatal("failed " + c[0])
			}
		}()
		exp, err := DoCoreLogic(c[0], rootEnv)
		if err != nil {
			if e, ok := err.(*SyntaxError); ok {
				if e.MsgCode != c[1] {
					t.Log(c[0])
					t.Fatal("failed "+testName+" test", e.MsgCode)
				}
			} else if e, ok := err.(*RuntimeError); ok {
				if e.MsgCode != c[1] {
					t.Log(c[0])
					t.Fatal("failed "+testName+" test", e.MsgCode)
				}
			} else {
				t.Log(c[0])
				t.Fatal("failed "+testName+" test", err.Error())
			}
		} else {
			if exp.String() != c[1] {
				t.Log(c[0])
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
		{"(list 10))", "E0003"},
		{"(a))", "E1008"},
		{"1)", "E0003"},
		{"#\\abc", "E0004"},
		{"\"A", "E0004"},
		{"A\"", "E0004"},
		{"\"", "E0004"},
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
func TestDoCoreLogic(t *testing.T) {
	testCode := [][]string{
		{"(define a 100) a", "100"},
		{"(define a 100) a (+ a 100)", "200"},
	}
	executeTest(testCode, "do_core_logic", t)
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
