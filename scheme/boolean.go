/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

import (
	"fmt"
	"reflect"
	"strconv"
)

// Boolean Type
type Boolean struct {
	Expression
	Value bool
	exp   string
}

func NewBoolean(v bool) *Boolean {
	b := new(Boolean)
	b.Value = v
	if v {
		b.exp = "#t"
	} else {
		b.exp = "#f"
	}
	return b
}
func (self *Boolean) String() string {
	return self.exp
}
func (self *Boolean) Print() {
	fmt.Print(self.String())
}
func (self *Boolean) isAtom() bool {
	return true
}
func (self *Boolean) clone() Expression {
	return NewBoolean(self.Value)
}
func (self *Boolean) equalValue(e Expression) bool {
	if v, ok := e.(*Boolean); ok {
		return self.Value == v.Value
	}
	return false
}
func doBoolean(exp []Expression, env *SimpleEnv, f func(bool) bool) (Expression, error) {
	if len(exp) != 1 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env,
		func(exp ...Expression) (Expression, error) {
			if v, ok := exp[0].(*Boolean); ok {
				return NewBoolean(f(v.Value)), nil
			}
			return NewBoolean(true), nil
		})
}
func doBooleanEq(exp []Expression, env *SimpleEnv) (Expression, error) {
	if len(exp) < 2 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env,
		func(exp ...Expression) (Expression, error) {
			var b bool
			if v, ok := exp[0].(*Boolean); ok {
				b = v.Value
			} else {
				return nil, NewRuntimeError("E1001", reflect.TypeOf(exp[0]).String())
			}
			for _, e := range exp[1:] {
				if v, ok := e.(*Boolean); ok {
					if b != v.Value {
						return NewBoolean(false), nil
					}
				} else {
					return nil, NewRuntimeError("E1001", reflect.TypeOf(e).String())
				}
			}
			return NewBoolean(true), nil
		})
}

// Build Global environement.
func buildBooleanFunc() {
	buildInFuncTbl["not"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return doBoolean(exp, env, func(b bool) bool { return !b })
	}
	buildInFuncTbl["boolean"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return doBoolean(exp, env, func(b bool) bool { return b })
	}
	buildInFuncTbl["boolean=?"] = doBooleanEq
}
