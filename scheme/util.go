/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

// eqv
func eqv(exp []Expression, env *SimpleEnv) (Expression, error) {

	if len(exp) != 2 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}

	x, err := eval(exp[0], env)
	if err != nil {
		return x, err
	}
	y, err := eval(exp[1], env)
	if err != nil {
		return y, err
	}
	if a, ok := x.(*Integer); ok {
		if b, ok := y.(*Integer); ok {
			return NewBoolean(a.Value == b.Value), nil
		}
	}
	if a, ok := x.(*Float); ok {
		if b, ok := y.(*Float); ok {
			return NewBoolean(a.Value == b.Value), nil
		}
	}
	if a, ok := x.(*Boolean); ok {
		if b, ok := y.(*Boolean); ok {
			return NewBoolean(a.Value == b.Value), nil
		}
	}
	if a, ok := x.(*Symbol); ok {
		if b, ok := y.(*Symbol); ok {
			return NewBoolean(a.Value == b.Value), nil
		}
	}
	if a, ok := x.(*Char); ok {
		if b, ok := y.(*Char); ok {
			return NewBoolean(a.Value == b.Value), nil
		}
	}
	if a, ok := x.(*String); ok {
		if b, ok := y.(*String); ok {
			return NewBoolean(a.Value == b.Value), nil
		}
	}
	return NewBoolean(false), nil
}

// Build Global environement.
func buildUtilFunc() {

	buildInFuncTbl["eqv?"] = eqv
	buildInFuncTbl["eq?"] = buildInFuncTbl["eqv?"]

	buildInFuncTbl["identity"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return eval(exp[0], env)
	}
	buildInFuncTbl["time"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		t0 := time.Now()
		e, err := eval(exp[0], env)
		t1 := time.Now()
		fmt.Println(t1.Sub(t0))
		return e, err
	}
	//srfi-98
	buildInFuncTbl["get-environment-variable"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		v, err := eval(exp[0], env)
		if err != nil {
			return v, err
		}
		s, ok := v.(*String)
		if !ok {
			return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
		}
		return NewString(os.Getenv(s.Value)), nil
	}
}
