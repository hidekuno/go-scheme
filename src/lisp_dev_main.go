/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package main

/*
ex.1
(time (append (list (let loop ((i 0)) (if (<= 300000 i) i (loop (+ 1 i))))) (list (let loop ((i 0)) (if (<= 300000 i) i (loop (+ 1 i)))))))
(time (go-append (list (let loop ((i 0)) (if (<= 300000 i) i (loop (+ 1 i))))) (list (let loop ((i 0)) (if (<= 300000 i) i (loop (+ 1 i)))))))

ex.2
(time (append (perm (iota 8) 8)(perm (iota 8) 8)))
(time (go-append (perm (iota 8) 8)(perm (iota 8) 8)))
*/

import (
	"runtime"
	"scheme"
	"scheme/experiment"
)

// Main
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	scheme.BuildFunc()
	experiment.BuildGoFunc()
	scheme.DoInteractive()
}
