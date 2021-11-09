package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
)

// 当前程序的版本
const version = "v0.0.1"

func main() {
	// 试试SO_REUSEPORT 的效果
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
		//ListenConfig: lc,
	})
	if err != nil {
		panic(err)
	}
	defer upg.Stop()

	// 为了演示方便，为程序启动强行加入 1s 的延时，并在日志中附上进程 pid
	time.Sleep(time.Second)
	log.SetPrefix(fmt.Sprintf("[PID: %d] ", os.Getpid()))

	// 监听系统的 SIGHUP 信号，以此信号触发进程重启
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP)
		for range sig {
			// 核心的 Upgrade 调用
			err := upg.Upgrade()
			if err != nil {
				log.Println("Upgrade failed:", err)
			}
		}
	}()

	// 注意必须使用 upg.Listen 对端口进行监听
	ln, err := upg.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Can't listen:", err)
	}

	// 创建一个简单的 http server，/version 返回当前的程序版本
	mux := http.NewServeMux()
	mux.HandleFunc("/version", func(rw http.ResponseWriter, r *http.Request) {
		log.Println(version)
		rw.Write([]byte(version + "\n"))
	})
	server := http.Server{
		Handler: mux,
	}

	// 照常启动 http server
	go func() {
		err := server.Serve(ln)
		if err != http.ErrServerClosed {
			log.Println("HTTP server:", err)
		}
	}()

	if err := upg.Ready(); err != nil {
		panic(err)
	}

	<-upg.Exit()

	// 给老进程的退出设置一个 30s 的超时时间，保证老进程的退出
	time.AfterFunc(30*time.Second, func() {
		log.Println("Graceful shutdown timed out")
		os.Exit(1)
	})

	// 等待 http server 的优雅退出
	server.Shutdown(context.Background())
}
