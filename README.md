Implementation of Scheme (subset version) by Go Lang
=================

## Overview
- Implemented a Lisp for Go Lang lessons. (It's Scheme base)
- As an implementation goal, we will provide an environment for easily operating a graphic program.

## Quality
- Level at which a simple program works
    - https://github.com/hidekuno/go-scheme/blob/master/scheme/lisp_test.go
- I confirmed that the SICP graphic language program works.
    - https://github.com/hidekuno/picture-language

## TEST & Run (CLI interpreter)
### Requirement
- go lang installed.

```
cd ${HOME}
git clone https://github.com/hidekuno/go-scheme
cd go-scheme/scheme
go test -v
go run cmd/lisp/main.go
```

## TEST & Run (Scheme on GUI for Draw)
### Requirement
- X server is running.

```
cd ${HOME}/go-scheme/draw
go test -v
go run cmd/lisp/main.go
```
- Then type "(draw-init)"

## TEST & Run (Scheme on Web API Server)
```
cd {HOME}/go-scheme/web
go test -v
go run cmd/api/main.go
```

## Build & Run (Scheme on Web Assembly)
```
cd {HOME}/go-scheme/web/wasm
cp /usr/local/go/misc/wasm/wasm_exec.js .
GOARCH=wasm GOOS=js go build -o lisp.wasm lisp_wasm.go
cd {HOME}/go-scheme/web
go run cmd/wasm/main.go
```

## Run on docker
### Requirement
- docker is running.
- X Server is running.(XQuartz 2.7.11 for mac)

### macOS
```
docker pull hidekuno/go-scheme
xhost +
docker run -it --name go-scheme -e DISPLAY=docker.for.mac.localhost:0 hidekuno/go-scheme /root/glisp
```

<img src="https://user-images.githubusercontent.com/22115777/68912921-e9619300-079c-11ea-976c-340252936eb1.png" width=50% height="50%">

### Linux
```
docker pull hidekuno/go-scheme
xhost +
docker run -it --name go-scheme -e DISPLAY=${HOSTIP}:0.0 hidekuno/go-scheme /root/glisp
```

### For environments where the X server is not running
```
docker pull hidekuno/go-scheme
docker run -it --name go-scheme hidekuno/go-scheme /root/lisp
```

<img src="https://user-images.githubusercontent.com/22115777/67071430-783eb800-f1bd-11e9-9a94-18c3b371ab39.png" width=80% height="80%">
