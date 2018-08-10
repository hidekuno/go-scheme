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

var (
	execFinished = false
)

func buildGtkFunc() {
	errorMsg["E2001"] = "Aleady Gtk Init"

	specialFuncTbl["draw-init"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
		if execFinished == true {
			return nil, NewRuntimeError("E2001")
		}
		go runDrawApp()

		specialFuncTbl["draw-clear"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
			drawClear()
			return NewNil(), nil
		}
		specialFuncTbl["draw-line"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
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
			drawLineLisp(point[0], point[1], point[2], point[3])
			return NewNil(), nil
		}
		specialFuncTbl["draw-imagefile"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
			if len(v) != 1 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(v)))
			}
			if s, ok := v[0].(*String); ok {
				drawImageFile(s.Value)
			} else {
				return nil, NewRuntimeError("E1003", reflect.TypeOf(v[0]).String())
			}
			return NewNil(), nil
		}
		execFinished = true
		return NewNil(), nil
	}
}

// Main
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	buildFunc()
	buildGtkFunc()

	cui := make(chan bool)
	go func() {
		doInteractive()
		cui <- true
	}()
	<-cui
}
