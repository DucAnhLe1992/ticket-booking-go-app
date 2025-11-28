# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/orders ./cmd/orders

# Final stage
FROM alpine:3.18
RUN addgroup -S app && adduser -S -G app app
COPY --from=builder /bin/orders /bin/orders
USER app
EXPOSE 3000
ENTRYPOINT ["/bin/orders"]
