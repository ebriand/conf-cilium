FROM golang:1.12.1

WORKDIR /go/api

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o heroes-api .

FROM scratch
COPY --from=0 /go/api/heroes-api /heroes-api
EXPOSE 80
ENTRYPOINT ["/heroes-api"]
