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
func (self *Integer) Add(n *Integer) *Integer {
	self.Value = self.Value.Add(self.Value, n.Value)
	return self
}
func (self *Integer) Sub(n *Integer) *Integer {
	self.Value = self.Value.Sub(self.Value, n.Value)
	return self
}
func (self *Integer) Mul(n *Integer) *Integer {
	self.Value = self.Value.Mul(self.Value, n.Value)
	return self
}
func (self *Integer) Div(n *Integer) *Integer {
	self.Value = self.Value.Div(self.Value, n.Value)
	return self
}
func (self *Integer) Equal(n *Integer) bool {
	return self.Value.Cmp(n.Value) == 0
}
func (self *Integer) Greater(n *Integer) bool {
	return self.Value.Cmp(n.Value) == 1
}
func (self *Integer) Less(n *Integer) bool {
	return self.Value.Cmp(n.Value) == -1
}
func (self *Integer) GreaterEqual(n *Integer) bool {
	return self.Value.Cmp(n.Value) >= 0
}
func (self *Integer) LessEqual(n *Integer) bool {
	return self.Value.Cmp(n.Value) <= 0
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
	t1 := ParseInteger("340282366920938463463374607431768211456")
	if t1 == nil {
		return
	}
	t2 := ParseInteger("100000000000000000000000000000000000000")
	if t2 == nil {
		return
	}
	t1 = t1.Add(t2)
	println(t1.Value.String())
	t2 = ParseInteger("2")
	if t2 == nil {
		return
	}
	t1 = t1.Add(t2)
	println(t1.Value.String())

	t2 = ParseInteger("10")
	if t2 == nil {
		return
	}
	t1 = t1.Sub(t2)
	println(t1.Value.String())

	t1 = t1.Mul(t2)
	println(t1.Value.String())

	t2 = ParseInteger("2")
	if t2 == nil {
		return
	}
	t1 = t1.Div(t2)
	println(t1.Value.String())

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
	println(t1.Greater(t2))
	println(t2.Less(t1))
}
