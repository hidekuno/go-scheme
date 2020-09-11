Go言語によるScheme(subset版)の実装
=================

## 概要
- Go言語手習いのため、Schemeの縮小版を実装した。
- 実装目標として、フラクタル図形プログラムを簡単に動作させるための環境を提供する。
- さらに、WebAPIを実装してみた。

## 完成度合い
- 簡単なプログラムが動くレベル
    - https://github.com/hidekuno/go-scheme/blob/master/src/scheme/lisp_test.go
- SICPの図形言語プログラムが動作するところまで確認した。
    - https://github.com/hidekuno/picture-language

## 変数
```
GOARCH=amd64
GOOS=linux
GOPATH=${HOME}/go-scheme:${HOME}/go
```

## 動かし方
### 動作条件
- dockerが動いていること
- Xサーバ(macの場合、XQuartz)が動いていること

### macOS
```
docker pull hidekuno/go-scheme
docker run -it --name go-scheme -e DISPLAY=docker.for.mac.localhost:0 hidekuno/go-scheme /root/lisp_draw_main
```

<img src="https://user-images.githubusercontent.com/22115777/68912921-e9619300-079c-11ea-976c-340252936eb1.png" width=80%>

### Linux
```
docker pull hidekuno/go-scheme
xhost +
docker run -it --name go-scheme -e DISPLAY=${HOSTIP}:0.0 hidekuno/go-scheme /root/lisp_draw_main
```
### Xサーバが動いていない環境向け(replのみ版)
```
docker pull hidekuno/go-scheme
docker run -it --name go-scheme hidekuno/go-scheme /root/lisp_main
```

<img src="https://user-images.githubusercontent.com/22115777/67071430-783eb800-f1bd-11e9-9a94-18c3b371ab39.png" width=80%>
