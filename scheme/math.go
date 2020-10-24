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
func mathImpl(exp []Expression, env *SimpleEnv, fn func(float64) float64) (Expression, error) {
	if len(exp) != 1 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env,
		func(exp ...Expression) (Expression, error) {
			if v, ok := exp[0].(*Float); ok {
				return NewFloat(fn(v.Value)), nil
			} else if v, ok := exp[0].(*Integer); ok {
				return NewFloat(fn((float64)(v.Value))), nil
			}
			return nil, NewRuntimeError("E1003", reflect.TypeOf(exp[0]).String())
		})
}

// Build Global environement.
func buildMathFunc() {

	buildInFuncTbl["sqrt"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return mathImpl(exp, env, math.Sqrt)
	}
	buildInFuncTbl["sin"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return mathImpl(exp, env, math.Sin)
	}
	buildInFuncTbl["cos"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return mathImpl(exp, env, math.Cos)
	}
	buildInFuncTbl["tan"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return mathImpl(exp, env, math.Tan)
	}
	buildInFuncTbl["asin"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return mathImpl(exp, env, math.Asin)
	}
	buildInFuncTbl["acos"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return mathImpl(exp, env, math.Acos)
	}
	buildInFuncTbl["atan"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return mathImpl(exp, env, math.Atan)
	}
	buildInFuncTbl["log"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return mathImpl(exp, env, math.Log)
	}
	buildInFuncTbl["exp"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return mathImpl(exp, env, math.Exp)
	}
	buildInFuncTbl["abs"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return mathImpl(exp, env, math.Abs)
	}
	buildInFuncTbl["truncate"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return mathImpl(exp, env, math.Trunc)
	}
	buildInFuncTbl["floor"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return mathImpl(exp, env, math.Floor)
	}
	buildInFuncTbl["ceiling"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return mathImpl(exp, env, math.Ceil)
	}
	buildInFuncTbl["round"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return mathImpl(exp, env, math.Round)
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
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if v, ok := exp[0].(*Integer); ok {
					return NewInteger(rand.Intn(v.Value)), nil
				}
				return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[0]).String())
			})
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
