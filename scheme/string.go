/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

import (
	"reflect"
	"strconv"
	"strings"
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

// Build Global environement.
func buildStringFunc() {

	buildInFuncTbl["string-append"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		ret := make([]string, 0, len(exp))
		for _, e := range exp {
			v, err := eval(e, env)
			if err != nil {
				return v, err
			}
			s, ok := v.(*String)
			if !ok {
				return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
			}
			ret = append(ret, s.Value)
		}
		return NewString(strings.Join(ret, "")), nil
	}
}
