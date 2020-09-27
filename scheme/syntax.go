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

// and or not
func logicalOperate(exp []Expression, env *SimpleEnv, bcond bool, bret bool) (Expression, error) {
	if len(exp) < 2 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	for _, e := range exp {
		b, err := eval(e, env)
		if err != nil {
			return nil, err
		}
		if _, ok := b.(*Boolean); !ok {
			return nil, NewRuntimeError("E1001", reflect.TypeOf(b).String())
		}
		if bcond == (b.(*Boolean)).Value {
			return NewBoolean(bcond), nil
		}
	}
	return NewBoolean(bret), nil
}

// Build Global environement.
func buildSyntaxFunc() {

	// syntax keyword implements
	buildInFuncTbl["if"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		e, err := eval(exp[0], env)
		if err != nil {
			return nil, err
		}

		b, ok := e.(*Boolean)
		if !ok {
			return nil, NewRuntimeError("E1001", reflect.TypeOf(exp).String())
		}
		if b.Value {
			return eval(exp[1], env)
		} else if 3 <= len(exp) {
			return eval(exp[2], env)
		}
		return NewNil(), nil
	}
	buildInFuncTbl["define"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		var (
			key *Symbol
			err error
			ok  bool
			e   Expression
		)
		if key, ok = exp[0].(*Symbol); ok {
			e, err = eval(exp[1], env)
			if err != nil {
				return nil, err
			}
		} else if l, ok := exp[0].(*List); ok {
			if len(l.Value) < 1 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			key, ok = l.Value[0].(*Symbol)
			if !ok {
				return nil, NewRuntimeError("E1004", reflect.TypeOf(exp[0]).String())
			}
			for _, p := range l.Value[1:] {
				e, ok := p.(*Symbol)
				if !ok {
					return nil, NewRuntimeError("E1004", reflect.TypeOf(e).String())
				}
			}
			e = NewFunction(env, NewList(l.Value[1:]), exp[1:], key.Value)
		} else {
			return nil, NewRuntimeError("E1004", reflect.TypeOf(exp[0]).String())
		}
		(*env).Regist(key.Value, e)
		return key, nil
	}
	buildInFuncTbl["lambda"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		// if l == (), internal list implements
		l, ok := exp[0].(*List)
		if !ok {
			return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
		}
		for _, p := range l.Value {
			e, ok := p.(*Symbol)
			if !ok {
				return nil, NewRuntimeError("E1004", reflect.TypeOf(e).String())
			}
		}
		return NewFunction(env, l, exp[1:], "lambda"), nil
	}
	buildInFuncTbl["set!"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		s, ok := exp[0].(*Symbol)
		if !ok {
			return nil, NewRuntimeError("E1004", reflect.TypeOf(exp[0]).String())
		}

		e, err := eval(exp[1], env)
		if err != nil {
			return exp[1], err
		}
		if _, ok := (*env).Find(s.Value); !ok {
			return e, NewRuntimeError("E1008", s.Value)
		}
		(*env).Set(s.Value, e)
		return e, nil
	}
	buildInFuncTbl["let"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		var letsym *Symbol
		var pname []Expression
		body := 1

		l, ok := exp[0].(*List)
		if ok && len(exp) < 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		if !ok {
			letsym, ok = exp[0].(*Symbol)
			if !ok {
				return nil, NewRuntimeError("E1004", reflect.TypeOf(exp[0]).String())
			}
			if len(exp) < 3 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			l, ok = exp[1].(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[1]).String())
			}
			body = 2
		}

		localEnv := Environment{}
		for _, let := range l.Value {
			r, ok := let.(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(let).String())
			}
			if len(r.Value) != 2 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(r.Value)))
			}
			sym, ok := (r.Value[0]).(*Symbol)
			if !ok {
				return nil, NewRuntimeError("E1004", reflect.TypeOf(r.Value[0]).String())
			}
			v, err := eval(r.Value[1], env)
			if err != nil {
				return nil, err
			}
			pname = append(pname, sym)
			localEnv[sym.Value] = v
		}

		nse := NewSimpleEnv(env, &localEnv)
		if letsym != nil {
			localEnv[letsym.Value] = NewFunction(nse, NewList(pname), exp[body:], letsym.Value)
		}
		var lastExp Expression
		for i := body; i < len(exp); i++ {
			if e, err := eval(exp[i], nse); err == nil {
				lastExp = e
			} else {
				return nil, err
			}
		}
		return lastExp, nil

	}
	buildInFuncTbl["and"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return logicalOperate(exp, env, false, true)
	}
	buildInFuncTbl["or"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return logicalOperate(exp, env, true, false)
	}
	buildInFuncTbl["not"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 1 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if _, ok := exp[0].(*Boolean); !ok {
					return nil, NewRuntimeError("E1001", reflect.TypeOf(exp[0]).String())
				}
				return NewBoolean(!(exp[0].(*Boolean)).Value), nil
			})
	}
	buildInFuncTbl["delay"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return NewPromise(env, exp[0]), nil
	}
	buildInFuncTbl["force"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		e, err := eval(exp[0], env)
		if err != nil {
			return nil, err
		}
		p, ok := e.(*Promise)
		if !ok {
			return e, nil
		} else {
			return eval(p.Body, p.Env)
		}
	}
	buildInFuncTbl["cond"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		for _, e := range exp {
			l, ok := e.(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(e).String())
			}
			if len(l.Value) < 2 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(l.Value)))
			}
			if _, ok := l.Value[0].(*List); ok {
				b, err := eval(l.Value[0], env)
				if err != nil {
					return nil, err
				}
				if b, ok := b.(*Boolean); ok {
					if b.Value {
						return evalMulti(l.Value[1:], env)
					}
				} else {
					return nil, NewRuntimeError("E1001")
				}
			} else if sym, ok := l.Value[0].(*Symbol); ok {
				if sym.Value == "else" {
					return evalMulti(l.Value[1:], env)
				} else {
					return nil, NewRuntimeError("E1012")
				}
			} else {
				return nil, NewRuntimeError("E1012")
			}
		}
		return NewNil(), nil
	}
	buildInFuncTbl["case"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}

		param := make([]Expression, 2)
		v, err := eval(exp[0], env)
		if err != nil {
			return v, err
		}
		param[0] = v

		if 2 <= len(exp) {
			for _, e := range exp[1:] {
				l, ok := e.(*List)
				if !ok {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(e).String())
				}
				if len(l.Value) < 2 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(l.Value)))
				}
				if p, ok := l.Value[0].(*List); ok {
					for _, e := range p.Value {
						v, err := eval(e, env)
						if err != nil {
							return nil, err
						}
						param[1] = v
						if b, err := eqv(param, env); err == nil {
							v, _ := b.(*Boolean)
							if v.Value == true {
								return evalMulti(l.Value[1:], env)
							}
						}
					}
				} else if sym, ok := l.Value[0].(*Symbol); ok {
					if sym.Value == "else" {
						return evalMulti(l.Value[1:], env)
					} else {
						return nil, NewRuntimeError("E1012")
					}
				} else {
					return nil, NewRuntimeError("E1012")
				}
			}
		}
		return NewNil(), nil
	}
	buildInFuncTbl["quote"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return exp[0], nil
	}
	buildInFuncTbl["begin"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {

		if len(exp) < 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return evalMulti(exp, env)
	}
}
