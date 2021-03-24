all: lisp.wasm

lisp.wasm: lisp_wasm.go
	GOARCH=wasm GOOS=js go build -o lisp.wasm lisp_wasm.go
clean:
	rm -f lisp.wasm wasm_exec.js
jscp:
	cp /usr/local/go/misc/wasm/wasm_exec.js .
