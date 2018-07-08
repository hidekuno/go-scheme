/*
   Go lang 4th study program.
   This is test program for mini scheme subset.

   ex.) go test lisp.go lisp_test.go

   hidekuno@gmail.com
*/
package main

import (
	"testing"
)
func check_hanoi(exp Expression) bool {
	if l, ok := exp.(*List); ok {
		if len(l.Value) == 7 {
			return true
		}
	}
	return false
}
func check_logic_matrix(exp Expression, items [][]int) bool {
	if l, ok := exp.(*List); ok {
		for i,e := range l.Value {
			if _, ok := e.(*List); !ok {
				return false
			}
			for j,r := range (e.(*List)).Value {
				if !check_logic_int(r, items[i][j]) {
					return false
				}
			}
		}
		return true
	} else {
		return false
	}
}
func check_logic_list(exp Expression, items []int) bool {

	if l, ok := exp.(*List); ok {
		for i,e := range l.Value {
			if !check_logic_int(e, items[i]) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func check_logic_int(exp Expression, v int) bool {
	if i, ok := exp.(*Integer); ok {
		if i.Value != v {
			return false
		}
	} else {
		return false
	}
	return true
}
var (
	program = []string {
		"(define test_list (list 36 27 14 19 2 8 7 6 0 9 3))",
		"(define counter (lambda () (let ((c 0)) (lambda () (set! c (+ 1 c)) c))))",
		"(define a (counter))",
		"(define b (counter))",
		"(define gcm (lambda (n m) (let ((mod (modulo n m))) (if (= 0 mod) m (gcm m mod)))))",
		"(define lcm (lambda (n m) (/(* n m)(gcm n m))))",
		"(define hanoi (lambda (from to work n) (if (>= 0 n)(list)(append (hanoi from work to (- n 1))(list (list (cons from to) n))(hanoi work to from (- n 1))))))",
		"(define prime (lambda (l) (if (> (car l)(sqrt (last l))) l (cons (car l)(prime (filter (lambda (n) (not (= 0 (modulo n (car l))))) (cdr l)))))))",
		"(define qsort (lambda (l pred) (if (null? l) l (append (qsort (filter (lambda (n) (pred n (car l))) (cdr l)) pred) (cons (car l) (qsort (filter (lambda (n) (not (pred n (car l))))(cdr l)) pred))))))",
		"(define comb (lambda (l n) (if (null? l) l (if (= n 1) (map (lambda (n) (list n)) l) (append (map (lambda (p) (cons (car l) p)) (comb (cdr l)(- n 1))) (comb (cdr l) n))))))",
		"(define delete (lambda (x l) (filter (lambda (n) (not (= x n))) l)))",
		"(define perm (lambda (l n)(if (>= 0 n) (list (list))(reduce (lambda (a b)(append a b))(map (lambda (x) (map (lambda (p) (cons x p)) (perm (delete x l)(- n 1)))) l)))))",
		"(define bubble-iter (lambda (x l)(if (or (null? l)(< x (car l)))(cons x l)(cons (car l)(bubble-iter x (cdr l))))))",
		"(define bsort (lambda (l)(if (null? l) l (bubble-iter (car l)(bsort (cdr l))))))",
	}
)
func Test_do_core_logic(t *testing.T) {
	var (
		exp Expression
	)
	build_func()
	root_env := NewSimpleEnv(nil,nil)
	for _, p := range program {
		exp, _ = do_core_logic(p,root_env)
	}

	exp, _ = do_core_logic("(let loop ((a 0)(r (list 1 2 3))) (if (null? r) a (loop (+ (car r) a)(cdr r))))", root_env)
	if (!check_logic_int(exp,6)) {
		t.Fatal("failed test")
	}
	exp, _ = do_core_logic("(a)",root_env)
	exp, _ = do_core_logic("(a)",root_env)
	exp, _ = do_core_logic("(a)",root_env)
	if (!check_logic_int(exp,3)) {
		t.Fatal("failed test")
	}
	exp, _ = do_core_logic("(b)",root_env)
	exp, _ = do_core_logic("(b)",root_env)
	if (!check_logic_int(exp,2)) {
		t.Fatal("failed test")
	}
	exp, _ = do_core_logic("(gcm 36 27)",root_env)
	if (!check_logic_int(exp,9)) {
		t.Fatal("failed test")
	}
	exp, _ = do_core_logic("(lcm 36 27)",root_env)
	if (!check_logic_int(exp,108)) {
		t.Fatal("failed test")
	}

	test_sort_data := []int {0, 2, 3, 6, 7, 8, 9, 14, 19, 27, 36}
	exp, _ = do_core_logic("(qsort test_list (lambda (a b)(< a b)))",root_env)
	if (!check_logic_list(exp, test_sort_data)) {
		t.Fatal("failed test")
	}
	exp, _ = do_core_logic("(bsort test_list)",root_env)
	if (!check_logic_list(exp, test_sort_data)) {
		t.Fatal("failed test")
	}

	prime_data := []int {2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31}
	exp, _ = do_core_logic("(prime (iota 30 2))", root_env)
	if (!check_logic_list(exp, prime_data)) {
		t.Fatal("failed test")
	}

	perm_data := [][]int{{1, 2}, {1, 3}, {2, 1}, {2, 3},{3, 1}, {3, 2}}
	exp, _ = do_core_logic("(perm (list 1 2 3) 2)", root_env)
	if (!check_logic_matrix(exp, perm_data)) {
		t.Fatal("failed test")
	}

	comb := [][]int{{1, 2}, {1, 3}, {2, 3}}
	exp, _ = do_core_logic("(comb (list 1 2 3) 2)", root_env)
	if (!check_logic_matrix(exp, comb)) {
		t.Fatal("failed test")
	}
	exp, _ = do_core_logic("(hanoi \"a\" \"b\" \"c\" 3)", root_env)
	if (!check_hanoi(exp)){
		t.Fatal("failed test")
	}
}
func Test_basic_opration(t *testing.T) {
	var (
		exp Expression
	)
	build_func()
	root_env := NewSimpleEnv(nil,nil)
	exp, _ = do_core_logic("(+ 1 1.5)",root_env)
	if (exp.(*Float)).Value != 2.5 {
		t.Fatal("failed test", (exp.(*Float)).Value)
	}
	exp,_ = do_core_logic("(let ((a 10)(b 20)(c 30)) (and (< a b)(< a c)(< a c)))",root_env)
	if (exp.(*Boolean)).Value != true {
		t.Fatal("failed test")
	}
	exp,_ = do_core_logic("(let ((a 10)(b 20)(c 30)) (or (= a b)(< b c)))",root_env)
	if (exp.(*Boolean)).Value != true {
		t.Fatal("failed test")
	}
}
