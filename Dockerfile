# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.24 AS builder
WORKDIR /src

# Cache module downloads
COPY go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source and build binary
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/msisdn-lookup ./main.go

# Runtime stage
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates \
    && adduser -D -H -s /sbin/nologin appuser

COPY --from=builder /out/msisdn-lookup /usr/local/bin/msisdn-lookup
COPY --from=builder /src/lookup/rules.json /app/rules.json

ENV LOOKUP_RULES_PATH=/app/rules.json
EXPOSE 9090
USER appuser
ENTRYPOINT ["/usr/local/bin/msisdn-lookup"]
