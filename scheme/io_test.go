/*
   Go lang 4th study program.
   This is test program for mini scheme subset.

   ex.) go test -v lisp.go lisp_test.go

   hidekuno@gmail.com
*/
package scheme

import (
	"bufio"
	"strings"
	"testing"
)

func readTest(testName string,
	testData string,
	resultData string,
	f func(*SimpleEnv, *bufio.Reader) (Expression, error),
	t *testing.T) {

	BuildFunc()
	rootEnv := NewSimpleEnv(nil, nil)
	exp, err := f(rootEnv, bufio.NewReader(strings.NewReader(testData)))
	if err != nil {
		t.Log(exp)
		t.Fatal("failed "+testName+" test", exp)
	}
	if exp.String() != resultData {
		t.Log(exp)
		t.Fatal("failed "+testName+" test", exp)
	}
}
func TestDisplay(t *testing.T) {
	testCode := [][]string{
		{"(display)", "E1007"},
	}
	executeTest(testCode, "display", t)
}
func TestNewLine(t *testing.T) {
	testCode := [][]string{
		{"(newline 10)", "E1007"},
	}
	executeTest(testCode, "newline", t)
}
func TestLoadFile(t *testing.T) {
	testCode := [][]string{
		{"(load-file)", "E1007"},
		{"(load-file 10)", "E1015"},
		{"(load-file a)", "E1008"},
		{"(load-file \"example/no.scm\")", "E1014"},
		{"(load-file \"/tmp\")", "E1016"},
		{"(load-file \"/etc/sudoers\")", "E9999"},
	}
	executeTest(testCode, "load-file", t)
}
func TestLoadUrl(t *testing.T) {
	testCode := [][]string{
		{"(load-url)", "E1007"},
		{"(load-url 10 10)", "E1007"},
		{"(load-url 10)", "E1015"},
		{"(load-url a)", "E1008"},
		{"(load-file \"https://raw.githubusercontent.com/hidekuno/go-scheme/master/base64.scm\")", "E1014"},
	}
	executeTest(testCode, "load-url", t)
}
func TestRead(t *testing.T) {
	readTest("read", "abcdef", "abcdef", read, t)

	testCode := [][]string{
		{"(read 1)", "E1007"},
	}
	executeTest(testCode, "read", t)
}
func TestReadChar(t *testing.T) {
	readTest("read-char", "a", "#\\a", readChar, t)
	testCode := [][]string{
		{"(read-char 1)", "E1007"},
	}
	executeTest(testCode, "read-char", t)
}
