package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func ck(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func main() {
	server_addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8981")
	ck(err)

	local_addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	ck(err)

	conn, err := net.DialUDP("udp", local_addr, server_addr)
	ck(err)

	defer conn.Close()
	i := 0
	for {
		msg := strconv.Itoa(i)
		i++
		txbuf := []byte("hello" + msg)
		_, err := conn.Write(txbuf)
		if err != nil {
			fmt.Println(msg, err)
		}
		time.Sleep(time.Second * 1)
		rxbuf := make([]byte, 65536)
		n, addr, err := conn.ReadFromUDP(rxbuf)
		fmt.Println("Received ", string(rxbuf[0:n]), " from ", addr)
	}
}
