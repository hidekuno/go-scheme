/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
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
		"E1017": "Not Case Gramar",
		"E9999": "System Panic",
	}
	tracer = log.New(os.Stderr, "", log.Lshortfile)
)

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
func tokenize(s string) ([]string, error) {
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
	if stringMode {
		return token, NewSyntaxError("E0004")
	}
	if DEBUG {
		for _, c := range token {
			fmt.Println(c)
		}
	}
	return token, nil
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

	var val Expression

	token, err := tokenize(program)
	if err != nil {
		return NewNil(), err
	}
	for {
		ast, c, err := parse(token)
		if err != nil {
			return nil, err
		}
		val, err = eval(ast, rootEnv)
		if err != nil {
			return nil, err
		}
		if c == len(token) {
			break
		} else {
			token = token[c:]
		}
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

// Build Global environement.
func BuildFunc() {
	buildInFuncTbl = map[string]EvalFunc{}

	buildMathFunc()
	buildListFunc()
	buildOperationFunc()
	buildSyntaxFunc()
	buildIoFunc()
	buildUtilFunc()
	buildStringFunc()
}

// add func
func AddBuildInFunc(name string, body EvalFunc) {
	buildInFuncTbl[name] = body
}

// add error message from other package
func AddErrorMsg(code string, value string) {
	errorMsg[code] = value
}

func DoEval(sexp Expression, env *SimpleEnv) (Expression, error) {
	return eval(sexp, env)
}
