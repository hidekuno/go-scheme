package main

/*
   Go lang 3rd study program.
   hidekuno@gmail.com

   wget https://raw.githubusercontent.com/golang/go/release-branch.go1.XX/misc/wasm/wasm_exec.js
   or
   cp /usr/local/go/misc/wasm/wasm_exec.js .

   GOARCH=wasm GOOS=js go build -o lisp.wasm lisp_wasm.go
*/
import (
	"syscall/js"

	"github.com/hidekuno/go-scheme/scheme"
)

var rootEnv *scheme.SimpleEnv

func eval(this js.Value, vs []js.Value) interface{} {

	document := js.Global().Get("document")
	text := document.Call("getElementById", "sExpression")
	result := document.Call("getElementById", "calcResult")

	exp, err := scheme.DoCoreLogic(text.Get("value").String(), rootEnv)
	if err != nil {
		result.Set("innerText", err.Error())
		println(err.Error())
		return err.Error()
	}
	result.Set("innerText", exp.String())
	return exp.String()
}

func main() {

	println("Wasm Lisp Initialized")
	scheme.BuildFunc()
	rootEnv = scheme.NewSimpleEnv(nil, nil)

	js.Global().Set("eval", js.FuncOf(eval))
	select {}
}
