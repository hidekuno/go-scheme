/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

type Number interface {
	Expression
	Add(Number) Number
	Sub(Number) Number
	Mul(Number) Number
	Div(Number) Number
	Equal(Number) bool
	Greater(Number) bool
	Less(Number) bool
	GreaterEqual(Number) bool
	LessEqual(Number) bool
}

func CreateNumber(exp Expression) (Number, error) {
	if v, ok := exp.(*Integer); ok {
		return NewInteger(v.Value), nil
	}
	if v, ok := exp.(*Float); ok {
		return NewFloat(v.Value), nil
	}
	if v, ok := exp.(*Rat); ok {
		return v, nil
	}
	return nil, NewRuntimeError("E1003", reflect.TypeOf(exp).String())
}
func toInt(n Number) *Integer {
	if v, ok := n.(*Integer); ok {
		return v
	} else if v, ok := n.(*Float); ok {
		return NewInteger(int(v.Value))
	}
	tracer.Fatal("toInt(): Impossible error\n")
	return NewInteger(1)
}
func toFloat(n Number) *Float {
	if v, ok := n.(*Integer); ok {
		return NewFloat(float64(v.Value))
	} else if v, ok := n.(*Float); ok {
		return v
	} else if v, ok := n.(*Rat); ok {
		return NewFloat(float64(v.Value[0]) / float64(v.Value[1]))
	}
	tracer.Fatal("toFloat(): Impossible error\n")
	return NewFloat(1.0)
}
func castNumber(x Number, y Number) (a Number, b Number) {
	a = x
	b = y
	if _, ok := x.(*Float); ok {
		if v, ok := y.(*Integer); ok {
			b = toFloat(v)
		}
		if v, ok := y.(*Rat); ok {
			b = toFloat(v)
		}
	}
	if v, ok := x.(*Integer); ok {
		if _, ok := y.(*Float); ok {
			a = toFloat(v)
		}
		if _, ok := y.(*Rat); ok {
			a = NewRat(v.Value, 1)
		}
	}
	if v, ok := x.(*Rat); ok {
		if _, ok := y.(*Float); ok {
			a = toFloat(v)
		}
		if w, ok := y.(*Integer); ok {
			b = NewRat(w.Value, 1)
		}
	}
	return a, b
}

// Integer Type
type Integer struct {
	Number
	Value int
}

func NewInteger(n int) *Integer {
	v := new(Integer)
	v.Value = n
	return v
}
func (self *Integer) Add(n Number) Number {
	self.Value += toInt(n).Value
	return self
}
func (self *Integer) Sub(n Number) Number {
	self.Value -= toInt(n).Value
	return self
}
func (self *Integer) Mul(n Number) Number {
	self.Value *= toInt(n).Value
	return self
}
func (self *Integer) Div(n Number) Number {
	v := toInt(n)
	if v.Value == 0 {
		panic(NewRuntimeError("E1013"))
	}
	if 0 == self.Value%v.Value {
		self.Value /= v.Value
		return self
	} else {
		return NewRat(self.Value, v.Value)
	}
}

func (self *Integer) Equal(n Number) bool {
	return self.Value == toInt(n).Value
}
func (self *Integer) Greater(n Number) bool {
	return self.Value > toInt(n).Value
}
func (self *Integer) Less(n Number) bool {
	return self.Value < toInt(n).Value
}
func (self *Integer) GreaterEqual(n Number) bool {
	return self.Value >= toInt(n).Value
}
func (self *Integer) LessEqual(n Number) bool {
	return self.Value <= toInt(n).Value
}
func (self *Integer) String() string {
	return strconv.Itoa(self.Value)
}
func (self *Integer) Print() {
	fmt.Print(self.Value)
}
func (self *Integer) isAtom() bool {
	return true
}
func (self *Integer) clone() Expression {
	return NewInteger(self.Value)
}
func (self *Integer) equalValue(e Expression) bool {
	if v, ok := e.(*Integer); ok {
		return self.Value == v.Value
	}
	return false
}

// Float Type
type Float struct {
	Number
	Value float64
}

func NewFloat(n float64) *Float {
	v := new(Float)
	v.Value = n
	return v
}
func (self *Float) Add(n Number) Number {
	self.Value += toFloat(n).Value
	return self
}
func (self *Float) Sub(n Number) Number {
	self.Value -= toFloat(n).Value
	return self
}
func (self *Float) Mul(n Number) Number {
	self.Value *= toFloat(n).Value
	return self
}
func (self *Float) Div(n Number) Number {
	self.Value /= toFloat(n).Value
	return self
}
func (self *Float) Equal(n Number) bool {
	return self.Value == toFloat(n).Value
}
func (self *Float) Greater(n Number) bool {
	return self.Value > toFloat(n).Value
}
func (self *Float) Less(n Number) bool {
	return self.Value < toFloat(n).Value
}
func (self *Float) GreaterEqual(n Number) bool {
	return self.Value >= toFloat(n).Value
}
func (self *Float) LessEqual(n Number) bool {
	return self.Value <= toFloat(n).Value
}
func (self *Float) String() string {
	return fmt.Sprint(self.Value)
}
func (self *Float) Print() {
	fmt.Print(self.Value)
}
func (self *Float) isAtom() bool {
	return true
}
func (self *Float) clone() Expression {
	return NewFloat(self.Value)
}
func (self *Float) equalValue(e Expression) bool {
	if v, ok := e.(*Float); ok {
		return self.Value == v.Value
	}
	return false
}
func (self *Float) FormatString(prec int) string {
	return strconv.FormatFloat(self.Value, 'f', prec, 64)
}
func (self *Float) LogFormatString(prec int) string {
	return strconv.FormatFloat(self.Value, 'e', prec, 64)
}

// Rat Type
type Rat struct {
	Number
	Value []int
}

func NewRat(m int, n int) *Rat {
	if n == 0 {
		panic(NewRuntimeError("E1013"))
	}
	v := new(Rat)
	v.Value = []int{m, n}
	v.calc()
	return v
}
func CreateRat(m int, n int) Number {
	v := NewRat(m, n)
	if v.Value[1] == 1 {
		return NewInteger(v.Value[0])
	}
	return v
}
func (self *Rat) Add(n Number) Number {
	if v, ok := n.(*Rat); ok {
		return CreateRat(self.Value[0]*v.Value[1]+v.Value[0]*self.Value[1], self.Value[1]*v.Value[1])
	}
	panic(NewRuntimeError("E1022"))
}
func (self *Rat) Sub(n Number) Number {
	if v, ok := n.(*Rat); ok {
		return CreateRat(self.Value[0]*v.Value[1]-v.Value[0]*self.Value[1], self.Value[1]*v.Value[1])
	}
	panic(NewRuntimeError("E1022"))
}
func (self *Rat) Mul(n Number) Number {
	if v, ok := n.(*Rat); ok {
		return CreateRat(self.Value[0]*v.Value[0], self.Value[1]*v.Value[1])
	}
	panic(NewRuntimeError("E1022"))
}
func (self *Rat) Div(n Number) Number {
	if v, ok := n.(*Rat); ok {
		return CreateRat(self.Value[0]*v.Value[1], self.Value[1]*v.Value[0])
	}
	panic(NewRuntimeError("E1022"))
}

func (self *Rat) Equal(n Number) bool {
	if v, ok := n.(*Rat); ok {
		return self.Value[0]*v.Value[1] == v.Value[0]*self.Value[1]
	}
	panic(NewRuntimeError("E1022"))
}
func (self *Rat) Greater(n Number) bool {
	if v, ok := n.(*Rat); ok {
		return self.Value[0]*v.Value[1] > v.Value[0]*self.Value[1]
	}
	panic(NewRuntimeError("E1022"))
}
func (self *Rat) Less(n Number) bool {
	if v, ok := n.(*Rat); ok {
		return self.Value[0]*v.Value[1] < v.Value[0]*self.Value[1]
	}
	panic(NewRuntimeError("E1022"))
}
func (self *Rat) GreaterEqual(n Number) bool {
	if v, ok := n.(*Rat); ok {
		return self.Value[0]*v.Value[1] >= v.Value[0]*self.Value[1]
	}
	panic(NewRuntimeError("E1022"))
}
func (self *Rat) LessEqual(n Number) bool {
	if v, ok := n.(*Rat); ok {
		return self.Value[0]*v.Value[1] <= v.Value[0]*self.Value[1]
	}
	panic(NewRuntimeError("E1022"))
}
func (self *Rat) String() string {
	return strconv.Itoa(self.Value[0]) + "/" + strconv.Itoa(self.Value[1])
}
func (self *Rat) Print() {
	fmt.Print(self.String())
}
func (self *Rat) isAtom() bool {
	return true
}
func (self *Rat) equalValue(e Expression) bool {
	if v, ok := e.(*Rat); ok {
		return self.Value[0]*v.Value[1] == v.Value[0]*self.Value[1]
	}
	return false
}
func (self *Rat) calc() {

	var gcm func(m, n int) int

	gcm = func(m, n int) int {
		mod := m % n
		if mod == 0 {
			return n
		}
		return gcm(n, mod)
	}
	g := gcm(self.Value[0], self.Value[1])
	if 0 > g {
		g = -g
	}
	sign := 1
	if self.Value[0]*self.Value[1] < 0 {
		sign = -1
	}
	self.Value[0] = self.Value[0] / g
	self.Value[1] = self.Value[1] / g

	self.Value[0] = int(math.Abs(float64(self.Value[0])))
	self.Value[1] = int(math.Abs(float64(self.Value[1])))
	self.Value[0] *= sign
}
func MakeRat(s string) Number {
	return MakeRatRadix(s, 10)
}
func MakeRatRadix(s string, r int) Number {
	rat := strings.Split(s, "/")
	if len(rat) == 2 {
		if m, err := strconv.ParseInt(rat[0], r, 0); err == nil {
			if n, err := strconv.ParseInt(rat[1], r, 0); err == nil {
				if 0 != n {
					return CreateRat(int(m), int(n))
				}
			}
		}
	}
	return nil
}
