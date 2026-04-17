FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go build -buildvcs=false -o gateway ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/gateway /app/gateway
COPY --from=builder /app/configs /app/configs
EXPOSE 8080
ENTRYPOINT ["/app/gateway", "-config", "/app/configs/config.example.yaml"]
