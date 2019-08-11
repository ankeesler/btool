FROM golang:buster
WORKDIR /etc/btool
COPY . /etc/btool
RUN go env && CGO_ENABLED=0 go build -o /usr/local/bin/btool ./cmd/btool
