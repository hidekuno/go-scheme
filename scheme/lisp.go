/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

import (
	"bufio"
	"bytes"
	crand "crypto/rand"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type (
	Environment        map[string]Expression
	EvalFunc           func([]Expression, *SimpleEnv) (Expression, error)
	EvaledAllParamFunc func(exp ...Expression) (Expression, error)
)

const (
	MaxLineSize = 1024
	DEBUG       = false
	PROMPT      = "scheme.go> "
)

var (
	//	builtinFuncTbl map[string]func(...Expression) (Expression, error)
	buildInFuncTbl map[string]EvalFunc
	errorMsg       = map[string]string{
		"E0001": "Unexpected EOF while reading",
		"E0002": "Unexpected ')' while reading",
		"E0003": "Extra close parenthesis `)'",
		"E0004": "Charactor syntax error",
		"E1001": "Not Boolean",
		"E1002": "Not Integer",
		"E1003": "Not Number",
		"E1004": "Not Symbol",
		"E1005": "Not List",
		"E1006": "Not Function",
		"E1007": "Not Enough Parameter Counts",
		"E1008": "Undefine variable",
		"E1009": "Not Enough Data Type",
		"E1010": "Not Promise",
		"E1011": "Not Enough List Length",
		"E1012": "Not Cond Gramar",
		"E1013": "Calculate A Division By Zero",
		"E1014": "Not Found Program File",
		"E1015": "Not String",
		"E1016": "Not Program File",
		"E9999": "System Panic",
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

// Error(syntax, runtime)
type SyntaxError struct {
	MsgCode           string
	SourceFileName    string
	SourceFileLineNum int
}

func (err *SyntaxError) Error() string {
	return errorMsg[err.MsgCode] + " (" + path.Base(err.SourceFileName) + ":" + strconv.Itoa(err.SourceFileLineNum) + ")"
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
	return errorMsg[err.MsgCode] + ": " + args + " (" + path.Base(err.SourceFileName) + ":" + strconv.Itoa(err.SourceFileLineNum) + ")"
}

func NewRuntimeError(text string, args ...string) error {
	_, sourceFileName, sourceFileLineNum, _ := runtime.Caller(1)
	return &RuntimeError{text, sourceFileName, sourceFileLineNum, args}
}

// eval type
type Expression interface {
	String() string
}

type Any interface{}
type Atom interface {
	Expression
	// Because Expression is different
	Dummy() Any
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

func CreateNumber(exp Expression) (Number, error) {
	if v, ok := exp.(*Integer); ok {
		return NewInteger(v.Value), nil
	}
	if v, ok := exp.(*Float); ok {
		return NewFloat(v.Value), nil
	}
	return nil, NewRuntimeError("E1003", reflect.TypeOf(exp).String())
}

// Integer Type
type Integer struct {
	Number
	Value int
}

func NewInteger(p int) *Integer {
	v := new(Integer)
	v.Value = p
	return v
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
	if v.Value == 0 {
		panic(NewRuntimeError("E1013"))
	}
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

func (self *Integer) String() string {
	return strconv.Itoa(self.Value)
}

// Float Type
type Float struct {
	Number
	Value float64
}

func NewFloat(p float64) *Float {
	v := new(Float)
	v.Value = p
	return v
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
func (self *Float) String() string {
	return fmt.Sprint(self.Value)
}
func (self *Float) FormatString(prec int) string {
	return strconv.FormatFloat(self.Value, 'f', prec, 64)
}
func (self *Float) LogFormatString(prec int) string {
	return strconv.FormatFloat(self.Value, 'e', prec, 64)
}

// Symbol Type
type Symbol struct {
	Atom
	Value string
}

func NewSymbol(token string) *Symbol {
	s := new(Symbol)
	s.Value = token
	return s
}

func (self *Symbol) String() string {
	return self.Value
}

// Boolean Type
type Boolean struct {
	Atom
	Value bool
	exp   string
}

func NewBoolean(v bool) *Boolean {
	b := new(Boolean)
	b.Value = v
	if v {
		b.exp = "#t"
	} else {
		b.exp = "#f"
	}
	return b
}
func (self *Boolean) String() string {
	return self.exp
}

// Character Type
type Char struct {
	Atom
	Value rune
	exp   string
}

func NewChar(v string) *Char {
	b := new(Char)
	b.exp = v
	b.Value = []rune(v)[2]
	return b
}
func (self *Char) String() string {
	return self.exp
}

// String Type
type String struct {
	Atom
	Value string
}

func NewString(p string) *String {
	v := new(String)
	v.Value = p
	return v
}

func (self *String) String() string {
	return "\"" + self.Value + "\""
}

// Nil Type
type Nil struct {
	Atom
	value string
}

func NewNil() *Nil {
	n := new(Nil)
	n.value = "nil"
	return n
}

func (self *Nil) String() string {
	return self.value
}

// List Type
type List struct {
	Expression
	Value []Expression
}

func NewList(exp []Expression) *List {
	l := new(List)
	l.Value = exp
	return l
}

func (self *List) String() string {
	var buffer bytes.Buffer
	var makeString func(*List)

	makeString = func(l *List) {
		buffer.WriteString("(")

		for _, i := range l.Value {
			if j, ok := i.(*List); ok {
				makeString(j)

			} else if j, ok := i.(Expression); ok {
				buffer.WriteString(j.String())
			}
			if i != l.Value[len(l.Value)-1] {
				buffer.WriteString(" ")
			}
		}
		buffer.WriteString(")")
	}
	makeString(self)
	return buffer.String()
}

// Pair Type
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
func (self *Pair) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("(")
	buffer.WriteString(self.Car.String())
	buffer.WriteString(" . ")
	buffer.WriteString(self.Cdr.String())
	buffer.WriteString(")")
	return buffer.String()
}

// BuildInFunc
type BuildInFunc struct {
	Atom
	Impl EvalFunc
	name string
}

func NewBuildInFunc(fn EvalFunc, key string) *BuildInFunc {
	f := new(BuildInFunc)
	f.Impl = fn
	f.name = key
	return f
}
func (self *BuildInFunc) String() string {
	return "Build In Function: " + self.name
}
func (self *BuildInFunc) Execute(exp []Expression, env *SimpleEnv) (Expression, error) {
	return self.Impl(exp, env)
}

// Function (lambda). Env is exists for closure.
type Function struct {
	Expression
	ParamName List
	Body      []Expression
	Env       *SimpleEnv
	Name      string
}

func NewFunction(parent *SimpleEnv, param *List, body []Expression, name string) *Function {
	self := new(Function)
	self.ParamName = *param
	self.Body = body
	self.Env = parent
	self.Name = name
	return self
}
func (self *Function) String() string {
	return "Function: "
}
func (self *Function) Execute(exp []Expression, env *SimpleEnv) (Expression, error) {

	// Bind lambda function' parameters.
	if len(self.ParamName.Value) != len(exp) {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	nse := NewSimpleEnv(self.Env, nil)
	idx := 0
	for _, v := range self.ParamName.Value {
		if sym, ok := v.(*Symbol); ok {
			if env != nil {
				v, err := eval(exp[idx], env)
				if err != nil {
					return nil, err
				}
				nse.Regist(sym.Value, v)
			} else {
				nse.Regist(sym.Value, exp[idx])
			}
			idx++
		}
	}
	var (
		result Expression
		err    error
	)
	for _, e := range self.Body {
		if body, ok := e.(*List); ok && self.Name != "lambda" {
			evalTailRecursion(nse, body, self.Name, self.ParamName.Value)
		}
		for {
			result, err = eval(e, nse)
			if err != nil {
				return nil, err
			}
			if _, ok := result.(*TailRecursion); !ok {
				break
			}
		}
	}
	return result, nil
}

// Promise
type Promise struct {
	Expression
	Body Expression
	Env  *SimpleEnv
}

func NewPromise(parent *SimpleEnv, body Expression) *Promise {
	fn := new(Promise)
	fn.Body = body
	fn.Env = NewSimpleEnv(parent, nil)
	return fn
}

func (self *Promise) String() string {
	return "Promise: "
}

// TailRecursion
type TailRecursion struct {
	Expression
	param    []Expression
	nameList []Expression
}

func NewTailRecursion(param []Expression, nameList []Expression) *TailRecursion {
	self := new(TailRecursion)
	self.param = param
	self.nameList = nameList
	return self
}

func (self *TailRecursion) SetParam(env *SimpleEnv) (Expression, error) {

	values := make([]Expression, 0, 8)
	for i, _ := range self.nameList {
		v, err := eval(self.param[i], env)
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}

	for i, c := range self.nameList {
		pname := c.(*Symbol)
		(*env).Set(pname.Value, values[i])
	}
	return self, nil
}

func (self *TailRecursion) String() string {
	return "TailRecursion"
}

// lex support  for  string
func tokenize(s string) []string {
	var token []string
	stringMode := false
	symbolName := make([]rune, 0, 1024)
	from := 0

	s = strings.NewReplacer("\t", " ", "\n", " ", "\r", " ").Replace(s)
	for i, c := range s {
		if stringMode {
			if c == '"' {
				if s[i-1] != '\\' {
					token = append(token, s[from:i+1])
					stringMode = false
				}
			}
		} else {
			if c == '"' {
				from = i
				stringMode = true
			} else if c == '(' {
				token = append(token, "(")
			} else if c == ')' {
				token = append(token, ")")
			} else if c == ' ' {
				// Nop
			} else {
				n := len(string(c)) // for multibyte
				symbolName = append(symbolName, c)
				if len(s)-n == i {
					token = append(token, string(symbolName))
				} else {
					switch s[i+n] {
					case '(', ')', ' ':
						token = append(token, string(symbolName))
						symbolName = make([]rune, 0, 1024)
					}
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
				count++
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
		return nil, 0, NewSyntaxError("E0003")
	} else {
		atomType, err := atom(token)
		return atomType, 1, err
	}
}

// Atom To "Integer, Float, Symbol"
func atom(token string) (Atom, error) {
	var (
		atom Atom
	)
	if ivalue, err := strconv.Atoi(token); err == nil {
		atom = NewInteger(ivalue)
	} else {
		if fvalue, err := strconv.ParseFloat(token, 64); err == nil {
			atom = NewFloat(fvalue)
		} else if f, ok := buildInFuncTbl[token]; ok {
			atom = NewBuildInFunc(f, token)
		} else {
			switch token {
			case "#t":
				atom = NewBoolean(true)
			case "#f":
				atom = NewBoolean(false)
			default:
				if strings.Index(token, "#\\") == 0 {
					whitespaceChar := map[string]byte{
						"#\\tab":     0x09,
						"#\\space":   0x20,
						"#\\newline": 0x0A,
						"#\\return":  0x0D,
					}
					if v, ok := whitespaceChar[token]; ok {
						char := NewChar(token)
						char.Value = rune(v)
						atom = char
					} else if len([]rune(token)) == 3 {
						atom = NewChar(token)
					} else {
						return nil, NewSyntaxError("E0004")
					}
				} else if (len(token) > 1) && (token[0] == '"') && (token[len(token)-1] == '"') {
					atom = NewString(token[1 : len(token)-1])
				} else {
					atom = NewSymbol(token)
				}
			}
		}

	}
	return atom, nil
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
			} else if v, ok := buildInFuncTbl[sym.Value]; ok {
				return NewBuildInFunc(v, sym.Value), nil
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
		if f, ok := v[0].(*BuildInFunc); ok {
			return f.Execute(v[1:], env)
		}
		proc, err := eval(v[0], env)
		if err != nil {
			return sexp, err
		}
		if f, ok := proc.(*BuildInFunc); ok {
			// (* (+ a 1) (+ b 2)),(if (= a b) "a" "b")
			return f.Execute(v[1:], env)
		} else if fn, ok := proc.(*Function); ok {
			// (proc 10 20)
			return fn.Execute(v[1:], env)
		} else {
			return sexp, NewRuntimeError("E1006", reflect.TypeOf(proc).String())
		}
	} else if te, ok := sexp.(*TailRecursion); ok {
		return te.SetParam(env)
	}
	return sexp, NewRuntimeError("E1009", reflect.TypeOf(sexp).String())
}

// eval tail recursion
func evalTailRecursion(env *SimpleEnv, body *List, label string, nameList []Expression) {

	if len(body.Value) == 0 {
		return
	}
	v := body.Value
	for i := 0; i < len(body.Value); i++ {
		if l, ok := v[i].(*List); ok {
			if len(l.Value) == 0 {
				continue
			}
			if sym, ok := l.Value[0].(*Symbol); ok {
				proc, err := eval(l.Value[0], env)
				if err != nil {
					return
				}
				if fn, ok := proc.(*Function); ok && fn.Name != "lambda" && label == fn.Name {
					v[i] = NewTailRecursion(l.Value[1:], nameList)
					continue
				}
				if sym.Value == "if" || sym.Value == "cond" || sym.Value == "else" {
					evalTailRecursion(env, l, label, nameList)
				}
			}
		}
	}
	return
}
func evalMulti(exp []Expression, env *SimpleEnv) (Expression, error) {
	var (
		v   Expression
		err error
	)
	for _, e := range exp {
		v, err = eval(e, env)
		if err != nil {
			return v, err
		}
	}
	return v, nil
}

// main logic
func DoCoreLogic(program string, rootEnv *SimpleEnv) (Expression, error) {

	token := tokenize(program)
	ast, c, err := parse(token)
	if err != nil {
		return nil, err
	}

	if c != len(token) {
		return nil, NewSyntaxError("E0003")
	}

	val, err := eval(ast, rootEnv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// CUI desu.
func countParenthesis(program string) bool {
	left := 0
	right := 0
	search := true

	for _, c := range program {
		if c == '"' && search {
			search = false
		} else if c == '"' && !search {
			search = true
		}
		if c == '(' && search {
			left++
		}
		if c == ')' && search {
			right++
		}
	}
	return left <= right
}

// CUI desu.
func DoInteractive() {
	rootEnv := NewSimpleEnv(nil, nil)
	repl(os.Stdin, rootEnv)
}

// Read-eval-print loop
func repl(stream *os.File, rootEnv *SimpleEnv) {
	program := make([]string, 0, 64)
	prompt := PROMPT
	reader := bufio.NewReaderSize(stream, MaxLineSize)

	for {
		if stream == os.Stdin {
			fmt.Print(prompt + " ")
		}
		b, _, err := reader.ReadLine()
		line := string(b)
		if err == io.EOF {
			break
		} else if line == "" {
			continue
		} else if line[0] == ';' {
			continue
		} else if line == "(quit)" {
			break
		}
		program = append(program, line)
		if !countParenthesis(strings.Join(program, " ")) {
			prompt = ""
			continue
		}
		val, err := DoCoreLogic(strings.Join(program, " "), rootEnv)
		if err != nil {
			fmt.Println(err.Error())
			goto FINISH
		}
		fmt.Println(val.String())
		if DEBUG {
			fmt.Print(reflect.TypeOf(val))
		}
	FINISH:
		program = make([]string, 0, 64)
		prompt = PROMPT
	}
}
func EvalCalcParam(exp []Expression, env *SimpleEnv, fn EvaledAllParamFunc) (Expression, error) {
	var args []Expression

	for _, e := range exp {
		c, err := eval(e, env)
		if err != nil {
			return e, err
		}
		args = append(args, c)
	}
	return fn(args...)
}

// addl, subl, imul, idiv
func calcOperate(calc func(Number, Number) Number, exp ...Expression) (Number, error) {
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

// gt,lt,ge,le
func cmpOperate(cmp func(Number, Number) bool, exp ...Expression) (*Boolean, error) {
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

// map,filter,reduce
func listFunc(lambda func(Expression, Expression, []Expression) ([]Expression, error), exp ...Expression) (Expression, error) {
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
	var vaList []Expression
	param := make([]Expression, 1)
	for _, param[0] = range l.Value {
		result, err := fn.Execute(param, nil)
		if err != nil {
			return nil, err
		}
		vaList, err = lambda(result, param[0], vaList)
		if err != nil {
			return nil, err
		}
	}
	return NewList(vaList), nil
}

// and or not
func logicalOperate(exp []Expression, env *SimpleEnv, bcond bool, bret bool) (Expression, error) {
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

// math skelton
func mathImpl(mathFunc func(float64) float64, exp ...Expression) (Expression, error) {
	if len(exp) != 1 {
		return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
	}
	if v, ok := exp[0].(*Float); ok {
		return NewFloat(mathFunc(v.Value)), nil
	} else if v, ok := exp[0].(*Integer); ok {
		return NewFloat(mathFunc((float64)(v.Value))), nil
	}
	return nil, NewRuntimeError("E1003", reflect.TypeOf(exp[0]).String())
}

// imul, skelton
func idivImpl(idivFunc func(int, int) int, exp ...Expression) (Expression, error) {
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
	if prm[1].Value == 0 {
		return nil, NewRuntimeError("E1013")
	}
	return NewInteger(idivFunc(prm[0].Value, prm[1].Value)), nil
}

// Build Global environement.
func BuildFunc() {
	buildInFuncTbl = map[string]EvalFunc{}

	buildInFuncTbl["+"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return calcOperate(func(a Number, b Number) Number { return a.Add(b) }, exp...)
			})
	}
	buildInFuncTbl["-"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return calcOperate(func(a Number, b Number) Number { return a.Sub(b) }, exp...)
			})
	}
	buildInFuncTbl["*"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return calcOperate(func(a Number, b Number) Number { return a.Mul(b) }, exp...)
			})
	}
	buildInFuncTbl["/"] = func(exp []Expression, env *SimpleEnv) (se Expression, e error) {
		// Not the best. But, Better than before.
		defer func() {
			if err := recover(); err != nil {
				if zero, ok := err.(*RuntimeError); ok {
					se = nil
					e = zero
				}
			}
		}()
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return calcOperate(func(a Number, b Number) Number { return a.Div(b) }, exp...)
			})
	}
	buildInFuncTbl["quotient"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return idivImpl(func(a int, b int) int { return a / b }, exp...)
			})
	}
	buildInFuncTbl["modulo"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return idivImpl(func(a int, b int) int { return a % b }, exp...)
			})
	}
	buildInFuncTbl[">"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return cmpOperate(func(a Number, b Number) bool { return a.Greater(b) }, exp...)
			})
	}
	buildInFuncTbl["<"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return cmpOperate(func(a Number, b Number) bool { return a.Less(b) }, exp...)
			})
	}
	buildInFuncTbl[">="] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return cmpOperate(func(a Number, b Number) bool { return a.GreaterEqual(b) }, exp...)
			})
	}
	buildInFuncTbl["<="] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return cmpOperate(func(a Number, b Number) bool { return a.LessEqual(b) }, exp...)
			})
	}
	buildInFuncTbl["="] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				return cmpOperate(func(a Number, b Number) bool { return a.Equal(b) }, exp...)
			})
	}
	buildInFuncTbl["not"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 1 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if _, ok := exp[0].(*Boolean); !ok {
					return nil, NewRuntimeError("E1001", reflect.TypeOf(exp[0]).String())
				}
				return NewBoolean(!(exp[0].(*Boolean)).Value), nil
			})
	}

	// list operator
	buildInFuncTbl["list"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				var l []Expression
				return NewList(append(l, exp...)), nil
			})
	}
	buildInFuncTbl["null?"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 1 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if l, ok := exp[0].(*List); ok {
					return NewBoolean(0 == len(l.Value)), nil
				} else {
					return NewBoolean(false), nil
				}
			})
	}
	buildInFuncTbl["length"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 1 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if l, ok := exp[0].(*List); ok {
					return NewInteger(len(l.Value)), nil
				} else {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
			})
	}
	buildInFuncTbl["car"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 1 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if l, ok := exp[0].(*List); ok {
					if len(l.Value) <= 0 {
						return nil, NewRuntimeError("E1011", strconv.Itoa(len(l.Value)))
					}
					return l.Value[0], nil
				} else if p, ok := exp[0].(*Pair); ok {
					return p.Car, nil
				} else {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
			})
	}
	buildInFuncTbl["cdr"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 1 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if l, ok := exp[0].(*List); ok {
					if len(l.Value) <= 0 {
						return nil, NewRuntimeError("E1011", strconv.Itoa(len(l.Value)))
					}
					return NewList(l.Value[1:]), nil
				} else if p, ok := exp[0].(*Pair); ok {
					return p.Cdr, nil
				} else {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
			})
	}
	buildInFuncTbl["cadr"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 1 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if l, ok := exp[0].(*List); ok {
					if len(l.Value) < 2 {
						return nil, NewRuntimeError("E1011", strconv.Itoa(len(l.Value)))
					}
					return l.Value[1], nil
				} else {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
			})
	}
	buildInFuncTbl["cons"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 2 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if _, ok := exp[1].(*List); ok {
					var args []Expression
					args = append(args, exp[0])
					return NewList(append(args, (exp[1].(*List)).Value...)), nil
				}
				return NewPair(exp[0], exp[1]), nil
			})
	}
	buildInFuncTbl["append"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) < 2 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				var expList []Expression
				for _, e := range exp {
					if v, ok := e.(*List); ok {
						expList = append(expList, v.Value...)
					} else {
						return nil, NewRuntimeError("E1005", reflect.TypeOf(e).String())
					}
				}
				return NewList(expList), nil
			})
	}
	buildInFuncTbl["last"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 1 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				if l, ok := exp[0].(*List); ok {
					if len(l.Value) <= 0 {
						return nil, NewRuntimeError("E1011", strconv.Itoa(len(l.Value)))
					}
					return l.Value[len(l.Value)-1], nil
				} else if p, ok := exp[0].(*Pair); ok {
					return p.Car, nil
				} else {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
				}
			})
	}
	buildInFuncTbl["reverse"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
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
			})
	}
	buildInFuncTbl["iota"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
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
			})
	}
	buildInFuncTbl["map"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 2 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				lambda := func(result Expression, item Expression, vaList []Expression) ([]Expression, error) {
					return append(vaList, result), nil
				}
				return listFunc(lambda, exp...)
			})
	}
	buildInFuncTbl["for-each"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
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
				param := make([]Expression, 1)
				for _, param[0] = range l.Value {
					_, err := fn.Execute(param, nil)
					if err != nil {
						return nil, err
					}
				}
				return NewNil(), nil
			})
	}
	buildInFuncTbl["filter"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 2 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				lambda := func(result Expression, item Expression, vaList []Expression) ([]Expression, error) {
					b, ok := result.(*Boolean)
					if !ok {
						return nil, NewRuntimeError("E1001", reflect.TypeOf(result).String())
					}
					if b.Value {
						return append(vaList, item), nil
					}
					return vaList, nil
				}
				return listFunc(lambda, exp...)
			})
	}
	buildInFuncTbl["reduce"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 3 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				fn, ok := exp[0].(*Function)
				if !ok {
					return nil, NewRuntimeError("E1006", reflect.TypeOf(exp[0]).String())
				}
				l, ok := exp[2].(*List)
				if !ok {
					return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[1]).String())
				}
				if len(l.Value) == 0 {
					return exp[1], nil
				}
				param := make([]Expression, len(fn.ParamName.Value))
				result := l.Value[0]
				for _, c := range l.Value[1:] {
					param[0] = result
					param[1] = c
					r, err := fn.Execute(param, nil)
					result = r
					if err != nil {
						return nil, err
					}
				}
				return result, nil
			})
	}
	buildInFuncTbl["sqrt"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) { return mathImpl(math.Sqrt, exp...) })
	}
	buildInFuncTbl["sin"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) { return mathImpl(math.Sin, exp...) })
	}
	buildInFuncTbl["cos"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) { return mathImpl(math.Cos, exp...) })
	}
	buildInFuncTbl["tan"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) { return mathImpl(math.Tan, exp...) })
	}
	buildInFuncTbl["atan"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) { return mathImpl(math.Atan, exp...) })
	}
	buildInFuncTbl["log"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) { return mathImpl(math.Log, exp...) })
	}
	buildInFuncTbl["exp"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) { return mathImpl(math.Exp, exp...) })
	}
	buildInFuncTbl["rand-init"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 0 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			return nil, NewRuntimeError("E9999")
		}
		rand.Seed(seed.Int64())
		return NewNil(), nil
	}
	buildInFuncTbl["rand-integer"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		if v, ok := exp[0].(*Integer); ok {
			return NewInteger(rand.Intn(v.Value)), nil
		}
		return nil, NewRuntimeError("E1002", reflect.TypeOf(exp[0]).String())
	}
	buildInFuncTbl["expt"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return EvalCalcParam(exp, env,
			func(exp ...Expression) (Expression, error) {
				if len(exp) != 2 {
					return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
				}
				f := []float64{0.0, 0.0}
				for i, e := range exp {
					if n, ok := e.(*Float); ok {
						f[i] = n.Value
					} else if n, ok := e.(*Integer); ok {
						f[i] = (float64)(n.Value)
					} else {
						return nil, NewRuntimeError("E1003", reflect.TypeOf(e).String())
					}
				}
				return NewFloat(math.Pow(f[0], f[1])), nil
			})
	}
	// syntax keyword implements
	buildInFuncTbl["if"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		e, err := eval(exp[0], env)
		if err != nil {
			return nil, err
		}

		b, ok := e.(*Boolean)
		if !ok {
			return nil, NewRuntimeError("E1001", reflect.TypeOf(exp).String())
		}
		if b.Value {
			return eval(exp[1], env)
		} else if 3 <= len(exp) {
			return eval(exp[2], env)
		}
		return NewNil(), nil
	}
	buildInFuncTbl["define"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		var (
			key *Symbol
			err error
			ok  bool
			e   Expression
		)
		if key, ok = exp[0].(*Symbol); ok {
			e, err = eval(exp[1], env)
			if err != nil {
				return nil, err
			}
		} else if l, ok := exp[0].(*List); ok {
			if len(l.Value) < 1 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			key, ok = l.Value[0].(*Symbol)
			if !ok {
				return nil, NewRuntimeError("E1004", reflect.TypeOf(exp[0]).String())
			}
			for _, p := range l.Value[1:] {
				e, ok := p.(*Symbol)
				if !ok {
					return nil, NewRuntimeError("E1004", reflect.TypeOf(e).String())
				}
			}
			e = NewFunction(env, NewList(l.Value[1:]), exp[1:], key.Value)
		} else {
			return nil, NewRuntimeError("E1004", reflect.TypeOf(exp[0]).String())
		}
		(*env).Regist(key.Value, e)
		return key, nil
	}
	buildInFuncTbl["lambda"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		// if l == (), internal list implements
		l, ok := exp[0].(*List)
		if !ok {
			return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[0]).String())
		}
		for _, p := range l.Value {
			e, ok := p.(*Symbol)
			if !ok {
				return nil, NewRuntimeError("E1004", reflect.TypeOf(e).String())
			}
		}
		return NewFunction(env, l, exp[1:], "lambda"), nil
	}
	buildInFuncTbl["set!"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		s, ok := exp[0].(*Symbol)
		if !ok {
			return nil, NewRuntimeError("E1004", reflect.TypeOf(exp[0]).String())
		}

		e, err := eval(exp[1], env)
		if err != nil {
			return exp[1], err
		}
		if _, ok := (*env).Find(s.Value); !ok {
			return e, NewRuntimeError("E1008", s.Value)
		}
		(*env).Set(s.Value, e)
		return e, nil
	}
	buildInFuncTbl["let"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		var letsym *Symbol
		var pname []Expression
		body := 1

		l, ok := exp[0].(*List)
		if ok && len(exp) < 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		if !ok {
			letsym, ok = exp[0].(*Symbol)
			if !ok {
				return nil, NewRuntimeError("E1004", reflect.TypeOf(exp[0]).String())
			}
			if len(exp) < 3 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
			}
			l, ok = exp[1].(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(exp[1]).String())
			}
			body = 2
		}

		localEnv := Environment{}
		for _, let := range l.Value {
			r, ok := let.(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(let).String())
			}
			if len(r.Value) != 2 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(r.Value)))
			}
			sym, ok := (r.Value[0]).(*Symbol)
			if !ok {
				return nil, NewRuntimeError("E1004", reflect.TypeOf(r.Value[0]).String())
			}
			v, err := eval(r.Value[1], env)
			if err != nil {
				return nil, err
			}
			pname = append(pname, sym)
			localEnv[sym.Value] = v
		}

		nse := NewSimpleEnv(env, &localEnv)
		if letsym != nil {
			localEnv[letsym.Value] = NewFunction(nse, NewList(pname), exp[body:], letsym.Value)
		}
		var lastExp Expression
		for i := body; i < len(exp); i++ {
			if e, err := eval(exp[i], nse); err == nil {
				lastExp = e
			} else {
				return nil, err
			}
		}
		return lastExp, nil

	}
	buildInFuncTbl["and"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return logicalOperate(exp, env, false, true)
	}
	buildInFuncTbl["or"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		return logicalOperate(exp, env, true, false)
	}
	buildInFuncTbl["delay"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return NewPromise(env, exp[0]), nil
	}
	buildInFuncTbl["force"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		e, err := eval(exp[0], env)
		if err != nil {
			return nil, err
		}
		p, ok := e.(*Promise)
		if !ok {
			return e, nil
		} else {
			return eval(p.Body, p.Env)
		}
	}
	buildInFuncTbl["identity"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return eval(exp[0], env)
	}
	buildInFuncTbl["cond"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		for _, e := range exp {
			l, ok := e.(*List)
			if !ok {
				return nil, NewRuntimeError("E1005", reflect.TypeOf(e).String())
			}
			if len(l.Value) < 2 {
				return nil, NewRuntimeError("E1007", strconv.Itoa(len(l.Value)))
			}
			if _, ok := l.Value[0].(*List); ok {
				b, err := eval(l.Value[0], env)
				if err != nil {
					return nil, err
				}
				if b, ok := b.(*Boolean); ok {
					if b.Value {
						return evalMulti(l.Value[1:], env)
					}
				} else {
					return nil, NewRuntimeError("E1001")
				}
			} else if sym, ok := l.Value[0].(*Symbol); ok {
				if sym.Value == "else" {
					return evalMulti(l.Value[1:], env)
				} else {
					return nil, NewRuntimeError("E1012")
				}
			} else {
				return nil, NewRuntimeError("E1012")
			}
		}
		return NewNil(), nil
	}
	buildInFuncTbl["quote"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return exp[0], nil
	}
	buildInFuncTbl["time"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		t0 := time.Now()
		e, err := eval(exp[0], env)
		t1 := time.Now()
		fmt.Println(t1.Sub(t0))
		return e, err
	}
	buildInFuncTbl["begin"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {

		if len(exp) < 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		return evalMulti(exp, env)
	}
	buildInFuncTbl["display"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		for _, e := range exp {
			v, err := eval(e, env)
			if err != nil {
				return v, err
			}
			fmt.Print(v.String())
		}
		return NewNil(), nil
	}
	buildInFuncTbl["newline"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 0 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		fmt.Println("")
		return NewNil(), nil
	}
	buildInFuncTbl["load-file"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {

		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		v, err := eval(exp[0], env)
		if err != nil {
			return v, err
		}
		filename, ok := v.(*String)
		if !ok {
			return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
		}
		f, err := os.Stat(filename.Value)
		if err != nil && os.IsNotExist(err) {
			return nil, NewRuntimeError("E1014")
		}
		if err != nil || !f.Mode().IsRegular() {
			return nil, NewRuntimeError("E1016")
		}

		fd, err := os.Open(filename.Value)
		if err != nil {
			return nil, NewRuntimeError("E9999")
		}
		defer func() { _ = fd.Close() }()
		repl(fd, env)
		return NewNil(), nil
	}
	//srfi-98
	buildInFuncTbl["get-environment-variable"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) != 1 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		v, err := eval(exp[0], env)
		if err != nil {
			return v, err
		}
		s, ok := v.(*String)
		if !ok {
			return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
		}
		return NewString(os.Getenv(s.Value)), nil
	}
	buildInFuncTbl["string-append"] = func(exp []Expression, env *SimpleEnv) (Expression, error) {
		if len(exp) < 2 {
			return nil, NewRuntimeError("E1007", strconv.Itoa(len(exp)))
		}
		ret := make([]string, 0, len(exp))
		for _, e := range exp {
			v, err := eval(e, env)
			if err != nil {
				return v, err
			}
			s, ok := v.(*String)
			if !ok {
				return nil, NewRuntimeError("E1015", reflect.TypeOf(exp[0]).String())
			}
			ret = append(ret, s.Value)
		}
		return NewString(strings.Join(ret, "")), nil
	}
}

// add error message from other package
func AddErrorMsg(code string, value string) {
	errorMsg[code] = value
}

// add func
func AddBuildInFunc(name string, body EvalFunc) {
	buildInFuncTbl[name] = body
}

func DoEval(sexp Expression, env *SimpleEnv) (Expression, error) {
	return eval(sexp, env)
}