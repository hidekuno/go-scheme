/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

import (
	"draw"
	"fmt"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/gtk"
	"reflect"
	"strconv"
)

var (
	execFinished = false
)

type Image struct {
	Expression
	Value *gdkpixbuf.Pixbuf
}

func NewImage(pixbuf *gdkpixbuf.Pixbuf) *Image {
	self := new(Image)
	self.Value = pixbuf
	return self
}
func (self *Image) Print() {
	fmt.Print("Pixbuf: ", self)
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
	errorMsg["E2001"] = "Aleady Gtk Init"
	errorMsg["E2002"] = "Cannot Read Image File"
	errorMsg["E2003"] = "Not Image"

	specialFuncTbl["draw-init"] = func(env *SimpleEnv, exps []Expression) (Expression, error) {
		if execFinished == true {
			return nil, NewRuntimeError("E2001")
		}
		pixmap, gdkwin, fg, bg := draw.BuildGtkApp()
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

		builtinFuncTbl["draw-clear"] = func(exp ...Expression) (Expression, error) {
			LispDrawClear()
			return NewNil(), nil
		}
		builtinFuncTbl["draw-line"] = func(exp ...Expression) (Expression, error) {
			var point [4]int
			if len(exp) != 4 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			for i, e := range exp {
				if p, ok := e.(*Integer); ok {
					point[i] = p.Value
				} else if p, ok := e.(*Float); ok {
					point[i] = int(p.Value)
				} else {
					return nil, NewRuntimeError("E1003", reflect.TypeOf(e).String())
				}
			}
			LispDrawLine(point[0], point[1], point[2], point[3])
			return NewNil(), nil
		}
		builtinFuncTbl["create-image-from-file"] = func(exp ...Expression) (Expression, error) {
			if len(exp) != 1 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			if s, ok := exp[0].(*String); ok {
				pixbuf, err := gdkpixbuf.NewPixbufFromFile(s.Value)
				if err != nil {
					return nil, NewRuntimeError("E2002")
				}
				return NewImage(pixbuf), nil
			}
			return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
		}
		builtinFuncTbl["draw-image"] = func(exp ...Expression) (Expression, error) {
			if len(exp) != 3 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			img, ok := exp[0].(*Image)
			if !ok {
				return nil, NewRuntimeError("E2003", reflect.TypeOf(exp[0]).String())
			}
			x, ok := exp[1].(*Integer)
			if !ok {
				return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[1]).String())
			}
			y, ok := exp[2].(*Integer)
			if !ok {
				return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[2]).String())
			}
			LispDrawImage(img.Value, x.Value, y.Value)
			return NewNil(), nil

		}
		builtinFuncTbl["scale-image"] = func(exp ...Expression) (Expression, error) {
			if len(exp) != 3 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			img, ok := exp[0].(*Image)
			if !ok {
				return nil, NewRuntimeError("E2003", reflect.TypeOf(exp[0]).String())
			}
			w, ok := exp[1].(*Integer)
			if !ok {
				return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[1]).String())
			}
			h, ok := exp[2].(*Integer)
			if !ok {
				return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[2]).String())
			}
			return img.Scale(w.Value, h.Value), nil
		}
		builtinFuncTbl["rotate90-image"] = func(exp ...Expression) (Expression, error) {
			if len(exp) != 1 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			if img, ok := exp[0].(*Image); ok {
				return img.RotateSimple90(), nil
			}
			return nil, NewRuntimeError("E2003", reflect.TypeOf(exp[0]).String())
		}
		builtinFuncTbl["rotate180-image"] = func(exp ...Expression) (Expression, error) {
			if len(exp) != 1 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			if img, ok := exp[0].(*Image); ok {
				return img.RotateSimple180(), nil
			}
			return nil, NewRuntimeError("E2003", reflect.TypeOf(exp[0]).String())
		}
		builtinFuncTbl["rotate270-image"] = func(exp ...Expression) (Expression, error) {
			if len(exp) != 1 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			if img, ok := exp[0].(*Image); ok {
				return img.RotateSimple270(), nil
			}
			return nil, NewRuntimeError("E2003", reflect.TypeOf(exp[0]).String())
		}
		execFinished = true
		return NewNil(), nil
	}
}
