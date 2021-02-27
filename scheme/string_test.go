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
func TestString(t *testing.T) {
	testCode := [][]string{
		{"(string #\\a)", "\"a\""},
		{"(string #\\A)", "\"A\""},
		{"(string #\\0)", "\"0\""},
		{"(string #\\9)", "\"9\""},

		{"(string)", "E1007"},
		{"(string 1 2)", "E1007"},
		{"(string 10)", "E1019"},
		{"(string a)", "E1008"},
	}
	executeTest(testCode, "string", t)
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
func TestListString(t *testing.T) {
	testCode := [][]string{
		{"(list->string (list))", "\"\""},
		{"(list->string (list #\\a #\\b #\\c))", "\"abc\""},

		{"(list->string)", "E1007"},
		{"(list->string (list #\\a #\\b)(list #\\a #\\b))", "E1007"},
		{"(list->string 10)", "E1005"},
		{"(list->string (list #\\a 10))", "E1019"},
		{"(list->string a)", "E1008"},
	}
	executeTest(testCode, "list->string", t)
}
func TestStringList(t *testing.T) {
	testCode := [][]string{
		{"(string->list \"\")", "()"},
		{"(string->list \"abc\")", "(#\\a #\\b #\\c)"},
		{"(string->list \"山田\")", "(#\\山 #\\田)"},

		{"(string->list)", "E1007"},
		{"(string->list \"a\" \"b\")", "E1007"},
		{"(string->list #\\a)", "E1015"},
		{"(string->list a)", "E1008"},
	}
	executeTest(testCode, "string->list", t)
}
func TestSubString(t *testing.T) {
	testCode := [][]string{
		{"(substring \"1234567890\" 1 2)", "\"2\""},
		{"(substring \"1234567890\" 1 3)", "\"23\""},
		{"(substring \"1234567890\" 0 10)", "\"1234567890\""},
		{"(substring \"山\" 0 1)", "\"山\""},
		{"(substring \"山1\" 0 2)", "\"山1\""},

		{"(substring)", "E1007"},
		{"(substring \"1234567890\" 1)", "E1007"},
		{"(substring \"1234567890\" 1 2 3)", "E1007"},
		{"(substring  1 2 3)", "E1015"},
		{"(substring \"1234567890\" #t 2)", "E1002"},
		{"(substring \"1234567890\" 0 #t)", "E1002"},
		{"(substring \"1234567890\" a 2)", "E1008"},
		{"(substring \"1234567890\" 0 a)", "E1008"},
		{"(substring \"1234567890\" -1 2)", "E1021"},
		{"(substring \"1234567890\" 0 -2)", "E1021"},
		{"(substring \"1234567890\" 0 11)", "E1021"},
		{"(substring \"1234567890\" 6 5)", "E1021"},
		{"(substring \"山\" 0 2)", "E1021"},
	}
	executeTest(testCode, "substring", t)
}
func TestSymbolString(t *testing.T) {
	testCode := [][]string{
		{"(symbol->string (quote abc))", "\"abc\""},

		{"(symbol->string)", "E1007"},
		{"(symbol->string (quote a) (quote b))", "E1007"},
		{"(symbol->string #t)", "E1004"},
	}
	executeTest(testCode, "symbol->string", t)
}
func TestStringSymbol(t *testing.T) {
	testCode := [][]string{
		{"(string->symbol \"abc\")", "abc"},

		{"(string->symbol)", "E1007"},
		{"(string->symbol \"abc\"  \"def\")", "E1007"},
		{"(string->symbol #t)", "E1015"},
	}
	executeTest(testCode, "string->symbol", t)
}
func TestMakeString(t *testing.T) {
	testCode := [][]string{
		{"(make-string 4 #\\a)", "\"aaaa\""},
		{"(make-string 4 #\\山)", "\"山山山山\""},

		{"(make-string)", "E1007"},
		{"(make-string 1)", "E1007"},
		{"(make-string 1 1 1)", "E1007"},
		{"(make-string #t #\\a)", "E1002"},
		{"(make-string -1 #\\a)", "E1021"},
		{"(make-string 4 a)", "E1008"},
		{"(make-string 4 #t)", "E1019"},
	}
	executeTest(testCode, "make-string", t)
}
func TestStringSplit(t *testing.T) {
	testCode := [][]string{
		{"(string-split  \"abc:def:g\"  #\\:)", "(\"abc\" \"def\" \"g\")", "(\"abc\" \"def\" \"g\")"},
		{"(string-split  \"abcdef\"  #\\,)", "(\"abcdef\")"},
		{"(string-split  \",abcdef\"  #\\,)", "(\"\" \"abcdef\")"},
		{"(string-split  \"abcdef,\"  #\\,)", "(\"abcdef\" \"\")"},
		{"(string-split  \"\"  #\\,)", "(\"\")"},

		{"(string-split)", "E1007"},
		{"(string-split 1 2 3)", "E1007"},
		{"(string-split #\\a #\\a)", "E1015"},
		{"(string-split \"\" \"\")", "E1019"},
		{"(string-split a #\\a)", "E1008"},
		{"(string-split \"\" a)", "E1008"},
	}
	executeTest(testCode, "string-split", t)
}
func TestStringJoin(t *testing.T) {
	testCode := [][]string{
		{"(string-join '(\"a\" \"b\" \"c\" \"d\" \"e\") \":\")", "\"a:b:c:d:e\""},
		{"(string-join '(\"a\" \"b\" \"c\" \"d\" \"e\") \"::\")", "\"a::b::c::d::e\""},
		{"(string-join '(\"a\") \"::\")", "\"a\""},
		{"(string-join '(\"\") \"::\")", "\"\""},

		{"(string-join)", "E1007"},
		{"(string-join 1 2 3)", "E1007"},
		{"(string-join  #\\a (list \"\" \"\"))", "E1005"},
		{"(string-join (list 1 \"a\" \"b\")  \",\")", "E1015"},
		{"(string-join (list \"a\" \"b\" 1) \",\")", "E1015"},
		{"(string-join a #\\a)", "E1008"},
		{"(string-join (list \"a\" \"b\"  \"c\") a)", "E1008"},
	}
	executeTest(testCode, "string-join", t)
}
func TestStringScan(t *testing.T) {
	testCode := [][]string{
		{"(string-scan \"abracadabra\" \"ada\")", "5"},
		{"(string-scan \"abracadabra\" #\\c)", "4"},
		{"(string-scan \"abracadabra\" \"aba\")", "#f"},
		{"(string-scan \"abracadabra\" #\\z)", "#f"},
		{"(string-scan \"1122\" #\\2)", "2"},

		{"(string-scan)", "E1007"},
		{"(string-scan \"abracadabra\")", "E1007"},
		{"(string-scan \"abracadabra\" \"aba\" \"aba\")", "E1007"},
		{"(string-scan 10  #\\z)", "E1015"},
		{"(string-scan \"abracadabra\" 10)", "E1009"},
		{"(string-scan a #\\2)", "E1008"},
		{"(string-scan \"1122\" a)", "E1008"},
	}
	executeTest(testCode, "string-scan", t)
}
func TestStringScanRight(t *testing.T) {
	testCode := [][]string{
		{"(string-scan-right \"abracadabra\" \"ada\")", "5"},
		{"(string-scan-right \"abracadabra\" #\\c)", "4"},
		{"(string-scan-right \"abracadabra\" \"aba\")", "#f"},
		{"(string-scan-right \"abracadabra\" #\\z)", "#f"},
		{"(string-scan-right \"1122\" #\\2)", "3"},

		{"(string-scan-right)", "E1007"},
		{"(string-scan-right \"abracadabra\")", "E1007"},
		{"(string-scan-right \"abracadabra\" \"aba\" \"aba\")", "E1007"},
		{"(string-scan-right 10  #\\z)", "E1015"},
		{"(string-scan-right \"abracadabra\" 10)", "E1009"},
		{"(string-scan-right a #\\2)", "E1008"},
		{"(string-scan-right \"1122\" a)", "E1008"},
	}
	executeTest(testCode, "string-scan-right", t)
}
