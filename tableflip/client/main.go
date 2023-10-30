package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	host = flag.String("host", "localhost", "host")
	port = flag.String("port", "8080", "port")
)

// The client is used to test the long lived socket case:
// long lived socket: opened by client, client use it to communicate with server for several times.
// long lived socket has some requirement for zero down time upgrade:
// - the client should check the network error and reopen the connection if necessary.
// - the server should stop response at application boundary defined by developer.
//
// nc client is used to test the short lived socket case:
// short lived socket: opened by client, read/write to it and close it.
// short lived socket has no requirement for zero down time upgrade.
func main() {
	flag.Parse()
	// rand.Seed(42)
	// var wg sync.WaitGroup

	// short connection
	if len(os.Args) > 1 {
		for {
			conn, err := net.Dial("tcp", *host+":"+*port)
			if err != nil {
				fmt.Println("Error connecting:", err)
				os.Exit(1)
			}
			// fmt.Println("Connecting to " + *host + ":" + *port)

			// wg.Add(1)
			for i := 0; i < 10; i++ {
				handleWrite(conn, i, nil)
				time.Sleep(time.Millisecond * 1)
				handleRead(conn, nil)
			}
			conn.Close()
		}
		// wg.Wait()
	} else {
		fmt.Println("no argument.")
	}
}

func handleWrite(conn net.Conn, idx int, wg *sync.WaitGroup) error {
	// defer wg.Done()

	_, err := conn.Write([]byte("hello " + strconv.Itoa(idx) + "\n"))
	if err != nil {
		fmt.Println("W: ", err)
		return err
	}

	return nil
}

func handleRead(conn net.Conn, wg *sync.WaitGroup) error {
	// defer wg.Done()
	reader := bufio.NewReader(conn)
	// for i := 1; i <= 10; i++ {
	_, err := reader.ReadString(byte('\n'))
	if err != nil {
		fmt.Println("R: ", err)
		return err
	}
	// fmt.Print(line)
	// }
	return nil
}
