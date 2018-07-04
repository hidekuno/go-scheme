/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type (
	Environment map[string]Expression
)

const (
	MAX_LINE_SIZE = 1024
	DEBUG         = false
)

var (
	builtin_func map[string]func(...Expression) (Expression, error)
	special_func map[string]func(*Environment, []Expression) (Expression, error)
	define_env   Environment
)

// Basic Data Type. (need
type SyntaxError struct {
	Msg string
}

func (err *SyntaxError) Error() string {
	return err.Msg
}

func NewSyntaxError(text string) error {
	return &SyntaxError{text}
}

type RuntimeError struct {
	Msg string
}

func (err *RuntimeError) Error() string {
	return err.Msg
}

func NewRuntimeError(text string) error {
	return &RuntimeError{text}
}

type Expression interface {
	Print()
}

type Any interface{}
type Atom interface {
	Expression
	// Because Expression is different
	Dummy() Any
}

type Symbol struct {
	Atom
	Value string
}

func NewSymbol(token string) *Symbol {
	s := new(Symbol)
	s.Value = token
	return s
}

func (self *Symbol) Print() {
	fmt.Print(self.Value)
}

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

type Integer struct {
	Number
	Value int
}

func NewInteger(p int) *Integer {
	v := new(Integer)
	v.Value = p
	return v
}

func (self *Integer) Print() {
	fmt.Print(self.Value)
}

func (self *Integer) Add(p Number) Number {
	v, _ := p.(*Integer)
	self.Value += v.Value
	return self
}
func (self *Integer) Sub(p Number) Number {
	v, _ := p.(*Integer)
	self.Value -= v.Value
	return self
}
func (self *Integer) Mul(p Number) Number {
	v, _ := p.(*Integer)
	self.Value *= v.Value
	return self
}
func (self *Integer) Div(p Number) Number {
	v, _ := p.(*Integer)
	self.Value /= v.Value
	return self
}

func (self *Integer) Equal(p Number) bool {
	v, _ := p.(*Integer)
	return self.Value == v.Value
}
func (self *Integer) Greater(p Number) bool {
	v, _ := p.(*Integer)
	return self.Value > v.Value
}
func (self *Integer) Less(p Number) bool {
	v, _ := p.(*Integer)
	return self.Value < v.Value
}
func (self *Integer) GreaterEqual(p Number) bool {
	v, _ := p.(*Integer)
	return self.Value >= v.Value
}
func (self *Integer) LessEqual(p Number) bool {
	v, _ := p.(*Integer)
	return self.Value <= v.Value
}

type Boolean struct {
	Atom
	Value bool
	name  string
}

func NewBoolean(v bool) *Boolean {
	b := new(Boolean)
	b.Value = v
	if v {
		b.name = "#t"
	} else {
		b.name = "#f"
	}
	return b
}

func (self *Boolean) Print() {
	fmt.Print(self.name)
}

type Float struct {
	Number
	Value float64
}

func NewFloat(p float64) *Float {
	v := new(Float)
	v.Value = p
	return v
}
func (self *Float) Print() {
	fmt.Print(self.Value)
}
func (self *Float) Add(p Number) Number {
	v, _ := p.(*Float)
	self.Value += v.Value
	return self
}
func (self *Float) Sub(p Number) Number {
	v, _ := p.(*Float)
	self.Value -= v.Value
	return self
}
func (self *Float) Mul(p Number) Number {
	v, _ := p.(*Float)
	self.Value *= v.Value
	return self
}
func (self *Float) Div(p Number) Number {
	v, _ := p.(*Float)
	self.Value /= v.Value
	return self
}
func (self *Float) Equal(p Number) bool {
	v, _ := p.(*Float)
	return self.Value == v.Value
}
func (self *Float) Greater(p Number) bool {
	v, _ := p.(*Float)
	return self.Value > v.Value
}
func (self *Float) Less(p Number) bool {
	v, _ := p.(*Float)
	return self.Value < v.Value
}
func (self *Float) GreaterEqual(p Number) bool {
	v, _ := p.(*Float)
	return self.Value >= v.Value
}
func (self *Float) LessEqual(p Number) bool {
	v, _ := p.(*Float)
	return self.Value <= v.Value
}

func CreateNumber(p Number) (Number, error) {
	if v, ok := p.(*Integer); ok {
		return NewInteger(v.Value), nil
	}
	if v, ok := p.(*Float); ok {
		return NewFloat(v.Value), nil
	}
	return nil, NewRuntimeError("Not Number")
}

type String struct {
	Atom
	Value string
}

func NewString(p string) *String {
	v := new(String)
	v.Value = strings.Replace(p, "\"", "", -1)
	return v
}

func (self *String) Print() {
	fmt.Print("\"" + self.Value + "\"")
}

type List struct {
	Expression
	Value []Expression
}

func NewList(el []Expression) *List {
	l := new(List)
	l.Value = el
	return l
}

func (self *List) Print() {
	var tprint func(*List)
	tprint = func(l *List) {
		fmt.Print("(")

		for _, i := range l.Value {
			if j, ok := i.(*List); ok {
				tprint(j)
			} else if j, ok := i.(Expression); ok {
				j.Print()
			}
			if i != l.Value[len(l.Value)-1] {
				fmt.Print(" ")
			}
		}
		fmt.Print(")")
	}
	tprint(self)
}

type Pair struct {
	Expression
	Car Expression
	Cdr Expression
}

func NewPair(car Expression, cdr Expression) *Pair {
	p := new(Pair)
	p.Car = car
	p.Cdr = cdr
	return p
}

func (self *Pair) Print() {
	fmt.Print("(")
	self.Car.Print()
	fmt.Print(" . ")
	self.Cdr.Print()
	fmt.Print(")")
}

type Operator struct {
	Expression
	Impl func(...Expression) (Expression, error)
}

func NewOperator(fn func(...Expression) (Expression, error)) *Operator {
	op := new(Operator)
	op.Impl = fn
	return op
}

func (self *Operator) Print() {
	fmt.Print("Operatotion or Builtin: ", self)
}

type Function struct {
	Expression
	ParamName Expression
	Body      []Expression
	Env       *Environment
}

func NewFunction(param Expression, body []Expression) *Function {
	fn := new(Function)
	fn.ParamName = param
	fn.Body = body
	fn.Env = &Environment{}
	return fn
}

func (self *Function) Print() {
	fmt.Print("Function: ", self)
}

// Bind lambda function' parameters.
func (self *Function) BindParam(env *Environment, values []Expression) (*Environment, error) {

	plist, _ := self.ParamName.(*List)
	local_env := Environment{}
	for key, _ := range *self.Env {
		local_env[key] = (*self.Env)[key]
	}

	idx := 0
	for _, i := range plist.Value {
		if sym, ok := i.(*Symbol); ok {
			if idx+1 > len(values) {
				return nil, NewRuntimeError("Not Enough ParamName Number")
			}
			if env != nil {
				v, err := eval(values[idx], env)
				if err != nil {
					return env, err
				}
				local_env[sym.Value] = v
			} else {
				local_env[sym.Value] = values[idx]
			}
			idx = idx + 1
		}
	}
	return &local_env, nil
}
func (self *Function) SetKeyNameEnv(env *Environment) {
	for key, _ := range *self.Env {
		if _, ok := (*env)[key]; ok {
			(*self.Env)[key] = (*env)[key]
		}
	}
}
func (self *Function) Execute(env *Environment) (Expression, error) {
	var (
		result Expression
		err    error
	)
	for _, exp := range self.Body {
		result, err = eval(exp, env)
		if err != nil {
			return exp, err
		}
		// (lambda () (let ((c 0)) (lambda () (set! c (+ 1 c))))))
		if closure, ok := result.(*Function); ok {
			closure.Env = env
		} else {
			self.SetKeyNameEnv(env)
		}
	}
	return result, nil
}

type LetLoop struct {
	Expression
	ParamName Expression
	Body      Expression
}

func NewLetLoop(param Expression, body Expression) *LetLoop {
	let := new(LetLoop)
	let.ParamName = param
	let.Body = body
	return let
}

func (self *LetLoop) Print() {
	fmt.Print("Let Macro: ", self)
}

// Parse from tokens,
func parse(line string) (Expression, error) {
	token := tokenize(line)
	ast, c, err := create_ast(token)

	if err != nil {
		return nil, err
	}
	if c != len(token) {
		err := NewSyntaxError("extra close parenthesis `)'")
		return nil, err
	}
	return ast, nil
}

// Tenuki lex.
func tokenize(s string) []string {
	s = strings.Replace(s, "(", " ( ", -1)
	s = strings.Replace(s, ")", " ) ", -1)
	token := strings.Fields(s)
	return token
}

// Create abstract syntax tree.
func create_ast(tokens []string) (Expression, int, error) {
	if len(tokens) == 0 {
		err := NewSyntaxError("unexpected EOF while reading")
		return nil, 0, err
	}
	token := tokens[0]
	tokens = tokens[1:]
	if "(" == token {
		var L []Expression

		count := 1
		for {
			exp, c, _ := create_ast(tokens)
			L = append(L, exp)
			tokens = tokens[c:]
			count = count + c

			if len(tokens) == 0 {
				err := NewSyntaxError("unexpected ')' while reading")
				return nil, 0, err
			}
			if tokens[0] == ")" {
				count = count + 1
				break
			}
		}
		item := NewList(L)
		return item, count, nil

	} else if ")" == token {
		err := NewSyntaxError("unexpected )")
		return nil, 0, err
	} else {
		return atom(token), 1, nil
	}
}

// Atom To "Integer, Float, Symbol"
func atom(token string) Atom {
	var (
		atom Atom
	)
	ivalue, err := strconv.Atoi(token)
	if err == nil {
		atom = NewInteger(ivalue)
	} else {
		fvalue, err := strconv.ParseFloat(token, 64)
		if err == nil {
			atom = NewFloat(fvalue)
		} else {
			switch token {
			case "#t":
				atom = NewBoolean(true)
			case "#f":
				atom = NewBoolean(false)
			default:
				if (len(token) > 1) && (token[0] == '"') && (token[len(token)-1] == '"') {
					atom = NewString(token)
				} else {
					atom = NewSymbol(token)
				}
			}
		}
	}
	return atom
}

// Evaluate an expression in an environment.
func eval(sexp Expression, env *Environment) (Expression, error) {

	if _, ok := sexp.(Atom); ok {
		if sym, ok := sexp.(*Symbol); ok {
			if v, ok := (*env)[sym.Value]; ok {
				return v, nil
			} else if v, ok := define_env[sym.Value]; ok {
				return v, nil
			} else if v, ok := builtin_func[sym.Value]; ok {
				return NewOperator(v), nil
			} else {
				return sexp, NewRuntimeError("Undefine Operator or variable: " + sym.Value)
			}
		} else {
			// 10,11.. ,etc
			return sexp, nil
		}
	} else if sl, ok := sexp.(*List); ok {
		v := sl.Value

		if sym, ok := v[0].(*Symbol); ok {
			if kfn, ok := special_func[sym.Value]; ok {
				return kfn(env, v)
			}
			proc, err := eval(v[0], env)
			if err != nil {
				return sexp, err
			}
			if op, ok := proc.(*Operator); ok {
				var args []Expression
				for _, exp := range v[1:] {
					e, err := eval(exp, env)
					if err != nil {
						return sexp, err
					}
					// adhoc
					if lambda, ok := e.(*Function); ok {
						lambda.Env = env
					}
					args = append(args, e)
				}
				return op.Impl(args...)

			} else if fn, ok := proc.(*Function); ok {
				// (proc 10 20)
				let, err := fn.BindParam(env, v[1:])
				if err != nil {
					return sexp, err
				}
				return fn.Execute(let)

			} else if let, ok := proc.(*LetLoop); ok {
				// (let loop ((a (list 1 2))) (if (null? a) "ok" (loop (cdr a))))
				l, _ := let.ParamName.(*List)

				for i, c := range l.Value {
					pname := c.(*Symbol)
					(*env)[pname.Value], err = eval(v[i+1], env)
					if err != nil {
						return sexp, err
					}
				}
				return eval(let.Body, env)
			}
		} else if slf, ok := v[0].(*List); ok {
			// ((lambda (a b) (+ a b)) 10 20)
			e, err := eval(slf, env)
			if err != nil {
				return sexp, err
			}
			fn, ok := e.(*Function)
			if !ok {
				return sexp, NewRuntimeError("Not Function")
			}
			// name binding
			ef, err := fn.BindParam(env, v[1:])
			if err != nil {
				return nil, err
			}
			return fn.Execute(ef)
		}
	}
	return sexp, NewRuntimeError("Undefine Data Type")
}

// CUI desu.
func do_interactive() {
	prompt := "scheme.go> "
	reader := bufio.NewReaderSize(os.Stdin, MAX_LINE_SIZE)
	local_env := Environment{}
	for {
		fmt.Print(prompt + " ")

		var line string
		b, _, err := reader.ReadLine()
		line = string(b)
		if err == io.EOF {
			break
		} else if line == "" {
			continue
		} else if line == "(quit)" {
			break
		}

		ast, err := parse(line)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		val, err := eval(ast, &local_env)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if DEBUG {
			fmt.Print(reflect.TypeOf(val))
		}
		val.Print()
		fmt.Print("\n")
	}
}

// Build Global environement.
func build_env() {

	builtin_func = map[string]func(...Expression) (Expression, error){}
	special_func = map[string]func(*Environment, []Expression) (Expression, error){}
	define_env = Environment{}

	// addl, subl,imul,idiv,modulo
	iter_calc := func(calc func(Number, Number) Number, exp ...Expression) (Number, error) {
		if 1 >= len(exp) {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		if _, ok := exp[0].(Number); !ok {
			return nil, NewRuntimeError("Not Number")
		}
		result, err := CreateNumber(exp[0].(Number))
		if err != nil {
			return nil, err
		}
		for _, i := range exp[1:] {
			prm, ok := i.(Number)
			if !ok {
				return nil, NewRuntimeError("Not Number")
			}

			if _, ok := result.(*Float); ok {
				if c, ok := i.(*Integer); ok {
					prm = NewFloat(float64(c.Value))
				}
			}
			if org, ok := result.(*Integer); ok {
				if _, ok := i.(*Float); ok {
					result = NewFloat(float64(org.Value))
				}
			}
			result = calc(result, prm)
		}
		return result, nil
	}
	builtin_func["+"] = func(exp ...Expression) (Expression, error) {
		return iter_calc(func(a Number, b Number) Number { return a.Add(b) }, exp...)
	}
	builtin_func["-"] = func(exp ...Expression) (Expression, error) {
		return iter_calc(func(a Number, b Number) Number { return a.Sub(b) }, exp...)
	}
	builtin_func["*"] = func(exp ...Expression) (Expression, error) {
		return iter_calc(func(a Number, b Number) Number { return a.Mul(b) }, exp...)
	}
	builtin_func["/"] = func(exp ...Expression) (Expression, error) {
		return iter_calc(func(a Number, b Number) Number { return a.Div(b) }, exp...)
	}
	builtin_func["modulo"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		var prm []*Integer
		for _, e := range exp {
			v, ok := e.(*Integer)
			if !ok {
				return nil, NewRuntimeError("Not Integer")
			}
			prm = append(prm, v)
		}
		return NewInteger(prm[0].Value % prm[1].Value), nil
	}
	// gt,lt,ge,le
	op_cmp := func(cmp func(Number, Number) bool, exp ...Expression) (*Boolean, error) {
		if 2 != len(exp) {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}

		result, err := CreateNumber(exp[0].(Number))
		if err != nil {
			return nil, err
		}
		prm, ok := exp[1].(Number)
		if !ok {
			return nil, NewRuntimeError("Not Integer")
		}
		if _, ok := result.(*Float); ok {
			if c, ok := prm.(*Integer); ok {
				prm = NewFloat(float64(c.Value))
			}
		}
		if org, ok := result.(*Integer); ok {
			if c, ok := exp[1].(*Float); ok {
				result = NewFloat(float64(org.Value))
				prm = c
			}
		}
		return NewBoolean(cmp(result, prm)), nil
	}
	builtin_func[">"] = func(exp ...Expression) (Expression, error) {
		return op_cmp(func(a Number, b Number) bool { return a.Greater(b) }, exp...)
	}
	builtin_func["<"] = func(exp ...Expression) (Expression, error) {
		return op_cmp(func(a Number, b Number) bool { return a.Less(b) }, exp...)
	}
	builtin_func[">="] = func(exp ...Expression) (Expression, error) {
		return op_cmp(func(a Number, b Number) bool { return a.GreaterEqual(b) }, exp...)
	}
	builtin_func["<="] = func(exp ...Expression) (Expression, error) {
		return op_cmp(func(a Number, b Number) bool { return a.LessEqual(b) }, exp...)
	}
	builtin_func["="] = func(exp ...Expression) (Expression, error) {
		return op_cmp(func(a Number, b Number) bool { return a.Equal(b) }, exp...)
	}
	builtin_func["not"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		if _, ok := exp[0].(*Boolean); !ok {
			return nil, NewRuntimeError("Not Boolean")
		}
		return NewBoolean(!(exp[0].(*Boolean)).Value), nil
	}

	// list operator
	builtin_func["list"] = func(exp ...Expression) (Expression, error) {
		var l []Expression
		return NewList(append(l, exp...)), nil
	}
	builtin_func["null?"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		if l, ok := exp[0].(*List); ok {
			return NewBoolean(0 == len(l.Value)), nil
		} else {
			return nil, NewRuntimeError("Not List")
		}
	}
	builtin_func["length"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		if l, ok := exp[0].(*List); ok {
			return NewInteger(len(l.Value)), nil
		} else {
			return nil, NewRuntimeError("Not List")
		}
	}
	builtin_func["car"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		if l, ok := exp[0].(*List); ok {
			if len(l.Value) <= 0 {
				return nil, NewRuntimeError("Not Enough Parameter Number")
			}
			return l.Value[0], nil
		} else if p, ok := exp[0].(*Pair); ok {
			return p.Car, nil
		} else {
			return nil, NewRuntimeError("Not List")
		}
	}
	builtin_func["cdr"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		if l, ok := exp[0].(*List); ok {
			if len(l.Value) <= 0 {
				var v []Expression
				return NewList(v), nil
			}
			return NewList(l.Value[1:]), nil
		} else if p, ok := exp[0].(*Pair); ok {
			return p.Cdr, nil
		} else {
			return nil, NewRuntimeError("Not List")
		}
	}
	builtin_func["cons"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		if _, ok := exp[1].(*List); ok {
			var args []Expression
			args = append(args, exp[0])
			return NewList(append(args, (exp[1].(*List)).Value...)), nil
		}
		return NewPair(exp[0], exp[1]), nil
	}
	builtin_func["append"] = func(exp ...Expression) (Expression, error) {
		if len(exp) < 2 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		var append_list []Expression
		for _, e := range exp {
			if v, ok := e.(*List); ok {
				append_list = append(append_list, v.Value...)
			} else {
				return nil, NewRuntimeError("Not List")
			}
		}
		return NewList(append_list), nil
	}
	builtin_func["last"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		if l, ok := exp[0].(*List); ok {
			if len(l.Value) <= 0 {
				return nil, NewRuntimeError("Not Enough Parameter Length")
			}
			return l.Value[len(l.Value)-1], nil
		} else if p, ok := exp[0].(*Pair); ok {
			return p.Car, nil
		} else {
			return nil, NewRuntimeError("Not List")
		}
	}
	builtin_func["iota"] = func(exp ...Expression) (Expression, error) {
		if len(exp) < 1 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		var l []Expression
		max, ok := exp[0].(*Integer)
		if !ok {
			return nil, NewRuntimeError("Not Integer")
		}
		start := 0
		if len(exp) == 2 {
			v, ok := exp[1].(*Integer)
			if !ok {
				return nil, NewRuntimeError("Not Integer")
			}
			start = v.Value
		}
		for i := start; i < start+max.Value; i++ {
			l = append(l, NewInteger(i))
		}

		return NewList(l), nil
	}
	// map,filter,reduce
	iter_func := func(lambda func(Expression, Expression, []Expression) ([]Expression, error), exp ...Expression) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		fn, ok := exp[0].(*Function)
		if !ok {
			return nil, NewRuntimeError("Not Function")
		}
		l, ok := exp[1].(*List)
		if !ok {
			return nil, NewRuntimeError("Not List")
		}
		var va_list []Expression
		param := make([]Expression, 1)
		for _, param[0] = range l.Value {
			let, _ := fn.BindParam(nil, param)
			result, err := fn.Execute(let)
			if err != nil {
				return nil, err
			}
			va_list, err = lambda(result, param[0], va_list)
			if err != nil {
				return nil, err
			}
		}
		return NewList(va_list), nil
	}
	builtin_func["map"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		lambda := func(result Expression, item Expression, va_list []Expression) ([]Expression, error) {
			return append(va_list, result), nil
		}
		return iter_func(lambda, exp...)
	}
	builtin_func["filter"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		lambda := func(result Expression, item Expression, va_list []Expression) ([]Expression, error) {
			b, ok := result.(*Boolean)
			if !ok {
				return nil, NewRuntimeError("Not Boolean")
			}
			if b.Value {
				return append(va_list, item), nil
			}
			return va_list, nil
		}
		return iter_func(lambda, exp...)
	}
	builtin_func["reduce"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		fn, ok := exp[0].(*Function)
		if !ok {
			return nil, NewRuntimeError("Not Function")
		}
		l, ok := exp[1].(*List)
		if !ok {
			return nil, NewRuntimeError("Not List")
		}

		param := make([]Expression, len((fn.ParamName.(*List)).Value))
		result := l.Value[0]
		var err error
		for _, c := range l.Value[1:] {
			param[0] = result
			param[1] = c
			let, _ := fn.BindParam(nil, param)
			result, err = fn.Execute(let)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	}
	// math skelton
	math_impl := func(math_func func(float64) float64, exp ...Expression) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		if v, ok := exp[0].(*Float); ok {
			return NewFloat(math_func(v.Value)), nil
		} else if v, ok := exp[0].(*Integer); ok {
			return NewFloat(math_func((float64)(v.Value))), nil
		}
		return nil, NewRuntimeError("Not Enough Type")
	}
	builtin_func["sqrt"] = func(exp ...Expression) (Expression, error) {
		return math_impl(math.Sqrt, exp...)
	}
	builtin_func["sin"] = func(exp ...Expression) (Expression, error) {
		return math_impl(math.Sin, exp...)
	}
	builtin_func["cos"] = func(exp ...Expression) (Expression, error) {
		return math_impl(math.Cos, exp...)
	}
	builtin_func["tan"] = func(exp ...Expression) (Expression, error) {
		return math_impl(math.Tan, exp...)
	}
	builtin_func["atan"] = func(exp ...Expression) (Expression, error) {
		return math_impl(math.Atan, exp...)
	}
	builtin_func["log"] = func(exp ...Expression) (Expression, error) {
		return math_impl(math.Log, exp...)
	}
	builtin_func["exp"] = func(exp ...Expression) (Expression, error) {
		return math_impl(math.Exp, exp...)
	}
	builtin_func["rand-integer"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		if v, ok := exp[0].(*Integer); ok {
			return NewInteger(rand.Intn(v.Value)), nil
		}
		return nil, NewRuntimeError("Not Integer")
	}
	// syntax keyword implements
	special_func["if"] = func(env *Environment, v []Expression) (Expression, error) {
		if len(v) != 4 {
			return nil, NewRuntimeError("Not Enough Parameter")
		}

		exp, err := eval(v[1], env)
		if err != nil {
			return nil, err
		}

		b, ok := exp.(*Boolean)
		if !ok {
			return nil, NewRuntimeError("Not Boolean")
		}

		if b.Value {
			return eval(v[2], env)
		} else {
			return eval(v[3], env)
		}
	}
	special_func["define"] = func(env *Environment, v []Expression) (Expression, error) {
		if len(v) != 3 {
			return nil, NewRuntimeError("Not Enough Parameter")
		}
		key, ok := v[1].(*Symbol)
		if !ok {
			return nil, NewRuntimeError("Not Symbol")
		}
		exp, err := eval(v[2], env)
		if err != nil {
			return nil, err
		}
		define_env[key.Value] = exp
		return key, nil
	}
	special_func["lambda"] = func(env *Environment, v []Expression) (Expression, error) {
		if len(v) < 3 {
			return nil, NewRuntimeError("Not Enough Parameter")
		}
		return NewFunction(v[1], v[2:]), nil
	}
	special_func["set!"] = func(env *Environment, v []Expression) (Expression, error) {
		if len(v) != 3 {
			return nil, NewRuntimeError("Not Enough Parameter")
		}
		s, ok := v[1].(*Symbol)
		if !ok {
			return nil, NewRuntimeError("Not Symbol")
		}

		exp, err := eval(v[2], env)
		if err != nil {
			return v[2], err
		}

		if _, ok := (*env)[s.Value]; ok {
			(*env)[s.Value] = exp
			return (*env)[s.Value], nil

		} else if _, ok := define_env[s.Value]; ok {
			define_env[s.Value] = exp
			return define_env[s.Value], nil
		} else {
			return exp, NewRuntimeError("Undefined Variable: " + s.Value)
		}
	}
	special_func["let"] = func(env *Environment, v []Expression) (Expression, error) {
		var letsym *Symbol
		var pname []Expression
		body := 2

		l, ok := v[1].(*List)
		if ok && len(v) <= 2 {
			return nil, NewRuntimeError("Not Enough Parameter")
		}

		if !ok {
			letsym, ok = v[1].(*Symbol)
			if !ok {
				return nil, NewRuntimeError("Not Symbol")
			}
			if len(v) <= 3 {
				return nil, NewRuntimeError("Not Enough Parameter")
			}
			l, ok = v[2].(*List)
			if !ok {
				return nil, NewRuntimeError("Not List")
			}
			body = 3
		}
		for _, let := range l.Value {
			r, ok := let.(*List)
			if !ok {
				return nil, NewRuntimeError("Not List")
			}
			if len(r.Value) != 2 {
				return nil, NewRuntimeError("Not Enough Parameter")
			}
			sym := (r.Value[0]).(*Symbol)
			v, err := eval(r.Value[1], env)
			if err != nil {
				return nil, err
			}
			pname = append(pname, sym)
			(*env)[sym.Value] = v
		}
		if letsym != nil {
			(*env)[letsym.Value] = NewLetLoop(NewList(pname), v[body])
		}
		return eval(v[body], env)
	}
	// and or not
	op_logical := func(env *Environment, exp []Expression, bcond bool, bret bool) (Expression, error) {
		if len(exp) <= 1 {
			return nil, NewRuntimeError("Not Enough Parameter Number")
		}
		for _, e := range exp {
			b, err := eval(e, env)
			if err != nil {
				return nil, err
			}
			if _, ok := b.(*Boolean); !ok {
				return nil, NewRuntimeError("Not Boolean")
			}
			if bcond == (b.(*Boolean)).Value {
				return NewBoolean(bcond), nil
			}
		}
		return NewBoolean(bret), nil
	}
	special_func["and"] = func(env *Environment, exp []Expression) (Expression, error) {
		return op_logical(env, exp[1:], false, true)
	}
	special_func["or"] = func(env *Environment, exp []Expression) (Expression, error) {
		return op_logical(env, exp[1:], true, false)
	}

}

// Main
func main() {
	build_env()
	do_interactive()
}
