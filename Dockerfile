FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main", "server", "-p", "3000", "-a", "0.0.0.0"]