/*
   Go lang 2nd study program.
   This is thing for go routine (multi threading program).

   +------------------------------------------+
   | go get github.com/mattn/go-gtk/gtk
   | go install github.com/mattn/go-gtk/gtk
   +------------------------------------------+

   hidekuno@gmail.com
*/
package main

import (
	"./fractal"
	"fmt"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/gtk"
	_ "runtime"
	"time"
)

const (
	KOCH_MAX       = 12
	TREE_MAX       = 20
	SIERPINSKI_MAX = 16
)

var draw_line_reentrant_lisp func(x0, y0, x1, y1 int)
var draw_clear func()
var draw_imagefile func(filename string)

func build_gtk_app() {

	var (
		pixmap *gdk.Pixmap
		fg_gc  *gdk.GC
		bg_gc  *gdk.GC
		gdkwin *gdk.Window
	)
	win := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("scheme.go")
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.Connect("destroy", gtk.MainQuit)

	//--------------------------------------------------------
	// GtkVBox
	//--------------------------------------------------------
	vbox := gtk.NewVBox(false, 1)
	vbox.SetBorderWidth(5)

	//--------------------------------------------------------
	// GtkMenuBar
	//--------------------------------------------------------
	draw_line_reentrant_lisp = func(x0, y0, x1, y1 int) {
		gdk.ThreadsEnter()
		pixmap.GetDrawable().DrawLine(fg_gc, x0, y0, x1, y1)
		gdkwin.Invalidate(nil, false)
		gdk.ThreadsLeave()
	}
	draw_clear = func() {
		pixmap.GetDrawable().DrawRectangle(bg_gc, true, 0, 0, -1, -1)
		gdkwin.Invalidate(nil, false)
	}
	draw_imagefile = func(filename string) {
		pixbuf, err := gdkpixbuf.NewPixbufFromFile(filename)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		pixmap.GetDrawable().DrawPixbuf(fg_gc, pixbuf, 0, 0, 0, 0, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)
		gdkwin.Invalidate(nil, false)
		gdkwin.GetDrawable().DrawDrawable(fg_gc, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
	}
	var draw_line_reentrant = func(x0, y0, x1, y1 int) {
		gdk.ThreadsEnter()
		pixmap.GetDrawable().DrawLine(fg_gc, x0, y0, x1, y1)
		gdk.ThreadsLeave()
	}
	var draw_reentrant = func(paint func()) {
		pixmap.GetDrawable().DrawRectangle(bg_gc, true, 0, 0, -1, -1)
		gdkwin.Invalidate(nil, false)

		stop := make(chan bool)
		go func() {
			stop <- false

			t0 := time.Now()
			paint()
			t1 := time.Now()
			fmt.Println(t1.Sub(t0))
			stop <- true
		}()
		go func() {
			defer close(stop)

			v := false
			for {
				time.Sleep(50 * time.Millisecond)
				select {
				case v = <-stop:
				default:
					v = false
				}
				gdkwin.Invalidate(nil, false)
				if v == true {
					break
				}
			}
		}()
	}

	var draw_line_single = func(x0, y0, x1, y1 int) { pixmap.GetDrawable().DrawLine(fg_gc, x0, y0, x1, y1) }
	var draw_single = func(paint func()) {
		pixmap.GetDrawable().DrawRectangle(bg_gc, true, 0, 0, -1, -1)
		paint()
		gdkwin.GetDrawable().DrawDrawable(fg_gc, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
	}

	menubar := gtk.NewMenuBar()
	cascademenu := gtk.NewMenuItemWithMnemonic("_File")
	menubar.Append(cascademenu)
	submenu := gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	var menuitem *gtk.MenuItem
	menuitem = gtk.NewMenuItemWithMnemonic("E_xit")
	menuitem.Connect("activate", func() {
		gtk.MainQuit()
	})
	submenu.Append(menuitem)

	cascademenu = gtk.NewMenuItemWithMnemonic("_Fractal")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	koch := fractal.CreateKoch(KOCH_MAX, draw_line_reentrant)
	menuitem = gtk.NewMenuItemWithMnemonic("_Koch")
	menuitem.Connect("activate", func() { draw_reentrant(koch) })
	submenu.Append(menuitem)

	koch_single := fractal.CreateKoch(KOCH_MAX, draw_line_single)
	menuitem = gtk.NewMenuItemWithMnemonic("_CSingleKoch")
	menuitem.Connect("activate", func() { draw_single(koch_single) })
	submenu.Append(menuitem)

	tree := fractal.CreateTree(TREE_MAX, draw_line_reentrant)
	menuitem = gtk.NewMenuItemWithMnemonic("_Tree")
	menuitem.Connect("activate", func() { draw_reentrant(tree) })
	submenu.Append(menuitem)

	sierpinski := fractal.CreateSierpinski(SIERPINSKI_MAX, draw_line_reentrant)
	menuitem = gtk.NewMenuItemWithMnemonic("_Sierpinski")
	menuitem.Connect("activate", func() { draw_reentrant(sierpinski) })
	submenu.Append(menuitem)

	cascademenu = gtk.NewMenuItemWithMnemonic("_Image")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	menuitem = gtk.NewMenuItemWithMnemonic("_Duke")
	menuitem.Connect("activate", func() {
		pixmap.GetDrawable().DrawRectangle(bg_gc, true, 0, 0, -1, -1)
		pixbuf, err := gdkpixbuf.NewPixbufFromFile("./images/duke.png")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		pixmap.GetDrawable().DrawPixbuf(fg_gc, pixbuf, 0, 0, 0, 0, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)
		gdkwin.GetDrawable().DrawDrawable(fg_gc, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
	})
	submenu.Append(menuitem)

	//--------------------------------------------------------
	// DrawingArea
	//--------------------------------------------------------
	canvas := gtk.NewDrawingArea()
	canvas.SetSizeRequest(720, 560)
	canvas.Connect("configure-event", func() {

		if fg_gc != nil {
			gdkwin.GetDrawable().DrawDrawable(fg_gc,
				pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
			return
		}
		if pixmap != nil {
			pixmap.Unref()
		}
		allocation := canvas.GetAllocation()
		pixmap = gdk.NewPixmap(canvas.GetWindow().GetDrawable(),
			allocation.Width,
			allocation.Height,
			24)
		bg_gc = gdk.NewGC(pixmap.GetDrawable())
		bg_gc.SetRgbFgColor(gdk.NewColor("whitesmoke"))
		pixmap.GetDrawable().DrawRectangle(bg_gc, true, 0, 0, -1, -1)

		fg_gc = gdk.NewGC(pixmap.GetDrawable())
		fg_gc.SetRgbFgColor(gdk.NewColor("black"))
	})
	canvas.Connect("expose-event", func() {
		if pixmap == nil {
			return
		}
		gdkwin.GetDrawable().DrawDrawable(fg_gc, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
	})
	//--------------------------------------------------------
	// Setting Layout
	//--------------------------------------------------------
	vbox.PackStart(menubar, false, false, 0)
	vbox.PackEnd(canvas, true, true, 0)
	win.Add(vbox)
	win.ShowAll()

	gdkwin = canvas.GetWindow()

}
func run_draw_app() {

	gdk.ThreadsInit()
	gtk.Init(nil)

	build_gtk_app()

	gtk.Main()
}
