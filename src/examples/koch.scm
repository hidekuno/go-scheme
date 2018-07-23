;
; this is a sample program, and this is drawing koch
;
; hidekuno@gmail.com
;

(define kcos (cs 60))
(define ksin (sn 60))

(define koch (lambda (x0 y0 x1 y1 c)
  (if (> c 1)
      (let (
            (xa (/ (+ (* x0 2) x1) 3))
            (ya (/ (+ (* y0 2) y1) 3))
            (xb (/ (+ (* x1 2) x0) 3))
            (yb (/ (+ (* y1 2) y0) 3)))
        (let ((yc (+ ya (+ (* (- xb xa) ksin) (* (- yb ya) kcos))))
              (xc (+ xa (- (* (- xb xa) kcos) (* (- yb ya) ksin)))))
        (koch x0 y0 xa  ya (- c 1))
        (koch xa ya xc  yc (- c 1))
        (koch xc yc xb  yb (- c 1))
        (koch xb yb x1  y1 (- c 1))))
      (draw_line x0 y0 x1 y1))))

(koch 259 0 34 390 4)
(koch 34 390 483 390 4)
(koch 483 390 259 0 4)
