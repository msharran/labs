# go-grpc

This is a POC for gRPC greeter service and client http
service that calls the gRPC server

## Installation 

Installation follows helm based k8s deployment. So
assuming a k8s cluster has been setup. you can 
spin up a local one based on docker using k3d or kind. 

The following instructions assume you are connected to the 
cluster and a namespace.

Run the following command to install both the
server and client

```bash
helm install go-grpc ./go-grpc

watch kubectl get pods
```

## Development

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
