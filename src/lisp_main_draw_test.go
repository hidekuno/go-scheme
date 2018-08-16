/*
   Go lang 4th study program.
   This is test program for mini scheme subset.

   ex.) go test -v lisp_main_draw.go draw.go lisp.go lisp_main_draw_test.go

   hidekuno@gmail.com
*/
package main

import (
	"testing"
)

func TestDraw(t *testing.T) {
	var (
		exp Expression
	)
	buildFunc()
	rootEnv := NewSimpleEnv(nil, nil)
	buildGtkFunc()

	exp, _ = doCoreLogic("(draw-init)", rootEnv)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: draw-init")
	}
	exp, _ = doCoreLogic("(draw-clear)", rootEnv)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: draw-clear")
	}
	exp, _ = doCoreLogic("(draw-line 100 100 200 200)", rootEnv)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: draw-line")
	}
	exp, _ = doCoreLogic("(create-image-from-file \"images/ch2-Z-G-30.gif\")", rootEnv)
	if _, ok := exp.(*Image); !ok {
		t.Fatal("failed test: create-image-from-file")
	}
	_, _ = doCoreLogic("(define img (create-image-from-file \"images/ch2-Z-G-30.gif\"))", rootEnv)

	exp, _ = doCoreLogic("(draw-image img 10 10)", rootEnv)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: draw-image")
	}
	exp, _ = doCoreLogic("(scale-image img 90 90)", rootEnv)
	if _, ok := exp.(*Image); !ok {
		t.Fatal("failed test: draw-image")
	}
	exp, _ = doCoreLogic("(rotate90-image img)", rootEnv)
	if _, ok := exp.(*Image); !ok {
		t.Fatal("failed test: rotate90-image")
	}
	exp, _ = doCoreLogic("(rotate180-image img)", rootEnv)
	if _, ok := exp.(*Image); !ok {
		t.Fatal("failed test: rotate180-image")
	}
	exp, _ = doCoreLogic("(rotate270-image img)", rootEnv)
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
		_, err := doCoreLogic(e[0], rootEnv)
		if err.(*RuntimeError).MsgCode != e[1] {
			t.Fatal("failed test: " + e[0])
		}
	}
}
