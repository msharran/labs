FROM golang:1.20.0-buster AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make

# Path: Dockerfile

FROM debian:buster-slim
WORKDIR /app
COPY --from=builder /app/bin/server /app/bin/server
COPY --from=builder /app/bin/client /app/bin/health-check-probe
EXPOSE 3456
CMD ["/app/bin/server", "-port", "3456"]
