;
; this is a sample program, and this is drawing koch
;
; hidekuno@gmail.com
;
(define (koch-demo n)
  (begin
    (koch 259 0 34 390 n)
    (koch 34 390 483 390 n)
    (koch 483 390 259 0 n)))

(define (sierpinski-demo n)
  (sierpinski 319 40 30 430 609 430 n))

(define (tree-demo n)
  (tree 300 400 300 300 n))
