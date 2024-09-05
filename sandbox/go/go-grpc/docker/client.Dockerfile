FROM golang:1.20.0-buster AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make

# Path: Dockerfile

FROM debian:buster-slim
WORKDIR /app
COPY --from=builder /app/bin/client /app/bin/client
EXPOSE 3455
ENTRYPOINT ["/app/bin/client"]
CMD ["-port", "3455"]
