/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

import (
	"reflect"
	"strconv"
)

// addl, subl, imul, idiv
func calcOperate(calc func(Number, Number) Number, exp ...Expression) (Number, error) {
	if 1 >= len(exp) {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	result, err := CreateNumber(exp[0])
	if err != nil {
		return nil, err
	}
	for _, i := range exp[1:] {
		prm, ok := i.(Number)
		if !ok {
			return nil, NewRuntimeError("E1003", reflect.TypeOf(i).String())
		}

		if _, ok := result.(*Float); ok {
			if c, ok := i.(*Integer); ok {
				prm = NewFloat(float64(c.Value))
			}
		}
		if org, ok := result.(*Integer); ok {
			if _, ok := i.(*Float); ok {
				result = NewFloat(float64(org.Value))
			}
		}
		result = calc(result, prm)
	}
	return result, nil
}

// gt,lt,ge,le
func cmpOperate(cmp func(Number, Number) bool, exp ...Expression) (*Boolean, error) {
	if 2 != len(exp) {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}

	result, err := CreateNumber(exp[0])
	if err != nil {
		return nil, err
	}
	prm, ok := exp[1].(Number)
	if !ok {
		return nil, NewRuntimeError("E1003", reflect.TypeOf(exp[1]).String())
	}
	if _, ok := result.(*Float); ok {
		if c, ok := prm.(*Integer); ok {
			prm = NewFloat(float64(c.Value))
		}
	}
	if org, ok := result.(*Integer); ok {
		if c, ok := exp[1].(*Float); ok {
			result = NewFloat(float64(org.Value))
			prm = c
		}
	}
	return NewBoolean(cmp(result, prm)), nil
}

// imul, skelton
func idivImpl(idivFunc func(int, int) int, exp ...Expression) (Expression, error) {
	if len(exp) != 2 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	var prm []*Integer
	for _, e := range exp {
		v, ok := e.(*Integer)
		if !ok {
			return nil, NewRuntimeError("E1002", reflect.TypeOf(e).String())
		}
		prm = append(prm, v)
	}
	if prm[1].Value == 0 {
		return nil, NewRuntimeError("E1013")
	}
	return NewInteger(idivFunc(prm[0].Value, prm[1].Value)), nil
}

// Build Global environement.
func buildOperationFunc() {

	buildInFuncTbl["+"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return calcOperate(func(a Number, b Number) Number { return a.Add(b) }, exp...)
			})
	}
	buildInFuncTbl["-"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return calcOperate(func(a Number, b Number) Number { return a.Sub(b) }, exp...)
			})
	}
	buildInFuncTbl["*"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return calcOperate(func(a Number, b Number) Number { return a.Mul(b) }, exp...)
			})
	}
	buildInFuncTbl["/"] = func(exp []Expression, env *SimpleEnv) (se Expression, e error) {
		// Not the best. But, Better than before.
		defer func() {
			if err := recover(); err != nil {
				if zero, ok := err.(*RuntimeError); ok {
					se = nil
					e = zero
				}
			}
		}()
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return calcOperate(func(a Number, b Number) Number { return a.Div(b) }, exp...)
			})
	}
	buildInFuncTbl["quotient"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return idivImpl(func(a int, b int) int { return a / b }, exp...)
			})
	}
	buildInFuncTbl["modulo"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return idivImpl(func(a int, b int) int { return a % b }, exp...)
			})
	}
	buildInFuncTbl[">"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return cmpOperate(func(a Number, b Number) bool { return a.Greater(b) }, exp...)
			})
	}
	buildInFuncTbl["<"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return cmpOperate(func(a Number, b Number) bool { return a.Less(b) }, exp...)
			})
	}
	buildInFuncTbl[">="] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return cmpOperate(func(a Number, b Number) bool { return a.GreaterEqual(b) }, exp...)
			})
	}
	buildInFuncTbl["<="] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return cmpOperate(func(a Number, b Number) bool { return a.LessEqual(b) }, exp...)
			})
	}
	buildInFuncTbl["="] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return cmpOperate(func(a Number, b Number) bool { return a.Equal(b) }, exp...)
			})
	}
}