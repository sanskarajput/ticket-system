# Build stage
FROM golang:1.26.4-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o ticket-system ./cmd/main.go

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache sqlite-libs ca-certificates

WORKDIR /app

COPY --from=builder /app/ticket-system .

EXPOSE 8080

RUN mkdir -p /app/database

CMD ["./ticket-system"]
