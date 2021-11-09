# Reference

-[Graceful shutdown of a TCP server in Go](https://eli.thegreenplace.net/2020/graceful-shutdown-of-a-tcp-server-in-go/)

# Demo steps

Run the following command to build the binaraies.

- `go build -o demo tcp.go`

Run the following command to start the server.

- `./demo`

Run the following command to start the client.

- `while true; do nc localhost 8080 -e echo hello; sleep 0.1; done`

The server side log is as following.

```sh
$ ./demo
[30377] 2021/11/09 13:22:43 listening on  127.0.0.1:8080
[30377] 2021/11/09 13:22:43 receive message:[hello]
[30377] 2021/11/09 13:22:43 receive message:[hello]
[30377] 2021/11/09 13:22:43 receive message:[hello]
[30377] 2021/11/09 13:22:43 receive message:[hello]
[30377] 2021/11/09 13:22:44 receive message:[hello]
```

Run the following command to restart the server. Note that the [PID] should be replaced with the previous PID.

- `kill -s HUP [PID]`

Watch the output of server side log. You should notice that the server doesn't report any error during the restart.

```sh
[30377] 2021/11/09 13:23:01 receive message:[hello]
[30377] 2021/11/09 13:23:01 receive message:[hello]
[30377] 2021/11/09 13:23:01 receive message:[hello]
[30720] 2021/11/09 13:23:01 listening on  127.0.0.1:8080
[30377] 2021/11/09 13:23:01 return from main goroutine.
[30377] 2021/11/09 13:23:01 receive message:[hello]
[30377] 2021/11/09 13:23:01 finish the old process.
ide@golangide:~/proj/examples/tableflip $ [30720] 2021/11/09 13:23:01 receive message:[hello]
[30720] 2021/11/09 13:23:01 receive message:[hello]
[30720] 2021/11/09 13:23:01 receive message:[hello]
[30720] 2021/11/09 13:23:02 receive message:[hello]
[30720] 2021/11/09 13:23:02 receive message:[hello]

```
