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

func buildGoFunc() {
	specialFuncTbl["go-append"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
		finish := make(chan bool, 2)
		var result []Expression
		var exp1 Expression
		var exp2 Expression

		go func() {
			t0 := time.Now()
			exp1, _ = eval(v[0], env)
			t1 := time.Now()
			fmt.Println("go-1", t1.Sub(t0))
			finish <- true
		}()
		go func() {
			t0 := time.Now()
			exp2, _ = eval(v[1], env)
			t1 := time.Now()
			fmt.Println("go-2", t1.Sub(t0))
			finish <- true
		}()
		<-finish
		<-finish
		result = append(result, (exp1.(*List)).Value...)
		result = append(result, (exp2.(*List)).Value...)
		return NewList(result), nil
	}
	specialFuncTbl["test-list"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
		time.Sleep(5 * time.Second)
		var result []Expression
		return NewList(result), nil
	}
}

// Main
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	buildFunc()
	buildGoFunc()
	doInteractive()
}
