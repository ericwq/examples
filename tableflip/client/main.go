package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var host = flag.String("host", "localhost", "host")
var port = flag.String("port", "8080", "port")

func main() {
	flag.Parse()
	rand.Seed(42)
	//var wg sync.WaitGroup

	// short connection
	if len(os.Args) > 1 {
		for {
			conn, err := net.Dial("tcp", *host+":"+*port)
			if err != nil {
				fmt.Println("Error connecting:", err)
				os.Exit(1)
			}
			defer conn.Close()
			//fmt.Println("Connecting to " + *host + ":" + *port)

			//wg.Add(1)
			handleWrite(conn, nil)
			//go handleRead(conn, &wg)
			time.Sleep(time.Millisecond * 100)
		}
		//wg.Wait()
	} else {
		fmt.Println("no argument.")

	}
}

func handleWrite(conn net.Conn, wg *sync.WaitGroup) {
	//defer wg.Done()

	_, e := conn.Write([]byte("hello " + strconv.Itoa(rand.Intn(9)) + "\n"))
	if e != nil {
		fmt.Println("Error to send message because of ", e.Error())
	}
}

func handleRead(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	reader := bufio.NewReader(conn)
	for i := 1; i <= 10; i++ {
		line, err := reader.ReadString(byte('\n'))
		if err != nil {
			fmt.Print("Error to read message because of ", err)
			return
		}
		fmt.Print(line)
	}
}
