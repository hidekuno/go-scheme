/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package draw

import (
	"draw"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/gtk"
	"reflect"
	"scheme"
	"strconv"
)

var (
	execFinished = false
)

type Image struct {
	scheme.Expression
	Value *gdkpixbuf.Pixbuf
}

func NewImage(pixbuf *gdkpixbuf.Pixbuf) *Image {
	self := new(Image)
	self.Value = pixbuf
	return self
}
func (self *Image) String() string {
	return "Pixbuf: "
}
func (self *Image) Scale(w, h int) *Image {
	return NewImage(self.Value.ScaleSimple(w, h, gdkpixbuf.INTERP_HYPER))
}
func (self *Image) RotateSimple90() *Image {
	return NewImage(self.Value.RotateSimple(gdkpixbuf.PIXBUF_ROTATE_COUNTERCLOCKWISE))
}
func (self *Image) RotateSimple180() *Image {
	return NewImage(self.Value.RotateSimple(gdkpixbuf.PIXBUF_ROTATE_UPSIDEDOWN))
}
func (self *Image) RotateSimple270() *Image {
	return NewImage(self.Value.RotateSimple(gdkpixbuf.PIXBUF_ROTATE_CLOCKWISE))
}

func BuildGtkFunc() {
	scheme.AddErrorMsg("E2001", "Aleady Gtk Init")
	scheme.AddErrorMsg("E2002", "Cannot Read Image File")
	scheme.AddErrorMsg("E2003", "Not Image")

	draw_init := func(exp ...scheme.Expression) (scheme.Expression, error) {
		if len(exp) != 0 {
			return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		if execFinished == true {
			return nil, scheme.NewRuntimeError("E2001")
		}
		pixmap, gdkwin, fg, bg := draw.BuildGtkApp("scheme.go")
		go gtk.Main()

		//--------------------------------------------------------
		// Function for lisp
		//--------------------------------------------------------
		LispDrawLine := func(x0, y0, x1, y1 int) {
			gdk.ThreadsEnter()
			pixmap.GetDrawable().DrawLine(fg, x0, y0, x1, y1)
			gdkwin.Invalidate(nil, false)
			gdk.ThreadsLeave()
		}
		LispDrawClear := func() {
			gdk.ThreadsEnter()
			pixmap.GetDrawable().DrawRectangle(bg, true, 0, 0, -1, -1)
			gdkwin.Invalidate(nil, false)
			gdk.ThreadsLeave()
		}
		LispDrawImage := func(pixbuf *gdkpixbuf.Pixbuf, x, y int) {
			gdk.ThreadsEnter()
			pixmap.GetDrawable().DrawPixbuf(fg, pixbuf, 0, 0, x, y, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)
			gdkwin.Invalidate(nil, false)
			gdkwin.GetDrawable().DrawDrawable(fg, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
			gdk.ThreadsLeave()
		}

		scheme.AddBuiltInFunc("draw-clear", func(exp ...scheme.Expression) (scheme.Expression, error) {
			if len(exp) != 0 {
				return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			LispDrawClear()
			return scheme.NewNil(), nil
		})
		scheme.AddBuiltInFunc("draw-line", func(exp ...scheme.Expression) (scheme.Expression, error) {
			var point [4]int
			if len(exp) != 4 {
				return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			for i, e := range exp {
				if p, ok := e.(*scheme.Integer); ok {
					point[i] = p.Value
				} else if p, ok := e.(*scheme.Float); ok {
					point[i] = int(p.Value)
				} else {
					return nil, scheme.NewRuntimeError("E1003", reflect.TypeOf(e).String())
				}
			}
			LispDrawLine(point[0], point[1], point[2], point[3])
			return scheme.NewNil(), nil
		})
		scheme.AddBuiltInFunc("create-image-from-file", func(exp ...scheme.Expression) (scheme.Expression, error) {
			if len(exp) != 1 {
				return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			if s, ok := exp[0].(*scheme.String); ok {
				pixbuf, err := gdkpixbuf.NewPixbufFromFile(s.Value)
				if err != nil {
					return nil, scheme.NewRuntimeError("E2002")
				}
				return NewImage(pixbuf), nil
			}
			return nil, scheme.NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
		})
		scheme.AddBuiltInFunc("draw-image", func(exp ...scheme.Expression) (scheme.Expression, error) {
			if len(exp) != 3 {
				return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			img, ok := exp[0].(*Image)
			if !ok {
				return nil, scheme.NewRuntimeError("E2003", reflect.TypeOf(exp[0]).String())
			}
			x, ok := exp[1].(*scheme.Integer)
			if !ok {
				return nil, scheme.NewRuntimeError("E1002", reflect.TypeOf(exp[1]).String())
			}
			y, ok := exp[2].(*scheme.Integer)
			if !ok {
				return nil, scheme.NewRuntimeError("E1002", reflect.TypeOf(exp[2]).String())
			}
			LispDrawImage(img.Value, x.Value, y.Value)
			return scheme.NewNil(), nil

		})
		scheme.AddBuiltInFunc("scale-image", func(exp ...scheme.Expression) (scheme.Expression, error) {
			if len(exp) != 3 {
				return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			img, ok := exp[0].(*Image)
			if !ok {
				return nil, scheme.NewRuntimeError("E2003", reflect.TypeOf(exp[0]).String())
			}
			w, ok := exp[1].(*scheme.Integer)
			if !ok {
				return nil, scheme.NewRuntimeError("E1002", reflect.TypeOf(exp[1]).String())
			}
			h, ok := exp[2].(*scheme.Integer)
			if !ok {
				return nil, scheme.NewRuntimeError("E1002", reflect.TypeOf(exp[2]).String())
			}
			return img.Scale(w.Value, h.Value), nil
		})
		scheme.AddBuiltInFunc("rotate90-image", func(exp ...scheme.Expression) (scheme.Expression, error) {
			if len(exp) != 1 {
				return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			if img, ok := exp[0].(*Image); ok {
				return img.RotateSimple90(), nil
			}
			return nil, scheme.NewRuntimeError("E2003", reflect.TypeOf(exp[0]).String())
		})
		scheme.AddBuiltInFunc("rotate180-image", func(exp ...scheme.Expression) (scheme.Expression, error) {
			if len(exp) != 1 {
				return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			if img, ok := exp[0].(*Image); ok {
				return img.RotateSimple180(), nil
			}
			return nil, scheme.NewRuntimeError("E2003", reflect.TypeOf(exp[0]).String())
		})
		scheme.AddBuiltInFunc("rotate270-image", func(exp ...scheme.Expression) (scheme.Expression, error) {
			if len(exp) != 1 {
				return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			if img, ok := exp[0].(*Image); ok {
				return img.RotateSimple270(), nil
			}
			return nil, scheme.NewRuntimeError("E2003", reflect.TypeOf(exp[0]).String())
		})
		execFinished = true
		return scheme.NewNil(), nil
	}
	scheme.AddBuiltInFunc("draw-init", draw_init)
}
