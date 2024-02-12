package main

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
)

func main() {
	s, err := syslog.New(syslog.LOG_WARNING|syslog.LOG_LOCAL7, "example")
	if err != nil {
		log.Fatal(err)
	}

	msg := fmt.Sprintf("syslog from pid=%d", os.Getpid())
	s.Info(msg)

	// the following is the result of the above log
	//
	// 2024-02-12T16:57:52.594611+08:00 openrc-nvide example[30777]: syslog from pid=30777
}
