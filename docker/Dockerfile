FROM centos:centos7
MAINTAINER hidekuno@gmail.com

RUN yum update -y|true
RUN yum install -y golang git gtk2-devel |true
RUN go get github.com/mattn/go-gtk/gtk
RUN go install github.com/mattn/go-gtk/gtk
ENV HOME /root
ENV GOPATH ${HOME}/go
ENV PATH ${GOPATH}/bin:$PATH
RUN go get github.com/mattn/go-gtk/gtk
RUN go install github.com/mattn/go-gtk/gtk
RUN go get github.com/gorilla/sessions
RUN go install github.com/gorilla/sessions
WORKDIR $HOME
RUN git clone https://github.com/hidekuno/go-scheme.git go-scheme
ENV GOPATH ${HOME}/go-scheme:${HOME}/go
