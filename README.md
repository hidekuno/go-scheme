Go言語によるScheme(subset版)の実装
=================

## 起動方法
### scheme単体
```
[kunohi@centos7-dev-docker src]$ go run lisp_main.go lisp.go 
scheme.go>
```
### scheme(グラフィックス処理付き)
```
[kunohi@centos7-dev-docker src]$ go run lisp_main_draw.go lisp.go draw.go
scheme.go>
```
### グラフィックス単体
```
[kunohi@centos7-dev-docker src]$ go run draw_main.go  draw.go
```

## 終了方法
```
scheme.go>  (quit)
[kunohi@centos7-dev-docker src]$ 
```

## テスト方法
### scheme単体
```
[kunohi@centos7-dev-docker src]$ go test -v lisp.go lisp_test.go
=== RUN   Test_lisp_sample_program
--- PASS: Test_lisp_sample_program (0.01s)
=== RUN   Test_math_func
--- PASS: Test_math_func (0.00s)
=== RUN   Test_list_func
--- PASS: Test_list_func (0.00s)
=== RUN   Test_basic_operation
--- PASS: Test_basic_operation (0.00s)
=== RUN   Test_err_case
--- PASS: Test_err_case (0.00s)
=== RUN   Test_interactive
--- PASS: Test_interactive (0.00s)
PASS
ok  	command-line-arguments	0.010s
[kunohi@centos7-dev-docker src]
```
### scheme グラフィックス単体
```
[kunohi@centos7-dev-docker src]$ go test -v lisp_main_draw.go lisp_main_draw_test.go  lisp.go
=== RUN   Test_draw
--- PASS: Test_draw (0.00s)
PASS
ok  	command-line-arguments	0.002s
[kunohi@centos7-dev-docker src]$ 
```

## ビルド方法
### scheme単体
```
[kunohi@centos7-dev-docker src]$ go build  -ldflags '-w -s' lisp.go lisp_main.go
[kunohi@centos7-dev-docker src]$ 
```
### scheme(グラフィックス処理付き)
```
[kunohi@centos7-dev-docker src]$ go build  -ldflags '-w -s'  lisp.go draw.go lisp_main_draw.go 
[kunohi@centos7-dev-docker src]$ 
```

### グラフィックス単体
```
[kunohi@centos7-dev-docker src]$ go build  -ldflags '-w -s'  lisp.go draw.go lisp_main_draw.go 
[kunohi@centos7-dev-docker src]$ 
```


## emacsでの設定(例)
```
(setq scheme-program-name "~/bin/lisp") 
```

### (注意点)環境の違いにより、"go test"が失敗する

Linux 32bit
```
[hideki@gentoo src]$ go run lisp.go 
scheme.go>  (atan 1)
0.7853981633974483
scheme.go>  (define pi (* 4 (atan 1)))
pi
scheme.go>  pi
3.141592653589793
scheme.go>  (tan (/ (* 45 pi) 180))
0.9999999999999999
scheme.go>  (quit)
[hideki@gentoo src]$ uname -a
Linux gentoo.mukogawa.or.jp 2.6.18-411.el5 #1 SMP Mon Jul 11 18:16:41 CDT 2016 i686 i686 i386 GNU/Linux
[hideki@gentoo src]$ 
```

Linux 64bit, MacOSX 64bit
```
macbookair:src hideki$ go run lisp.go 
scheme.go>  (define pi (* 4 (atan 1)))
pi
scheme.go>  pi
3.141592653589793
scheme.go>  (tan (/ (* 45 pi) 180))
1
scheme.go>  (quit)
macbookair:src hideki$ uname -a
Darwin macbookair.local 17.6.0 Darwin Kernel Version 17.6.0: Tue May  8 15:22:16 PDT 2018; root:xnu-4570.61.1~1/RELEASE_X86_64 x86_64
macbookair:src hideki$ 
```
