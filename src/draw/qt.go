/*
   Go lang 2nd study program.
   This is thing for go routine (multi threading program).

   hidekuno@gmail.com
*/
package draw

import (
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"

	"fractal"
)

type Canvas struct {
	*widgets.QWidget
	painter *gui.QPainter
	draw    func()
}

func NewCanvas() *Canvas {
	canvas := &Canvas{widgets.NewQWidget(nil, 0), nil, nil}
	canvas.ConnectPaintEvent(canvas.paintEvent)
	return canvas
}
func (self *Canvas) DrawLine(x0, y0, x1, y1 int) {
	self.painter.DrawLine3(x0, y0, x1, y1)
}

func (self *Canvas) paintEvent(event *gui.QPaintEvent) {
	self.painter = gui.NewQPainter2(self.BackingStore().PaintDevice())
	if self.draw != nil {
		self.draw()
	}
	self.painter.End()
}

func (self *Canvas) DrawImage() {
	image := gui.NewQImage()
	image.Load(SampleImage, "png")

	tf := gui.NewQTransform()
	w := (float64)(image.Width())
	h := (float64)(image.Height())
	tf.Translate(w, h)
	// tf.Rotate(180.0, core.Qt__XAxis)
	tf.Rotate(180.0, core.Qt__YAxis)
	self.painter.SetTransform(tf, true)
	self.painter.DrawImage(core.NewQRectF4(0, 0, w, h), image, core.NewQRectF(), core.Qt__AutoColor)
}
func (self *Canvas) ChangeDrawble(p func()) {
	self.draw = p
	self.Update()
}

func BuildQtApp(titleName string) {

	qApp := widgets.NewQApplication(len(os.Args), os.Args)
	window := widgets.NewQMainWindow(nil, 0)
	window.SetObjectName(titleName)
	window.SetWindowTitle(titleName)
	window.Resize(core.NewQSize2(720, 560))

	menubar := window.MenuBar()
	fileMenu := menubar.AddMenu2("&File")
	exitAction := fileMenu.AddAction("&Exit")
	exitAction.SetShortcut(gui.QKeySequence_FromString("Ctrl+Q", 0))
	exitAction.ConnectTriggered(func(bool) { qApp.Quit() })

	fractalMenu := menubar.AddMenu2("Fractal")
	kochAction := fractalMenu.AddAction("Koch")
	treeAction := fractalMenu.AddAction("Tree")
	sierpinskiAction := fractalMenu.AddAction("Sierpinski")

	canvas := NewCanvas()
	kochAction.ConnectTriggered(func(bool) {
		canvas.ChangeDrawble(fractal.CreateKoch(KochMax, canvas.DrawLine))
	})
	treeAction.ConnectTriggered(func(bool) {
		canvas.ChangeDrawble(fractal.CreateTree(TreeMax, canvas.DrawLine))
	})
	sierpinskiAction.ConnectTriggered(func(bool) {
		canvas.ChangeDrawble(fractal.CreateSierpinski(SierpinskiMax, canvas.DrawLine))
	})

	imageMenu := menubar.AddMenu2("Image")
	glendaAction := imageMenu.AddAction("Glenda")
	glendaAction.ConnectTriggered(func(bool) {
		canvas.ChangeDrawble(func() {
			image := gui.NewQImage()
			image.Load(SampleImage, "png")

			tf := gui.NewQTransform()
			w := (float64)(image.Width())
			h := (float64)(image.Height())
			tf.Translate(w, h)
			// tf.Rotate(180.0, core.Qt__XAxis)
			tf.Rotate(180.0, core.Qt__YAxis)
			canvas.painter.SetTransform(tf, true)
			canvas.painter.DrawImage(core.NewQRectF4(0, 0, w, h), image, core.NewQRectF(), core.Qt__AutoColor)
		})
	})
	window.SetCentralWidget(canvas)
	window.Show()
}
