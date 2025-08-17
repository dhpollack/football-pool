FROM golang:1.23

WORKDIR /app

COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

COPY . .

RUN go build -o /football-pool ./cmd/main.go
RUN go build -o /seed ./cmd/seed/main.go # Add this line

EXPOSE 8080

CMD ["/football-pool"]
