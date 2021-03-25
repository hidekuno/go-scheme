/*
   Go lang 4th study program.
   This is test program for mini scheme subset.

   ex.) go test -v lisp_dev.go lisp_dev_test.go

   hidekuno@gmail.com
*/
package experiment

import (
	"testing"

	"github.com/hidekuno/go-scheme/scheme"
)

func executeTest(testCode [][]string, testName string, t *testing.T) {

	rootEnv := scheme.NewSimpleEnv(nil, nil)

	for i, c := range testCode {
		exp, err := scheme.DoCoreLogic(c[0], rootEnv)
		if err != nil {
			if e, ok := err.(*scheme.SyntaxError); ok {
				if e.MsgCode != c[1] {
					t.Log(i)
					t.Fatal("failed "+testName+" test", e.MsgCode)
				}
			} else if e, ok := err.(*scheme.RuntimeError); ok {
				if e.MsgCode != c[1] {
					t.Log(i)
					t.Fatal("failed "+testName+" test", e.MsgCode)
				}
			} else {
				t.Log(i)
				t.Fatal("failed "+testName+" test", err.Error())
			}
		} else {
			if exp.String() != c[1] {
				t.Log(i)
				t.Fatal("failed "+testName+" test", exp)
			}
		}
	}
}
func TestBuildFunc(t *testing.T) {
	scheme.BuildFunc()
	BuildGoFunc()
}
func TestGoApppend(t *testing.T) {
	testCode := [][]string{
		{"(go-append (list 1 2 3)(list 4 5 6))", "(1 2 3 4 5 6)"},
		{"(go-append)", "E1007"},
		{"(go-append 10 20 30)", "E1007"},
		{"(go-append 10 (list 1 2 3))", "E1005"},
		{"(go-append (list 1 2 3) 10)", "E1005"},
	}
	executeTest(testCode, "go-append", t)

}
