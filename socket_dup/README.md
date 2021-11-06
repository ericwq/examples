# Reference

- [The evolution of reuseport in the Linux kernel](https://programmer.group/the-evolution-of-reuseport-in-the-linux-kernel.html)
- [跨进程复制socket](https://cong.im/post/linux/%E8%B7%A8%E8%BF%9B%E7%A8%8B%E5%A4%8D%E5%88%B6socket/)
- [Linux网络编程“惊群”问题总结](https://www.cnblogs.com/Anker/p/7071849.html)
- [accept 与 epoll 惊群](https://pureage.info/2015/12/22/thundering-herd.html)

# Demo steps

Run the following command to build the binaraies.

- `go build -o demo main.go`

Run the following command to start the first server.

- `./demo`

You should get the following output:

```sh
$ ./demo 
listen on:  [::]:7000
Listening on /tmp/unix_socket_tcp
[49451] 2021-11-06 21:43:55.0125327 +0800 CST m=+4.472036001 recv msg is: hello
[49451] 2021-11-06 21:43:55.1211736 +0800 CST m=+4.580667201 recv msg is: hello
```

Run the following command to start the client. The client will keep sending `hello` message to localhost:7000

- `while true; do nc localhost 7000 -e echo hello; sleep 0.1; done`

Run the following command to start the second server.

- `./demo 2`

You should get the following output:

```sh
$ ./demo 2
recv fd 7
listen on:  [::]:7000
[49591] 2021-11-06 21:44:02.6632161 +0800 CST m=+1.219149001 recv msg is: hello
[49591] 2021-11-06 21:44:02.8883711 +0800 CST m=+1.444308201 recv msg is: hello
[49591] 2021-11-06 21:44:03.0006343 +0800 CST m=+1.556868001 recv msg is: hello
```

You can run more servers, if you like it.

```sh
$ ./demo 3
recv fd 7
listen on:  [::]:7000
[49778] 2021-11-06 21:44:13.1593171 +0800 CST m=+1.555800001 recv msg is: hello
[49778] 2021-11-06 21:44:13.6133832 +0800 CST m=+2.009869001 recv msg is: hello
[49778] 2021-11-06 21:44:13.717499 +0800 CST m=+2.113980301 recv msg is: hello
```

Now, watch the output from the client and 3 servers. The same socket is shared by 3 processes.
