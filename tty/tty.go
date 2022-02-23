// Companion code for the Linux terminals blog series: https://dev.to/napicella/linux-terminals-tty-pty-and-shell-192e
// I have simplified the code to highlight the interesting bits for the purpose of the blog post:
// - windows resizing is not addressed
// - client does not catch signals (CTRL + C, etc.) to gracefully close the tcp connection
// 
// Build: go build -o remote main.go
// In one terminal run: ./remote -server
// In another terminal run: ./remote 
// 
// Run on multiple machines:
// In the client function, replace the loopback address with IP of the machine, then rebuild
// Beware the unecrypted TCP connection!
package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"net"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

var isServer *bool

func init()  {
	isServer = flag.Bool("server", false, "")
}

func server() error {
	// Create command
	c := exec.Command("bash")

	// Start the command with a pty.
	ptmx, e := pty.Start(c)
	if e != nil {
		return e
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	return listen(ptmx)
}

func listen(ptmx *os.File) error {
	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, e := net.Listen("tcp", ":8081")
	if e != nil {
		return e
	}
	// accept connection on port
	conn, e := ln.Accept()
	if e != nil {
		return e
	}

	go func() { _, _ = io.Copy(ptmx, conn) }()
	_, e = io.Copy(conn, ptmx)
	return e
}

func client() error {
	// connect to this socket
	conn, e := net.Dial("tcp", "127.0.0.1:8081")
	if e != nil {
		return e
	}

	// MakeRaw put the terminal connected to the given file descriptor into raw
	// mode and returns the previous state of the terminal so that it can be
	// restored.
	oldState, e := terminal.MakeRaw(int(os.Stdin.Fd()))
	if e != nil {
		return e
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.


	go func() { _, _ = io.Copy(os.Stdout, conn) }()
	_, e = io.Copy(conn, os.Stdin)
	fmt.Println("Bye!")

	return e
}

func clientAndServer() error {
	flag.Parse()
	if isServer != nil && *isServer {
		fmt.Println("Starting server mode")
		return server()
	} else {
		fmt.Println("Starting client mode")
		return client()
	}
}

func main() {
	if e := clientAndServer(); e != nil {
		fmt.Println(e)
	}
}
