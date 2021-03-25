/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package main

import (
	"runtime"
	"local.packages/draw"
	"github.com/hidekuno/go-scheme/scheme"
)

// Main
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	scheme.BuildFunc()
	draw.BuildGtkFunc()
	scheme.DoInteractive()
}
