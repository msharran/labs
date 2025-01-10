# Redis Server

## Start server

```sh
zig run src/main.zig
```

## Test server

```sh
echo -n '*1\r\n$4\r\nPING\r\n' | nc localhost 6379
```


![ScreenShot](./images/ping.png)
