FROM golang:1.20.0-buster AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest && \
    mv $GOPATH/bin/grpcurl /app/bin/grpcurl

# Path: Dockerfile

FROM debian:buster-slim
WORKDIR /app
COPY --from=builder /app/bin/server /app/bin/server
COPY --from=builder /app/bin/grpcurl /app/bin/grpcurl
EXPOSE 3456
CMD ["/app/bin/server", "-port", "3456"]
