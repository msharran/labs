# Readme

Expose ingress

```bash
$ k port-forward svc/external-ingress-nginx-chart-c05923f4-controller 8080:80
Forwarding from 127.0.0.1:8080 -> 80
Forwarding from [::1]:8080 -> 80
Handling connection for 8080
Handling connection for 8080
Handling connection for 8080
Handling connection for 8080
```

Hit from random host -> 404

```bash
$ curl http://127.0.0.1:8080
<html>
<head><title>404 Not Found</title></head>
<body>
<center><h1>404 Not Found</h1></center>
<hr><center>nginx</center>
...
```


Hit from foo.example.com -> nginx welcome page

```bash
$ curl -H "Host: foo.example.com" http://127.0.0.1:8080
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
...
```
