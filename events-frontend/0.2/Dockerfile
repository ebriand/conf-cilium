FROM golang:1.12.1

WORKDIR /go/frontend

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o events-frontend .

FROM scratch
COPY --from=0 /go/frontend/events-frontend /events-frontend
COPY templates templates
EXPOSE 80
ENTRYPOINT ["/events-frontend"]
