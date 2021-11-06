package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// 掩饰了，复制的listen socket 可以并行，轮流响应客户的请求。
// 没有参数时，启动一个TCP服务，同时打开UnixSocket来监听请求
// 带参数时，从UnixSocket获取信息，进行处理：启动TCP服务

func main() {
	tcpSrv := NewTcpSrv()
	if len(os.Args) <= 1 {
		if err := tcpSrv.Init(); err != nil {
			fmt.Println("tcp srv init fail, err is", err)
			return
		}
		if err := tcpSrv.Start(); err != nil {
			fmt.Println("tcp srv start fail, err is", err)
			return
		}
		//// 迁移listen
		if err := tcpSrv.SendListenerWithUnixSocket(); err != nil {
			fmt.Println("send listener with unix socket fail, err is", err)
			return
		}

		// 迁移conn
		//if err := tcpSrv.SendConnWithUnixSocket(); err != nil {
		//	fmt.Println("send listener with unix socket fail, err is", err)
		//	return
		//}

	} else {
		//// 迁移listen
		if err := tcpSrv.RecvListenerFromUnixSocket(); err != nil {
			fmt.Println("recv listener with unix socket fail, err is", err)
			return
		}

		// 迁移conn
		//if err := tcpSrv.RecvConnFromUnixSocket(); err != nil {
		//	fmt.Println("recv listener with unix socket fail, err is", err)
		//	return
		//}
	}

	select {}
}

type TcpSrv struct {
	listener *net.TCPListener
	conns    map[string]*net.TCPConn
}

func NewTcpSrv() *TcpSrv {
	return &TcpSrv{
		conns: make(map[string]*net.TCPConn),
	}
}

func (t *TcpSrv) Init() error {
	listener, err := net.Listen("tcp", ":7000")
	if err != nil {
		return err
	}

	t.listener = listener.(*net.TCPListener)
	return nil
}

func (t *TcpSrv) Start() error {
	fmt.Println("listen on: ", t.listener.Addr())
	go func() {
		for {
			conn, err := t.listener.Accept()
			if err != nil {
				fmt.Println("accept fail, err msg is", err)
				continue
			}
			go t.clientSrv(conn)
			storeConn := conn.(*net.TCPConn)
			t.conns[conn.RemoteAddr().String()] = storeConn
		}
	}()
	return nil
}

func (t *TcpSrv) StartWithListenSocket(listener *net.TCPListener) error {
	fmt.Println("listen on: ", listener.Addr())
	go func() {
		for {
			c, err := listener.Accept()
			if err != nil {
				fmt.Println("accept fail, err msg is", err)
				continue
			}
			go t.clientSrv(c)
		}
	}()
	return nil
}

func (t *TcpSrv) clientSrv(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	pids := fmt.Sprintf("[%d] ", os.Getpid())

	for {
		time.Sleep(1 * time.Second)

		nRead, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("read msg fail, err is", err)
			return
		}

		// after receive the message, print it and stop the processing.
		fmt.Print(pids, time.Now(), " recv msg is: ", string(buf[:nRead]))
		return
		/*
			if _, err := conn.Write(buf[:nRead]); err != nil {
				fmt.Println("write msg fail, err is", err)
				return
			}
		*/
	}
}

func (t *TcpSrv) SendListenerWithUnixSocket() error {
	_ = os.Remove("/tmp/unix_socket_tcp")
	addr, err := net.ResolveUnixAddr("unix", "/tmp/unix_socket_tcp")
	if err != nil {
		fmt.Println("Cannot resolve unix addr: " + err.Error())
		return err
	}

	listener, err := net.ListenUnix("unix", addr)
	if err != nil {
		fmt.Println("Cannot listen to unix domain socket: " + err.Error())
		return err
	}
	fmt.Println("Listening on", listener.Addr())

	go func() {
		for {
			c, err := listener.Accept()
			if err != nil {
				fmt.Println("Accept: " + err.Error())
				return
			}

			file, _ := t.listener.File()
			buf := make([]byte, 1)
			buf[0] = 0
			rights := syscall.UnixRights(int(file.Fd()))
			_, _, err = c.(*net.UnixConn).WriteMsgUnix(buf, rights, nil)
			if err != nil {
				fmt.Println("synchronize listen socket fail, err is", err.Error())
			}
		}
	}()

	return nil
}

func (t *TcpSrv) RecvListenerFromUnixSocket() error {
	connInterface, err := net.Dial("unix", "/tmp/unix_socket_tcp")
	if err != nil {
		fmt.Println("net dial unix fail", err.Error())
		return err
	}
	defer func() {
		_ = connInterface.Close()
	}()

	unixConn := connInterface.(*net.UnixConn)

	b := make([]byte, 1)
	oob := make([]byte, 32)
	for {
		err = unixConn.SetWriteDeadline(time.Now().Add(time.Minute * 3))
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		n, oobn, _, _, err := unixConn.ReadMsgUnix(b, oob)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if n != 1 || b[0] != 0 {
			if n != 1 {
				fmt.Printf("recv fd type error: %d\n", n)
			} else {
				fmt.Println("init finish")
			}
			return err
		}

		scms, err := unix.ParseSocketControlMessage(oob[0:oobn])
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if len(scms) != 1 {
			fmt.Printf("recv fd num != 1 : %d\n", len(scms))
			return err
		}
		fds, err := unix.ParseUnixRights(&scms[0])
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if len(fds) != 1 {
			fmt.Printf("recv fd num != 1 : %d\n", len(fds))
			return err
		}
		fmt.Printf("recv fd %d\n", fds[0])
		// 这里需要把file close， 不然每次重启都会多复制一个socket
		file := os.NewFile(uintptr(fds[0]), "fd-from-old")
		conn, err := net.FileListener(file)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		_ = file.Close()
		//fmt.Println(conn)

		lc := conn.(*net.TCPListener)
		go t.StartWithListenSocket(lc)
	}
}

func (t *TcpSrv) SendConnWithUnixSocket() error {
	_ = os.Remove("/tmp/unix_socket_tcp")
	addr, err := net.ResolveUnixAddr("unix", "/tmp/unix_socket_tcp")
	if err != nil {
		fmt.Println("Cannot resolve unix addr: " + err.Error())
		return err
	}

	listener, err := net.ListenUnix("unix", addr)
	if err != nil {
		fmt.Println("Cannot listen to unix domain socket: " + err.Error())
		return err
	}
	fmt.Println("Listening on", listener.Addr())

	go func() {
		for {
			c, err := listener.Accept()
			if err != nil {
				fmt.Println("Accept: " + err.Error())
				return
			}
			for _, conn := range t.conns {
				file, _ := conn.File()
				buf := make([]byte, 1)
				buf[0] = 0
				rights := syscall.UnixRights(int(file.Fd()))
				_, _, err = c.(*net.UnixConn).WriteMsgUnix(buf, rights, nil)
				if err != nil {
					fmt.Println("synchronize listen socket fail, err is", err.Error())
				}
			}
		}
	}()

	return nil
}

func (t *TcpSrv) RecvConnFromUnixSocket() error {
	connInterface, err := net.Dial("unix", "/tmp/unix_socket_tcp")
	if err != nil {
		fmt.Println("net dial unix fail", err.Error())
		return err
	}
	defer func() {
		_ = connInterface.Close()
	}()

	unixConn := connInterface.(*net.UnixConn)

	b := make([]byte, 1)
	oob := make([]byte, 32)
	for {
		err = unixConn.SetWriteDeadline(time.Now().Add(time.Minute * 3))
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		n, oobn, _, _, err := unixConn.ReadMsgUnix(b, oob)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if n != 1 || b[0] != 0 {
			if n != 1 {
				fmt.Printf("recv fd type error: %d\n", n)
			} else {
				fmt.Println("init finish")
			}
			return err
		}

		scms, err := unix.ParseSocketControlMessage(oob[0:oobn])
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if len(scms) != 1 {
			fmt.Printf("recv fd num != 1 : %d\n", len(scms))
			return err
		}
		fds, err := unix.ParseUnixRights(&scms[0])
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if len(fds) != 1 {
			fmt.Printf("recv fd num != 1 : %d\n", len(fds))
			return err
		}
		fmt.Printf("recv fd %d\n", fds[0])
		// 这里需要把file close， 不然每次重启都会多复制一个socket
		file := os.NewFile(uintptr(fds[0]), "fd-from-old")
		conn, err := net.FileConn(file)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		_ = file.Close()
		fmt.Println(conn)
		t.clientSrv(conn)
	}
}
