FROM golang:1.24-alpine AS builder

WORKDIR /src/app

COPY app/go.mod app/go.sum* ./
RUN go mod download

COPY app/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app ./cmd/app

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /out/app ./app

EXPOSE 8080

CMD ["./app"]
