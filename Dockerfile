FROM centos:centos7
MAINTAINER hidekuno@gmail.com

RUN yum update -y|true
RUN yum install -y epel-release |true
RUN yum install -y golang golang-bin git gtk2-devel |true
ENV HOME /root
ENV GOPATH ${HOME}/go
ENV PATH ${GOPATH}/bin:${PATH}
WORKDIR $HOME
RUN go get github.com/mattn/go-gtk/gtk
RUN go install github.com/mattn/go-gtk/gtk
RUN go get github.com/gorilla/sessions
RUN go install github.com/gorilla/sessions
RUN go get golang.org/x/net/websocket
RUN go install golang.org/x/net/websocket
RUN git clone https://github.com/hidekuno/go-scheme
ENV GOPATH ${HOME}/go-scheme:${HOME}/go
WORKDIR $HOME/go-scheme/src
RUN go build  -ldflags '-w -s' lisp_main.go
RUN go build  -ldflags '-w -s' lisp_draw_main.go
RUN go build  -ldflags '-w -s' web_api_main.go
RUN go build  -ldflags '-w -s' web_wasm_main.go
RUN GOARCH=wasm GOOS=js go build -o wasm/lisp.wasm lisp_wasm.go
