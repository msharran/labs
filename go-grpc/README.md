# go-grpc

This is a POC for gRPC greeter service and client http
service that calls the gRPC server

## Installation 

Installation follows helm based k8s deployment. So
assuming a k8s cluster has been setup. you can 
spin up a local one based on docker using k3d or kind. 

The following instructions assume you are connected to the 
cluster and a namespace.

You can install using both Kustomize and Helm.
Kustomize is recommended

### Using Kustomize (RECOMMENDED)

Install

```bash
make ./kustomize
```

Uninstall

```bash
make ./kustomize ACTION=uninstall
```

### Using Helm

Install

```bash
make helm-install
```

Uninstall

```bash
make helm-uninstall 
```
