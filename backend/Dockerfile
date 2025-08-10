FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

COPY . .

RUN go build -o /football-pool ./cmd/main.go

EXPOSE 8080

CMD ["/football-pool"]
