/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package main

import (
	"reflect"
	"runtime"
	"strconv"
)

func build_gtk_func() {

	special_func["draw-init"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
		go run_draw_app()

		special_func["draw-clear"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
			draw_clear()
			return NewNil(), nil
		}
		special_func["draw-line"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
			var point [4]int
			if len(v) != 4 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(v)))
			}

			for i, sexp := range v {
				e, err := eval(sexp, env)
				if err != nil {
					return nil, err
				}
				if p, ok := e.(*Integer); ok {
					point[i] = p.Value
				} else if p, ok := e.(*Float); ok {
					point[i] = int(p.Value)
				} else {
					return nil, NewRuntimeError("E1003", reflect.TypeOf(e).String())
				}
			}
			draw_line_reentrant_lisp(point[0], point[1], point[2], point[3])
			return NewNil(), nil
		}
		special_func["draw-imagefile"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
			if len(v) != 1 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(v)))
			}
			if s, ok := v[0].(*String); ok {
				draw_imagefile(s.Value)
			} else {
				return nil, NewRuntimeError("E1003", reflect.TypeOf(v[0]).String())
			}
			return NewNil(), nil
		}
		return NewNil(), nil
	}
}

// Main
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	build_func()
	build_gtk_func()

	cui_ch := make(chan bool)
	go func() {
		do_interactive()
		cui_ch <- true
	}()
	<-cui_ch
}
