/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package main

import (
	"runtime"

	"local.packages/draw"
	"github.com/mattn/go-gtk/gtk"
)

// Main
func main() {
	runtime.GOMAXPROCS(1)
	draw.BuildGtkApp("draw-demo")
	gtk.Main()
}
