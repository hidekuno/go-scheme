/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// Build Global environement.
func buildIoFunc() {

	buildInFuncTbl["display"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		for _, e := range exp {
			v, err := eval(e, env)
			if err != nil {
				return v, err
			}
			fmt.Print(v.String())
		}
		return NewNil(), nil
	}
	buildInFuncTbl["newline"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 0 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		fmt.Println("")
		return NewNil(), nil
	}
	buildInFuncTbl["load-file"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {

		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		v, err := eval(exp[0], env)
		if err != nil {
			return v, err
		}
		filename, ok := v.(*String)
		if !ok {
			return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
		}
		f, err := os.Stat(filename.Value)
		if err != nil && os.IsNotExist(err) {
			return nil, NewRuntimeError("E1014")
		}
		if err != nil || !f.Mode().IsRegular() {
			return nil, NewRuntimeError("E1016")
		}

		fd, err := os.Open(filename.Value)
		if err != nil {
			return nil, NewRuntimeError("E9999")
		}
		defer func() { _ = fd.Close() }()
		repl(fd, env)
		return NewNil(), nil
	}

}
