FROM golang:1.16-alpine AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go get ./...
COPY cmd/api/*.go ./
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o /api

FROM alpine:latest
COPY --from=build /api /api
EXPOSE 8080
CMD [ "/api" ]
