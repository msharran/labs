FROM golang:1.22.0 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o tmp-file-exporter .

FROM debian:stretch-slim
WORKDIR /app
COPY --from=builder /app/tmp-file-exporter /app/
EXPOSE 9900
CMD ["/app/tmp-file-exporter"]
