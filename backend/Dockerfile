FROM golang:1.21 AS builder

WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /backend cmd/api/main.go

FROM alpine:latest

COPY --from=builder /backend /backend

CMD ["/backend"]