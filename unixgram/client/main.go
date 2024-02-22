package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("unixgram", "/tmp/aprilsh_test")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// defer os.RemoveAll("/tmp/unixdomaincli")

	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("hello %d", i)
		_, err = conn.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
	}
}
