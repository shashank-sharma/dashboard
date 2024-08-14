# Cache gogcc alpine
FROM --platform=$BUILDPLATFORM golang:1.21-alpine as gogcc

# Accept the target architecture as a build argument
ARG TARGETARCH

ENV GOOS=linux
ENV CGO_ENABLED=0

RUN apk update && apk add --no-cache \
        gcc \
        musl-dev

# Set the GOARCH environment variable based on the TARGETARCH build argument
ENV GOARCH=$TARGETARCH

# Build the binary
FROM gogcc as builder

WORKDIR /app

# Download dependencies
COPY go.mod ./
COPY go.sum ./

RUN go mod download

# Build /app/bin
COPY . .

RUN go build -ldflags="-s -w" -o bin -v ./cmd/dashboard/main.go

# Serve the binary with pb_public
FROM alpine:latest as bin

RUN apk update && apk add --no-cache \
        ca-certificates

WORKDIR /app/
COPY pb_public ./pb_public
COPY --from=builder /app/bin .

EXPOSE 8080

CMD ["/app/bin", "serve", "--http=0.0.0.0:8090"]
