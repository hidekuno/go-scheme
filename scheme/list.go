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
	"sort"
	"strconv"
)

type Sequence interface {
	GetValue() []Expression
}

func buildString(seq Sequence) string {
	var buffer bytes.Buffer
	var makeString func(Sequence)

	makeString = func(seq Sequence) {
		buffer.WriteString("(")

		for i, e := range seq.GetValue() {
			if j, ok := e.(Sequence); ok {
				makeString(j)

			} else if j, ok := e.(Expression); ok {
				buffer.WriteString(j.String())
			}
			if i != len(seq.GetValue())-1 {
				buffer.WriteString(" ")
			}
		}
		buffer.WriteString(")")
	}
	makeString(seq)
	return buffer.String()
}

// List Type
type List struct {
	Expression
	Sequence
	Value []Expression
}

func NewList(exp []Expression) *List {
	l := new(List)
	l.Value = exp
	return l
}
func (self *List) String() string {
	return buildString(self)
}
func (self *List) Print() {
	fmt.Print(self.String())
}
func (self *List) isAtom() bool {
	return false
}
func (self *List) clone() Expression {
	return NewList(self.Value)
}
func (self *List) equalValue(e Expression) bool {
	// Not Support this method
	return false
}
func (self *List) GetValue() []Expression {
	return self.Value
}

// Vector Type
type Vector struct {
	Expression
	Sequence
	List
}

func NewVector(exp []Expression) *Vector {
	v := new(Vector)
	v.Value = exp
	return v
}
func (self *Vector) String() string {
	return "#" + buildString(self)
}
func (self *Vector) Print() {
	fmt.Print(self.String())
}
func (self *Vector) clone() Expression {
	return NewVector(self.Value)
}
func (self *Vector) isAtom() bool {
	return false
}
func (self *Vector) equalValue(e Expression) bool {
	// Not Support this method
	return false
}
func (self *Vector) GetValue() []Expression {
	return self.Value
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
func (self *Pair) Print() {
	fmt.Print(self.String())
}
func (self *Pair) isAtom() bool {
	return true
}
func (self *Pair) clone() Expression {
	return NewPair(self.Car, self.Cdr)
}
func (self *Pair) equalValue(e Expression) bool {
	// Not Support this method
	return false
}
func MakeQuotedValue(fn Expression, l []Expression, result Expression) *List {
	size := 4
	if len(l) > size {
		size = len(l) + 1
	}

	sexp := NewList(make([]Expression, 0, size))
	sexp.Value = append(sexp.Value, fn)

	if result != nil {
		quote := NewList(make([]Expression, 2))
		quote.Value[0] = NewBuildInFunc(Quote, "quote")

		if _, ok := result.(*List); ok {
			quote.Value[1] = result
			sexp.Value = append(sexp.Value, quote)

		} else if _, ok := result.(*Symbol); ok {
			quote.Value[1] = result
			sexp.Value = append(sexp.Value, quote)

		} else {
			sexp.Value = append(sexp.Value, result)
		}
	}
	for _, e := range l {
		quote := NewList(make([]Expression, 2))
		quote.Value[0] = NewBuildInFunc(Quote, "quote")

		if _, ok := e.(*List); ok {
			quote.Value[1] = e
			sexp.Value = append(sexp.Value, quote)

		} else if _, ok := e.(*Symbol); ok {
			quote.Value[1] = e
			sexp.Value = append(sexp.Value, quote)

		} else {
			sexp.Value = append(sexp.Value, e)
		}
	}
	return sexp
}

// map,filter
func doListFunc(exp []Expression, env *SimpleEnv,
	lambda func(Expression, Expression, []Expression) ([]Expression, error)) (Expression, error) {
	if len(exp) != 2 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env,
		func(exp ...Expression) (Expression, error) {
			l, ok := exp[1].(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[1]).String())
			}
			var result []Expression

			for _, e := range l.Value {
				sexp := MakeQuotedValue(exp[0], []Expression{e}, nil)
				v, err := eval(sexp, env)
				if err != nil {
					return nil, err
				}
				result, err = lambda(e, v, result)
				if err != nil {
					return nil, err
				}
			}
			return NewList(result), nil
		})

}
func subList(exp []Expression, env *SimpleEnv, fn func(*List, *Integer) (int, int)) (Expression, error) {
	if len(exp) != 2 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env,
		func(exp ...Expression) (Expression, error) {

			l, ok := exp[0].(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
			}
			n, ok := exp[1].(*Integer)
			if !ok {
				return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[1]).String())
			}
			if n.Value < 0 || len(l.Value) < n.Value {
				return nil, NewRuntimeError("E1011", strconv.Itoa(n.Value))
			}

			x, y := fn(l, n)
			result := make([]Expression, y-x)
			copy(result, l.Value[x:y])
			return NewList(result), nil
		})
}
func defaultSortCallback(l, m Expression) bool {
	if x, ok := l.(Number); ok {
		if y, ok := m.(Number); ok {
			a, b := castNumber(x, y)
			return a.LessEqual(b)
		}
	}
	if x, ok := l.(*Char); ok {
		if y, ok := m.(*Char); ok {
			return x.Value < y.Value
		}
	}
	if x, ok := l.(*String); ok {
		if y, ok := m.(*String); ok {
			return x.Value < y.Value
		}
	}
	return false
}
func sortCallback(l, m, exp Expression, env *SimpleEnv) bool {
	sexp := make([]Expression, 0)
	sexp = append(sexp, exp)
	sexp = append(sexp, l)
	sexp = append(sexp, m)
	v, err := eval(NewList(sexp), env)
	if err != nil {
		return false
	}
	if b, ok := v.(*Boolean); ok {
		return b.Value
	}
	return false
}

func copySort(exp []Expression,
	env *SimpleEnv,
	sortImpl func(interface{}, func(int, int) bool)) (Expression, error) {

	if 1 > len(exp) || 2 < len(exp) {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env,
		func(exp ...Expression) (Expression, error) {
			l, ok := exp[0].(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
			}
			m := make([]Expression, len(l.Value))
			copy(m, l.Value)
			return doSort(env, NewList(m), sortImpl, exp)
		})

}
func effectSort(exp []Expression,
	env *SimpleEnv,
	sortImpl func(interface{}, func(int, int) bool)) (Expression, error) {

	if 1 > len(exp) || 2 < len(exp) {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env,
		func(exp ...Expression) (Expression, error) {
			l, ok := exp[0].(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
			}
			return doSort(env, l, sortImpl, exp)
		})
}
func doSort(env *SimpleEnv,
	l *List,
	sortImpl func(interface{}, func(int, int) bool),
	exp []Expression) (Expression, error) {

	if len(exp) == 1 {
		sortImpl(l.Value, func(i, j int) bool {
			return defaultSortCallback(l.Value[i], l.Value[j])
		})

	} else if len(exp) == 2 {
		if _, ok := exp[1].(*BuildInFunc); !ok {
			if _, ok := exp[1].(*Function); !ok {
				return nil, NewRuntimeError("E1006", reflect.TypeOf(exp[1]).String())
			}
		}
		sortImpl(l.Value, func(i, j int) bool {
			return sortCallback(l.Value[i], l.Value[j], exp[1], env)
		})
	}
	return l, nil
}
func merge(exp []Expression, env *SimpleEnv) (Expression, error) {
	merge_iter := func(l *List, m *List) (Expression, error) {
		v := make([]Expression, 0)
		i, j := 0, 0

		for {
			if len(l.Value) <= i || len(m.Value) <= j {
				break
			}
			if x, ok := l.Value[i].(Number); ok {
				if y, ok := m.Value[j].(Number); ok {
					a, b := castNumber(x, y)
					if a.LessEqual(b) {
						v = append(v, l.Value[i])
						i++
					} else {
						v = append(v, m.Value[j])
						j++
					}
					continue
				}
			} else if x, ok := l.Value[i].(*Char); ok {
				if y, ok := l.Value[j].(*Char); ok {
					if x.Value < y.Value {
						v = append(v, l.Value[i])
						i++
					} else {
						v = append(v, m.Value[j])
						j++
					}
					continue
				}
			} else if x, ok := l.Value[i].(*String); ok {
				if y, ok := l.Value[j].(*String); ok {
					if x.Value < y.Value {
						v = append(v, l.Value[i])
						i++
					} else {
						v = append(v, m.Value[j])
						j++
					}
					continue
				}
			}
			v = append(v, l.Value[i])
			i++
		}
		v = append(v, l.Value[i:]...)
		v = append(v, m.Value[j:]...)
		return NewList(v), nil
	}
	merge_iter_by := func(l *List, m *List, exp Expression) (Expression, error) {
		v := make([]Expression, 0)
		i, j := 0, 0

		if _, ok := exp.(*BuildInFunc); !ok {
			if _, ok := exp.(*Function); !ok {
				return nil, NewRuntimeError("E1006", reflect.TypeOf(exp).String())
			}
		}
		for {
			if len(l.Value) <= i || len(m.Value) <= j {
				break
			}
			sexp := make([]Expression, 0)
			sexp = append(sexp, exp)
			sexp = append(sexp, l.Value[i])
			sexp = append(sexp, m.Value[j])
			e, err := eval(NewList(sexp), env)
			if err != nil {
				v = append(v, l.Value[i])
				i++
			} else if b, ok := e.(*Boolean); ok {
				if b.Value {
					v = append(v, l.Value[i])
					i++
				} else {
					v = append(v, m.Value[j])
					j++
				}
			} else {
				return nil, NewRuntimeError("E1001", reflect.TypeOf(exp).String())
			}
		}
		v = append(v, l.Value[i:]...)
		v = append(v, m.Value[j:]...)
		return NewList(v), nil
	}
	if 2 > len(exp) || 3 < len(exp) {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env,
		func(exp ...Expression) (Expression, error) {

			l, ok := exp[0].(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
			}
			m, ok := exp[1].(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[1]).String())
			}
			if len(exp) == 3 {
				return merge_iter_by(l, m, exp[2])
			} else {
				return merge_iter(l, m)
			}
		})
}
func isSorted(exp []Expression, env *SimpleEnv) (Expression, error) {
	if 1 > len(exp) || 2 < len(exp) {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env,
		func(exp ...Expression) (Expression, error) {
			l, ok := exp[0].(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
			}
			var b bool
			if len(exp) == 1 {
				b = sort.SliceIsSorted(l.Value, func(i, j int) bool {
					return defaultSortCallback(l.Value[i], l.Value[j])
				})
			} else if len(exp) == 2 {
				if _, ok := exp[1].(*BuildInFunc); !ok {
					if _, ok := exp[1].(*Function); !ok {
						return nil, NewRuntimeError("E1006", reflect.TypeOf(exp[1]).String())
					}
				}
				b = sort.SliceIsSorted(l.Value, func(i, j int) bool {
					return sortCallback(l.Value[i], l.Value[j], exp[1], env)
				})
			}
			return NewBoolean(b), nil
		})
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
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if l, ok := exp[0].(*List); ok {
					return NewBoolean(0 == len(l.Value)), nil
				} else {
					return NewBoolean(false), nil
				}
			})
	}
	buildInFuncTbl["length"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if l, ok := exp[0].(*List); ok {
					return NewInteger(len(l.Value)), nil
				} else {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
			})
	}
	buildInFuncTbl["car"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
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
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
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
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
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
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if _, ok := exp[1].(*List); ok {
					var args []Expression
					args = append(args, exp[0])
					return NewList(append(args, (exp[1].(*List)).Value...)), nil
				}
				return NewPair(exp[0], exp[1]), nil
			})
	}
	buildInFuncTbl["append"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				var l []Expression
				for _, e := range exp {
					if v, ok := e.(*List); ok {
						l = append(l, v.Value...)
					} else {
						return nil, NewRuntimeError("E1005", reflect.TypeOf(e).String())
					}
				}
				return NewList(l), nil
			})
	}
	buildInFuncTbl["append!"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				var l *List
				if v, ok := exp[0].(*List); ok {
					l = v
				} else {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
				for _, e := range exp[1:] {
					if v, ok := e.(*List); ok {
						l.Value = append(l.Value, v.Value...)
					} else {
						return nil, NewRuntimeError("E1005", reflect.TypeOf(e).String())
					}
				}
				return l, nil
			})
	}
	buildInFuncTbl["last"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
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
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
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
		if len(exp) < 1 || 3 < len(exp) {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
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
		return doListFunc(exp, env,
			func(org Expression,
				value Expression,
				result []Expression) ([]Expression, error) {
				return append(result, value), nil
			})
	}
	buildInFuncTbl["for-each"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if _, err := doListFunc(exp, env,
			func(org Expression,
				value Expression,
				result []Expression) ([]Expression, error) {
				return nil, nil
			}); err != nil {
			return nil, err
		}
		return NewNil(), nil
	}
	buildInFuncTbl["filter"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return doListFunc(exp, env,
			func(org Expression,
				value Expression,
				result []Expression) ([]Expression, error) {
				b, ok := value.(*Boolean)
				if !ok {
					return nil, NewRuntimeError("E1001", reflect.TypeOf(value).String())
				}
				if b.Value {
					return append(result, org), nil
				}
				return result, nil
			})
	}
	buildInFuncTbl["reduce"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 3 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				l, ok := exp[2].(*List)
				if !ok {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[1]).String())
				}
				if len(l.Value) == 0 {
					return exp[1], nil
				}
				result := l.Value[0]
				for _, e := range l.Value[1:] {
					sexp := MakeQuotedValue(exp[0], []Expression{e}, result)
					v, err := eval(sexp, env)
					if err != nil {
						return nil, err
					}
					result = v
				}
				return result, nil
			})
	}
	buildInFuncTbl["make-list"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				size, ok := exp[0].(*Integer)
				if !ok {
					return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[0]).String())
				}
				if size.Value < 0 {
					return nil, NewRuntimeError("E1011", strconv.Itoa(size.Value))
				}
				l := make([]Expression, 0, size.Value)
				for i := 0; i < size.Value; i++ {
					l = append(l, exp[1].clone())
				}
				return NewList(l), nil
			})
	}
	buildInFuncTbl["take"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return subList(exp, env,
			func(l *List, n *Integer) (x int, y int) {
				x = 0
				y = n.Value
				return x, y
			})

	}
	buildInFuncTbl["drop"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return subList(exp, env,
			func(l *List, n *Integer) (x int, y int) {
				x = n.Value
				y = len(l.Value)
				return x, y
			})
	}
	buildInFuncTbl["delete"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				m, ok := exp[1].(*List)
				if !ok {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[1]).String())
				}
				l := make([]Expression, 0, len(m.Value))
				for _, e := range m.Value {

					if exp[0].equalValue(e) {
						continue
					}
					l = append(l, e)
				}
				return NewList(l), nil
			})
	}
	buildInFuncTbl["delete!"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				l, ok := exp[1].(*List)
				if !ok {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[1]).String())
				}
				m := make([]Expression, len(l.Value))
				copy(m, l.Value)
				l.Value = make([]Expression, 0, len(m))
				for _, e := range m {
					if exp[0].equalValue(e) {
						continue
					}
					l.Value = append(l.Value, e)
				}
				return exp[1], nil
			})
	}
	buildInFuncTbl["list-ref"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				l, ok := exp[0].(*List)
				if !ok {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
				n, ok := exp[1].(*Integer)
				if !ok {
					return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[1]).String())
				}

				if n.Value < 0 || len(l.Value) <= n.Value {
					return nil, NewRuntimeError("E1011", strconv.Itoa(n.Value))
				}
				return l.Value[n.Value], nil
			})
	}
	buildInFuncTbl["list-set!"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 3 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				l, ok := exp[0].(*List)
				if !ok {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
				n, ok := exp[1].(*Integer)
				if !ok {
					return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[1]).String())
				}

				if n.Value < 0 || len(l.Value) <= n.Value {
					return nil, NewRuntimeError("E1011", strconv.Itoa(n.Value))
				}
				l.Value[n.Value] = exp[2]

				return NewNil(), nil
			})
	}
	buildInFuncTbl["set-car!"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				l, ok := exp[0].(*List)
				if !ok {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
				if len(l.Value) <= 0 {
					return nil, NewRuntimeError("E1011")
				}
				l.Value[0] = exp[1]
				return NewNil(), nil
			})
	}
	buildInFuncTbl["set-cdr!"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				l, ok := exp[0].(*List)
				if !ok {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
				if len(l.Value) <= 0 {
					return nil, NewRuntimeError("E1011")
				}
				tmp := l.Value[0]
				l.Value = make([]Expression, 0)
				l.Value = append(l.Value, tmp)

				if m, ok := exp[1].(*List); ok {
					l.Value = append(l.Value, m.Value...)
				} else {
					l.Value = append(l.Value, exp[1])
				}
				return NewNil(), nil
			})
	}
	buildInFuncTbl["sort"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return copySort(exp, env, sort.Slice)
	}
	buildInFuncTbl["sort!"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return effectSort(exp, env, sort.Slice)
	}
	buildInFuncTbl["stable-sort"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return copySort(exp, env, sort.SliceStable)
	}
	buildInFuncTbl["stable-sort!"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return effectSort(exp, env, sort.SliceStable)
	}
	buildInFuncTbl["merge"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return merge(exp, env)
	}
	buildInFuncTbl["sorted?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isSorted(exp, env)
	}

	// list operator
	buildInFuncTbl["vector"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				var l []Expression
				return NewVector(append(l, exp...)), nil
			})
	}
}
