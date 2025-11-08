FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN GIN_MODE=release go test -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
