FROM golang:1.12.1

WORKDIR /go/frontend

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o events-frontend .

FROM debian:9
RUN apt update && apt install -y curl
COPY --from=0 /go/frontend/events-frontend /events-frontend
COPY templates templates
COPY assets assets
EXPOSE 80
ENTRYPOINT ["/events-frontend"]
