/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package main

import (
	"runtime"
	"scheme"
	"scheme/draw"
)

// Main
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	scheme.BuildFunc()
	draw.BuildGtkFunc()
	scheme.DoInteractive()
}
