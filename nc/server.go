package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func ck(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

// check listening udp port on system
// % netstat -ul
//
// send udp request and read reply
// % echo "hello world" | nc localhost 8981 -u -w 1
//
// send udp request to remote host
// % ssh ide@localhost  "echo 'open aprilsh' | nc localhost 8981 -u -w 1"
func serve(port string) (done chan bool) {
	local_addr, err := net.ResolveUDPAddr("udp", port)
	ck(err)
	conn, err := net.ListenUDP("udp", local_addr)
	ck(err)
	buf := make([]byte, 65536)
	done = make(chan bool)

	fmt.Printf("listening on %s\n", port)
	go func() {
		defer conn.Close()
		for {
			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("Error: ", err)
				continue
			}
			fmt.Printf("Received %q from %s\n", strings.TrimSpace(string(buf[0:n])), addr)
			// rx <- buf[0:n]
			conn.WriteToUDP([]byte("#"), addr) // add response prefix
			conn.WriteToUDP(buf[0:n], addr)
		}
	}()
	return
}

func main() {
	done := serve(":8981")
	<-done
}
