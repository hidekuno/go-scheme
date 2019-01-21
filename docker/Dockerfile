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
RUN mkdir go-scheme
ENV GOPATH ${HOME}/go-scheme:${HOME}/go
