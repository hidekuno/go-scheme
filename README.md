Go言語によるScheme(subset版)の実装
=================

## 開発環境
| Item   | Ver. |備考|
|--------|--------|--------|
| OS     | CentOS | draw系を使わなければ特になし|
| Gtk+   | 2.24.31||
| golang   | 1.9.4||
| go-gtk | release-0.1|https://github.com/mattn/go-gtk|

## インストール手順
### 環境変数のsetup
```
export GOARCH="amd64"
export GOOS="linux"
export GOPATH=${HOME}/go
export PATH=${GOPATH}/bin:$PATH
```
### go-gtkのinstall
```
go get github.com/mattn/go-gtk/gtk
go install github.com/mattn/go-gtk/gtk
```
### 本体のinstall
```
cd ${dokoka}
git clone https://github.com/hidekuno/go-scheme.git  go-scheme
```

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

## グラフィックス関連の使用方法
### 描画用Windowの起動
```
[kunohi@centos7-dev-docker src]$ ./lisp 
scheme.go>  (draw-init)
nil
scheme.go>  
```
### 描画用Windowのクリア
```
scheme.go>  (draw-clear)
nil
scheme.go>  
```

### 線を引く
```
scheme.go>  (draw-line 100 100 200 200)
nil
scheme.go>  
```

### コッホ曲線を描画するプログラム
```
(draw-init)
(define pi (*(atan 1)4))
(define (cs angle)(cos (/(* pi angle)180)))
(define (sn angle)(sin (/(* pi angle)180)))
(define koch 
  (lambda (x0 y0 x1 y1 c)
    (let ((kcos (cs 60))
          (ksin (sn 60)))
      (if (> c 1)
          (let (
                (xa (/ (+ (* x0 2) x1) 3))
                (ya (/ (+ (* y0 2) y1) 3))
                (xb (/ (+ (* x1 2) x0) 3))
                (yb (/ (+ (* y1 2) y0) 3)))
            (let ((yc (+ ya (+ (* (- xb xa) ksin) (* (- yb ya) kcos))))
                  (xc (+ xa (- (* (- xb xa) kcos) (* (- yb ya) ksin)))))
              (koch x0 y0 xa  ya (- c 1))
              (koch xa ya xc  yc (- c 1))
              (koch xc yc xb  yb (- c 1))
              (koch xb yb x1  y1 (- c 1))))
          (draw-line x0 y0 x1 y1)))))
```
#### 実行例
```
(koch 259 0 34 390 4)
(koch 34 390 483 390 4)
(koch 483 390 259 0 4)
```

![image](https://user-images.githubusercontent.com/4899700/42983927-89247a4c-8c24-11e8-82e7-5c2ac3f47e37.png)

### ツリーカーブ曲線を描画するプログラム
```
(define (tree x0 y0 x1 y1 c)
  (let ((tcos (cs 15))
        (tsin (sn 45))
        (alpha 0.6))
    (let ((ya (+ y1  (*    tsin (- x1 x0) alpha) (*    tcos (- y1 y0) alpha)))
          (xa (+ x1  (*    tcos (- x1 x0) alpha) (* -1 tsin (- y1 y0) alpha)))
          (yb (+ y1  (* -1 tsin (- x1 x0) alpha) (*    tcos (- y1 y0) alpha)))
          (xb (+ x1  (*    tcos (- x1 x0) alpha) (*    tsin (- y1 y0) alpha))))
      (draw-line x0 y0 x1 y1)
      (if (>= 0 c)
          ((lambda () (draw-line x1 y1 xa ya) (draw-line x1 y1 xb yb)))
          ((lambda () (tree x1 y1 xa ya (- c 1))(tree x1  y1  xb  yb (- c 1))))))))
```
#### 実行例
```
(tree 300 400 300 300 12)
```

![image](https://user-images.githubusercontent.com/4899700/42988528-dfc3149a-8c37-11e8-8b72-0d8afe921ac3.png)

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
