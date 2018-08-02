/*
   Go lang 2nd study program.
   This is thing for go routine (multi threading program).

   hidekuno@gmail.com
*/
package fractal

func CreateSierpinski(scale int, drawLine func(x0, y0, x1, y1 int)) func() {

	var draw func(x0, y0, x1, y1, x2, y2, c int)

	draw = func(x0, y0, x1, y1, x2, y2, c int) {
		if c > 1 {
			xx0 := (int)((x0 + x1) / 2.0)
			yy0 := (int)((y0 + y1) / 2.0)
			xx1 := (int)((x1 + x2) / 2.0)
			yy1 := (int)((y1 + y2) / 2.0)
			xx2 := (int)((x2 + x0) / 2.0)
			yy2 := (int)((y2 + y0) / 2.0)

			draw(x0, y0, xx0, yy0, xx2, yy2, c-1)
			draw(x1, y1, xx0, yy0, xx1, yy1, c-1)
			draw(x2, y2, xx2, yy2, xx1, yy1, c-1)

		} else {
			drawLine(x0, y0, x1, y1)
			drawLine(x1, y1, x2, y2)
			drawLine(x2, y2, x0, y0)
		}
	}
	return func() {
		draw(319, 40, 30, 430, 609, 430, scale)
	}
}
