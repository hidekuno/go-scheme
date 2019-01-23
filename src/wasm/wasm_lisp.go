package main

/*
   Go lang 3rd study program.
   hidekuno@gmail.com

   GOARCH=wasm GOOS=js go build -o lisp.wasm wasm_lisp.go
*/
import (
	"bytes"
	"scheme"
	"strconv"
	"syscall/js"
)

var rootEnv *scheme.SimpleEnv

func eval(i []js.Value) {
	document := js.Global().Get("document")
	text := document.Call("getElementById", "sExpression")
	result := document.Call("getElementById", "calcResult")

	exp, err := scheme.DoCoreLogic(text.Get("value").String(), rootEnv)
	if err != nil {
		result.Set("innerText", err.Error())
		println(err.Error())
		return
	}
	var buffer bytes.Buffer
	var make_string func(scheme.Expression)

	make_string = func(exp scheme.Expression) {
		if l, ok := exp.(*scheme.List); ok {
			buffer.WriteString("(")
			for _, i := range l.Value {
				if j, ok := i.(*scheme.List); ok {
					make_string(j)
				} else if j, ok := i.(scheme.Atom); ok {
					make_string(j)
				}
				if i != l.Value[len(l.Value)-1] {
					buffer.WriteString(" ")
				}
			}
			buffer.WriteString(")")
		} else if j, ok := exp.(*scheme.Integer); ok {
			buffer.WriteString(strconv.Itoa(j.Value))
		} else if j, ok := exp.(*scheme.Float); ok {
			buffer.WriteString(strconv.FormatFloat(j.Value, 'f', 8, 64))
		}
	}
	make_string(exp)
	result.Set("innerText", buffer.String())
}

func main() {

	println("Wasm Lisp Initialized")
	scheme.BuildFunc()
	rootEnv = scheme.NewSimpleEnv(nil, nil)

	js.Global().Set("eval", js.NewCallback(eval))
	select {}
}
