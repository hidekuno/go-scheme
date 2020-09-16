/*
   Go lang 2nd study program.
   This is thing for go routine (multi threading program).

   hidekuno@gmail.com
*/
package fractal

import (
	"os"
	"math"
)

func CreateKoch(scale int, drawLine func(x0, y0, x1, y1 int)) func() {

	cos60 := math.Cos((math.Pi * 60) / 180)
	sin60 := math.Sin((math.Pi * 60) / 180)
	var draw func(x0, y0, x1, y1, c int)

	if os.Getenv("DISPLAY") ==  "docker.for.mac.localhost:0" {
		scale = 4
	}
	draw = func(x0, y0, x1, y1, c int) {
		if c > 1 {

			xa := ((x0*2 + x1) / 3)
			ya := ((y0*2 + y1) / 3)
			xb := ((x1*2 + x0) / 3)
			yb := ((y1*2 + y0) / 3)

			yc := ya + int(float64(xb-xa)*sin60+float64(yb-ya)*cos60)
			xc := xa + int(float64(xb-xa)*cos60-float64(yb-ya)*sin60)

			draw(x0, y0, xa, ya, c-1)
			draw(xa, ya, xc, yc, c-1)
			draw(xc, yc, xb, yb, c-1)
			draw(xb, yb, x1, y1, c-1)

		} else {
			drawLine(x0, y0, x1, y1)
		}
	}
	return func() {
		draw(259, 0, 34, 390, scale)
		draw(34, 390, 483, 390, scale)
		draw(483, 390, 259, 0, scale)
	}
}
