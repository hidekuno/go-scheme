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
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type (
	Environment map[string]Expression
)

const (
	MAX_LINE_SIZE = 1024
	DEBUG         = false
	PROMPT        = "scheme.go> "
)

var (
	builtin_func map[string]func(...Expression) (Expression, error)
	special_func map[string]func(*SimpleEnv, []Expression) (Expression, error)
	error_msg    = map[string]string{
		"E0001": "Unexpected EOF while reading",
		"E0002": "Unexpected ')' while reading",
		"E0003": "Extra close parenthesis `)'",
		"E1001": "Not Boolean",
		"E1002": "Not Integer",
		"E1003": "Not Number",
		"E1004": "Not Symbol",
		"E1005": "Not List",
		"E1006": "Not Function",
		"E1007": "Not Enough Parameter Counts",
		"E1008": "Undefine variable",
		"E1009": "Not Enough Data Type",
	}
	tracer = log.New(os.Stderr, "", log.Lshortfile)
)

// env structure
type SimpleEnv struct {
	EnvTable *Environment
	Parent   *SimpleEnv
}

func NewSimpleEnv(parent *SimpleEnv, et *Environment) *SimpleEnv {
	v := new(SimpleEnv)
	v.Parent = parent
	if et != nil {
		v.EnvTable = et
	} else {
		env := Environment{}
		v.EnvTable = &env
	}
	return v
}
func (self *SimpleEnv) Find(key string) (Expression, bool) {
	if v, ok := (*self.EnvTable)[key]; ok {
		return v, true
	}
	for c := self.Parent; c != nil; c = c.Parent {
		if v, ok := (*c.EnvTable)[key]; ok {
			return v, true
		}
	}
	return nil, false
}
func (self *SimpleEnv) Regist(key string, exp Expression) {
	(*self.EnvTable)[key] = exp
}
func (self *SimpleEnv) Set(key string, exp Expression) {
	if _, ok := (*self.EnvTable)[key]; ok {
		(*self.EnvTable)[key] = exp
		return
	}
	for c := self.Parent; c != nil; c = c.Parent {
		if _, ok := (*c.EnvTable)[key]; ok {
			(*c.EnvTable)[key] = exp
			return
		}
	}
}

// Basic Data Type. (need
type SyntaxError struct {
	MsgCode           string
	SourceFileName    string
	SourceFileLineNum int
}

func (err *SyntaxError) Error() string {
	return error_msg[err.MsgCode] + " (" + path.Base(err.SourceFileName) + ":" + strconv.Itoa(err.SourceFileLineNum) + ")"
}

func NewSyntaxError(text string) error {
	_, sourceFileName, sourceFileLineNum, _ := runtime.Caller(1)
	return &SyntaxError{text, sourceFileName, sourceFileLineNum}
}

type RuntimeError struct {
	MsgCode           string
	SourceFileName    string
	SourceFileLineNum int
	Args              []string
}

func (err *RuntimeError) Error() string {

	args := ""
	for i, e := range err.Args {
		if i != 0 {
			args = args + ","
		}
		args = args + e
	}
	return error_msg[err.MsgCode] + ": " + args + " (" + path.Base(err.SourceFileName) + ":" + strconv.Itoa(err.SourceFileLineNum) + ")"
}

func NewRuntimeError(text string, args ...string) error {
	_, sourceFileName, sourceFileLineNum, _ := runtime.Caller(1)
	return &RuntimeError{text, sourceFileName, sourceFileLineNum, args}
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

func CreateNumber(p Expression) (Number, error) {
	if v, ok := p.(*Integer); ok {
		return NewInteger(v.Value), nil
	}
	if v, ok := p.(*Float); ok {
		return NewFloat(v.Value), nil
	}
	return nil, NewRuntimeError("E1003", reflect.TypeOf(p).String())
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

type SpecialFunc struct {
	Expression
	Impl func(*SimpleEnv, []Expression) (Expression, error)
}

func NewSpecialFunc(fn func(*SimpleEnv, []Expression) (Expression, error)) *SpecialFunc {
	sf := new(SpecialFunc)
	sf.Impl = fn
	return sf
}
func (self *SpecialFunc) Print() {
	fmt.Print("Special Functon ex. if: ", self)
}
func (self *SpecialFunc) Execute(env *SimpleEnv, exps []Expression) (Expression, error) {
	return self.Impl(env, exps)
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

func (self *Operator) Execute(env *SimpleEnv, exps []Expression) (Expression, error) {
	var args []Expression

	for _, exp := range exps {
		e, err := eval(exp, env)
		if err != nil {
			return exp, err
		}
		args = append(args, e)
	}
	return self.Impl(args...)
}

type Function struct {
	Expression
	ParamName List
	Body      []Expression
	Env       *SimpleEnv
}

func NewFunction(parent *SimpleEnv, param *List, body []Expression) *Function {
	fn := new(Function)
	fn.ParamName = *param
	fn.Body = body
	fn.Env = NewSimpleEnv(parent, nil)
	return fn
}

func (self *Function) Print() {
	fmt.Print("Function: ", self)
}

// Bind lambda function' parameters.
func (self *Function) Execute(env *SimpleEnv, values []Expression) (Expression, error) {
	local_env := Environment{}
	idx := 0
	if len(self.ParamName.Value) != len(values) {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(values)))
	}
	for _, i := range self.ParamName.Value {
		if sym, ok := i.(*Symbol); ok {
			if env != nil {
				v, err := eval(values[idx], env)
				if err != nil {
					return nil, err
				}
				local_env[sym.Value] = v
			} else {
				local_env[sym.Value] = values[idx]
			}
			idx = idx + 1
		}
	}
	self.Env = NewSimpleEnv(self.Env, &local_env)
	var (
		result Expression
		err    error
	)
	for _, exp := range self.Body {
		result, err = eval(exp, self.Env)
		if err != nil {
			return exp, err
		}
	}
	return result, nil
}

type LetLoop struct {
	Expression
	ParamName List
	Body      Expression
}

func NewLetLoop(param *List, body Expression) *LetLoop {
	let := new(LetLoop)
	let.ParamName = *param
	let.Body = body
	return let
}

func (self *LetLoop) Print() {
	fmt.Print("Let Macro: ", self)
}
func (self *LetLoop) Execute(env *SimpleEnv, v []Expression) (Expression, error) {

	for i, c := range self.ParamName.Value {
		pname := c.(*Symbol)
		data, err := eval(v[i], env)
		if err != nil {
			return nil, err
		}
		(*env).Set(pname.Value, data)

	}
	return eval(self.Body, env)
}

// lex support  for  string
func tokenize(s string) []string {
	var token []string

	string_mode := false
	string_data := ""

	rep := regexp.MustCompile(`^ *`)
	s = rep.ReplaceAllString(s, "")
	rep = regexp.MustCompile(` *$`)
	s = rep.ReplaceAllString(s, "")

	from := 0
	for i, c := range s {
		if string_mode {
			if c == '"' {
				token = append(token, s[from:i+1])
				string_mode = false
			}
		} else {
			if c == '"' {
				from = i
				string_mode = true
			} else if c == '(' {
				token = append(token, "(")
			} else if c == ')' {
				token = append(token, ")")
			} else if c == ' ' {
				// Nop
			} else {
				string_data = string_data + s[i:i+1]
				if (len(s)-1 == i) || (s[i+1] == ')') || (s[i+1] == ' ') {
					token = append(token, string_data)
					string_data = ""
				}
			}
		}
	}
	if DEBUG {
		for _, c := range token {
			fmt.Println(c)
		}
	}
	return token
}

// Tenuki lex.
func tokenize_easy(s string) []string {
	s = strings.Replace(s, "(", " ( ", -1)
	s = strings.Replace(s, ")", " ) ", -1)
	token := strings.Fields(s)
	return token
}

// Create abstract syntax tree.
func parse(tokens []string) (Expression, int, error) {

	if len(tokens) == 0 {
		return nil, 0, NewSyntaxError("E0001")
	}
	token := tokens[0]
	tokens = tokens[1:]

	if "(" == token {
		if len(tokens) <= 0 {
			return nil, 0, NewSyntaxError("E0001")
		}
		var L []Expression

		count := 1
		for {
			if tokens[0] == ")" {
				count = count + 1
				break
			}
			exp, c, err := parse(tokens)
			if err != nil {
				return nil, c, err
			}
			L = append(L, exp)
			tokens = tokens[c:]
			count = count + c

			if len(tokens) == 0 {
				return nil, 0, NewSyntaxError("E0002")
			}
		}
		item := NewList(L)
		return item, count, nil

	} else if ")" == token {
		return nil, 0, NewSyntaxError("E0002")
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
func eval(sexp Expression, env *SimpleEnv) (Expression, error) {
	if DEBUG {
		fmt.Print(reflect.TypeOf(sexp))
	}
	if _, ok := sexp.(Atom); ok {
		if sym, ok := sexp.(*Symbol); ok {
			if v, ok := (*env).Find(sym.Value); ok {
				return v, nil
			} else if v, ok := builtin_func[sym.Value]; ok {
				return NewOperator(v), nil
			} else if v, ok := special_func[sym.Value]; ok {
				return NewSpecialFunc(v), nil
			} else {
				return sexp, NewRuntimeError("E1008", sym.Value)
			}
		} else {
			// 10,11.. ,"a", "B", ,etc
			return sexp, nil
		}
	} else if sl, ok := sexp.(*List); ok {
		if len(sl.Value) == 0 {
			return sexp, nil
		}
		v := sl.Value
		if _, ok := v[0].(*Symbol); ok {
			proc, err := eval(v[0], env)
			if err != nil {
				return sexp, err
			}
			if sf, ok := proc.(*SpecialFunc); ok {
				// (if (= a b) "a" "b")
				return sf.Execute(env, v[1:])

			} else if op, ok := proc.(*Operator); ok {
				// (* (+ a 1) (+ b 2))
				return op.Execute(env, v[1:])

			} else if fn, ok := proc.(*Function); ok {
				// (proc 10 20)
				return fn.Execute(env, v[1:])

			} else if let, ok := proc.(*LetLoop); ok {
				// (let loop ((a (list 1 2 3))(b 0))
				//   (if (null? a) b (loop (cdr a)(+ b (car a)))))
				return let.Execute(env, v[1:])
			}
		} else if slf, ok := v[0].(*List); ok {
			// ((lambda (a b) (+ a b)) 10 20)
			e, err := eval(slf, env)
			if err != nil {
				return sexp, err
			}
			fn, ok := e.(*Function)
			if !ok {
				return sexp, NewRuntimeError("E1006", reflect.TypeOf(e).String())
			}
			// execute
			return fn.Execute(env, v[1:])
		}
	}
	return sexp, NewRuntimeError("E1009", reflect.TypeOf(sexp).String())
}

// main logic
func do_core_logic(program string, root_env *SimpleEnv) (Expression, error) {

	token := tokenize_easy(program)
	ast, c, err := parse(token)
	if err != nil {
		return nil, err
	}

	if c != len(token) {
		return nil, NewSyntaxError("E0003")
	}

	val, err := eval(ast, root_env)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// CUI desu.
func count_parenthesis(program string) bool {
	left_parenthesis := 0
	right_parenthesis := 0
	search_mode := true

	for _, c := range program {
		if c == '"' && search_mode {
			search_mode = false
		} else if c == '"' && !search_mode {
			search_mode = true
		}
		if c == '(' && search_mode {
			left_parenthesis = left_parenthesis + 1
		}
		if c == ')' && search_mode {
			right_parenthesis = right_parenthesis + 1
		}
	}
	return left_parenthesis <= right_parenthesis
}

// CUI desu.
func do_interactive() {
	var program string

	prompt := PROMPT
	reader := bufio.NewReaderSize(os.Stdin, MAX_LINE_SIZE)
	root_env := NewSimpleEnv(nil, nil)
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

		program = program + " " + line
		if !count_parenthesis(program) {
			prompt = ""
			continue
		}
		val, err := do_core_logic(program, root_env)
		if err != nil {
			fmt.Println(err.Error())
			goto FINISH
		}
		val.Print()
		fmt.Print("\n")
		if DEBUG {
			fmt.Print(reflect.TypeOf(val))
		}
	FINISH:
		program = ""
		prompt = PROMPT
	}
}

// Build Global environement.
func build_func() {
	builtin_func = map[string]func(...Expression) (Expression, error){}
	special_func = map[string]func(*SimpleEnv, []Expression) (Expression, error){}

	// addl, subl,imul,idiv,modulo
	iter_calc := func(calc func(Number, Number) Number, exp ...Expression) (Number, error) {
		if 1 >= len(exp) {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		result, err := CreateNumber(exp[0])
		if err != nil {
			return nil, err
		}
		for _, i := range exp[1:] {
			prm, ok := i.(Number)
			if !ok {
				return nil, NewRuntimeError("E1003", reflect.TypeOf(i).String())
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
	builtin_func["quotient"] = builtin_func["/"]
	builtin_func["modulo"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		var prm []*Integer
		for _, e := range exp {
			v, ok := e.(*Integer)
			if !ok {
				return nil, NewRuntimeError("E1002", reflect.TypeOf(e).String())
			}
			prm = append(prm, v)
		}
		return NewInteger(prm[0].Value % prm[1].Value), nil
	}
	// gt,lt,ge,le
	op_cmp := func(cmp func(Number, Number) bool, exp ...Expression) (*Boolean, error) {
		if 2 != len(exp) {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}

		result, err := CreateNumber(exp[0])
		if err != nil {
			return nil, err
		}
		prm, ok := exp[1].(Number)
		if !ok {
			return nil, NewRuntimeError("E1003", reflect.TypeOf(exp[1]).String())
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
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		if _, ok := exp[0].(*Boolean); !ok {
			return nil, NewRuntimeError("E1001", reflect.TypeOf(exp[0]).String())
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
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		if l, ok := exp[0].(*List); ok {
			return NewBoolean(0 == len(l.Value)), nil
		} else {
			return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
		}
	}
	builtin_func["length"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		if l, ok := exp[0].(*List); ok {
			return NewInteger(len(l.Value)), nil
		} else {
			return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
		}
	}
	builtin_func["car"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		if l, ok := exp[0].(*List); ok {
			if len(l.Value) <= 0 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(l.Value)))
			}
			return l.Value[0], nil
		} else if p, ok := exp[0].(*Pair); ok {
			return p.Car, nil
		} else {
			return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
		}
	}
	builtin_func["cdr"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
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
			return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
		}
	}
	builtin_func["cons"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
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
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		var append_list []Expression
		for _, e := range exp {
			if v, ok := e.(*List); ok {
				append_list = append(append_list, v.Value...)
			} else {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(e).String())
			}
		}
		return NewList(append_list), nil
	}
	builtin_func["last"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		if l, ok := exp[0].(*List); ok {
			if len(l.Value) <= 0 {
				return nil, NewRuntimeError("E1007")
			}
			return l.Value[len(l.Value)-1], nil
		} else if p, ok := exp[0].(*Pair); ok {
			return p.Car, nil
		} else {
			return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
		}
	}
	builtin_func["reverse"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		if l, ok := exp[0].(*List); ok {
			if len(l.Value) <= 1 {
				return l, nil
			}
			args := make([]Expression, len(l.Value))
			idx := len(l.Value) - 1
			for _, c := range l.Value {
				args[idx] = c
				idx = idx - 1
			}
			return NewList(args), nil
		} else {
			return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
		}
	}
	builtin_func["iota"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 1 && len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		var l []Expression
		max, ok := exp[0].(*Integer)
		if !ok {
			return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[0]).String())
		}
		start := 0
		if len(exp) == 2 {
			v, ok := exp[1].(*Integer)
			if !ok {
				return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[1]).String())
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
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		fn, ok := exp[0].(*Function)
		if !ok {
			return nil, NewRuntimeError("E1006", reflect.TypeOf(exp[0]).String())
		}
		l, ok := exp[1].(*List)
		if !ok {
			return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[1]).String())
		}
		var va_list []Expression
		param := make([]Expression, 1)
		for _, param[0] = range l.Value {
			result, err := fn.Execute(nil, param)

			va_list, err = lambda(result, param[0], va_list)
			if err != nil {
				return nil, err
			}
		}
		return NewList(va_list), nil
	}
	builtin_func["map"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		lambda := func(result Expression, item Expression, va_list []Expression) ([]Expression, error) {
			return append(va_list, result), nil
		}
		return iter_func(lambda, exp...)
	}
	builtin_func["filter"] = func(exp ...Expression) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		lambda := func(result Expression, item Expression, va_list []Expression) ([]Expression, error) {
			b, ok := result.(*Boolean)
			if !ok {
				return nil, NewRuntimeError("E1001", reflect.TypeOf(result).String())
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
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		fn, ok := exp[0].(*Function)
		if !ok {
			return nil, NewRuntimeError("E1006", reflect.TypeOf(exp[0]).String())
		}
		l, ok := exp[1].(*List)
		if !ok {
			return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[1]).String())
		}

		param := make([]Expression, len(fn.ParamName.Value))
		result := l.Value[0]

		for _, c := range l.Value[1:] {
			param[0] = result
			param[1] = c
			r, err := fn.Execute(nil, param)
			result = r
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	}
	// math skelton
	math_impl := func(math_func func(float64) float64, exp ...Expression) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		if v, ok := exp[0].(*Float); ok {
			return NewFloat(math_func(v.Value)), nil
		} else if v, ok := exp[0].(*Integer); ok {
			return NewFloat(math_func((float64)(v.Value))), nil
		}
		return nil, NewRuntimeError("E1003", reflect.TypeOf(exp[0]).String())
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
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		if v, ok := exp[0].(*Integer); ok {
			return NewInteger(rand.Intn(v.Value)), nil
		}
		return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[0]).String())
	}
	// syntax keyword implements
	special_func["if"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
		if len(v) != 3 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(v)))
		}
		exp, err := eval(v[0], env)
		if err != nil {
			return nil, err
		}

		b, ok := exp.(*Boolean)
		if !ok {
			return nil, NewRuntimeError("E1001", reflect.TypeOf(exp).String())
		}
		if b.Value {
			return eval(v[1], env)
		} else {
			return eval(v[2], env)
		}
	}
	special_func["define"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
		if len(v) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(v)))
		}
		key, ok := v[0].(*Symbol)
		if !ok {
			return nil, NewRuntimeError("E1004", reflect.TypeOf(v[0]).String())
		}
		exp, err := eval(v[1], env)
		if err != nil {
			return nil, err
		}
		(*env).Regist(key.Value, exp)
		return key, nil
	}
	special_func["lambda"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
		if len(v) < 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(v)))
		}
		// if l == (), internal list implements
		l, ok := v[0].(*List)
		if !ok {
			return nil, NewRuntimeError("E1005", reflect.TypeOf(v[0]).String())
		}
		return NewFunction(env, l, v[1:]), nil
	}
	special_func["set!"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
		if len(v) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(v)))
		}
		s, ok := v[0].(*Symbol)
		if !ok {
			return nil, NewRuntimeError("E1004", reflect.TypeOf(v[0]).String())
		}

		exp, err := eval(v[1], env)
		if err != nil {
			return v[1], err
		}
		if _, ok := (*env).Find(s.Value); !ok {
			return exp, NewRuntimeError("E1008", s.Value)
		}
		(*env).Set(s.Value, exp)
		return exp, nil
	}
	special_func["let"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
		var letsym *Symbol
		var pname []Expression
		body := 1

		l, ok := v[0].(*List)
		if ok && len(v) < 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(v)))
		}
		if !ok {
			letsym, ok = v[0].(*Symbol)
			if !ok {
				return nil, NewRuntimeError("E1004", reflect.TypeOf(v[0]).String())
			}
			if len(v) < 3 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(v)))
			}
			l, ok = v[1].(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(v[1]).String())
			}
			body = 2
		}

		local_env := Environment{}
		for _, let := range l.Value {
			r, ok := let.(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(let).String())
			}
			if len(r.Value) != 2 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(r.Value)))
			}
			sym := (r.Value[0]).(*Symbol)
			v, err := eval(r.Value[1], env)
			if err != nil {
				return nil, err
			}
			pname = append(pname, sym)
			local_env[sym.Value] = v
		}
		if letsym != nil {
			(*env).Regist(letsym.Value, NewLetLoop(NewList(pname), v[body]))
		}
		return eval(v[body], NewSimpleEnv(env, &local_env))
	}
	// and or not
	op_logical := func(env *SimpleEnv, exp []Expression, bcond bool, bret bool) (Expression, error) {
		if len(exp) < 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		for _, e := range exp {
			b, err := eval(e, env)
			if err != nil {
				return nil, err
			}
			if _, ok := b.(*Boolean); !ok {
				return nil, NewRuntimeError("E1001", reflect.TypeOf(b).String())
			}
			if bcond == (b.(*Boolean)).Value {
				return NewBoolean(bcond), nil
			}
		}
		return NewBoolean(bret), nil
	}
	special_func["and"] = func(env *SimpleEnv, exp []Expression) (Expression, error) {
		return op_logical(env, exp, false, true)
	}
	special_func["or"] = func(env *SimpleEnv, exp []Expression) (Expression, error) {
		return op_logical(env, exp, true, false)
	}
}

// Main
func main() {
	build_func()
	do_interactive()
}
