FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./cmd/http/main.go ./

RUN go build -o /http-server

EXPOSE 8080

CMD ["/http-server"]
