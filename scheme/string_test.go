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
		{"(number->string 1/3)", "\"1/3\""},
		{"(number->string 3735927486 2)", "\"11011110101011011011101010111110\""},
		{"(number->string 3735927486 8)", "\"33653335276\""},
		{"(number->string 3735927486 10)", "\"3735927486\""},
		{"(number->string 3735927486 16)", "\"deadbabe\""},
		{"(number->string 3735927486 36)", "\"1ps9w3i\""},

		{"(number->string)", "E1007"},
		{"(number->string 10 20 30)", "E1007"},
		{"(number->string #f)", "E1003"},
		{"(number->string #f 10)", "E1003"},
		{"(number->string 100 1)", "E1021"},
		{"(number->string 100 37)", "E1021"},
		{"(number->string a)", "E1008"},
		{"(number->string 10 a)", "E1008"},
	}
	executeTest(testCode, "number->string", t)
}
func TestStringNumber(t *testing.T) {
	testCode := [][]string{
		{"(string->number \"123\")", "123"},
		{"(string->number \"10.5\")", "10.5"},
		{"(string->number \"1/3\")", "1/3"},
		{"(string->number \"10000\" 2)", "16"},
		{"(string->number \"012\" 8)", "10"},
		{"(string->number \"123\" 10)", "123"},
		{"(string->number \"ab\" 16)", "171"},

		{"(string->number)", "E1007"},
		{"(string->number \"123\" \"10.5\" 10)", "E1007"},
		{"(string->number 100)", "E1015"},
		{"(string->number 100 10)", "E1015"},
		{"(string->number 100 #f)", "E1002"},
		{"(string->number 100 1)", "E1021"},
		{"(string->number 100 37)", "E1021"},
		{"(string->number a)", "E1008"},
		{"(string->number 10 a)", "E1008"},
		{"(string->number \"ab\" 2)", "#f"},
		{"(string->number \"ab\" 8)", "#f"},
		{"(string->number \"ab\" 10)", "#f"},
		{"(string->number \"/1\")", "#f"},
		{"(string->number \"1/3/2\")", "#f"},
		{"(string->number \"1/0\")", "#f"},
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
func TestStringVector(t *testing.T) {
	testCode := [][]string{
		{"(string->vector \"\")", "#()"},
		{"(string->vector \"abc\")", "#(#\\a #\\b #\\c)"},
		{"(string->vector \"山田\")", "#(#\\山 #\\田)"},
		{"(string->vector)", "E1007"},
		{"(string->vector \"a\" \"b\")", "E1007"},
		{"(string->vector #\\a)", "E1015"},
		{"(string->vector a)", "E1008"},
	}
	executeTest(testCode, "string->vector", t)
}
func TestVectorString(t *testing.T) {
	testCode := [][]string{
		{"(vector->string (vector))", "\"\""},
		{"(vector->string (vector #\\a #\\b #\\c))", "\"abc\""},
		{"(vector->string)", "E1007"},
		{"(vector->string (list #\\a #\\b)(list #\\a #\\b))", "E1007"},
		{"(vector->string 10)", "E1022"},
		{"(vector->string (vector #\\a 10))", "E1019"},
		{"(vector->string a)", "E1008"},
	}
	executeTest(testCode, "vector->string", t)
}
func TestStringReverse(t *testing.T) {
	testCode := [][]string{
		{"(string-reverse \"1234567890\")", "\"0987654321\""},
		{"(string-reverse \"1234567890\" 3)", "\"0987654\""},
		{"(string-reverse \"1234567890\" 3 8)", "\"87654\""},
		{"(string-reverse  \"山川\")", "\"川山\""},
		{"(string-reverse)", "E1007"},
		{"(string-reverse \"abcdefghijklmn\" 2 3 4)", "E1007"},
		{"(string-reverse 10)", "E1015"},
		{"(string-reverse \"abcdefghijklmn\" #\\a 1)", "E1002"},
		{"(string-reverse \"abcdefghijklmn\" 1 #\\a)", "E1002"},
		{"(string-reverse \"abcdefghijklmn\" 1 20)", "E1021"},
		{"(string-reverse \"abcdefghijklmn\" 6 5)", "E1021"},
	}
	executeTest(testCode, "string-reverse", t)
}
func TestStringUpcase(t *testing.T) {
	testCode := [][]string{
		{"(string-upcase \"ab1012cd\")", "\"AB1012CD\""},
		{"(string-upcase \"abcd\" 1)", "\"BCD\""},
		{"(string-upcase \"abcd\" 1 2)", "\"B\""},
		{"(string-upcase)", "E1007"},
		{"(string-upcase \"abcdefghijklmn\" 2 3 4)", "E1007"},
		{"(string-upcase 10)", "E1015"},
		{"(string-upcase \"abcdefghijklmn\" #\\a 1)", "E1002"},
		{"(string-upcase \"abcdefghijklmn\" 1 #\\a)", "E1002"},
		{"(string-upcase \"abcdefghijklmn\" 1 20)", "E1021"},
		{"(string-upcase \"abcdefghijklmn\" 6 5)", "E1021"},
	}
	executeTest(testCode, "string-upcase", t)
}

func TestStringDowncase(t *testing.T) {
	testCode := [][]string{
		{"(string-downcase \"AB1012CD\")", "\"ab1012cd\""},
		{"(string-downcase \"ABCD\" 1)", "\"bcd\""},
		{"(string-downcase \"ABCD\" 1 2)", "\"b\""},
		{"(string-downcase)", "E1007"},
		{"(string-downcase \"ABCDEFGHIJKLMN\" 2 3 4)", "E1007"},
		{"(string-downcase 10)", "E1015"},
		{"(string-downcase \"ABCDEFGHIJKLMN\" 1 #\\a)", "E1002"},
		{"(string-downcase \"ABCDEFGHIJKLMN\" 1 #\\a)", "E1002"},
		{"(string-downcase \"ABCDEFGHIJKLMN\" 1 20)", "E1021"},
		{"(string-downcase \"ABCDEFGHIJKLMN\" 6 5)", "E1021"},
	}
	executeTest(testCode, "string-downcase", t)
}
func TestStringIndex(t *testing.T) {
	testCode := [][]string{
		{"(string-index \"abcdefghijlklmn\" #\\a)", "0"},
		{"(string-index \"abcdefghijlklmn\" #\\e)", "4"},
		{"(string-index \"abcdefghijlklmn\" #\\z)", "#f"},
		{"(string-index \"abcdefghijlklmn\" #\\c 2)", "2"},
		{"(string-index \"abcdefghijlklmn\" #\\d 2 8)", "3"},
		{"(string-index \"abcdefghijlklmn\" #\\k 2 8)", "#f"},
		{"(string-index \"abcdefghijlklcn\" #\\n 0 14)", "#f"},
		{"(string-index \"abcdefghijlklcn\" #\\n 0 15)", "14"},
		{"(string-index)", "E1007"},
		{"(string-index \"abcdefghijklmn\" 2 3 4 5)", "E1007"},
		{"(string-index 10 #\\a)", "E1015"},
		{"(string-index \"abc\" 10)", "E1019"},
		{"(string-index \"abcdefghijklmn\" #\\a #\\a 1)", "E1002"},
		{"(string-index \"abcdefghijklmn\" #\\a  1 #\\a)", "E1002"},
		{"(string-index \"abcdefghijklmn\" #\\a 1 20)", "E1021"},
		{"(string-index \"abcdefghijklmn\" #\\a 6 5)", "E1021"},
	}
	executeTest(testCode, "string-index", t)
}
func TestStringIndexRight(t *testing.T) {
	testCode := [][]string{
		{"(string-index-right \"abcdefghijlklmn\" #\\a)", "0"},
		{"(string-index-right \"abcdefghijlklmn\" #\\z)", "#f"},
		{"(string-index-right \"abcdefghijlklcn\" #\\c 2)", "13"},
		{"(string-index-right \"abcdefghijlklcn\" #\\n 2 14)", "#f"},
		{"(string-index-right \"abcdefghijlklcn\" #\\n 2 15)", "14"},
		{"(string-index-right)", "E1007"},
		{"(string-index-right \"abcdefghijklmn\" 2 3 4 5)", "E1007"},
		{"(string-index-right 10 #\\a)", "E1015"},
		{"(string-index-right \"abc\" 10)", "E1019"},
		{"(string-index-right \"abcdefghijklmn\" #\\a #\\a 1)", "E1002"},
		{"(string-index-right \"abcdefghijklmn\" #\\a  1 #\\a)", "E1002"},
		{"(string-index-right \"abcdefghijklmn\" #\\a 1 20)", "E1021"},
		{"(string-index-right \"abcdefghijklmn\" #\\a 6 5)", "E1021"},
	}
	executeTest(testCode, "string-index-right", t)
}
func TestStringDelete(t *testing.T) {
	testCode := [][]string{
		{"(string-delete \"abcdefghijlklcn\" #\\a)", "\"bcdefghijlklcn\""},
		{"(string-delete \"abcdefghijlklcn\" #\\a 3)", "\"defghijlklcn\""},
		{"(string-delete \"abcdefghijlklcn\" #\\n 13)", "\"c\""},
		{"(string-delete \"abcdefghijlklcn\" #\\n 14)", "\"\""},
		{"(string-delete \"abcdefghijlklcn\" #\\a 3 4)", "\"d\""},
		{"(string-delete \"aaaaaaaaaaaaaaaaaaaa\" #\\a 3 4)", "\"\""},
		{"(string-delete \"abcdefghijlklcn\" #\\a 3 9)", "\"defghi\""},
		{"(string-delete \"abcdefghijlklcn\" #\\h 3 9)", "\"defgi\""},
		{"(string-delete)", "E1007"},
		{"(string-delete \"abcdefghijklmn\" 2 3 4 5)", "E1007"},
		{"(string-delete 10 #\\a)", "E1015"},
		{"(string-delete \"abc\" 10)", "E1019"},
		{"(string-delete \"abcdefghijklmn\" #\\a #\\a 1)", "E1002"},
		{"(string-delete \"abcdefghijklmn\" #\\a  1 #\\a)", "E1002"},
		{"(string-delete \"abcdefghijklmn\" #\\a 1 20)", "E1021"},
		{"(string-delete \"abcdefghijklmn\" #\\a 6 5)", "E1021"},
	}
	executeTest(testCode, "string-delete", t)
}
func TestStringTrim(t *testing.T) {
	testCode := [][]string{
		{"(string-trim  \"  ad  \")", "\"ad  \""},
		{"(string-trim \"ada\" #\\a)", "\"da\""},
		{"(string-trim)", "E1007"},
		{"(string-trim \"abcdefghijklmn\" 2 3 4)", "E1007"},
		{"(string-trim 10 #\\a)", "E1015"},
		{"(string-trim \"abcdefghijklmn\" 2)", "E1019"},
	}
	executeTest(testCode, "string-trim", t)
}
func TestStringTrimRight(t *testing.T) {
	testCode := [][]string{
		{"(string-trim-right  \"  ad  \")", "\"  ad\""},
		{"(string-trim-right \"ada\" #\\a)", "\"ad\""},
		{"(string-trim-right)", "E1007"},
		{"(string-trim-right \"abcdefghijklmn\" 2 3 4)", "E1007"},
		{"(string-trim-right 10 #\\a)", "E1015"},
		{"(string-trim-right \"abcdefghijklmn\" 2)", "E1019"},
	}
	executeTest(testCode, "string-trim-right", t)
}
func TestStringTrimBoth(t *testing.T) {
	testCode := [][]string{
		{"(string-trim-both  \"  ad  \")", "\"ad\""},
		{"(string-trim-both \"ada\" #\\a)", "\"d\""},
		{"(string-trim-both)", "E1007"},
		{"(string-trim-both \"abcdefghijklmn\" 2 3 4)", "E1007"},
		{"(string-trim-both 10 #\\a)", "E1015"},
		{"(string-trim-both \"abcdefghijklmn\" 2)", "E1019"},
	}
	executeTest(testCode, "string-trim-both", t)
}
func TestStringTake(t *testing.T) {
	testCode := [][]string{
		{"(string-take \"1234567890\" 0)", "\"\""},
		{"(string-take \"1山2\" 2)", "\"1山\""},
		{"(string-take \"1234567890\" 10)", "\"1234567890\""},
		{"(string-take)", "E1007"},
		{"(string-take 2 3 4)", "E1007"},
		{"(string-take \"123456\" #\\a)", "E1002"},
		{"(string-take \"abcdefghijklmn\" -1)", "E1021"},
		{"(string-take \"abcdefghijklmn\" 15)", "E1021"},
	}
	executeTest(testCode, "string-take", t)
}
func TestStringTakeRigth(t *testing.T) {
	testCode := [][]string{
		{"(string-take-right \"1234567890\" 0)", "\"\""},
		{"(string-take-right \"1山2\" 2)", "\"山2\""},
		{"(string-take-right \"1234567890\" 10)", "\"1234567890\""},
		{"(string-take-right)", "E1007"},
		{"(string-take-right 2 3 4)", "E1007"},
		{"(string-take-right \"123456\" #\\a)", "E1002"},
		{"(string-take-right \"abcdefghijklmn\" -1)", "E1021"},
		{"(string-take-right \"abcdefghijklmn\" 15)", "E1021"},
	}
	executeTest(testCode, "string-take-right", t)
}
func TestStringDrop(t *testing.T) {
	testCode := [][]string{
		{"(string-drop \"1234567890\" 0)", "\"1234567890\""},
		{"(string-drop \"1山2\" 1)", "\"山2\""},
		{"(string-drop \"1234567890\" 10)", "\"\""},
		{"(string-drop)", "E1007"},
		{"(string-drop 2 3 4)", "E1007"},
		{"(string-drop \"123456\" #\\a)", "E1002"},
		{"(string-drop \"abcdefghijklmn\" -1)", "E1021"},
		{"(string-drop \"abcdefghijklmn\" 15)", "E1021"},
	}
	executeTest(testCode, "string-drop", t)
}
func TestStringDropRigth(t *testing.T) {
	testCode := [][]string{
		{"(string-drop-right \"1234567890\" 0)", "\"1234567890\""},
		{"(string-drop-right \"1山2\" 1)", "\"1山\""},
		{"(string-drop-right \"1234567890\" 10)", "\"\""},
		{"(string-drop-right)", "E1007"},
		{"(string-drop-right 2 3 4)", "E1007"},
		{"(string-drop-right \"123456\" #\\a)", "E1002"},
		{"(string-drop-right \"abcdefghijklmn\" -1)", "E1021"},
		{"(string-drop-right \"abcdefghijklmn\" 15)", "E1021"},
	}
	executeTest(testCode, "string-drop-right", t)
}
