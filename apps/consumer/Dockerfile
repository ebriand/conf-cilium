FROM golang:1.11 as builder

WORKDIR /usr/src/app
COPY main.go /usr/src/app

RUN go mod init consumer

RUN go build .

ENTRYPOINT [ "/usr/src/app/consumer" ]
