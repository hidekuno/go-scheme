package main

/*
   Go lang 3rd study program.
   hidekuno@gmail.com

   GOARCH=wasm GOOS=js go build -o wasm/lisp.wasm lisp_wasm.go
*/
import (
	"scheme"
	"syscall/js"
)

var rootEnv *scheme.SimpleEnv

func eval(this js.Value,vs []js.Value) interface{}{

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
