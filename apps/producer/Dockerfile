FROM golang:1.11 as builder

RUN git clone https://github.com/edenhill/librdkafka.git
RUN cd librdkafka && \
  ./configure --prefix /usr && \
  make && \
  make install

WORKDIR /usr/src/app
COPY main.go /usr/src/app

RUN go mod init producer

RUN go build .

ENTRYPOINT [ "/usr/src/app/producer" ]
