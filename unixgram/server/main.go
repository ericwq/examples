package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"
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

		conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(2)))
		n, err := conn.Read(buf[:])
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				slog.Info("uxServe read timeout.")
				continue
			}
		}
		fmt.Printf("got from client: %s\n", string(buf[:n]))
	}
}
