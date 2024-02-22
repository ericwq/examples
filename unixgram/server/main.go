package main

import (
	"fmt"
	"net"
	"os"
)

const sockAddr = "/tmp/aprilsh_test"
const network = "unixgram"

func main() {

	addr, err := net.ResolveUnixAddr(network, sockAddr)
	conn, err := net.ListenUnixgram(network, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	defer os.RemoveAll(sockAddr)

	for i := 0; i < 5; i++ {
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			panic(err)
		}
		fmt.Printf("got from client: %s\n", string(buf[:n]))
	}
}
