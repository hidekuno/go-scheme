/*
   Go lang 4th study program.
   This is test program for mini scheme subset.

   ex.) go test -v lisp.go lisp_test.go

   hidekuno@gmail.com
*/
package scheme

import (
	"testing"
)

func TestCharEq(t *testing.T) {
	testCode := [][]string{
		{"(char=? #\\a #\\a)", "#t"},
		{"(char=? #\\a #\\b)", "#f"},

		{"(char=?)", "E1007"},
		{"(char=? #\\a)", "E1007"},
		{"(char=? #\\a #\\b #\\c)", "E1007"},
		{"(char=? #\\a 10)", "E1019"},
		{"(char=? 10 #\\a)", "E1019"},
		{"(char=? #\\a a)", "E1008"},
	}
	executeTest(testCode, "char=?", t)
}

func TestCharLess(t *testing.T) {
	testCode := [][]string{
		{"(char<? #\\a #\\b)", "#t"},
		{"(char<? #\\b #\\a)", "#f"},

		{"(char<?)", "E1007"},
		{"(char<? #\\a)", "E1007"},
		{"(char<? #\\a #\\b #\\c)", "E1007"},
		{"(char<? #\\a 10)", "E1019"},
		{"(char<? 10 #\\a)", "E1019"},
		{"(char<? #\\a a)", "E1008"},
	}
	executeTest(testCode, "char<?", t)
}
func TestCharThan(t *testing.T) {
	testCode := [][]string{
		{"(char>? #\\b #\\a)", "#t"},
		{"(char>? #\\a #\\b)", "#f"},

		{"(char>?)", "E1007"},
		{"(char>? #\\a)", "E1007"},
		{"(char>? #\\a #\\b #\\c)", "E1007"},
		{"(char>? #\\a 10)", "E1019"},
		{"(char>? 10 #\\a)", "E1019"},
		{"(char>? #\\a a)", "E1008"},
	}
	executeTest(testCode, "char>?", t)
}
func TestCharLessEq(t *testing.T) {
	testCode := [][]string{
		{"(char<=? #\\a #\\b)", "#t"},
		{"(char<=? #\\a #\\a)", "#t"},
		{"(char<=? #\\b #\\a)", "#f"},

		{"(char<=?)", "E1007"},
		{"(char<=? #\\a)", "E1007"},
		{"(char<=? #\\a #\\b #\\c)", "E1007"},
		{"(char<=? #\\a 10)", "E1019"},
		{"(char<=? 10 #\\a)", "E1019"},
		{"(char<=? #\\a a)", "E1008"},
	}
	executeTest(testCode, "char<=?", t)
}
func TestCharThanEq(t *testing.T) {
	testCode := [][]string{
		{"(char>=? #\\b #\\a)", "#t"},
		{"(char>=? #\\a #\\a)", "#t"},
		{"(char>=? #\\a #\\b)", "#f"},

		{"(char>=?)", "E1007"},
		{"(char>=? #\\a)", "E1007"},
		{"(char>=? #\\a #\\b #\\c)", "E1007"},
		{"(char>=? #\\a 10)", "E1019"},
		{"(char>=? 10 #\\a)", "E1019"},
		{"(char>=? #\\a a)", "E1008"},
	}
	executeTest(testCode, "char>=?", t)
}
func TestCharCaseIgnoreEq(t *testing.T) {
	testCode := [][]string{
		{"(char-ci=? #\\a #\\A)", "#t"},
		{"(char-ci=? #\\A #\\a)", "#t"},
		{"(char-ci=? #\\a #\\B)", "#f"},

		{"(char-ci=?)", "E1007"},
		{"(char-ci=? #\\a)", "E1007"},
		{"(char-ci=? #\\a #\\b #\\c)", "E1007"},
		{"(char-ci=? #\\a 10)", "E1019"},
		{"(char-ci=? 10 #\\a)", "E1019"},
		{"(char-ci=? #\\a a)", "E1008"},
	}
	executeTest(testCode, "char-ci=?", t)
}

func TestCharCaseIgnoreLess(t *testing.T) {
	testCode := [][]string{
		{"(char-ci<? #\\a #\\C)", "#t"},
		{"(char-ci<? #\\C #\\a)", "#f"},

		{"(char-ci<?)", "E1007"},
		{"(char-ci<? #\\a)", "E1007"},
		{"(char-ci<? #\\a #\\b #\\c)", "E1007"},
		{"(char-ci<? #\\a 10)", "E1019"},
		{"(char-ci<? 10 #\\a)", "E1019"},
		{"(char-ci<? #\\a a)", "E1008"},
	}
	executeTest(testCode, "char-ci<?", t)
}
func TestCharCaseIgnoreThan(t *testing.T) {
	testCode := [][]string{
		{"(char-ci>? #\\C #\\a)", "#t"},
		{"(char-ci>? #\\a #\\C)", "#f"},

		{"(char-ci>?)", "E1007"},
		{"(char-ci>? #\\a)", "E1007"},
		{"(char-ci>? #\\a #\\b #\\c)", "E1007"},
		{"(char-ci>? #\\a 10)", "E1019"},
		{"(char-ci>? 10 #\\a)", "E1019"},
		{"(char-ci>? #\\a a)", "E1008"},
	}
	executeTest(testCode, "char-ci>?", t)
}
func TestCharCaseIgnoreEqLess(t *testing.T) {
	testCode := [][]string{
		{"(char-ci<=? #\\a #\\C)", "#t"},
		{"(char-ci<=? #\\C #\\C)", "#t"},
		{"(char-ci<=? #\\C #\\a)", "#f"},

		{"(char-ci<=?)", "E1007"},
		{"(char-ci<=? #\\a)", "E1007"},
		{"(char-ci<=? #\\a #\\b #\\c)", "E1007"},
		{"(char-ci<=? #\\a 10)", "E1019"},
		{"(char-ci<=? 10 #\\a)", "E1019"},
		{"(char-ci<=? #\\a a)", "E1008"},
	}
	executeTest(testCode, "char-ci<=?", t)
}

func TestCharCaseIgnoreEqThan(t *testing.T) {
	testCode := [][]string{
		{"(char-ci>=? #\\C #\\a)", "#t"},
		{"(char-ci>=? #\\C #\\C)", "#t"},
		{"(char-ci>=? #\\a #\\C)", "#f"},

		{"(char-ci>=?)", "E1007"},
		{"(char-ci>=? #\\a)", "E1007"},
		{"(char-ci>=? #\\a #\\b #\\c)", "E1007"},
		{"(char-ci>=? #\\a 10)", "E1019"},
		{"(char-ci>=? 10 #\\a)", "E1019"},
		{"(char-ci>=? #\\a a)", "E1008"},
	}
	executeTest(testCode, "char-ci>=?", t)
}

func TestIntegerChar(t *testing.T) {
	testCode := [][]string{
		{"(integer->char 65)", "#\\A"},
		{"(integer->char 23665)", "#\\山"},

		{"(integer->char)", "E1007"},
		{"(integer->char 23 665)", "E1007"},
		{"(integer->char #\\a)", "E1002"},
		{"(integer->char a)", "E1008"},
	}
	executeTest(testCode, "integer->char", t)
}
func TestCharInteger(t *testing.T) {
	testCode := [][]string{
		{"(char->integer #\\A)", "65"},
		{"(char->integer #\\山)", "23665"},

		{"(char->integer)", "E1007"},
		{"(char->integer #\\a #\\b)", "E1007"},
		{"(char->integer 999)", "E1019"},
		{"(char->integer a)", "E1008"},
	}
	executeTest(testCode, "char->integer", t)
}
