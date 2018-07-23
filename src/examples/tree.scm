;
; this is a sample program, and this is drawing tree curve
;
; hidekuno@gmail.com
;

(define tcos (cs 15))
(define tsin (sn 45))
(define alpha 0.6)

(define (tree x0 y0 x1 y1 c)
  (let ((ya (+ y1  (*    tsin (- x1 x0) alpha) (*    tcos (- y1 y0) alpha)))
        (xa (+ x1  (*    tcos (- x1 x0) alpha) (* -1 tsin (- y1 y0) alpha)))
        (yb (+ y1  (* -1 tsin (- x1 x0) alpha) (*    tcos (- y1 y0) alpha)))
        (xb (+ x1  (*    tcos (- x1 x0) alpha) (*    tsin (- y1 y0) alpha))))
    (draw_line x0 y0 x1 y1)
    (if (>= 0 c)
        ((lambda () (draw_line x1 y1 xa ya) (draw_line x1 y1 xb yb)))
        ((lambda () (tree x1 y1 xa ya (- c 1))(tree x1  y1  xb  yb (- c 1)))))))

(tree 300 400 300 300 12)
