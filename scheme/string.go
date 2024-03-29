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
func (self *String) Print() {
	fmt.Print(self.Value)
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

func innerSubString(exp []Expression, env *SimpleEnv, s string) (*string, error) {
	r := [2]int{0, utf8.RuneCountInString(s)}

	for i, e := range exp {
		v, ok := e.(*Integer)
		if !ok {
			return nil, NewRuntimeError("E1002", reflect.TypeOf(e).String())
		}

		r[i] = v.Value
	}
	from, to := r[0], r[1]
	if from < 0 || to > utf8.RuneCountInString(s) || from > to {
		return nil, NewRuntimeError("E1021", strconv.Itoa(from), strconv.Itoa(to))
	}

	ret := string([]rune(s)[from:to])
	return &ret, nil
}
func getStartEnd(exp []Expression, env *SimpleEnv, s string) (int, int, error) {

	r := [2]int{0, utf8.RuneCountInString(s)}

	for i, e := range exp {
		v, ok := e.(*Integer)
		if !ok {
			return r[0], r[1], NewRuntimeError("E1002", reflect.TypeOf(e).String())
		}

		r[i] = v.Value
	}
	from, to := r[0], r[1]
	if from < 0 || to > utf8.RuneCountInString(s) || from > to {
		return from, to, NewRuntimeError("E1021", strconv.Itoa(from), strconv.Itoa(to))
	}
	return from, to, nil
}
func stringCase(exp []Expression, env *SimpleEnv, caseFunc func(s string) string) (Expression, error) {

	s, ok := exp[0].(*String)
	if !ok {
		return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
	}
	ss, err := innerSubString(exp[1:], env, s.Value)
	if err != nil {
		return nil, err
	}
	return NewString(caseFunc(*ss)), nil
}
func stringIndex(exp []Expression, env *SimpleEnv, index func(s, chars string) int) (Expression, error) {

	s, ok := exp[0].(*String)
	if !ok {
		return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
	}
	c, ok := exp[1].(*Char)
	if !ok {
		return nil, NewRuntimeError("E1019", reflect.TypeOf(exp[1]).String())
	}
	from, to, err := getStartEnd(exp[2:], env, s.Value)
	if err != nil {
		return nil, err
	}
	idx := index(s.Value, string(c.Value))
	if idx == -1 {
		return NewBoolean(false), nil
	} else if (from <= idx) && (idx < to) {
		return NewInteger(idx), nil
	} else {
		return NewBoolean(false), nil
	}
}
func stringTrim(exp []Expression, env *SimpleEnv, trim func(s, chars string) string) (Expression, error) {

	s, ok := exp[0].(*String)
	if !ok {
		return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
	}
	if len(exp) == 1 {
		return NewString(trim(s.Value, " ")), nil
	} else {
		c, ok := exp[1].(*Char)
		if !ok {
			return nil, NewRuntimeError("E1019", reflect.TypeOf(exp[1]).String())
		}
		return NewString(trim(s.Value, string(c.Value))), nil
	}
}
func stringScan(exp []Expression, env *SimpleEnv, index func(s, chars string) int) (Expression, error) {
	if len(exp) != 2 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
		x, ok := exp[0].(*String)
		if !ok {
			return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
		}
		sep := ""
		if s, ok := exp[1].(*String); ok {
			sep = s.Value
		} else if c, ok := exp[1].(*Char); ok {
			// strings.IndexRune is exists, but there is not strings.LastIndexRune
			sep = string(c.Value)
		} else {
			return nil, NewRuntimeError("E1009")
		}

		i := index(x.Value, sep)
		if i >= 0 {
			return NewInteger(i), nil
		} else {
			return NewBoolean(false), nil
		}
	})
}
func stringCompare(exp []Expression, env *SimpleEnv, operate func(string, string) bool) (Expression, error) {
	if len(exp) != 2 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
		x, ok := exp[0].(*String)
		if !ok {
			return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
		}
		y, ok := exp[1].(*String)
		if !ok {
			return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[1]).String())
		}
		return NewBoolean(operate(x.Value, y.Value)), nil
	})
}
func stringLength(exp []Expression, env *SimpleEnv, fn func(string) int) (Expression, error) {
	if len(exp) != 1 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
		x, ok := exp[0].(*String)
		if !ok {
			return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
		}

		return NewInteger(fn(x.Value)), nil
	})
}
func numberString(exp Expression, env *SimpleEnv, r int) (Expression, error) {
	if _, ok := exp.(Number); !ok {
		return nil, NewRuntimeError("E1003", reflect.TypeOf(exp).String())
	}
	if r == 10 {
		return NewString(exp.String()), nil
	}
	if i, ok := exp.(*Integer); ok {
		return NewString(strconv.FormatInt(int64(i.Value), r)), nil
	} else {
		return NewString(exp.String()), nil
	}
}
func stringNumber(exp Expression, env *SimpleEnv, r int) (Expression, error) {
	s, ok := exp.(*String)
	if !ok {
		return nil, NewRuntimeError("E1015", reflect.TypeOf(exp).String())
	}
	if i, err := strconv.ParseInt(s.Value, r, 0); err == nil {
		return NewInteger(int(i)), nil
	} else if f, err := strconv.ParseFloat(s.Value, 64); err == nil {
		return NewFloat(f), nil
	} else {
		rat := MakeRatRadix(s.Value, r)
		if rat != nil {
			return rat, nil
		}
	}
	return NewBoolean(false), nil
}
func stringRange(exp []Expression, env *SimpleEnv, rangeFunc func(r []rune, n int) []rune) (Expression, error) {

	s, ok := exp[0].(*String)
	if !ok {
		return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
	}
	n, ok := exp[1].(*Integer)
	if !ok {
		return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[0]).String())
	}
	if 0 > n.Value || utf8.RuneCountInString(s.Value) < n.Value {
		return nil, NewRuntimeError("E1021", n.String())
	}
	return NewString(string(rangeFunc([]rune(s.Value), n.Value))), nil
}

// Build Global environement.
func buildStringFunc() {
	buildInFuncTbl["string-append"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {

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
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
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
	buildInFuncTbl["string"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			c, ok := exp[0].(*Char)
			if !ok {
				return nil, NewRuntimeError("E1019", reflect.TypeOf(exp[0]).String())
			}
			return NewString(string(c.Value)), nil
		})
	}
	buildInFuncTbl["string=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return stringCompare(exp, env, func(x string, y string) bool { return x == y })
	}
	buildInFuncTbl["string<?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return stringCompare(exp, env, func(x string, y string) bool { return x < y })
	}
	buildInFuncTbl["string>?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return stringCompare(exp, env, func(x string, y string) bool { return x > y })
	}
	buildInFuncTbl["string<=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return stringCompare(exp, env, func(x string, y string) bool { return x <= y })
	}
	buildInFuncTbl["string>=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return stringCompare(exp, env, func(x string, y string) bool { return x >= y })
	}
	buildInFuncTbl["string-ci=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return stringCompare(exp, env, func(x string, y string) bool { return strings.ToLower(x) == strings.ToLower(y) })
	}
	buildInFuncTbl["string-ci<?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return stringCompare(exp, env, func(x string, y string) bool { return strings.ToLower(x) < strings.ToLower(y) })
	}
	buildInFuncTbl["string-ci>?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return stringCompare(exp, env, func(x string, y string) bool { return strings.ToLower(x) > strings.ToLower(y) })
	}
	buildInFuncTbl["string-ci<=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return stringCompare(exp, env, func(x string, y string) bool { return strings.ToLower(x) <= strings.ToLower(y) })
	}
	buildInFuncTbl["string-ci>=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return stringCompare(exp, env, func(x string, y string) bool { return strings.ToLower(x) >= strings.ToLower(y) })
	}
	buildInFuncTbl["string-length"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return stringLength(exp, env, func(x string) int { return utf8.RuneCountInString(x) })
	}
	buildInFuncTbl["string-size"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return stringLength(exp, env, func(x string) int { return len(x) })
	}
	buildInFuncTbl["number->string"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return doRadix(exp, env, numberString)
	}
	buildInFuncTbl["string->number"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return doRadix(exp, env, stringNumber)
	}
	buildInFuncTbl["list->string"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			var buffer bytes.Buffer

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
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
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
		if len(exp) != 3 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			s, ok := exp[0].(*String)
			if !ok {
				return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
			}

			ss, err := innerSubString(exp[1:], env, s.Value)
			if err != nil {
				return nil, err
			}
			return NewString(*ss), nil
		})
	}
	buildInFuncTbl["symbol->string"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			s, ok := exp[0].(*Symbol)
			if !ok {
				return nil, NewRuntimeError("E1004", reflect.TypeOf(exp[0]).String())

			}
			return NewString(s.Value), nil
		})
	}
	buildInFuncTbl["string->symbol"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			s, ok := exp[0].(*String)
			if !ok {
				return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())

			}
			return NewSymbol(s.Value), nil
		})
	}
	buildInFuncTbl["make-string"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			var buffer bytes.Buffer
			n, ok := exp[0].(*Integer)
			if !ok {
				return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[0]).String())

			}
			c, ok := exp[1].(*Char)
			if !ok {
				return nil, NewRuntimeError("E1019", reflect.TypeOf(exp[1]).String())

			}
			if n.Value < 0 {
				return nil, NewRuntimeError("E1021", n.String())
			}
			for i := 0; i < n.Value; i++ {
				buffer.WriteRune(c.Value)
			}
			return NewString(buffer.String()), nil
		})
	}
	buildInFuncTbl["string-split"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			s, ok := exp[0].(*String)
			if !ok {
				return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())

			}
			c, ok := exp[1].(*Char)
			if !ok {
				return nil, NewRuntimeError("E1019", reflect.TypeOf(exp[1]).String())

			}
			v := strings.Split(s.Value, string(c.Value))
			l := make([]Expression, 0, len(v))

			for _, e := range v {
				l = append(l, NewString(e))
			}
			return NewList(l), nil
		})
	}
	buildInFuncTbl["string-join"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			l, ok := exp[0].(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
			}
			s, ok := exp[1].(*String)
			if !ok {
				return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
			}
			v := make([]string, 0, len(l.Value))
			for _, e := range l.Value {
				s, ok := e.(*String)
				if !ok {
					return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
				}
				v = append(v, s.Value)
			}
			return NewString(strings.Join(v, s.Value)), nil
		})
	}
	buildInFuncTbl["string-scan"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return stringScan(exp, env, strings.Index)
	}
	buildInFuncTbl["string-scan-right"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return stringScan(exp, env, strings.LastIndex)
	}
	buildInFuncTbl["vector->string"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			var buffer bytes.Buffer

			l, ok := exp[0].(*Vector)
			if !ok {
				return nil, NewRuntimeError("E1022", reflect.TypeOf(exp[0]).String())
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
	buildInFuncTbl["string->vector"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			s, ok := exp[0].(*String)
			if !ok {
				return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
			}
			l := make([]Expression, 0, len(s.Value))
			for _, c := range s.Value {
				l = append(l, NewCharFromRune(rune(c)))
			}
			return NewVector(l), nil
		})
	}
	buildInFuncTbl["string-reverse"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 1 || 3 < len(exp) {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			s, ok := exp[0].(*String)
			if !ok {
				return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())

			}
			ss, err := innerSubString(exp[1:], env, s.Value)
			if err != nil {
				return nil, err
			}
			// golang is not exists string reverse api
			runes := []rune(*ss)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			return NewString(string(runes)), nil
		})
	}
	buildInFuncTbl["string-upcase"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 1 || 3 < len(exp) {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringCase(exp, env, strings.ToUpper)
		})
	}
	buildInFuncTbl["string-downcase"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 1 || 3 < len(exp) {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringCase(exp, env, strings.ToLower)
		})
	}
	buildInFuncTbl["string-index"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 1 || 4 < len(exp) {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringIndex(exp, env, strings.Index)
		})
	}
	buildInFuncTbl["string-index-right"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 1 || 4 < len(exp) {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringIndex(exp, env, strings.LastIndex)
		})
	}
	buildInFuncTbl["string-delete"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 1 || 4 < len(exp) {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {

			s, ok := exp[0].(*String)
			if !ok {
				return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
			}
			c, ok := exp[1].(*Char)
			if !ok {
				return nil, NewRuntimeError("E1019", reflect.TypeOf(exp[1]).String())
			}
			ss, err := innerSubString(exp[2:], env, s.Value)
			if err != nil {
				return nil, err
			}
			runes := make([]rune, 0, len(*ss))
			for _, e := range []rune(*ss) {
				if e != c.Value {
					runes = append(runes, e)
				}
			}
			return NewString(string(runes)), nil
		})
	}
	buildInFuncTbl["string-trim"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 1 || 2 < len(exp) {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringTrim(exp, env, strings.TrimLeft)
		})
	}
	buildInFuncTbl["string-trim-right"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 1 || 2 < len(exp) {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringTrim(exp, env, strings.TrimRight)
		})
	}
	buildInFuncTbl["string-trim-both"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 1 || 2 < len(exp) {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringTrim(exp, env, strings.Trim)
		})
	}
	buildInFuncTbl["string-take"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringRange(exp, env, func(r []rune, n int) []rune { return r[:n] })
		})
	}
	buildInFuncTbl["string-take-right"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringRange(exp, env, func(r []rune, n int) []rune { return r[len(r)-n:] })
		})
	}
	buildInFuncTbl["string-drop"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringRange(exp, env, func(r []rune, n int) []rune { return r[n:] })
		})
	}
	buildInFuncTbl["string-drop-right"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			return stringRange(exp, env, func(r []rune, n int) []rune { return r[:len(r)-n] })
		})
	}
}
