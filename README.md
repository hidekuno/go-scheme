Go言語によるScheme(subset版)の実装
=================

## 概要
- Go言語手習いのため、Schemeの縮小版を実装した。
- 実装目標として、フラクタル図形プログラムを簡単に動作させるための環境を提供する。
- さらに、WebAPIを実装してみた。

## 完成度合い
- 簡単なプログラム(quick-sort, merge-sort, permutations,combination,Eratosthenes等)が動くレベル  
https://github.com/hidekuno/go-scheme/blob/master/src/scheme/lisp_test.go

- SICPの図形言語プログラムが動作するところまで確認した。

<img src="https://user-images.githubusercontent.com/22115777/44436239-11406600-a5ef-11e8-9860-0b3f73350114.png" width=50%>

## 開発環境
| Item   | Ver. |備考|
|--------|--------|--------|
| OS     | CentOS | draw系を使わなければ特になし|
| Gtk+   | 2.24.31||
| golang   | 1.9.4||
| go-gtk | release-0.1|https://github.com/mattn/go-gtk|
| gorilla/sessions|v1.1.2|https://github.com/gorilla/sessions|

## 変数
```
GOARCH=amd64
GOOS=linux
GOPATH=${HOME}/go-scheme:${HOME}/go
```

## その他
インストールの方法、動かし方などは下記を参照

https://github.com/hidekuno/go-scheme/wiki
