/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

import (
	"encoding/binary"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

func isType(exp []Expression, env *SimpleEnv, cmp func(Expression) bool) (Expression, error) {
	if len(exp) != 1 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	e, err := eval(exp[0], env)
	if err != nil {
		return nil, err
	}
	return NewBoolean(cmp(e)), nil
}
func isTypeOfInteger(exp []Expression, env *SimpleEnv, cmp func(int) bool) (Expression, error) {
	if len(exp) != 1 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	e, err := eval(exp[0], env)
	if err != nil {
		return nil, err
	}
	n, ok := e.(*Integer)
	if !ok {
		return nil, NewRuntimeError("E1002", reflect.TypeOf(e).String())
	}
	return NewBoolean(cmp(n.Value)), nil
}
func isTypeOfNumber(exp []Expression, env *SimpleEnv, cmp func(Number) bool) (Expression, error) {
	if len(exp) != 1 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	e, err := eval(exp[0], env)
	if err != nil {
		return nil, err
	}
	n, err := CreateNumber(e)
	if err != nil {
		return nil, err
	}
	return NewBoolean(cmp(n)), nil
}

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
	buildInFuncTbl["even?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isTypeOfInteger(exp, env, func(n int) bool { return n%2 == 0 })
	}
	buildInFuncTbl["odd?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isTypeOfInteger(exp, env, func(n int) bool { return n%2 != 0 })
	}
	buildInFuncTbl["zero?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isTypeOfNumber(exp, env, func(n Number) bool { return (toInt(n)).Value == 0 })
	}
	buildInFuncTbl["positive?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isTypeOfNumber(exp, env, func(n Number) bool { return (toInt(n)).Value > 0 })
	}
	buildInFuncTbl["negative?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isTypeOfNumber(exp, env, func(n Number) bool { return (toInt(n)).Value < 0 })
	}
	buildInFuncTbl["list?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isType(exp, env, func(e Expression) bool { return reflect.TypeOf(e) == reflect.TypeOf(&List{}) })
	}
	buildInFuncTbl["pair?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isType(exp, env, func(e Expression) bool { return reflect.TypeOf(e) == reflect.TypeOf(&Pair{}) })
	}
	buildInFuncTbl["char?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isType(exp, env, func(e Expression) bool { return reflect.TypeOf(e) == reflect.TypeOf(&Char{}) })
	}
	buildInFuncTbl["string?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isType(exp, env, func(e Expression) bool { return reflect.TypeOf(e) == reflect.TypeOf(&String{}) })
	}
	buildInFuncTbl["integer?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isType(exp, env, func(e Expression) bool { return reflect.TypeOf(e) == reflect.TypeOf(&Integer{}) })
	}
	buildInFuncTbl["number?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isType(exp, env, func(e Expression) bool {
			return (reflect.TypeOf(e) == reflect.TypeOf(&Integer{})) || (reflect.TypeOf(e) == reflect.TypeOf(&Float{}))
		})
	}
	buildInFuncTbl["procedure?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isType(exp, env, func(e Expression) bool {
			return (reflect.TypeOf(e) == reflect.TypeOf(&Function{})) || (reflect.TypeOf(e) == reflect.TypeOf(&BuildInFunc{}))
		})
	}
	buildInFuncTbl["symbol?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isType(exp, env, func(e Expression) bool { return reflect.TypeOf(e) == reflect.TypeOf(&Symbol{}) })
	}
	buildInFuncTbl["boolean?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isType(exp, env, func(e Expression) bool { return reflect.TypeOf(e) == reflect.TypeOf(&Boolean{}) })
	}
	buildInFuncTbl["native-endian"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 0 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		var x uint64 = 1
		buf := make([]byte, binary.MaxVarintLen64)
		binary.PutUvarint(buf, x)

		if buf[0] == 1 {
			return NewSymbol("little-endian"), nil
		} else {
			return NewSymbol("big-endian"), nil
		}
		return nil, NewRuntimeError("E9999")
	}
}
