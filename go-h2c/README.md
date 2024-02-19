# H2C (HTTP2 over clear text) to serve both gRPC and HTTP requests

## Sample gRPC and http requests

```bash
sharranm@2184-X1 ~/p/p/l/go-h2c> grpcurl -plaintext -d '{"name":"sharran"}' localhost:3000 helloworld.Greeter.SayHello
{
  "message": "Hello, sharran"
}
sharranm@2184-X1 ~/p/p/l/go-h2c> curl http://localhost:3000/hello/sharran
Hello, sharran!%
```

## Server output

```bash
sharranm@2184-X1 ~/p/p/l/go-h2c> go run main.go                                                                                                                 main!
2024/02/19 14:38:07 GET /hello/sharran
2024/02/19 14:38:26 Received: sharran
2024/02/19 14:38:38 Received: sharran
2024/02/19 14:38:47 GET /hello/sharran
```
