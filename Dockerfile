FROM ubuntu:18.04 as builder
MAINTAINER hidekuno@gmail.com

ENV HOME /root
RUN apt-get update && apt-get -y install git curl libgtk2.0-dev |true
RUN curl -O https://dl.google.com/go/go1.13.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.13.linux-amd64.tar.gz|true
ENV GOPATH ${HOME}/go
ENV PATH ${PATH}:/usr/local/go/bin:${GOPATH}/bin
RUN go get github.com/mattn/go-gtk/gtk && go install github.com/mattn/go-gtk/gtk

WORKDIR $HOME
ENV GOPATH ${HOME}/go
RUN go get github.com/hidekuno/go-scheme && git clone https://github.com/hidekuno/picture-language

WORKDIR $GOPATH/github.com/hidekuno/go-scheme
RUN go build -o lisp -ldflags '-w -s' cmd/lisp/main.go && go build -o glisp -ldflags '-w -s' cmd/lisp/draw/main.go

FROM ubuntu:18.04 as go-scheme
MAINTAINER hidekuno@gmail.com

RUN apt-get update && apt-get -y install libgtk2.0-0
COPY --from=builder /root/go/src/github.com/hidekuno/go-scheme/lisp /root/
COPY --from=builder /root/go/src/github.com/hidekuno/go-scheme/glisp /root/
COPY --from=builder /root/picture-language/z-learning/go.scm /root/
COPY --from=builder /root/picture-language/fractal/ /root/picture-language/fractal/
COPY --from=builder /root/picture-language/sicp/ /root/picture-language/sicp/
