FROM golang:1.12.1

WORKDIR /go/frontend

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o internal-frontend .

FROM scratch
COPY --from=0 /go/frontend/internal-frontend /internal-frontend
COPY templates templates
COPY assets assets
EXPOSE 80
ENTRYPOINT ["/internal-frontend"]
