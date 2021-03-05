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

// Continuation
type Continuation struct {
	Expression
	cont  Expression
	Value Expression
	Name  string
}

func NewContinuation(cont Expression) *Continuation {
	c := new(Continuation)
	c.cont = cont
	c.Value = nil
	return c
}
func (self *Continuation) String() string {
	return "Continuation: "
}
func (self *Continuation) isAtom() bool {
	return true
}
func (self *Continuation) equalValue(e Expression) bool {
	// Not Support this method
	return false
}
func (self *Continuation) clone() Expression {
	return NewContinuation(self.cont)
}
func (err *Continuation) Error() string {
	return "Continuation"
}
func (self *Continuation) Execute(exp []Expression, env *SimpleEnv) (Expression, error) {
	if len(exp) != 2 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	v, err := eval(exp[1], env)
	if err != nil {
		return nil, err
	}
	self.Value = v
	if s, ok := exp[0].(*Symbol); ok {
		self.Name = s.Value
	}
	return nil, self
}

func Quote(exp []Expression, env *SimpleEnv) (Expression, error) {
	if len(exp) != 1 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return exp[0], nil
}
func Begin(exp []Expression, env *SimpleEnv) (Expression, error) {
	if len(exp) < 1 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	var (
		v   Expression
		err error
	)
	for _, e := range exp {
		v, err = eval(e, env)
		if err != nil {
			return v, err
		}
	}
	return v, nil
}

// and or not
func doLogicalOperate(exp []Expression, env *SimpleEnv, bcond bool, bret bool) (Expression, error) {
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
		if len(exp) < 2 {
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

		if len(exp) == 0 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
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
		return doLogicalOperate(exp, env, false, true)
	}
	buildInFuncTbl["or"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return doLogicalOperate(exp, env, true, false)
	}
	buildInFuncTbl["not"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
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
						return Begin(l.Value[1:], env)
					}
				} else {
					return nil, NewRuntimeError("E1001")
				}
			} else if sym, ok := l.Value[0].(*Symbol); ok {
				if sym.Value == "else" {
					return Begin(l.Value[1:], env)
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
				if len(l.Value) < 1 {
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
								if len(l.Value) < 2 {
									return l.Value[0], nil
								}
								return Begin(l.Value[1:], env)
							}
						}
					}
				} else if sym, ok := l.Value[0].(*Symbol); ok {
					if sym.Value == "else" {
						if len(l.Value) < 2 {
							return NewInteger(0), nil
						}
						return Begin(l.Value[1:], env)
					} else {
						return nil, NewRuntimeError("E1017")
					}
				} else {
					return nil, NewRuntimeError("E1017")
				}
			}
		}
		return NewNil(), nil
	}
	buildInFuncTbl["apply"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}

		e, err := eval(exp[1], env)
		if err != nil {
			return nil, err
		}
		l, ok := e.(*List)
		if !ok {
			return nil, NewRuntimeError("E1005", reflect.TypeOf(l).String())
		}
		sexp := MakeQuotedValue(exp[0], l.Value, nil)
		return eval(sexp, env)

	}
	buildInFuncTbl["quote"] = Quote
	buildInFuncTbl["begin"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {

		if len(exp) < 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return Begin(exp, env)
	}
	buildInFuncTbl["do"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {

		if len(exp) < 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		l, ok := exp[0].(*List)
		if !ok {
			return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
		}
		if len(l.Value) < 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		localEnv := Environment{}
		param := make([]string, 0, len(l.Value))
		update := make([]Expression, 0, len(l.Value))
		for _, e := range l.Value {
			f, ok := e.(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
			}
			if len(f.Value) != 3 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(f.Value)))
			}
			if s, ok := f.Value[0].(*Symbol); ok {
				v, err := eval(f.Value[1], env)
				if err != nil {
					return nil, err
				}
				localEnv[s.Value] = v
				param = append(param, s.Value)
			} else {
				return nil, NewRuntimeError("E1004", reflect.TypeOf(f.Value[0]).String())
			}
			update = append(update, f.Value[2])
		}

		l, ok = exp[1].(*List)
		if !ok {
			return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[1]).String())
		}
		if len(l.Value) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(l.Value)))
		}
		cond := make([]Expression, 0, 2)
		for _, c := range l.Value {
			cond = append(cond, c)
		}
		nse := NewSimpleEnv(env, &localEnv)
		for {
			// eval condition
			e, err := eval(cond[0], nse)
			if err != nil {
				return nil, err
			}
			if b, ok := e.(*Boolean); ok {
				if b.Value {
					if e, err := eval(cond[1], nse); err == nil {
						return e, nil
					} else {
						return nil, err
					}
				}
			} else {
				return nil, NewRuntimeError("E1001", reflect.TypeOf(e).String())
			}

			// eval body
			for i := 2; i <= len(exp)-1; i++ {
				_, err := eval(exp[i], nse)
				if err != nil {
					return nil, err
				}
			}

			// eval step
			result := make([]Expression, 0, len(update))
			for _, u := range update {
				v, err := eval(u, nse)
				if err != nil {
					return nil, err
				}
				result = append(result, v)
			}
			for i, v := range result {
				nse.Regist(param[i], v)
			}
		}
	}
	buildInFuncTbl["call/cc"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		sexp := make([]Expression, 0)
		sexp = append(sexp, NewContinuation(nil))
		e, err := eval(exp[0], env)
		if err != nil {
			return nil, NewRuntimeError("E1006")
		}
		if fn, ok := e.(*Function); ok {
			return fn.Execute(sexp, env)
		}
		return nil, NewRuntimeError("E1006")
	}
}
