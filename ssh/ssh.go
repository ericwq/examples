package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

func main() {
	user := "ide"
	password := "password"
	host := "localhost"
	port := 22
	timeout := 10

	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)

	// only password authentication
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: time.Duration(timeout) * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	// connet to ssh
	fmt.Println("ssh dial.")
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		fmt.Printf("create connection failed. %s\n", err)
	}

	// create session
	fmt.Println("\nnew session.")
	if session, err = client.NewSession(); err != nil {
		fmt.Printf("create session failed. %s\n", err)
	}

	// run remote command
	fmt.Println("\nrun cmd.")
	out, err := session.Output("pwd")
	if err != nil {
		fmt.Println("run :", err)
	} else {
		fmt.Printf("cmd output:%s", out)
	}

	// fmt.Println("\nclose session.")
	// err = session.Close()
	// if err != nil {
	// 	fmt.Println("close session error=", err)
	// }

	time.Sleep(15 * time.Second)

	fmt.Println("\nclose client.")
	err = client.Close()
	if err != nil {
		fmt.Println("close client error=", err)
	}
}
