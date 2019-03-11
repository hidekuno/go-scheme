/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package main

import (
	"draw/prototype"
	"github.com/therecipe/qt/widgets"
)

// Main
func main() {
	draw.BuildQtApp("draw-demo")
	widgets.QApplication_Exec()
}
