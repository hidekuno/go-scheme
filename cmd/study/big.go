package main

import (
	"math"
	"math/big"
	"strconv"
)

const INTUNIT = 18

type Integer struct {
	Value *big.Int
}

func NewInteger(n *big.Int) *Integer {
	v := new(Integer)
	v.Value = n
	return v
}
func (self *Integer) String() string {
	return self.Value.String()
}
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
func ParseInteger(s string) *Integer {
	unit := big.NewInt(int64(math.Pow(10.0, 18.0)))

	n := Reverse(s)
	total := big.NewInt(0)
	mul := big.NewInt(1)
	for i := 0; i < len(s); i += INTUNIT {
		l := i + INTUNIT
		if len(s)-i < INTUNIT {
			l = len(s)
		}
		if ivalue, err := strconv.Atoi(Reverse(n[i:l])); err == nil {
			b := big.NewInt(int64(ivalue))
			b = b.Mul(b, mul)
			total = total.Add(total, b)
		} else {
			return nil
		}
		mul = mul.Mul(mul, unit)
	}
	return NewInteger(total)
}
func main() {
	if testData := ParseInteger("340282366920938463463374607431768211456"); testData != nil {
		println(testData.Value.String())
	}
	if testData := ParseInteger("100000000000000000000000000000000000000"); testData != nil {
		println(testData.Value.String())
	}
	if testData := ParseInteger("100000000000000000"); testData != nil {
		println(testData.Value.String())
	}
	if testData := ParseInteger("1000000000000000000"); testData != nil {
		println(testData.Value.String())
	}
}
