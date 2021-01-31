/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func read(env *SimpleEnv, reader *bufio.Reader) (Expression, error) {

	program := make([]string, 0, 64)
	for {
		b, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := string(b)
		if line == "" {
			continue
		}
		program = append(program, line)
		if !countParenthesis(strings.Join(program, " ")) {
			continue
		}
		token, err := tokenize(strings.Join(program, " "))
		if err != nil {
			return nil, err
		}
		ast, _, err := parse(token)
		if err != nil {
			return nil, err
		}
		return ast, nil
	}

	// EOF
	return NewNil(), nil
}
func read_char(env *SimpleEnv, reader *bufio.Reader) (Expression, error) {

	b, _, err := reader.ReadLine()
	c := "#\\"
	if len(b) == 0 {
		c = c + "newline"
	} else {
		if b[0] == ' ' {
			c = c + "space"
		} else {
			c = c + string(b[0])
		}
	}
	ast, _, err := parse([]string{c})
	if err != nil {
		return nil, err
	}
	return ast, nil
}

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
			v.Print()
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
	buildInFuncTbl["load-url"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		v, err := eval(exp[0], env)
		if err != nil {
			return v, err
		}
		url, ok := v.(*String)
		if !ok {
			return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
		}
		resp, err := http.Get(url.Value)
		if err != nil {
			return nil, NewRuntimeError("E1014")
		}
		if resp.StatusCode != http.StatusOK {
			return nil, NewRuntimeError("E1014")
		}

		defer resp.Body.Close()
		repl(resp.Body, env)

		return NewNil(), nil
	}
	buildInFuncTbl["read"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 0 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		// EOF
		return read(env, bufio.NewReaderSize(os.Stdin, MaxLineSize))
	}
	buildInFuncTbl["read-char"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 0 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return read_char(env, bufio.NewReaderSize(os.Stdin, MaxLineSize))
	}
}
