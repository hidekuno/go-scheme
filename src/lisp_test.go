package main

import (
	"testing"
)

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

func TestLisp(t *testing.T) {
	var (
		exp Expression
	)
	build_func()
	root_env := NewSimpleEnv(nil,nil)
	program := []string {
		"(define counter (lambda () (let ((c 0)) (lambda () (set! c (+ 1 c)) c))))",
		"(define a (counter))",
		"(define b (counter))",		
		"(define gcm (lambda (n m) (let ((mod (modulo n m))) (if (= 0 mod) m (gcm m mod)))))",
		"(define lcm (lambda (n m) (/(* n m)(gcm n m))))",
		"(define hanoi (lambda (from to work n) (if (>= 0 n) (list) (append (hanoi from work to (- n 1)) (list (list (cons from to) n)) (hanoi work to from (- n 1))))))",
		"(define prime (lambda (l) (if (> (car l)(sqrt (last l))) l (cons (car l)(prime (filter (lambda (n) (not (= 0 (modulo n (car l))))) (cdr l)))))))",
		"(define qsort (lambda (l pred) (if (null? l) l (append (qsort (filter (lambda (n) (pred n (car l))) (cdr l)) pred) (cons (car l) (qsort (filter (lambda (n) (not (pred n (car l))))(cdr l)) pred))))))",
		"(define comb (lambda (l n) (if (null? l) l (if (= n 1) (map (lambda (n) (list n)) l) (append (map (lambda (p) (cons (car l) p)) (comb (cdr l)(- n 1))) (comb (cdr l) n))))))",
		"(define delete (lambda (x l) (filter (lambda (n) (not (= x n))) l)))",
		"(define perm (lambda (l n)(if (>= 0 n) (list (list))(reduce (lambda (a b)(append a b))(map (lambda (x) (map (lambda (p) (cons x p)) (perm (delete x l)(- n 1)))) l)))))",
		"(define bubble-iter (lambda (x l)(if (or (null? l)(< x (car l)))(cons x l)(cons (car l)(bubble-iter x (cdr l))))))",
		"(define bsort (lambda (l)(if (null? l) l (bubble-iter (car l)(bsort (cdr l))))))" }

	for _, p := range program {
		exp, _ = do_core_logic(p,root_env)
	}
	exp, _ = do_core_logic("(let loop ((a 0)(r (list 1 2 3))) (if (null? r) a (loop (+ (car r) a)(cdr r))))", root_env)
	if (!check_logic_int(exp,6)) {
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
}
