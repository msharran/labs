# go-grpc

This is a POC for gRPC greeter service and client using `go`

## Usage

1. Start the gRPC server by running the following command

```
$ make serve
```

2. Say hello using gRPC client. Open a new terminal and run the below command

```bash
$ make say_hello
```

You will get a response like this

```
reply: message:"Hello, Sharran"
```
