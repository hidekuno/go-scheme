/*
   Go lang 4th study program.
   This is test program for mini scheme subset.

   ex.) go test -v lisp_go draw.go lisp.go lisp_draw_test.go

   hidekuno@gmail.com
*/
package draw

import (
	"scheme"
	"testing"
)

func executeTest(testCode [][]string, testName string, t *testing.T) {

	rootEnv := scheme.NewSimpleEnv(nil, nil)

	for i, c := range testCode {
		exp, err := scheme.DoCoreLogic(c[0], rootEnv)
		if err != nil {
			if e, ok := err.(*scheme.SyntaxError); ok {
				if e.MsgCode != c[1] {
					t.Log(i)
					t.Fatal("failed "+testName+" test", e.MsgCode)
				}
			} else if e, ok := err.(*scheme.RuntimeError); ok {
				if e.MsgCode != c[1] {
					t.Log(i)
					t.Fatal("failed "+testName+" test", e.MsgCode)
				}
			} else {
				t.Log(i)
				t.Fatal("failed "+testName+" test", err.Error())
			}
		} else {
			if exp.String() != c[1] {
				t.Log(i)
				t.Fatal("failed "+testName+" test", exp)
			}
		}
	}
}
func TestBuildFunc(t *testing.T) {
	scheme.BuildFunc()
	BuildGtkFunc()
}
func TestDrawInit(t *testing.T) {
	testCode := [][]string{
		{"(draw-init)", "nil"},

		{"(draw-init 10)", "E1007"},
		{"(draw-init)", "E2001"},
	}
	executeTest(testCode, "draw-init", t)

}
func TestDrawClear(t *testing.T) {
	testCode := [][]string{
		{"(draw-clear)", "nil"},

		{"(draw-clear 10)", "E1007"},
	}
	executeTest(testCode, "draw-clear", t)
}
func TestDrawLine(t *testing.T) {
	testCode := [][]string{
		{"(draw-line 100 100 200 200)", "nil"},

		{"(draw-line)", "E1007"},
		{"(draw-line 100 100 200)", "E1007"},
		{"(draw-line 100 100 200 200 100)", "E1007"},
		{"(draw-line #t 100 200 200)", "E1003"},
		{"(draw-line 100 100 200 #t)", "E1003"},
	}
	executeTest(testCode, "draw-line", t)
}
func TestCreateImageFromFile(t *testing.T) {
	testCode := [][]string{
		{"(create-image-from-file \"../../../images/duke.png\")", "Pixbuf: "},

		{"(create-image-from-file)", "E1007"},
		{"(create-image-from-file \"hoge.gif\")", "E2002"},
		{"(create-image-from-file 10)", "E1015"},
	}
	executeTest(testCode, "create-image-from-file", t)
}
func TestDrawImage(t *testing.T) {
	testCode := [][]string{
		{"(define img (create-image-from-file \"../../../images/duke.png\"))", "img"},
		{"(draw-image img 10 10)", "nil"},

		{"(draw-image 10 10)", "E1007"},
		{"(draw-image 10 10 10 10)", "E1007"},
		{"(draw-image 10 10 10)", "E2003"},
		{"(draw-image img #t 10)", "E1002"},
		{"(draw-image img 10 #t)", "E1002"},
	}
	executeTest(testCode, "draw-image", t)
}
func TestScaleImage(t *testing.T) {
	testCode := [][]string{
		{"(define img (create-image-from-file \"../../../images/duke.png\"))", "img"},
		{"(scale-image img 90 90)", "Pixbuf: "},

		{"(scale-image 10 10)", "E1007"},
		{"(scale-image 10 10 10 10)", "E1007"},
		{"(scale-image 10 10 10)", "E2003"},
		{"(scale-image img #t 10)", "E1002"},
		{"(scale-image img 10 #t)", "E1002"},
	}
	executeTest(testCode, "scale-image", t)
}
func TestRotate90Image(t *testing.T) {
	testCode := [][]string{
		{"(define img (create-image-from-file \"../../../images/duke.png\"))", "img"},
		{"(rotate90-image img)", "Pixbuf: "},

		{"(rotate90-image)", "E1007"},
		{"(rotate90-image img 10)", "E1007"},
		{"(rotate90-image #t)", "E2003"},
	}
	executeTest(testCode, "rotate90-image", t)
}
func TestRotate180Image(t *testing.T) {
	testCode := [][]string{
		{"(define img (create-image-from-file \"../../../images/duke.png\"))", "img"},
		{"(rotate180-image img)", "Pixbuf: "},

		{"(rotate180-image)", "E1007"},
		{"(rotate180-image img 10)", "E1007"},
		{"(rotate180-image #t)", "E2003"},
	}
	executeTest(testCode, "rotate180-image", t)
}
func TestRotate270Image(t *testing.T) {
	testCode := [][]string{
		{"(define img (create-image-from-file \"../../../images/duke.png\"))", "img"},
		{"(rotate270-image img)", "Pixbuf: "},

		{"(rotate270-image)", "E1007"},
		{"(rotate270-image img 10)", "E1007"},
		{"(rotate270-image #t)", "E2003"},
	}
	executeTest(testCode, "rotate270-image", t)
}
