package main

import (
	"flag"
	"fmt"
	"log"
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
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGUSR2)
		for s := range sig {
			switch s {
			case syscall.SIGHUP:
				//log.Println("got message SIGHUP.")
				err := upg.Upgrade()
				if err != nil {
					log.Println("upgrade failed:", err)
				}
			case syscall.SIGUSR2:
				//log.Println("got message SIGUSER2.")
				upg.Stop()
				return
			}
		}
	}()

	ln, err := upg.Fds.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalln("Can't listen:", err)
	}

	var wg sync.WaitGroup
	quit := make(chan struct{})

	wg.Add(1)
	go func() {
		defer func() {
			//log.Println("stop listening.")
			ln.Close()
			wg.Done()
		}()

		log.Println("listening on ", ln.Addr())

		for {
			select {
			case <-quit:
				return
			default:
			}
			c, err := ln.Accept()
			if err != nil {
				log.Println("Accept return error, stop accept. err:", err)
				return
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				c.SetDeadline(time.Now().Add(time.Second))

				//read message from the client and print it on screen
				var buf []byte = make([]byte, 32)
				_, err := c.Read(buf)
				if err != nil {
					log.Println("read from connection error:", err)
				} else {
					msg := strings.ReplaceAll(string(buf), "\n", "")
					log.Printf("receive message:[%s]", msg)
				}
				//c.Write([]byte("It is a mistake to think you can solve any major problems just with potatoes.\n"))
				c.Close()
			}()
		}
	}()

	if err := upg.Ready(); err != nil {
		panic(err)
	}
	<-upg.Exit()
	//log.Println("receive from exitC channel.")

	close(quit)
	//log.Println("quit the listening.")

	wg.Wait()
	log.Println("finish the old process.")

	//log.Println("Exit() pause for a moment.")
	//time.Sleep(time.Second * 2)
}
