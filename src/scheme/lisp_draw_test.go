/*
   Go lang 4th study program.
   This is test program for mini scheme subset.

   ex.) go test -v lisp_go draw.go lisp.go lisp_draw_test.go

   hidekuno@gmail.com
*/
package scheme

import (
	"testing"
)

func TestDraw(t *testing.T) {

	var exp Expression

	BuildFunc()
	rootEnv := NewSimpleEnv(nil, nil)
	BuildGtkFunc()

	exp, _ = DoCoreLogic("(draw-init)", rootEnv)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: draw-init")
	}
	exp, _ = DoCoreLogic("(draw-clear)", rootEnv)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: draw-clear")
	}
	exp, _ = DoCoreLogic("(draw-line 100 100 200 200)", rootEnv)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: draw-line")
	}
	exp, _ = DoCoreLogic("(create-image-from-file \"../../images/ch2-Z-G-30.gif\")", rootEnv)
	if _, ok := exp.(*Image); !ok {
		t.Fatal("failed test: create-image-from-file")
	}
	_, _ = DoCoreLogic("(define img (create-image-from-file \"../../images/ch2-Z-G-30.gif\"))", rootEnv)

	exp, _ = DoCoreLogic("(draw-image img 10 10)", rootEnv)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: draw-image")
	}
	exp, _ = DoCoreLogic("(scale-image img 90 90)", rootEnv)
	if _, ok := exp.(*Image); !ok {
		t.Fatal("failed test: draw-image")
	}
	exp, _ = DoCoreLogic("(rotate90-image img)", rootEnv)
	if _, ok := exp.(*Image); !ok {
		t.Fatal("failed test: rotate90-image")
	}
	exp, _ = DoCoreLogic("(rotate180-image img)", rootEnv)
	if _, ok := exp.(*Image); !ok {
		t.Fatal("failed test: rotate180-image")
	}
	exp, _ = DoCoreLogic("(rotate270-image img)", rootEnv)
	if _, ok := exp.(*Image); !ok {
		t.Fatal("failed test: rotate270-image")
	}
	testCode := [][]string{
		{"(draw-init)", "E2001"},
		{"(draw-line)", "E1007"},
		{"(draw-line 100 100 200)", "E1007"},
		{"(draw-line 100 100 200 200 100)", "E1007"},
		{"(draw-line #t 100 200 200)", "E1003"},
		{"(draw-line 100 100 200 #t)", "E1003"},
		{"(create-image-from-file)", "E1007"},
		{"(create-image-from-file \"hoge.gif\")", "E2002"},
		{"(create-image-from-file 10)", "E1015"},
		{"(draw-image 10 10)", "E1007"},
		{"(draw-image 10 10 10 10)", "E1007"},
		{"(draw-image 10 10 10)", "E2003"},
		{"(draw-image img #t 10)", "E1002"},
		{"(draw-image img 10 #t)", "E1002"},

		{"(scale-image 10 10)", "E1007"},
		{"(scale-image 10 10 10 10)", "E1007"},
		{"(scale-image 10 10 10)", "E2003"},
		{"(scale-image img #t 10)", "E1002"},
		{"(scale-image img 10 #t)", "E1002"},
		{"(rotate90-image)", "E1007"},
		{"(rotate90-image img 10)", "E1007"},
		{"(rotate90-image #t)", "E2003"},
		{"(rotate180-image)", "E1007"},
		{"(rotate180-image img 10)", "E1007"},
		{"(rotate180-image #t)", "E2003"},
		{"(rotate270-image)", "E1007"},
		{"(rotate270-image img 10)", "E1007"},
		{"(rotate270-image #t)", "E2003"},
	}
	for _, e := range testCode {
		_, err := DoCoreLogic(e[0], rootEnv)
		if err.(*RuntimeError).MsgCode != e[1] {
			t.Fatal("failed test: " + e[0])
		}
	}
}