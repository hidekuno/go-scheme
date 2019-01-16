/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package experiment

import (
	"fmt"
	"scheme"
	"strconv"
	"time"
)

func BuildGoFunc() {

	scheme.AddSpecialFunc("go-append", func(env *scheme.SimpleEnv, exp []scheme.Expression) (scheme.Expression, error) {
		if len(exp) < 2 {
			return nil, scheme.NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}

		finish := make(chan bool, 2)
		var exp1 scheme.Expression
		var exp2 scheme.Expression

		go func() {
			t0 := time.Now()
			exp1, _ = scheme.DoEval(exp[0], env)
			t1 := time.Now()
			fmt.Println("go-1", t1.Sub(t0))
			finish <- true
		}()
		go func() {
			t0 := time.Now()
			exp2, _ = scheme.DoEval(exp[1], copyEnv(env))
			t1 := time.Now()
			fmt.Println("go-2", t1.Sub(t0))
			finish <- true
		}()
		<-finish
		<-finish
		return scheme.NewList(append((exp1.(*scheme.List)).Value, (exp2.(*scheme.List)).Value...)), nil
	})
}

// support for multi threading
func copyEnv(env *scheme.SimpleEnv) *scheme.SimpleEnv {
	env2 := scheme.NewSimpleEnv(nil, nil)

	for key, _ := range *(env.EnvTable) {

		if fn, ok := ((*env.EnvTable)[key]).(*scheme.Function); ok {
			fn2 := scheme.NewFunction(env2, &fn.ParamName, fn.Body, fn.Name)
			env2.Regist(key, fn2)
		}
	}
	return env2
}
