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
func calcOperate(exp []Expression, env *SimpleEnv, calc func(Number, Number) Number, x int) (Expression, error) {
	if 1 > len(exp) {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {

		result, err := CreateNumber(exp[0])
		if err != nil {
			return nil, err
		}
		if 1 == len(exp) {
			return calc(NewInteger(x), result), nil
		} else {
			for _, e := range exp[1:] {
				prm, ok := e.(Number)
				if !ok {
					return nil, NewRuntimeError("E1003", reflect.TypeOf(e).String())
				}
				result, prm = castNumber(result, prm)
				result = calc(result, prm)
			}
		}
		return result, nil
	})
}

// gt,lt,ge,le
func cmpOperate(exp []Expression, env *SimpleEnv, cmp func(Number, Number) bool) (Expression, error) {
	if 2 != len(exp) {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {

		result, err := CreateNumber(exp[0])
		if err != nil {
			return nil, err
		}
		prm, ok := exp[1].(Number)
		if !ok {
			return nil, NewRuntimeError("E1003", reflect.TypeOf(exp[1]).String())
		}

		result, prm = castNumber(result, prm)
		return NewBoolean(cmp(result, prm)), nil
	})
}

// max,min
func selectOne(exp []Expression, env *SimpleEnv, cmp func(Number, Number) bool) (Expression, error) {
	if 1 > len(exp) {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {

		result, err := CreateNumber(exp[0])
		if err != nil {
			return nil, err
		}
		for _, e := range exp[1:] {
			prm, ok := e.(Number)
			if !ok {
				return nil, NewRuntimeError("E1003", reflect.TypeOf(e).String())
			}
			result, prm = castNumber(result, prm)
			if cmp(result, prm) {
				result = prm
			}
		}
		return result, nil
	})
}

// imul, skelton
func idivImpl(exp []Expression, env *SimpleEnv, fn func(int, int) int) (Expression, error) {
	if len(exp) != 2 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
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
		return NewInteger(fn(prm[0].Value, prm[1].Value)), nil
	})
}

// ash
func shift(exp []Expression, env *SimpleEnv) (Expression, error) {
	if len(exp) != 2 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
		x, ok := exp[0].(*Integer)
		if !ok {
			return nil, NewRuntimeError("E1002", reflect.TypeOf(x).String())
		}
		y, ok := exp[1].(*Integer)
		if !ok {
			return nil, NewRuntimeError("E1002", reflect.TypeOf(y).String())
		}

		if y.Value > 0 {
			return NewInteger(x.Value << y.Value), nil
		} else {
			return NewInteger(x.Value >> (-1 * y.Value)), nil
		}
	})
}

// and,or,xor
func calcLogic(exp []Expression, env *SimpleEnv, calc func(a *Integer, b *Integer) int) (Expression, error) {

	if 0 >= len(exp) {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
		result, ok := exp[0].(*Integer)
		if !ok {
			return nil, NewRuntimeError("E1002", reflect.TypeOf(result).String())
		}
		for _, e := range exp[1:] {
			prm, ok := e.(*Integer)
			if !ok {
				return nil, NewRuntimeError("E1002", reflect.TypeOf(e).String())
			}
			result = NewInteger(calc(result, prm))
		}
		return result, nil
	})
}

// lognot
func lognot(exp []Expression, env *SimpleEnv) (Expression, error) {
	if len(exp) != 1 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
		v, ok := exp[0].(*Integer)
		if !ok {
			return nil, NewRuntimeError("E1002", reflect.TypeOf(v).String())
		}
		return NewInteger(^v.Value), nil
	})
}

// Build Global environement.
func buildOperationFunc() {

	buildInFuncTbl["+"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return calcOperate(exp, env, func(a Number, b Number) Number { return a.Add(b) }, 0)
	}
	buildInFuncTbl["-"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return calcOperate(exp, env, func(a Number, b Number) Number { return a.Sub(b) }, 0)
	}
	buildInFuncTbl["*"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return calcOperate(exp, env, func(a Number, b Number) Number { return a.Mul(b) }, 1)
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
		return calcOperate(exp, env, func(a Number, b Number) Number { return a.Div(b) }, 1)
	}
	buildInFuncTbl["ash"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return shift(exp, env)
	}
	buildInFuncTbl["logand"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return calcLogic(exp, env, func(a *Integer, b *Integer) int { return a.Value & b.Value })
	}
	buildInFuncTbl["logior"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return calcLogic(exp, env, func(a *Integer, b *Integer) int { return a.Value | b.Value })
	}
	buildInFuncTbl["logxor"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return calcLogic(exp, env, func(a *Integer, b *Integer) int { return a.Value ^ b.Value })
	}
	buildInFuncTbl["lognot"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return lognot(exp, env)
	}
	buildInFuncTbl["max"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return selectOne(exp, env, func(a Number, b Number) bool { return b.Greater(a) })
	}
	buildInFuncTbl["min"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return selectOne(exp, env, func(a Number, b Number) bool { return b.Less(a) })
	}
	buildInFuncTbl["quotient"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return idivImpl(exp, env, func(a int, b int) int { return a / b })
	}
	buildInFuncTbl["modulo"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return idivImpl(exp, env, func(a int, b int) int { return a % b })
	}
	buildInFuncTbl[">"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return cmpOperate(exp, env, func(a Number, b Number) bool { return a.Greater(b) })
	}
	buildInFuncTbl["<"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return cmpOperate(exp, env, func(a Number, b Number) bool { return a.Less(b) })
	}
	buildInFuncTbl[">="] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return cmpOperate(exp, env, func(a Number, b Number) bool { return a.GreaterEqual(b) })
	}
	buildInFuncTbl["<="] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return cmpOperate(exp, env, func(a Number, b Number) bool { return a.LessEqual(b) })
	}
	buildInFuncTbl["="] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return cmpOperate(exp, env, func(a Number, b Number) bool { return a.Equal(b) })
	}
}
