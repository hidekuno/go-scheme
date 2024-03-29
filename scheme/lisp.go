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
		"E1018": "Not Format Gramar",
		"E1019": "Not Char",
		"E1020": "Not Rat",
		"E1021": "Out Of Range",
		"E1022": "Not Vector",
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
	isAtom() bool
	clone() Expression
	equalValue(Expression) bool
	Print()
}

// Symbol Type
type Symbol struct {
	Expression
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
func (self *Symbol) Print() {
	fmt.Print(self.Value)
}
func (self *Symbol) isAtom() bool {
	return true
}
func (self *Symbol) clone() Expression {
	return NewSymbol(self.Value)
}
func (self *Symbol) equalValue(e Expression) bool {
	if v, ok := e.(*Symbol); ok {
		return self.Value == v.Value
	}
	return false
}

// Nil Type
type Nil struct {
	Expression
	value string
}

func NewNil() *Nil {
	n := new(Nil)
	n.value = "#<nil>"
	return n
}

func (self *Nil) String() string {
	return self.value
}
func (self *Nil) Print() {
	fmt.Print(self.value)
}
func (self *Nil) isAtom() bool {
	return true
}
func (self *Nil) clone() Expression {
	return NewNil()
}
func (self *Nil) equalValue(e Expression) bool {
	// Not Support this method
	return false
}

// BuildInFunc
type BuildInFunc struct {
	Expression
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
func (self *BuildInFunc) Print() {
	fmt.Print(self.String())
}
func (self *BuildInFunc) Execute(exp []Expression, env *SimpleEnv) (Expression, error) {
	return self.Impl(exp, env)
}
func (self *BuildInFunc) isAtom() bool {
	return true
}
func (self *BuildInFunc) clone() Expression {
	return NewBuildInFunc(self.Impl, self.name)
}
func (self *BuildInFunc) equalValue(e Expression) bool {
	// Not Support this method
	return false
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
func (self *Function) Print() {
	fmt.Print(self.String())
}
func (self *Function) isAtom() bool {
	return true
}
func (self *Function) clone() Expression {
	return NewFunction(self.Env, &self.ParamName, self.Body, self.Name)
}
func (self *Function) equalValue(e Expression) bool {
	// Not Support this method
	return false
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
			v, err := eval(exp[idx], env)
			if err != nil {
				return nil, err
			}
			nse.Regist(sym.Value, v)
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
			if result, err = eval(e, nse); err != nil {
				if c, ok := err.(*Continuation); ok {
					if len(exp) == 1 {
						if s, ok := self.ParamName.Value[0].(*Symbol); ok {
							if c.Name == s.Value {
								return c.Value, nil
							}
						}
					}
				}
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
func (self *Promise) Print() {
	fmt.Print(self.String())
}
func (self *Promise) isAtom() bool {
	return false
}
func (self *Promise) clone() Expression {
	return NewPromise(self.Env, self.Body)
}
func (self *Promise) equalValue(e Expression) bool {
	// Not Support this method
	return false
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
func (self *TailRecursion) Print() {
	fmt.Print(self.String())
}
func (self *TailRecursion) isAtom() bool {
	return false
}
func (self *TailRecursion) clone() Expression {
	return NewTailRecursion(self.param, self.nameList)
}
func (self *TailRecursion) equalValue(e Expression) bool {
	// Not Support this method
	return false
}

// lex support  for  string
func tokenize(s string) ([]string, error) {
	var tokens []string
	stringMode := false
	quoteMode := false
	vectorMode := false
	tokenName := make([]rune, 0, 1024)
	from := 0
	left := 0
	right := 0

	rb := []rune(strings.NewReplacer("\t", " ", "\n", " ", "\r", " ").Replace(s))
	for i, c := range rb {
		if stringMode {
			if c == '"' {
				if rb[i-1] != '\\' {
					tokens = append(tokens, string(rb[from:i+1]))
					stringMode = false
				}
			}
		} else {
			if c == '"' {
				from = i
				stringMode = true
			} else if c == '(' {
				tokens = append(tokens, "(")
				if vectorMode == true {
					tokens = append(tokens, "vector")
					vectorMode = false
				}
				if quoteMode == true {
					left++
				}
			} else if c == ')' {
				tokens = append(tokens, ")")
				if quoteMode == true {
					right++
				}
				if quoteMode == true && left == right {
					tokens = append(tokens, ")")
					quoteMode = false
				}
			} else if c == '\'' {
				tokens = append(tokens, "(")
				tokens = append(tokens, "quote")

				quoteMode = true
			} else if c == ' ' {
				//Nop
			} else {
				if len(rb)-1 == i {
					tokenName = append(tokenName, c)
					tokens = append(tokens, string(tokenName))
					if quoteMode == true {
						tokens = append(tokens, ")")
						quoteMode = false
					}
				} else {
					if c == '#' && rb[i+1] == '(' {
						vectorMode = true
						continue
					}
					tokenName = append(tokenName, c)
					switch rb[i+1] {
					case '(', ')', ' ':
						tokens = append(tokens, string(tokenName))
						tokenName = make([]rune, 0, 1024)
						if quoteMode == true && left == right {
							tokens = append(tokens, ")")
							quoteMode = false
						}
					}
				}
			}
		}
	}
	if stringMode {
		return tokens, NewSyntaxError("E0004")
	}
	if DEBUG {
		fmt.Println(tokens)
	}
	return tokens, nil
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
func atom(token string) (Expression, error) {

	var (
		atom Expression
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
					if v, ok := whitespaceChar[token]; ok {
						atom = NewCharFromRune(rune(v))
					} else if len([]rune(token)) == 3 {
						atom = NewChar(token)
					} else {
						return nil, NewSyntaxError("E0004")
					}
				} else if (len(token) > 1) && (token[0] == '"') && (token[len(token)-1] == '"') {
					atom = NewString(token[1 : len(token)-1])
				} else {

					if rat := MakeRat(token); rat != nil {
						atom = rat
					} else {
						atom = NewSymbol(token)
					}
				}
			}
		}

	}
	return atom, nil
}

// Evaluate an expression in an environment.
func eval(sexp Expression, env *SimpleEnv) (Expression, error) {
	if DEBUG {
		fmt.Println(reflect.TypeOf(sexp))
	}
	if sexp.isAtom() {
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
		if b, ok := v[0].(*BuildInFunc); ok {
			return b.Execute(v[1:], env)
		}
		if f, ok := v[0].(*Function); ok {
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
		} else if c, ok := proc.(*Continuation); ok {
			// (proc 10 20)
			return c.Execute(v[0:], env)
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

// main logic
func DoCoreLogic(program string, rootEnv *SimpleEnv) (Expression, error) {

	var val Expression

	token, err := tokenize(program)
	if err != nil {
		return nil, err
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
func repl(stream io.Reader, rootEnv *SimpleEnv) {
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
	buildCharFunc()
	buildBooleanFunc()
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
