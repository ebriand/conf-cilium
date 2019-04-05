FROM golang:1.12.1

WORKDIR /go/frontend

COPY go.mod go.sum ./
RUN go mod download

COPY main.go .
RUN CGO_ENABLED=0 go build -o internal-frontend .

FROM scratch
COPY --from=0 /go/frontend/internal-frontend /internal-frontend
EXPOSE 80
ENTRYPOINT ["/internal-frontend"]
