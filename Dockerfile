# Multi-stage build
FROM golang:1.23.0 AS builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /app/server /server
COPY .env .env
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
