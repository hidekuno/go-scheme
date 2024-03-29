/*
   Go lang 2nd study program.
   This is thing for go routine (multi threading program).

   +------------------------------------------+
   | go get github.com/mattn/go-gtk/gtk
   | go install github.com/mattn/go-gtk/gtk
   +------------------------------------------+

   hidekuno@gmail.com
*/
package draw

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/gtk"
	"local.packages/fractal"
)

const (
	ScreenWidth  = 560
	ScreenHeight = 560
)

func BuildGtkApp(titleName string) (*gdk.Pixmap, *gdk.Window, *gdk.GC, *gdk.GC) {

	var (
		pixmap *gdk.Pixmap
		gdkwin *gdk.Window
		fg     *gdk.GC
		bg     *gdk.GC
	)
	//--------------------------------------------------------
	// Init etc...
	//--------------------------------------------------------
	gdk.ThreadsInit()
	gtk.Init(nil)

	win := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	win.SetTitle(titleName)
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
	if titleName == "scheme.go" {
		cascademenu.SetSensitive(false)
	}

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
	menuitem.SetSensitive(false)

	menuitem = gtk.NewMenuItemWithMnemonic("_Tree")
	menuitem.Connect("activate", func() { drawReEntrant(fractal.CreateTree(TreeMax, drawLineReEntrant)) })
	submenu.Append(menuitem)

	menuitem = gtk.NewMenuItemWithMnemonic("_Sierpinski")
	menuitem.Connect("activate", func() { drawReEntrant(fractal.CreateSierpinski(SierpinskiMax, drawLineReEntrant)) })
	submenu.Append(menuitem)
	if titleName == "scheme.go" {
		cascademenu.SetSensitive(false)
	}

	cascademenu = gtk.NewMenuItemWithMnemonic("_Image")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	menuitem = gtk.NewMenuItemWithMnemonic("_Simple")
	menuitem.Connect("activate", func() {
		pixmap.GetDrawable().DrawRectangle(bg, true, 0, 0, -1, -1)
		orgPixbuf, err := gdkpixbuf.NewPixbufFromFile(filepath.Join(os.Getenv("GOPATH"), SampleImage))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		type ImageSample struct {
			Scale int
			Angle gdkpixbuf.PixbufRotation
		}
		samples := []ImageSample{
			{180, gdkpixbuf.PIXBUF_ROTATE_NONE},
			{120, gdkpixbuf.PIXBUF_ROTATE_COUNTERCLOCKWISE},
			{80, gdkpixbuf.PIXBUF_ROTATE_UPSIDEDOWN},
			{52, gdkpixbuf.PIXBUF_ROTATE_CLOCKWISE},
			{34, gdkpixbuf.PIXBUF_ROTATE_NONE},
		}
		w, h := 0, 0
		for _, rec := range samples {
			pixbuf := orgPixbuf.ScaleSimple(rec.Scale, rec.Scale, gdkpixbuf.INTERP_HYPER).RotateSimple(rec.Angle)
			pixmap.GetDrawable().DrawPixbuf(fg, pixbuf, 0, 0, w, h, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)
			gdkwin.Invalidate(nil, false)
			w, h = w+pixbuf.GetWidth(), h+pixbuf.GetHeight()
		}
	})
	submenu.Append(menuitem)
	if titleName == "scheme.go" {
		cascademenu.SetSensitive(false)
	}
	//--------------------------------------------------------
	// DrawingArea
	//--------------------------------------------------------
	canvas := gtk.NewDrawingArea()
	canvas.SetSizeRequest(ScreenWidth, ScreenHeight)

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
	return pixmap, gdkwin, fg, bg
}
