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

		go func() {
			t0 := time.Now()
			exp1, _ = eval(v[0], env)
			t1 := time.Now()
			fmt.Println("go-1", t1.Sub(t0))
			finish_ch <- true
		}()
		go func() {
			t0 := time.Now()
			exp2, _ = eval(v[1], env)
			t1 := time.Now()
			fmt.Println("go-2", t1.Sub(t0))
			finish_ch <- true
		}()
		<-finish_ch
		<-finish_ch
		append_list = append(append_list, (exp1.(*List)).Value...)
		append_list = append(append_list, (exp2.(*List)).Value...)
		return NewList(append_list), nil
	}
	special_func["test-list"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
		time.Sleep(5 * time.Second)
		var append_list []Expression
		return NewList(append_list), nil
	}
}

// Main
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	build_func()
	build_go_func()
	do_interactive()
}
