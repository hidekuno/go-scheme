Go言語によるScheme(subset版)の実装
=================

## 起動方法
```
[kunohi@centos7-dev-docker src]$ go run lisp.go
scheme.go>
```

## 終了方法
```
scheme.go>  (quit)
[kunohi@centos7-dev-docker src]$ 
```

## テスト方法
```
[kunohi@centos7-dev-docker src]$ go test -v lisp.go lisp_test.go 
=== RUN   Test_lisp_sample_program
--- PASS: Test_lisp_sample_program (0.00s)
=== RUN   Test_math_func
--- PASS: Test_math_func (0.00s)
=== RUN   Test_list_func
--- PASS: Test_list_func (0.00s)
=== RUN   Test_basic_opration
--- PASS: Test_basic_opration (0.00s)
=== RUN   Test_err_case
--- PASS: Test_err_case (0.00s)
=== RUN   Test_interactive
--- PASS: Test_interactive (0.00s)
	lisp_test.go:544: 3.5
	lisp_test.go:544: 30
	lisp_test.go:544: a
	lisp_test.go:544: #t
	lisp_test.go:544: "ABC"
	lisp_test.go:544: (1 2 3 (4 5))
	lisp_test.go:544: (1 . 2)
PASS
ok  	command-line-arguments	0.013s
[kunohi@centos7-dev-docker src]
```

## ビルド方法
```
[kunohi@centos7-dev-docker src]$ go build lisp.go
[kunohi@centos7-dev-docker src]$ 
```

## emacsでの設定(例)
```
(setq scheme-program-name "~/bin/lisp") 
```
