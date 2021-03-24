FROM ubuntu:18.04 as builder
MAINTAINER hidekuno@gmail.com

ENV HOME /root
RUN apt-get update && apt-get -y install git curl libgtk2.0-dev |true
RUN curl -O https://dl.google.com/go/go1.16.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.16.linux-amd64.tar.gz|true
ENV PATH ${PATH}:/usr/local/go/bin

WORKDIR $HOME
RUN git clone https://github.com/hidekuno/go-scheme && git clone https://github.com/hidekuno/picture-language

WORKDIR $HOME/go-scheme/scheme
RUN go build -o lisp -ldflags '-w -s' cmd/lisp/main.go

WORKDIR $HOME/go-scheme/draw
RUN go build -o glisp -ldflags '-w -s' cmd/lisp/main.go

FROM ubuntu:18.04 as go-scheme
MAINTAINER hidekuno@gmail.com

RUN apt-get update && apt-get -y install libgtk2.0-0
COPY --from=builder /root/go-scheme/scheme/lisp /root/
COPY --from=builder /root/go-scheme/draw/glisp /root/
COPY --from=builder /root/picture-language/sicp.scm /root/
COPY --from=builder /root/picture-language/sicp/ /root/picture-language/sicp/
