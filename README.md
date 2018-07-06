Go言語によるScheme(subset版)の実装
=================

# 動作確認したプログラム
```
(let loop ((a 0)(r (list 1 2 3))) (if (null? r) a (loop (+ (car r) a)(cdr r))))
(define counter (lambda () (let ((c 0)) (lambda () (set! c (+ 1 c))))))
(define gcm (lambda (n m) (let ((mod (modulo n m))) (if (= 0 mod) m (gcm m mod)))))
(define lcm (lambda (n m) (/(* n m)(gcm n m))))
(define hanoi (lambda (from to work n) (if (>= 0 n) (list) (append (hanoi from work to (- n 1)) (list (list (cons from to) n)) (hanoi work to from (- n 1))))))
(define prime (lambda (l) (if (> (car l)(sqrt (last l))) l (cons (car l)(prime (filter (lambda (n) (not (= 0 (modulo n (car l))))) (cdr l)))))))
(define qsort (lambda (l pred) (if (null? l) l (append (qsort (filter (lambda (n) (pred n (car l))) (cdr l)) pred) (cons (car l) (qsort (filter (lambda (n) (not (pred n (car l))))(cdr l)) pred))))))
(define comb (lambda (l n) (if (null? l) l (if (= n 1) (map (lambda (n) (list n)) l) (append (map (lambda (p) (cons (car l) p)) (comb (cdr l)(- n 1))) (comb (cdr l) n))))))
(define delete (lambda (x l) (filter (lambda (n) (not (= x n))) l)))
(define perm (lambda (l n)(if (>= 0 n) (list (list))(reduce (lambda (a b)(append a b))(map (lambda (x) (map (lambda (p) (cons x p)) (perm (delete x l)(- n 1)))) l)))))
(define bubble-iter (lambda (x l)(if (or (null? l)(< x (car l)))(cons x l)(cons (car l)(bubble-iter x (cdr l))))))
(define bsort (lambda (l)(if (null? l) l (bubble-iter (car l)(bsort (cdr l))))))
```
