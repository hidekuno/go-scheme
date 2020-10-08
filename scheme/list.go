/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

import (
	"bytes"
	"reflect"
	"strconv"
)

// List Type
type List struct {
	Expression
	Value []Expression
}

func NewList(exp []Expression) *List {
	l := new(List)
	l.Value = exp
	return l
}

func (self *List) String() string {
	var buffer bytes.Buffer
	var makeString func(*List)

	makeString = func(l *List) {
		buffer.WriteString("(")

		for _, i := range l.Value {
			if j, ok := i.(*List); ok {
				makeString(j)

			} else if j, ok := i.(Expression); ok {
				buffer.WriteString(j.String())
			}
			if i != l.Value[len(l.Value)-1] {
				buffer.WriteString(" ")
			}
		}
		buffer.WriteString(")")
	}
	makeString(self)
	return buffer.String()
}

// Pair Type
type Pair struct {
	Expression
	Car Expression
	Cdr Expression
}

func NewPair(car Expression, cdr Expression) *Pair {
	p := new(Pair)
	p.Car = car
	p.Cdr = cdr
	return p
}
func (self *Pair) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("(")
	buffer.WriteString(self.Car.String())
	buffer.WriteString(" . ")
	buffer.WriteString(self.Cdr.String())
	buffer.WriteString(")")
	return buffer.String()
}

// map,filter,reduce
func listFunc(lambda func(Expression, Expression, []Expression) ([]Expression, error), env *SimpleEnv, exp ...Expression) (Expression, error) {
	if len(exp) != 2 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	l, ok := exp[1].(*List)
	if !ok {
		return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[1]).String())
	}
	var result []Expression
	sexp := NewList(make([]Expression, 2))
	sexp.Value[0] = exp[0]

	quote := NewList(make([]Expression, 2))
	quote.Value[0] = NewBuildInFunc(buildInFuncTbl["quote"], "quote")

	for _, e := range l.Value {
		if _, ok = e.(*List); ok {
			quote.Value[1] = e
			sexp.Value[1] = quote

		} else if _, ok = e.(*Symbol); ok {
			sexp.Value[1] = e
			sexp.Value[1] = quote
		} else {
			sexp.Value[1] = e
		}
		v, err := eval(sexp, env)
		if err != nil {
			return nil, err
		}
		result, err = lambda(sexp.Value[1], v, result)
		if err != nil {
			return nil, err
		}
	}
	return NewList(result), nil
}

// Build Global environement.
func buildListFunc() {

	// list operator
	buildInFuncTbl["list"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				var l []Expression
				return NewList(append(l, exp...)), nil
			})
	}
	buildInFuncTbl["null?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 1 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if l, ok := exp[0].(*List); ok {
					return NewBoolean(0 == len(l.Value)), nil
				} else {
					return NewBoolean(false), nil
				}
			})
	}
	buildInFuncTbl["length"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 1 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if l, ok := exp[0].(*List); ok {
					return NewInteger(len(l.Value)), nil
				} else {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
			})
	}
	buildInFuncTbl["car"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 1 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if l, ok := exp[0].(*List); ok {
					if len(l.Value) <= 0 {
						return nil, NewRuntimeError("E1011", strconv.Itoa(len(l.Value)))
					}
					return l.Value[0], nil
				} else if p, ok := exp[0].(*Pair); ok {
					return p.Car, nil
				} else {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
			})
	}
	buildInFuncTbl["cdr"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 1 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if l, ok := exp[0].(*List); ok {
					if len(l.Value) <= 0 {
						return nil, NewRuntimeError("E1011", strconv.Itoa(len(l.Value)))
					}
					return NewList(l.Value[1:]), nil
				} else if p, ok := exp[0].(*Pair); ok {
					return p.Cdr, nil
				} else {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
			})
	}
	buildInFuncTbl["cadr"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 1 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if l, ok := exp[0].(*List); ok {
					if len(l.Value) < 2 {
						return nil, NewRuntimeError("E1011", strconv.Itoa(len(l.Value)))
					}
					return l.Value[1], nil
				} else {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
			})
	}
	buildInFuncTbl["cons"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 2 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if _, ok := exp[1].(*List); ok {
					var args []Expression
					args = append(args, exp[0])
					return NewList(append(args, (exp[1].(*List)).Value...)), nil
				}
				return NewPair(exp[0], exp[1]), nil
			})
	}
	buildInFuncTbl["append"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) < 2 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				var expList []Expression
				for _, e := range exp {
					if v, ok := e.(*List); ok {
						expList = append(expList, v.Value...)
					} else {
						return nil, NewRuntimeError("E1005", reflect.TypeOf(e).String())
					}
				}
				return NewList(expList), nil
			})
	}
	buildInFuncTbl["last"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 1 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if l, ok := exp[0].(*List); ok {
					if len(l.Value) <= 0 {
						return nil, NewRuntimeError("E1011", strconv.Itoa(len(l.Value)))
					}
					return l.Value[len(l.Value)-1], nil
				} else if p, ok := exp[0].(*Pair); ok {
					return p.Car, nil
				} else {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
			})
	}
	buildInFuncTbl["reverse"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 1 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if l, ok := exp[0].(*List); ok {
					if len(l.Value) <= 1 {
						return l, nil
					}
					args := make([]Expression, len(l.Value))
					idx := len(l.Value) - 1
					for _, c := range l.Value {
						args[idx] = c
						idx = idx - 1
					}
					return NewList(args), nil
				} else {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
			})
	}
	buildInFuncTbl["iota"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) < 1 || 3 < len(exp) {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				var l []Expression
				param := [3]int{0, 0, 1}
				for i := 0; i < len(exp); i++ {
					v, ok := exp[i].(*Integer)
					if !ok {
						return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[i]).String())
					}
					param[i] = v.Value
				}
				max, start, step := param[0], param[1], param[2]
				v := start
				for i := start; i < start+max; i++ {
					l = append(l, NewInteger(v))
					v += step
				}
				return NewList(l), nil
			})
	}
	buildInFuncTbl["map"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 2 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				lambda := func(org Expression, value Expression, result []Expression) ([]Expression, error) {
					return append(result, value), nil
				}
				return listFunc(lambda, env, exp...)
			})
	}
	buildInFuncTbl["for-each"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 2 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}

				fn, ok := exp[0].(*Function)
				if !ok {
					return nil, NewRuntimeError("E1006", reflect.TypeOf(exp[0]).String())
				}
				l, ok := exp[1].(*List)
				if !ok {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[1]).String())
				}
				param := make([]Expression, 1)
				for _, param[0] = range l.Value {
					_, err := fn.Execute(param, nil)
					if err != nil {
						return nil, err
					}
				}
				return NewNil(), nil
			})
	}
	buildInFuncTbl["filter"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 2 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				lambda := func(org Expression, value Expression, result []Expression) ([]Expression, error) {
					b, ok := value.(*Boolean)
					if !ok {
						return nil, NewRuntimeError("E1001", reflect.TypeOf(value).String())
					}
					if b.Value {
						return append(result, org), nil
					}
					return result, nil
				}
				return listFunc(lambda, env, exp...)
			})
	}
	buildInFuncTbl["reduce"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 3 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				fn, ok := exp[0].(*Function)
				if !ok {
					return nil, NewRuntimeError("E1006", reflect.TypeOf(exp[0]).String())
				}
				l, ok := exp[2].(*List)
				if !ok {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[1]).String())
				}
				if len(l.Value) == 0 {
					return exp[1], nil
				}
				param := make([]Expression, len(fn.ParamName.Value))
				result := l.Value[0]
				for _, c := range l.Value[1:] {
					param[0] = result
					param[1] = c
					r, err := fn.Execute(param, nil)
					result = r
					if err != nil {
						return nil, err
					}
				}
				return result, nil
			})
	}

}
