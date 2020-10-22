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

func TestStringAppend(t *testing.T) {
	testCode := [][]string{
		{"(string-append \"a\" \"b\")", "\"ab\""},
		{"(string-append \"a\" \"b\"  \"c\")", "\"abc\""},
		{"(string-append)", "E1007"},
		{"(string-append 10)", "E1007"},
		{"(string-append a b)", "E1008"},
		{"(string-append \"a\" 10)", "E1015"},
	}
	executeTest(testCode, "string-append", t)
}
func TestFormat(t *testing.T) {
	testCode := [][]string{
		{"(format \"~D\" 10)", "\"10\""},
		{"(format \"~d\" 10)", "\"10\""},
		{"(format \"~X\" 10)", "\"A\""},
		{"(format \"~x\" 10)", "\"a\""},
		{"(format \"~O\" 10)", "\"12\""},
		{"(format \"~o\" 10)", "\"12\""},
		{"(format \"~B\" 10)", "\"1010\""},
		{"(format \"~b\" 10)", "\"1010\""},
		{"(define a \"~D\")", "a"},
		{"(define b 100)", "b"},
		{"(format a b)", "\"100\""},

		{"(format)", "E1007"},
		{"(format \"~B\")", "E1007"},
		{"(format \"~B\" 10 12)", "E1007"},
		{"(format 10 12)", "E1015"},
		{"(format \"~A\" #f)", "E1002"},
		{"(format \"~A\" 10)", "E1018"},
	}
	executeTest(testCode, "format", t)
}
func TestStringEq(t *testing.T) {
	testCode := [][]string{
		{"(string=? \"abc\" \"abc\")", "#t"},
		{"(string=? \"abc\" \"ABC\")", "#f"},

		{"(string=?)", "E1007"},
		{"(string=? \"abc\")", "E1007"},
		{"(string=? \"abc\" \"ABC\" \"DEF\")", "E1007"},
		{"(string=? \"abc\" 10)", "E1015"},
		{"(string=? 10 \"abc\")", "E1015"},
		{"(string=? \"abc\" a)", "E1008"},
	}
	executeTest(testCode, "string=?", t)
}

func TestStringLess(t *testing.T) {
	testCode := [][]string{

		{"(string<? \"1234\" \"9\")", "#t"},
		{"(string<? \"9\" \"1234\")", "#f"},
		{"(string<?)", "E1007"},
		{"(string<? \"abc\")", "E1007"},
		{"(string<? \"abc\" \"ABC\" \"DEF\")", "E1007"},
		{"(string<? \"abc\" 10)", "E1015"},
		{"(string<? 10 \"abc\")", "E1015"},
		{"(string<? \"abc\" a)", "E1008"},
	}
	executeTest(testCode, "string<?", t)
}
func TestStringThan(t *testing.T) {
	testCode := [][]string{

		{"(string>? \"9\" \"1234\")", "#t"},
		{"(string>? \"1234\" \"9\")", "#f"},
		{"(string>?)", "E1007"},
		{"(string>? \"abc\")", "E1007"},
		{"(string>? \"abc\" \"ABC\" \"DEF\")", "E1007"},
		{"(string>? \"abc\" 10)", "E1015"},
		{"(string>? 10 \"abc\")", "E1015"},
		{"(string>? \"abc\" a)", "E1008"},
	}
	executeTest(testCode, "string>?", t)
}
func TestStringLessEq(t *testing.T) {
	testCode := [][]string{
		{"(string<=? \"1234\" \"9\")", "#t"},
		{"(string<=? \"1234\" \"1234\")", "#t"},
		{"(string<=? \"9\" \"1234\")", "#f"},
		{"(string<=?)", "E1007"},
		{"(string<=? \"abc\")", "E1007"},
		{"(string<=? \"abc\" \"ABC\" \"DEF\")", "E1007"},
		{"(string<=? \"abc\" 10)", "E1015"},
		{"(string<=? 10 \"abc\")", "E1015"},
		{"(string<=? \"abc\" a)", "E1008"},
	}
	executeTest(testCode, "", t)
}
func TestStringThanEq(t *testing.T) {
	testCode := [][]string{
		{"(string>=?  \"9\" \"1234\")", "#t"},
		{"(string>=?  \"1234\" \"1234\")", "#t"},
		{"(string>=?  \"1234\" \"9\")", "#f"},
		{"(string>=?)", "E1007"},
		{"(string>=? \"abc\")", "E1007"},
		{"(string>=? \"abc\" \"ABC\" \"DEF\")", "E1007"},
		{"(string>=? \"abc\" 10)", "E1015"},
		{"(string>=? 10 \"abc\")", "E1015"},
		{"(string>=? \"abc\" a)", "E1008"},
	}
	executeTest(testCode, "string>=?", t)
}
func TestStringCaseIgnoreEq(t *testing.T) {
	testCode := [][]string{
		{"(string-ci=? \"Abc\" \"aBc\")", "#t"},
		{"(string-ci=? \"abc\" \"ABC\")", "#t"},
		{"(string-ci=? \"abcd\" \"ABC\")", "#f"},

		{"(string-ci=?)", "E1007"},
		{"(string-ci=? \"abc\")", "E1007"},
		{"(string-ci=? \"abc\" \"ABC\" \"DEF\")", "E1007"},
		{"(string-ci=? \"abc\" 10)", "E1015"},
		{"(string-ci=? 10 \"abc\")", "E1015"},
		{"(string-ci=? \"abc\" a)", "E1008"},
	}
	executeTest(testCode, "string-ci=?", t)
}

func TestStringCaseIgnoreLess(t *testing.T) {
	testCode := [][]string{
		{"(string-ci<? \"abc\" \"DEF\")", "#t"},
		{"(string-ci<? \"DEF\" \"abc\")", "#f"},
		{"(string-ci<?)", "E1007"},
		{"(string-ci<? \"abc\")", "E1007"},
		{"(string-ci<? \"abc\" \"ABC\" \"DEF\")", "E1007"},
		{"(string-ci<? \"abc\" 10)", "E1015"},
		{"(string-ci<? 10 \"abc\")", "E1015"},
		{"(string-ci<? \"abc\" a)", "E1008"},
	}
	executeTest(testCode, "string-ci<?", t)
}
func TestStringCaseIgnoreThan(t *testing.T) {
	testCode := [][]string{
		{"(string-ci>? \"DEF\" \"abc\")", "#t"},
		{"(string-ci>? \"abc\" \"DEF\")", "#f"},
		{"(string-ci>?)", "E1007"},
		{"(string-ci>? \"abc\")", "E1007"},
		{"(string-ci>? \"abc\" \"ABC\" \"DEF\")", "E1007"},
		{"(string-ci>? \"abc\" 10)", "E1015"},
		{"(string-ci>? 10 \"abc\")", "E1015"},
		{"(string-ci>? \"abc\" a)", "E1008"},
	}
	executeTest(testCode, "string-ci>?", t)
}
func TestStringCaseIgnoreEqLess(t *testing.T) {
	testCode := [][]string{
		{"(string-ci<=? \"abc\" \"DEF\")", "#t"},
		{"(string-ci<=? \"DEF\" \"abc\")", "#f"},
		{"(string-ci<=? \"Abc\" \"aBC\")", "#t"},

		{"(string-ci<=?)", "E1007"},
		{"(string-ci<=? \"abc\")", "E1007"},
		{"(string-ci<=? \"abc\" \"ABC\" \"DEF\")", "E1007"},
		{"(string-ci<=? \"abc\" 10)", "E1015"},
		{"(string-ci<=? 10 \"abc\")", "E1015"},
		{"(string-ci<=? \"abc\" a)", "E1008"},
	}
	executeTest(testCode, "string-ci<=?", t)
}

func TestStringCaseIgnoreEqThan(t *testing.T) {
	testCode := [][]string{
		{"(string-ci>=? \"abc\" \"DEF\")", "#f"},
		{"(string-ci>=? \"DEF\" \"abc\")", "#t"},
		{"(string-ci>=? \"Abc\" \"aBC\")", "#t"},

		{"(string-ci>=?)", "E1007"},
		{"(string-ci>=? \"abc\")", "E1007"},
		{"(string-ci>=? \"abc\" \"ABC\" \"DEF\")", "E1007"},
		{"(string-ci>=? \"abc\" 10)", "E1015"},
		{"(string-ci>=? 10 \"abc\")", "E1015"},
		{"(string-ci>=? \"abc\" a)", "E1008"},
	}
	executeTest(testCode, "string-ci>=?", t)
}
func TestStringLength(t *testing.T) {
	testCode := [][]string{
		{"(string-length \"\")", "0"},
		{"(string-length \"1234567890\")", "10"},
		{"(string-length \"山\")", "1"},

		{"(string-length)", "E1007"},
		{"(string-length \"1234\" \"12345\")", "E1007"},
		{"(string-length 1000)", "E1015"},
	}
	executeTest(testCode, "string-length", t)
}
func TestStringSize(t *testing.T) {
	testCode := [][]string{
		{"(string-size \"\")", "0"},
		{"(string-size \"1234567890\")", "10"},
		{"(string-size \"山\")", "3"},

		{"(string-size)", "E1007"},
		{"(string-size \"1234\" \"12345\")", "E1007"},
		{"(string-size 1000)", "E1015"},
	}
	executeTest(testCode, "string-size", t)
}
func TestNumberString(t *testing.T) {
	testCode := [][]string{
		{"(number->string 10)", "\"10\""},
		{"(number->string 10.5)", "\"10.5\""},

		{"(number->string)", "E1007"},
		{"(number->string 10 20)", "E1007"},
		{"(number->string #f)", "E1003"},
		{"(number->string a)", "E1008"},
	}
	executeTest(testCode, "number->string", t)
}
func TestStringNumber(t *testing.T) {
	testCode := [][]string{

		{"(string->number \"123\")", "123"},
		{"(string->number \"10.5\")", "10.5"},

		{"(string->number)", "E1007"},
		{"(string->number \"123\" \"10.5\")", "E1007"},
		{"(string->number 100)", "E1015"},
		{"(string->number \"/1\")", "E1003"},
		{"(string->number \"1/3/2\")", "E1003"},
		{"(string->number a)", "E1008"},
	}
	executeTest(testCode, "string->number", t)
}
