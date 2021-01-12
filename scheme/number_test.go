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

func TestRatAtom(t *testing.T) {
	testCode := [][]string{
		{"5/6", "5/6"},
		{"5/1", "5"},
		{"0/5", "0"},
		{"5/0", "E1008"},
	}
	executeTest(testCode, "rat_plus", t)
}
func TestRatPlus(t *testing.T) {
	testCode := [][]string{
		{"(+ (/ 1 2)(/ 1 3))", "5/6"},
		{"(+ (/ 1 2) 3)", "7/2"},
		{"(+ 3 (/ 1 2) )", "7/2"},
		{"(+ 0.5 (/ 1 2) )", "1"},
		{"(+ (/ 1 2) 0.5)", "1"},
	}
	executeTest(testCode, "rat_plus", t)
}
func TestRatMinus(t *testing.T) {
	testCode := [][]string{
		{"(- (/ 1 2)(/ 1 3))", "1/6"},
		{"(- 3 (/ 1 2))", "5/2"},
		{"(- (/ 1 2) 3)", "-5/2"},
		{"(- 0.15 (/ 1 2))", "-0.35"},
		{"(- (/ 1 2) 0.15)", "0.35"},
	}
	executeTest(testCode, "rat_minus", t)
}
func TestRatMulti(t *testing.T) {
	testCode := [][]string{
		{"(* (/ 1 2)(/ 1 4))", "1/8"},
		{"(* (/ 1 2) 3)", "3/2"},
		{"(* 3 (/ 1 2))", "3/2"},
		{"(* 0.25 (/ 1 2))", "0.125"},
		{"(* (/ 1 2) 0.25)", "0.125"},
	}
	executeTest(testCode, "rat_multi", t)
}
func TestRatDiv(t *testing.T) {
	testCode := [][]string{
		{"(/ 4 3)", "4/3"},
		{"(/ -1 3)", "-1/3"},
		{"(/ 1 -3)", "-1/3"},
		{"(/ -1 -3)", "1/3"},
		{"(+ (/ -1 3)(/ 1 3))", "0"},
	}
	executeTest(testCode, "rat_div", t)
}

func TestRatEq(t *testing.T) {
	testCode := [][]string{
		{"(= (/ 4 2) 2)", "#t"},
		{"(= 2 (/ 4 2))", "#t"},
		{"(= (/ 4 8) (/ 2 4))", "#t"},
		{"(= (/ 3 8) (/ 2 4))", "#f"},
		{"(= 0.5 (/ 1 2))", "#t"},
		{"(= (/ 3 2) 1.5)", "#t"},
		{"(= 0.6 (/ 1 2))", "#f"},
		{"(= (/ 3 2) 1.6)", "#f"},
	}
	executeTest(testCode, "rat_eq", t)
}
func TestRatThan(t *testing.T) {
	testCode := [][]string{
		{"(> (/ 7 2) 3)", "#t"},
		{"(> (/ 1 2) 0.3)", "#t"},
		{"(> 4 (/ 7 2))", "#t"},
		{"(> 1.6 (/ 3 2))", "#t"},
		{"(> (/ 3 4)(/ 4 8))", "#t"},
	}
	executeTest(testCode, "rat_than", t)
}
func TestRatLess(t *testing.T) {
	testCode := [][]string{
		{"(< 3 (/ 7 2))", "#t"},
		{"(< 0.3 (/ 1 2))", "#t"},
		{"(< (/ 7 2) 4)", "#t"},
		{"(< (/ 4 8)(/ 3 4))", "#t"},
	}
	executeTest(testCode, "rat_less", t)
}
func TestRatThanEq(t *testing.T) {
	testCode := [][]string{
		{"(>= (/ 7 2) 3)", "#t"},
		{"(>= (/ 1 2) 0.3)", "#t"},
		{"(>= 4 (/ 7 2))", "#t"},
		{"(>= 1.6 (/ 3 2))", "#t"},
		{"(>= (/ 3 4) (/ 4 8))", "#t"},
	}
	executeTest(testCode, "rat_than_eq", t)
}
func TestRatLessEq(t *testing.T) {
	testCode := [][]string{
		{"(<= 3 (/ 7 2))", "#t"},
		{"(<= 0.3 (/ 1 2))", "#t"},
		{"(<= (/ 7 2) 4)", "#t"},
		{"(<= (/ 4 8) (/ 3 4))", "#t"},
	}
	executeTest(testCode, "rat_less_eq", t)
}
