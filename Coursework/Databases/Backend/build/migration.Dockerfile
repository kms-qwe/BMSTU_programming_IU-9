FROM golang:1.24-alpine AS builder

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

FROM alpine:3.21

WORKDIR /migrations

COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY migrations /migrations
COPY build/migration.sh /migration.sh

RUN chmod +x /migration.sh

ENTRYPOINT ["/migration.sh"]
