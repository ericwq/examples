# Reference

- [Graceful shutdown of a TCP server in Go](https://eli.thegreenplace.net/2020/graceful-shutdown-of-a-tcp-server-in-go/)

## Demo steps

Run the following command to build the binaraies.

- `go build -o demo tcp.go`

Run the following command to start the server.

- `./demo`

Run the following command to start the client.

- `while true; do nc localhost 8080 -e echo hello; sleep 0.1; done`

The server side log is as following.

```sh
$ ./demo
[10185] 2021/11/11 23:34:59 listening on  127.0.0.1:8080
[10185] 2021/11/11 23:35:38 receive message from [127.0.0.1:43393]:[hello]
[10185] 2021/11/11 23:35:38 receive message from [127.0.0.1:42929]:[hello]
[10185] 2021/11/11 23:35:38 receive message from [127.0.0.1:33273]:[hello]
```

Run the following command to restart the server. Note that the [PID] should be replaced with the previous PID.

- `kill -s HUP [PID]`

Watch the output of server side log. You should notice that the server doesn't report any error during the restart.

```sh
[10185] 2021/11/11 23:37:00 receive message from [127.0.0.1:34589]:[hello]
[10185] 2021/11/11 23:37:00 receive message from [127.0.0.1:36367]:[hello]
[10185] 2021/11/11 23:37:00 receive from exitC channel.
[10573] 2021/11/11 23:37:00 listening on  127.0.0.1:8080
[10185] 2021/11/11 23:37:00 quit the listening.
[10185] 2021/11/11 23:37:00 stop listening.
[10185] 2021/11/11 23:37:00 finish the old process.
[10573] 2021/11/11 23:37:00 receive message from [127.0.0.1:37935]:[hello]
[10573] 2021/11/11 23:37:00 receive message from [127.0.0.1:34721]:[hello]
[10573] 2021/11/11 23:37:00 receive message from [127.0.0.1:34245]:[hello]
```

Run the following command to stop the server. Note that the [PID] should be replaced with the right PID.

- `kill -s QUIT [PID]`

Watch the output of server side log.

```sh
[10573] 2021/11/11 23:39:55 receive message from [127.0.0.1:44181]:[hello]
[10573] 2021/11/11 23:39:55 receive message from [127.0.0.1:37603]:[hello]
[10573] 2021/11/11 23:39:55 receive message from [127.0.0.1:33321]:[hello]
[10573] 2021/11/11 23:39:55 got message SIGQUIT.
[10573] 2021/11/11 23:39:55 receive from exitC channel.
[10573] 2021/11/11 23:39:55 quit the listening.
[10573] 2021/11/11 23:39:56 stop listening.
[10573] 2021/11/11 23:39:56 finish the old process.
```

## Understand tableflip code

The key design of tableflip is the parent process shares the same sockets withe the child process.

![tableflip.001.png](images/tableflip.001.png)
