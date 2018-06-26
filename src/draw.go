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
	"github.com/mattn/go-gtk/gtk"
	"runtime"
	"time"
)

const (
	KOCH_MAX       = 12
	TREE_MAX       = 20
	SIERPINSKI_MAX = 16
)

func build_gtk_app() {

	var (
		pixmap *gdk.Pixmap
		fg_gc  *gdk.GC
		bg_gc  *gdk.GC
		gdkwin *gdk.Window
	)
	win := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("フラクタル図サンプル")
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
	var draw_line_reentrant = func(x0, y0, x1, y1 int) {

		gdk.ThreadsEnter()

		pixmap.GetDrawable().DrawLine(fg_gc, x0, y0, x1, y1)

		var rec gdk.Rectangle
		rec.Height = y1 - y0
		rec.Width = x1 - x0
		rec.X = x0
		rec.Y = y0

		if rec.Height < 0 {
			rec.Height = rec.Height * -1
		}
		if rec.Width < 0 {
			rec.Width = rec.Width * -1
		}
		rec.Height += 2
		rec.Width += 2
		if x1 < x0 {
			rec.X = x1
		}
		if y1 < y0 {
			rec.Y = y1
		}
		// gdkwin.Invalidate(&rec, false)
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

	cascademenu = gtk.NewMenuItemWithMnemonic("_View")
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
		gdkwin.GetDrawable().DrawDrawable(fg_gc,
			pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
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
func main() {
	runtime.GOMAXPROCS(1)

	gdk.ThreadsInit()
	gtk.Init(nil)

	build_gtk_app()

	gtk.Main()
}
