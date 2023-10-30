package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
)

func main() {
	Example_tcpServer()
}

// This shows how to use the Upgrader
// with a listener based service.
func Example_tcpServer() {
	var (
		listenAddr = flag.String("listen", "localhost:8080", "`Address` to listen on")
		pidFile    = flag.String("pid-file", "", "`Path` to pid file")
	)

	flag.Parse()
	log.SetPrefix(fmt.Sprintf("[%d] ", os.Getpid()))

	/*
		var lc = &net.ListenConfig{
			Control: func(network, address string, c syscall.RawConn) error {
				var opErr error
				if err := c.Control(func(fd uintptr) {
					opErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
				}); err != nil {
					return err
				}
				return opErr
			},
		}
	*/

	upg, err := tableflip.New(tableflip.Options{
		PIDFile: *pidFile,
	})
	if err != nil {
		panic(err)
	}
	defer upg.Stop()

	// Do an upgrade on SIGHUP
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGQUIT)
		for s := range sig {
			switch s {
			case syscall.SIGHUP:
				// log.Println("got message SIGHUP.")
				err := upg.Upgrade()
				if err != nil {
					log.Println("upgrade failed:", err)
				}
			case syscall.SIGQUIT:
				log.Println("got message SIGQUIT.")
				upg.Stop()
				return
			}
		}
	}()

	ln, err := upg.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalln("Can't listen:", err)
	}

	var wg sync.WaitGroup
	quit := make(chan struct{})

	wg.Add(1)
	go func() {
		defer func() {
			log.Println("stop listening.")
			ln.Close()
			wg.Done()
		}()

		log.Println("listening on ", ln.Addr())
		lis := ln.(*net.TCPListener)

		for {
			select {
			case <-quit:
				return
			default:
				// set listener accept time out: 2 second
				lis.SetDeadline(time.Now().Add(time.Second * 2))
				conn, err := lis.Accept()
				if err != nil {
					if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
						// log.Println("accept time out")
					} else {
						log.Println("accept error", err)
					}
				} else {
					wg.Add(1)
					go func() {
						handleConnection(conn, quit)
						wg.Done()
					}()
				}
			}
		}
	}()

	if err := upg.Ready(); err != nil {
		panic(err)
	}
	<-upg.Exit()
	log.Println("receive from exitC channel.")

	close(quit)
	log.Println("quit the listening.")

	wg.Wait()
	log.Println("finish the old process.")

	// log.Println("Exit() pause for a moment.")
	// time.Sleep(time.Second * 2)
}

func handleConnection(conn net.Conn, quit chan struct{}) {
	defer conn.Close()

	var buf []byte = make([]byte, 32)
	// log.Println("handle connection")
	count := 10
ReadLoop:
	for i := 0; i < count; i++ {
		select {
		case <-quit:
			break ReadLoop
		default:

			// set read time out
			conn.SetDeadline(time.Now().Add(200 * time.Millisecond))

			// read message from the client and print it on screen
			n, err := conn.Read(buf)
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					// log.Println("read time out")
					continue ReadLoop
				} else if err != io.EOF {
					log.Println("read error", err)
					return
				}
			}

			if n == 0 {
				log.Printf("done connection.")
				return
			}
			msg := strings.ReplaceAll(string(buf), "\n", "")
			conn.Write(buf)
			log.Printf("receive message from [%s]:[%s]", conn.RemoteAddr().String(), msg)
			// log.Printf("receive message from [%s]:[%s - modified]", conn.RemoteAddr().String(), msg)
		}
	}
	log.Printf("done w/ message from [%s]", conn.RemoteAddr().String())
}
