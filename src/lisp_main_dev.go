/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package main

import (
	"fmt"
	"runtime"
	"time"
)

func build_go_func() {
	special_func["go-append"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
		finish_ch := make(chan bool, 2)
		var append_list []Expression
		var exp1 Expression
		var exp2 Expression
		var err error
		go func() {
			exp1, err = eval(v[0], env)
			finish_ch <- true
		}()
		go func() {
			exp2, err = eval(v[1], env)
			finish_ch <- true
		}()
		<-finish_ch
		<-finish_ch
		append_list = append(append_list, (exp1.(*List)).Value...)
		append_list = append(append_list, (exp2.(*List)).Value...)
		return NewList(append_list), nil
	}
	special_func["time"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
		t0 := time.Now()
		if exp, err := eval(v[0], env); err != nil {
			return exp, err
		}
		t1 := time.Now()
		fmt.Println(t1.Sub(t0))
		return NewNil(), nil
	}
}

// Main
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	build_func()
	build_go_func()
	do_interactive()
}
