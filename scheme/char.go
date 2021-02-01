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
	"unicode"
)

var (
	whitespaceChar = map[string]byte{
		"#\\tab":     0x09,
		"#\\space":   0x20,
		"#\\newline": 0x0A,
		"#\\return":  0x0D,
	}
)

// Character Type
type Char struct {
	Expression
	Value rune
	exp   string
}

func NewChar(v string) *Char {
	b := new(Char)
	b.exp = v
	b.Value = []rune(v)[2]
	return b
}
func NewCharFromRune(c rune) *Char {
	b := new(Char)
	b.Value = c

	for k, v := range whitespaceChar {
		if c == rune(v) {
			b.exp = k
			return b
		}
	}
	if unicode.IsPrint(b.Value) {
		b.exp = "#\\" + string(c)
	} else {
		b.exp = "#\\non-printable-char"
	}
	return b
}
func (self *Char) String() string {
	return self.exp
}
func (self *Char) Print() {
	fmt.Printf("%c", self.Value)
}
func (self *Char) isAtom() bool {
	return true
}
func (self *Char) clone() Expression {
	return NewChar(self.exp)
}
func (self *Char) equalValue(e Expression) bool {
	if v, ok := e.(*Char); ok {
		return self.Value == v.Value
	}
	return false
}
func charCompare(exp []Expression, env *SimpleEnv, cmp func(rune, rune) bool) (Expression, error) {
	if len(exp) != 2 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
		x, ok := exp[0].(*Char)
		if !ok {
			return nil, NewRuntimeError("E1019", reflect.TypeOf(exp[0]).String())
		}
		y, ok := exp[1].(*Char)
		if !ok {
			return nil, NewRuntimeError("E1019", reflect.TypeOf(exp[1]).String())
		}
		return NewBoolean(cmp(x.Value, y.Value)), nil
	})
}
func isCharKind(exp []Expression, env *SimpleEnv, fn func(rune) bool) (Expression, error) {
	if len(exp) != 1 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
		x, ok := exp[0].(*Char)
		if !ok {
			return nil, NewRuntimeError("E1019", reflect.TypeOf(exp[0]).String())
		}
		return NewBoolean(fn(x.Value)), nil
	})
}
func convertChar(exp []Expression, env *SimpleEnv, fn func(rune) rune) (Expression, error) {
	if len(exp) != 1 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
		x, ok := exp[0].(*Char)
		if !ok {
			return nil, NewRuntimeError("E1019", reflect.TypeOf(exp[0]).String())
		}
		return NewCharFromRune(fn(x.Value)), nil
	})
}
func buildCharFunc() {
	buildInFuncTbl["char=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return charCompare(exp, env, func(x rune, y rune) bool { return x == y })
	}
	buildInFuncTbl["char<?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return charCompare(exp, env, func(x rune, y rune) bool { return x < y })
	}
	buildInFuncTbl["char>?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return charCompare(exp, env, func(x rune, y rune) bool { return x > y })
	}
	buildInFuncTbl["char<=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return charCompare(exp, env, func(x rune, y rune) bool { return x <= y })
	}
	buildInFuncTbl["char>=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return charCompare(exp, env, func(x rune, y rune) bool { return x >= y })
	}
	buildInFuncTbl["char-ci=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return charCompare(exp, env, func(x rune, y rune) bool { return unicode.ToLower(x) == unicode.ToLower(y) })
	}
	buildInFuncTbl["char-ci<?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return charCompare(exp, env, func(x rune, y rune) bool { return unicode.ToLower(x) < unicode.ToLower(y) })
	}
	buildInFuncTbl["char-ci>?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return charCompare(exp, env, func(x rune, y rune) bool { return unicode.ToLower(x) > unicode.ToLower(y) })
	}
	buildInFuncTbl["char-ci<=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return charCompare(exp, env, func(x rune, y rune) bool { return unicode.ToLower(x) <= unicode.ToLower(y) })
	}
	buildInFuncTbl["char-ci>=?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return charCompare(exp, env, func(x rune, y rune) bool { return unicode.ToLower(x) >= unicode.ToLower(y) })
	}
	buildInFuncTbl["char-alphabetic?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isCharKind(exp, env, func(x rune) bool { return unicode.IsLetter(x) })
	}
	buildInFuncTbl["char-numeric?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isCharKind(exp, env, func(x rune) bool { return unicode.IsNumber(x) })
	}
	buildInFuncTbl["char-whitespace?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isCharKind(exp, env, func(x rune) bool { return unicode.IsSpace(x) })
	}
	buildInFuncTbl["char-upper-case?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isCharKind(exp, env, func(x rune) bool { return unicode.IsUpper(x) })
	}
	buildInFuncTbl["char-lower-case?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return isCharKind(exp, env, func(x rune) bool { return unicode.IsLower(x) })
	}
	buildInFuncTbl["integer->char"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			x, ok := exp[0].(*Integer)
			if !ok {
				return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[0]).String())
			}
			return NewCharFromRune(rune(x.Value)), nil
		})

	}
	buildInFuncTbl["char->integer"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return EvalCalcParam(exp, env, func(exp ...Expression) (Expression, error) {
			c, ok := exp[0].(*Char)
			if !ok {
				return nil, NewRuntimeError("E1019", reflect.TypeOf(exp[0]).String())
			}
			return NewInteger(int(c.Value)), nil
		})

	}
	buildInFuncTbl["char-upcase"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return convertChar(exp, env, func(x rune) rune { return unicode.ToUpper(x) })
	}
	buildInFuncTbl["char-downcase"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return convertChar(exp, env, func(x rune) rune { return unicode.ToLower(x) })
	}
}
