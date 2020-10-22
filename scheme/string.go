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
	"strings"
	"unicode/utf8"
)

// String Type
type String struct {
	Expression
	Value string
}

func NewString(p string) *String {
	v := new(String)
	v.Value = p
	return v
}

func (self *String) String() string {
	return "\"" + self.Value + "\""
}
func (self *String) isAtom() bool {
	return true
}
func (self *String) clone() Expression {
	return NewString(self.Value)
}
func (self *String) equalValue(e Expression) bool {
	if v, ok := e.(*String); ok {
		return self.Value == v.Value
	}
	return false
}
func strcmp(operate func(string, string) bool, exp ...Expression) (Expression, error) {
	if len(exp) != 2 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	x, ok := exp[0].(*String)
	if !ok {
		return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
	}
	y, ok := exp[1].(*String)
	if !ok {
		return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[1]).String())
	}
	return NewBoolean(operate(x.Value, y.Value)), nil
}
func strlen(fn func(string) int, exp ...Expression) (Expression, error) {
	if len(exp) != 1 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	x, ok := exp[0].(*String)
	if !ok {
		return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
	}
	return NewInteger(fn(x.Value)), nil
}

// Build Global environement.
func buildStringFunc() {
	buildInFuncTbl["string-append"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			if len(exp) < 2 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			ret := make([]string, 0, len(exp))
			for _, e := range exp {
				s, ok := e.(*String)
				if !ok {
					return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
				}
				ret = append(ret, s.Value)
			}
			return NewString(strings.Join(ret, "")), nil
		})
	}
	buildInFuncTbl["format"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			if len(exp) != 2 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			f, ok := exp[0].(*String)
			if !ok {
				return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
			}
			n, ok := exp[1].(*Integer)
			if !ok {
				return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[1]).String())
			}

			s := ""
			switch f.Value {
			case "~d", "~D":
				s = fmt.Sprintf("%d", n.Value)
			case "~o", "~O":
				s = fmt.Sprintf("%o", n.Value)
			case "~b", "~B":
				s = fmt.Sprintf("%b", n.Value)
			case "~x", "~X":
				s = fmt.Sprintf("%"+string(f.Value[1]), n.Value)
			default:
				return nil, NewRuntimeError("E1018")
			}
			return NewString(s), nil
		})
	}
	buildInFuncTbl["string=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return strcmp(func(x string, y string) bool { return x == y }, exp...)
		})
	}
	buildInFuncTbl["string<?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return strcmp(func(x string, y string) bool { return x < y }, exp...)
		})
	}
	buildInFuncTbl["string>?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return strcmp(func(x string, y string) bool { return x > y }, exp...)
		})
	}
	buildInFuncTbl["string<=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return strcmp(func(x string, y string) bool { return x <= y }, exp...)
		})
	}
	buildInFuncTbl["string>=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return strcmp(func(x string, y string) bool { return x >= y }, exp...)
		})
	}
	buildInFuncTbl["string-ci=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return strcmp(func(x string, y string) bool { return strings.ToLower(x) == strings.ToLower(y) }, exp...)
		})
	}
	buildInFuncTbl["string-ci<?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return strcmp(func(x string, y string) bool { return strings.ToLower(x) < strings.ToLower(y) }, exp...)
		})
	}
	buildInFuncTbl["string-ci>?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return strcmp(func(x string, y string) bool { return strings.ToLower(x) > strings.ToLower(y) }, exp...)
		})
	}
	buildInFuncTbl["string-ci<=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return strcmp(func(x string, y string) bool { return strings.ToLower(x) <= strings.ToLower(y) }, exp...)
		})
	}
	buildInFuncTbl["string-ci>=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return strcmp(func(x string, y string) bool { return strings.ToLower(x) >= strings.ToLower(y) }, exp...)
		})
	}
	buildInFuncTbl["string-length"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return strlen(func(x string) int { return utf8.RuneCountInString(x) }, exp...)
		})
	}
	buildInFuncTbl["string-size"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return strlen(func(x string) int { return len(x) }, exp...)
		})
	}
}
