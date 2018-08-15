/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package main

import (
	"fmt"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"reflect"
	"runtime"
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
	return NewImage(self.Value.RotateSimple(gdkpixbuf.PIXBUF_ROTATE_NONE))
}

func buildGtkFunc() {
	errorMsg["E2001"] = "Aleady Gtk Init"
	errorMsg["E2002"] = "Cannot Read Image File"
	errorMsg["E2003"] = "Not Image"

	specialFuncTbl["draw-init"] = func(env *SimpleEnv, exps []Expression) (Expression, error) {
		if execFinished == true {
			return nil, NewRuntimeError("E2001")
		}
		go runDrawApp()

		builtinFuncTbl["draw-clear"] = func(exp ...Expression) (Expression, error) {
			drawClear()
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
			drawLineLisp(point[0], point[1], point[2], point[3])
			return NewNil(), nil
		}
		specialFuncTbl["draw-imagefile"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
			if len(v) != 1 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(v)))
			}
			if s, ok := v[0].(*String); ok {
				drawImageFile(s.Value)
			} else {
				return nil, NewRuntimeError("E1003", reflect.TypeOf(v[0]).String())
			}
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
			if len(exp) != 1 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			if img, ok := exp[0].(*Image); ok {
				drawImage(img.Value)
				return NewNil(), nil
			}
			return nil, NewRuntimeError("E2003", reflect.TypeOf(exp[0]).String())
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

// Main
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	buildFunc()
	buildGtkFunc()

	cui := make(chan bool)
	go func() {
		doInteractive()
		cui <- true
	}()
	<-cui
}
