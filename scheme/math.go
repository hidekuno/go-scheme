/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"strconv"
)

// math skelton
func mathImpl(mathFunc func(float64) float64, exp ...Expression) (Expression, error) {
	if len(exp) != 1 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	if v, ok := exp[0].(*Float); ok {
		return NewFloat(mathFunc(v.Value)), nil
	} else if v, ok := exp[0].(*Integer); ok {
		return NewFloat(mathFunc((float64)(v.Value))), nil
	}
	return nil, NewRuntimeError("E1003", reflect.TypeOf(exp[0]).String())
}

// Build Global environement.
func buildMathFunc() {

	buildInFuncTbl["sqrt"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) { return mathImpl(math.Sqrt, exp...) })
	}
	buildInFuncTbl["sin"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) { return mathImpl(math.Sin, exp...) })
	}
	buildInFuncTbl["cos"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) { return mathImpl(math.Cos, exp...) })
	}
	buildInFuncTbl["tan"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) { return mathImpl(math.Tan, exp...) })
	}
	buildInFuncTbl["atan"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) { return mathImpl(math.Atan, exp...) })
	}
	buildInFuncTbl["log"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) { return mathImpl(math.Log, exp...) })
	}
	buildInFuncTbl["exp"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) { return mathImpl(math.Exp, exp...) })
	}
	buildInFuncTbl["rand-init"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 0 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			return nil, NewRuntimeError("E9999")
		}
		rand.Seed(seed.Int64())
		return NewNil(), nil
	}
	buildInFuncTbl["rand-integer"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		if v, ok := exp[0].(*Integer); ok {
			return NewInteger(rand.Intn(v.Value)), nil
		}
		return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[0]).String())
	}
	buildInFuncTbl["expt"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 2 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				f := []float64{0.0, 0.0}
				for i, e := range exp {
					if n, ok := e.(*Float); ok {
						f[i] = n.Value
					} else if n, ok := e.(*Integer); ok {
						f[i] = (float64)(n.Value)
					} else {
						return nil, NewRuntimeError("E1003", reflect.TypeOf(e).String())
					}
				}
				return NewFloat(math.Pow(f[0], f[1])), nil
			})
	}
}
