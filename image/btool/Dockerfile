FROM golang:buster

RUN apt-get update && apt-get install gdb valgrind man manpages manpages-dev emacs -y

WORKDIR /etc/btool
COPY . /etc/btool
RUN go env && CGO_ENABLED=0 go build -o /usr/local/bin/btool ./cmd/btool
