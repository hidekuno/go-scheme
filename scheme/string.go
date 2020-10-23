/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

import (
	"bytes"
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
func stringCompare(operate func(string, string) bool, exp ...Expression) (Expression, error) {
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
func stringLength(fn func(string) int, exp ...Expression) (Expression, error) {
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
			return stringCompare(func(x string, y string) bool { return x == y }, exp...)
		})
	}
	buildInFuncTbl["string<?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringCompare(func(x string, y string) bool { return x < y }, exp...)
		})
	}
	buildInFuncTbl["string>?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringCompare(func(x string, y string) bool { return x > y }, exp...)
		})
	}
	buildInFuncTbl["string<=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringCompare(func(x string, y string) bool { return x <= y }, exp...)
		})
	}
	buildInFuncTbl["string>=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringCompare(func(x string, y string) bool { return x >= y }, exp...)
		})
	}
	buildInFuncTbl["string-ci=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringCompare(func(x string, y string) bool { return strings.ToLower(x) == strings.ToLower(y) }, exp...)
		})
	}
	buildInFuncTbl["string-ci<?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringCompare(func(x string, y string) bool { return strings.ToLower(x) < strings.ToLower(y) }, exp...)
		})
	}
	buildInFuncTbl["string-ci>?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringCompare(func(x string, y string) bool { return strings.ToLower(x) > strings.ToLower(y) }, exp...)
		})
	}
	buildInFuncTbl["string-ci<=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringCompare(func(x string, y string) bool { return strings.ToLower(x) <= strings.ToLower(y) }, exp...)
		})
	}
	buildInFuncTbl["string-ci>=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringCompare(func(x string, y string) bool { return strings.ToLower(x) >= strings.ToLower(y) }, exp...)
		})
	}
	buildInFuncTbl["string-length"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringLength(func(x string) int { return utf8.RuneCountInString(x) }, exp...)
		})
	}
	buildInFuncTbl["string-size"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringLength(func(x string) int { return len(x) }, exp...)
		})
	}
	buildInFuncTbl["number->string"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			if len(exp) != 1 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			if _, ok := exp[0].(Number); !ok {
				return nil, NewRuntimeError("E1003", reflect.TypeOf(exp[0]).String())
			}
			return NewString(exp[0].String()), nil
		})
	}
	buildInFuncTbl["string->number"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			if len(exp) != 1 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			s, ok := exp[0].(*String)
			if !ok {
				return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())

			}
			if i, err := strconv.Atoi(s.Value); err == nil {
				return NewInteger(i), nil
			} else if f, err := strconv.ParseFloat(s.Value, 64); err == nil {
				return NewFloat(f), nil
			}
			return nil, NewRuntimeError("E1003", s.Value)
		})
	}
	buildInFuncTbl["list->string"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			var buffer bytes.Buffer

			if len(exp) != 1 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			l, ok := exp[0].(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())

			}
			for _, e := range l.Value {
				c, ok := e.(*Char)
				if !ok {
					return nil, NewRuntimeError("E1019", reflect.TypeOf(e).String())
				}
				buffer.WriteRune(c.Value)
			}
			return NewString(buffer.String()), nil
		})
	}
	buildInFuncTbl["string->list"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			if len(exp) != 1 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			s, ok := exp[0].(*String)
			if !ok {
				return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())

			}
			l := make([]Expression, 0, len(s.Value))
			for _, c := range s.Value {
				l = append(l, NewCharFromRune(rune(c)))
			}
			return NewList(l), nil
		})
	}
	buildInFuncTbl["substring"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			if len(exp) != 3 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			s, ok := exp[0].(*String)
			if !ok {
				return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())

			}
			from, ok := exp[1].(*Integer)
			if !ok {
				return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[1]).String())

			}
			to, ok := exp[2].(*Integer)
			if !ok {
				return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[2]).String())

			}
			if from.Value < 0 || to.Value > utf8.RuneCountInString(s.Value) || from.Value > to.Value {
				return nil, NewRuntimeError("E1021", from.String(), to.String())
			}
			return NewString(
				string(
					[]rune(s.Value)[from.Value:to.Value])), nil
		})
	}
	buildInFuncTbl["symbol->string"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {

			if len(exp) != 1 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			s, ok := exp[0].(*Symbol)
			if !ok {
				return nil, NewRuntimeError("E1004", reflect.TypeOf(exp[0]).String())

			}
			return NewString(s.Value), nil
		})
	}
	buildInFuncTbl["string->symbol"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {

			if len(exp) != 1 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			s, ok := exp[0].(*String)
			if !ok {
				return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())

			}
			return NewSymbol(s.Value), nil
		})
	}
}
