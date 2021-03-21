/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package draw

import (
	"reflect"
	"strconv"

	"github.com/hidekuno/go-scheme/draw"
	"github.com/hidekuno/go-scheme/scheme"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/gtk"
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

	drawInit := func(exp []scheme.Expression, env *scheme.SimpleEnv) (scheme.Expression, error) {
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
		lispDrawLine := func(x0, y0, x1, y1 int) {
			gdk.ThreadsEnter()
			pixmap.GetDrawable().DrawLine(fg, x0, y0, x1, y1)
			gdkwin.Invalidate(nil, false)
			gdk.ThreadsLeave()
		}
		lispDrawClear := func() {
			gdk.ThreadsEnter()
			pixmap.GetDrawable().DrawRectangle(bg, true, 0, 0, -1, -1)
			gdkwin.Invalidate(nil, false)
			gdk.ThreadsLeave()
		}
		lispDrawImage := func(pixbuf *gdkpixbuf.Pixbuf, x, y int) {
			gdk.ThreadsEnter()
			pixmap.GetDrawable().DrawPixbuf(fg, pixbuf, 0, 0, x, y, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)
			gdkwin.Invalidate(nil, false)
			gdkwin.GetDrawable().DrawDrawable(fg, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
			gdk.ThreadsLeave()
		}
		scheme.AddBuildInFunc("draw-clear", func(exp []scheme.Expression, env *scheme.SimpleEnv) (scheme.Expression, error) {
			if len(exp) != 0 {
				return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			lispDrawClear()
			return scheme.NewNil(), nil
		})
		scheme.AddBuildInFunc("draw-line", func(exp []scheme.Expression, env *scheme.SimpleEnv) (scheme.Expression, error) {
			var setParam = func(e scheme.Expression) (int, error) {
				if p, ok := e.(*scheme.Integer); ok {
					return p.Value, nil
				} else if p, ok := e.(*scheme.Float); ok {
					return int(p.Value), nil
				} else {
					return -1, scheme.NewRuntimeError("E1003", reflect.TypeOf(e).String())
				}
			}
			var err error
			return scheme.EvalCalcParam(exp, env,
				func(exp ...scheme.Expression) (scheme.Expression, error) {
					var point [4]int
					if len(exp) != 2 && len(exp) != 4 {
						return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
					}
					if len(exp) == 2 {
						i := 0
						for _, e := range exp {
							if p, ok := e.(*scheme.Pair); ok {
								if point[i], err = setParam(p.Car); err != nil {
									return nil, err
								}
								if point[i+1], err = setParam(p.Cdr); err != nil {
									return nil, err
								}
							} else {
								return nil, scheme.NewRuntimeError("E1005", reflect.TypeOf(e).String())
							}
							i = i + 2
						}
					} else {
						for i, e := range exp {
							if point[i], err = setParam(e); err != nil {
								return nil, err
							}
						}
					}
					lispDrawLine(point[0], point[1], point[2], point[3])
					return scheme.NewNil(), nil
				})
		})
		scheme.AddBuildInFunc("create-image-from-file", func(exp []scheme.Expression, env *scheme.SimpleEnv) (scheme.Expression, error) {
			return scheme.EvalCalcParam(exp, env,
				func(exp ...scheme.Expression) (scheme.Expression, error) {
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
		})
		scheme.AddBuildInFunc("draw-image", func(exp []scheme.Expression, env *scheme.SimpleEnv) (scheme.Expression, error) {
			return scheme.EvalCalcParam(exp, env,
				func(exp ...scheme.Expression) (scheme.Expression, error) {
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
					lispDrawImage(img.Value, x.Value, y.Value)
					return scheme.NewNil(), nil
				})
		})
		scheme.AddBuildInFunc("scale-image", func(exp []scheme.Expression, env *scheme.SimpleEnv) (scheme.Expression, error) {
			return scheme.EvalCalcParam(exp, env,
				func(exp ...scheme.Expression) (scheme.Expression, error) {
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
		})
		scheme.AddBuildInFunc("rotate90-image", func(exp []scheme.Expression, env *scheme.SimpleEnv) (scheme.Expression, error) {
			return scheme.EvalCalcParam(exp, env,
				func(exp ...scheme.Expression) (scheme.Expression, error) {
					if len(exp) != 1 {
						return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
					}
					if img, ok := exp[0].(*Image); ok {
						return img.RotateSimple90(), nil
					}
					return nil, scheme.NewRuntimeError("E2003", reflect.TypeOf(exp[0]).String())
				})
		})
		scheme.AddBuildInFunc("rotate180-image", func(exp []scheme.Expression, env *scheme.SimpleEnv) (scheme.Expression, error) {
			return scheme.EvalCalcParam(exp, env,
				func(exp ...scheme.Expression) (scheme.Expression, error) {
					if len(exp) != 1 {
						return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
					}
					if img, ok := exp[0].(*Image); ok {
						return img.RotateSimple180(), nil
					}
					return nil, scheme.NewRuntimeError("E2003", reflect.TypeOf(exp[0]).String())
				})
		})
		scheme.AddBuildInFunc("rotate270-image", func(exp []scheme.Expression, env *scheme.SimpleEnv) (scheme.Expression, error) {
			return scheme.EvalCalcParam(exp, env,
				func(exp ...scheme.Expression) (scheme.Expression, error) {
					if len(exp) != 1 {
						return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
					}
					if img, ok := exp[0].(*Image); ok {
						return img.RotateSimple270(), nil
					}
					return nil, scheme.NewRuntimeError("E2003", reflect.TypeOf(exp[0]).String())
				})
		})
		scheme.AddBuildInFunc("gtk-major-version", func(exp []scheme.Expression, env *scheme.SimpleEnv) (scheme.Expression, error) {
			if len(exp) != 0 {
				return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			return scheme.NewInteger(int(gtk.MajorVersion())), nil
		})
		scheme.AddBuildInFunc("gtk-minor-version", func(exp []scheme.Expression, env *scheme.SimpleEnv) (scheme.Expression, error) {
			if len(exp) != 0 {
				return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			return scheme.NewInteger(int(gtk.MinorVersion())), nil
		})
		scheme.AddBuildInFunc("gtk-micro-version", func(exp []scheme.Expression, env *scheme.SimpleEnv) (scheme.Expression, error) {
			if len(exp) != 0 {
				return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			return scheme.NewInteger(int(gtk.MicroVersion())), nil
		})
		scheme.AddBuildInFunc("screen-width", func(exp []scheme.Expression, env *scheme.SimpleEnv) (scheme.Expression, error) {
			if len(exp) != 0 {
				return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			return scheme.NewInteger(draw.ScreenWidth - 2), nil
		})
		scheme.AddBuildInFunc("screen-height", func(exp []scheme.Expression, env *scheme.SimpleEnv) (scheme.Expression, error) {
			if len(exp) != 0 {
				return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			return scheme.NewInteger(draw.ScreenHeight - 2), nil
		})

		execFinished = true
		return scheme.NewNil(), nil
	}
	scheme.AddBuildInFunc("draw-init", drawInit)
}
