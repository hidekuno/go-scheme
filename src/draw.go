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
	"time"
)

const (
	KochMax       = 12
	TreeMax       = 20
	SierpinskiMax = 16
)

var drawLineLisp func(x0, y0, x1, y1 int)
var drawClear func()
var drawImageFile func(filename string)

func buildGtkApp() {

	var (
		pixmap *gdk.Pixmap
		fg     *gdk.GC
		bg     *gdk.GC
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
	drawLineLisp = func(x0, y0, x1, y1 int) {
		gdk.ThreadsEnter()
		pixmap.GetDrawable().DrawLine(fg, x0, y0, x1, y1)
		gdkwin.Invalidate(nil, false)
		gdk.ThreadsLeave()
	}
	drawClear = func() {
		pixmap.GetDrawable().DrawRectangle(bg, true, 0, 0, -1, -1)
		gdkwin.Invalidate(nil, false)
	}
	drawImageFile = func(filename string) {

		pixbuf, err := gdkpixbuf.NewPixbufFromFileAtScale(filename, -1, -1, true)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		pixmap.GetDrawable().DrawPixbuf(fg, pixbuf, 0, 0, 0, 0, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)
		gdkwin.Invalidate(nil, false)
		gdkwin.GetDrawable().DrawDrawable(fg, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
	}
	var drawLineReEntrant = func(x0, y0, x1, y1 int) {
		gdk.ThreadsEnter()
		pixmap.GetDrawable().DrawLine(fg, x0, y0, x1, y1)
		gdk.ThreadsLeave()
	}
	var drawReEntrant = func(paint func()) {
		pixmap.GetDrawable().DrawRectangle(bg, true, 0, 0, -1, -1)
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

	var drawLineSingle = func(x0, y0, x1, y1 int) { pixmap.GetDrawable().DrawLine(fg, x0, y0, x1, y1) }
	var drawSingle = func(paint func()) {
		pixmap.GetDrawable().DrawRectangle(bg, true, 0, 0, -1, -1)
		paint()
		gdkwin.GetDrawable().DrawDrawable(fg, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
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

	menuitem = gtk.NewMenuItemWithMnemonic("_Koch")
	menuitem.Connect("activate", func() { drawReEntrant(fractal.CreateKoch(KochMax, drawLineReEntrant)) })
	submenu.Append(menuitem)

	menuitem = gtk.NewMenuItemWithMnemonic("_CSingleKoch")
	menuitem.Connect("activate", func() { drawSingle(fractal.CreateKoch(KochMax, drawLineSingle)) })
	submenu.Append(menuitem)

	menuitem = gtk.NewMenuItemWithMnemonic("_Tree")
	menuitem.Connect("activate", func() { drawReEntrant(fractal.CreateTree(TreeMax, drawLineReEntrant)) })
	submenu.Append(menuitem)

	menuitem = gtk.NewMenuItemWithMnemonic("_Sierpinski")
	menuitem.Connect("activate", func() { drawReEntrant(fractal.CreateSierpinski(SierpinskiMax, drawLineReEntrant)) })
	submenu.Append(menuitem)

	cascademenu = gtk.NewMenuItemWithMnemonic("_Image")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	menuitem = gtk.NewMenuItemWithMnemonic("_Duke")
	menuitem.Connect("activate", func() {
		pixmap.GetDrawable().DrawRectangle(bg, true, 0, 0, -1, -1)
		pixbuf, err := gdkpixbuf.NewPixbufFromFile("./images/duke.png")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		pixmap.GetDrawable().DrawPixbuf(fg, pixbuf, 0, 0, 0, 0, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)
		gdkwin.GetDrawable().DrawDrawable(fg, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
	})
	submenu.Append(menuitem)

	//--------------------------------------------------------
	// DrawingArea
	//--------------------------------------------------------
	canvas := gtk.NewDrawingArea()
	canvas.SetSizeRequest(720, 560)
	canvas.Connect("configure-event", func() {

		if fg != nil {
			gdkwin.GetDrawable().DrawDrawable(fg,
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
		bg = gdk.NewGC(pixmap.GetDrawable())
		bg.SetRgbFgColor(gdk.NewColor("whitesmoke"))
		pixmap.GetDrawable().DrawRectangle(bg, true, 0, 0, -1, -1)

		fg = gdk.NewGC(pixmap.GetDrawable())
		fg.SetRgbFgColor(gdk.NewColor("black"))
	})
	canvas.Connect("expose-event", func() {
		if pixmap == nil {
			return
		}
		gdkwin.GetDrawable().DrawDrawable(fg, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
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
func runDrawApp() {

	gdk.ThreadsInit()
	gtk.Init(nil)

	buildGtkApp()

	gtk.Main()
}
