FROM golang:1.21 AS builder

WORKDIR /backend
COPY go.* .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /backend/cmd/api/main cmd/api/main.go

FROM alpine:latest

COPY --from=builder /backend/cmd/api/main /backend/main

CMD ["/backend/main"]