/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

import (
	"fmt"
	"strconv"
	"time"
)

func BuildGoFunc() {
	specialFuncTbl["go-append"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
		if len(v) < 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(v)))
		}

		finish := make(chan bool, 2)
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
			exp2, _ = eval(v[1], copyEnv(env))
			t1 := time.Now()
			fmt.Println("go-2", t1.Sub(t0))
			finish <- true
		}()
		<-finish
		<-finish
		return NewList(append((exp1.(*List)).Value, (exp2.(*List)).Value...)), nil
	}
	specialFuncTbl["test-list"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
		time.Sleep(5 * time.Second)
		var result []Expression
		return NewList(result), nil
	}
}

// support for multi threading
func copyEnv(env *SimpleEnv) *SimpleEnv {
	env2 := NewSimpleEnv(nil, nil)

	for key, _ := range *(env.EnvTable) {

		if fn, ok := ((*env.EnvTable)[key]).(*Function); ok {
			fn2 := NewFunction(env2, &fn.ParamName, fn.Body, fn.Name)
			env2.Regist(key, fn2)
		}
	}
	return env2
}
