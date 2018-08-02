/*
   Go lang 2nd study program.
   This is thing for go routine (multi threading program).

   hidekuno@gmail.com
*/
package fractal

import (
	"math"
)

func CreateTree(scale int, drawLine func(x0, y0, x1, y1 int)) func() {

	cs := math.Cos((math.Pi * 15) / 180)
	sn := math.Sin((math.Pi * 45) / 180)
	alpha := 0.6
	var draw func(x0, y0, x1, y1, c int)

	draw = func(x0, y0, x1, y1, c int) {

		drawLine(x0, y0, x1, y1)
		ya := y1 + int(sn*float64(x1-x0)*alpha+cs*float64(y1-y0)*alpha)
		xa := x1 + int(cs*float64(x1-x0)*alpha-sn*float64(y1-y0)*alpha)

		yb := y1 + int(-sn*float64(x1-x0)*alpha+cs*float64(y1-y0)*alpha)
		xb := x1 + int(cs*float64(x1-x0)*alpha+sn*float64(y1-y0)*alpha)

		if 0 >= c {
			drawLine(x1, y1, xa, ya)
			drawLine(x1, y1, xb, yb)
		} else {
			draw(x1, y1, xa, ya, c-1)
			draw(x1, y1, xb, yb, c-1)
		}
	}
	return func() { draw(300, 400, 300, 300, scale) }
}
