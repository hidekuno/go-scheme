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

func run_draw_app()                          {}
func draw_clear()                            {}
func draw_line_reentrant(x0, y0, x1, y1 int) {}

func Test_draw(t *testing.T) {
	var (
		exp Expression
	)
	build_func()
	root_env := NewSimpleEnv(nil, nil)
	build_gtk_func()
	exp, _ = do_core_logic("(draw_init)", root_env)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: draw_init")
	}
	exp, _ = do_core_logic("(draw_clear)", root_env)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: draw_clear")
	}
	exp, _ = do_core_logic("(draw_line 100 100 200 200)", root_env)
	if _, ok := exp.(*Nil); !ok {
		t.Fatal("failed test: draw_line")
	}
	test_code := [][]string{
		{"(draw_line)", "E1007"},
		{"(draw_line 100 100 200)", "E1007"},
		{"(draw_line 100 100 200 200 100)", "E1007"},
		{"(draw_line #t 100 200 200)", "E1003"},
		{"(draw_line 100 100 200 #t)", "E1003"},
	}
	for _, e := range test_code {
		_, err := do_core_logic(e[0], root_env)
		if err.(*RuntimeError).MsgCode != e[1] {
			t.Fatal("failed test: " + e[0])
		}
	}
}
