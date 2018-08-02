/*
   Go lang 4th study program.
   This is test program for mini scheme subset.

   ex.) go test -v lisp.go lisp_test.go

   hidekuno@gmail.com
*/
package main

import (
	"testing"
)

func run_draw_app()                               {}
func draw_clear()                                 {}
func draw_line_reentrant_lisp(x0, y0, x1, y1 int) {}
func draw_imagefile(filename string)              {}

func TestDraw(t *testing.T) {
	var (
		exp Expression
	)
	build_func()
	rootEnv := NewSimpleEnv(nil, nil)
	build_gtk_func()
	exp, _ = do_core_logic("(draw-init)", rootEnv)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: draw-init")
	}
	exp, _ = do_core_logic("(draw-clear)", rootEnv)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: draw-clear")
	}
	exp, _ = do_core_logic("(draw-line 100 100 200 200)", rootEnv)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: draw-line")
	}
	exp, _ = do_core_logic("(draw-imagefile \"./images/duke.png\")", rootEnv)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: draw-imagefile")
	}
	testCode := [][]string{
		{"(draw-line)", "E1007"},
		{"(draw-line 100 100 200)", "E1007"},
		{"(draw-line 100 100 200 200 100)", "E1007"},
		{"(draw-line #t 100 200 200)", "E1003"},
		{"(draw-line 100 100 200 #t)", "E1003"},
		{"(draw-imagefile)", "E1007"},
		{"(draw-imagefile \"a.gif\" \"b.gif\")", "E1007"},
		{"(draw-imagefile #t)", "E1003"},
	}
	for _, e := range testCode {
		_, err := do_core_logic(e[0], rootEnv)
		if err.(*RuntimeError).MsgCode != e[1] {
			t.Fatal("failed test: " + e[0])
		}
	}
}
