(define frame (make-frame (make-vect 0 0) (make-vect 1 0) (make-vect 0 1)))
(define segments (list
                     (make-segment(make-point 0.35 0.15) (make-point 0.4 0))
		     (make-segment(make-point 0.65 0.15) (make-point 0.6 0))
		     (make-segment(make-point 0.35 0.15) (make-point 0.4 0.35))
		     (make-segment(make-point 0.65 0.15) (make-point 0.6 0.35))
		     (make-segment(make-point 0.6 0.35)  (make-point 0.75 0.35))
		     (make-segment(make-point 0.4 0.35)  (make-point 0.3 0.35))
		     (make-segment(make-point 0.75 0.35) (make-point 1 0.65))
		     (make-segment(make-point 0.6 0.55)  (make-point 1 0.85))
		     (make-segment(make-point 0.6 0.55)  (make-point 0.75 1))
		     (make-segment(make-point 0.5 0.7)   (make-point 0.6 1))
		     (make-segment(make-point 0.3 0.35)  (make-point 0.15 0.4))
		     (make-segment(make-point 0.3 0.4)   (make-point 0.15 0.6))
		     (make-segment(make-point 0.15 0.4)  (make-point 0 0.15))
		     (make-segment(make-point 0.15 0.6)  (make-point 0 0.35))
		     (make-segment(make-point 0.3 0.4)   (make-point 0.35 0.5))
		     (make-segment(make-point 0.35 0.5)  (make-point 0.25 1))
		     (make-segment(make-point 0.5 0.7)   (make-point 0.4 1))
                  ))
(define wave (segments->painter segments))
(define draw-line-from-point (lambda (p1 p2) (draw-line (* (x-point p1) 600)(* (y-point p1) 400)(* (x-point p2) 600)(* (y-point p2) 400))))
((square-limit wave 4) frame)
