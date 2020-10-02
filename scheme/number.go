/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

import (
	"fmt"
	"reflect"
	"strconv"
)

type Number interface {
	Atom
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
	}

	tracer.Fatal("toFloat(): Impossible error\n")
	return NewFloat(1.0)
}
func castNumber(x Number, y Number) (a Number, b Number) {
	a = x
	b = y
	if _, ok := x.(*Float); ok {
		if v, ok := y.(*Integer); ok {
			b = NewFloat(float64(v.Value))
		}
	}
	if v, ok := x.(*Integer); ok {
		if _, ok := y.(*Float); ok {
			a = NewFloat(float64(v.Value))
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
	self.Value /= v.Value
	return self
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
func (self *Float) FormatString(prec int) string {
	return strconv.FormatFloat(self.Value, 'f', prec, 64)
}
func (self *Float) LogFormatString(prec int) string {
	return strconv.FormatFloat(self.Value, 'e', prec, 64)
}
