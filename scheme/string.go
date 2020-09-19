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
