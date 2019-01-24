package main

/*
   Go lang 3rd study program.
   hidekuno@gmail.com

   GOARCH=wasm GOOS=js go build -o lisp.wasm wasm_lisp.go
*/
import (
	"scheme"
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
	result.Set("innerText", exp.String())
}

func main() {

	println("Wasm Lisp Initialized")
	scheme.BuildFunc()
	rootEnv = scheme.NewSimpleEnv(nil, nil)

	js.Global().Set("eval", js.NewCallback(eval))
	select {}
}
