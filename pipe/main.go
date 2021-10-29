package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	//done := make(chan struct{})

	var isChild bool
	env := os.Environ()
	fmt.Println("check env.")
	if env != nil && len(env[0]) > 4 {

		if strings.HasPrefix(env[0], "STOP") {
			isChild = true
		}
	}

	fmt.Printf("isChild is %t\n", isChild)

	var process *os.Process
	if !isChild {
		inR, inW, _ := os.Pipe()
		outR, outW, _ := os.Pipe()

		fmt.Println("#1 start process.")
		process, _ = os.StartProcess(os.Args[0], os.Args[1:], &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr, inW, outR},
			Env:   []string{"STOP=START"},
		})

		fmt.Println("#1 write to outW.")
		go writeTo(outW)
		fmt.Println("#1 read from inR.")
		readFrom(inR, process)

	} else {
		fmt.Println("#2 open inW.")
		inW := os.NewFile(3, "inW")

		fmt.Println("#2 open outR.")
		outR := os.NewFile(4, "outR")

		fmt.Println("#2 write to inW.")
		go writeTo(inW)
		fmt.Println("#2 read from outR. ")
		readFrom(outR, process)

	}
}

func readFrom(outR *os.File, process *os.Process) {
	scanner := bufio.NewScanner(outR)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	/*
		if process != nil {
			process.Signal(os.Kill)
		}
	*/
	fmt.Println("finish")
}

func writeTo(inW *os.File) {
	writer := bufio.NewWriter(inW)
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		writer.WriteString("date\n")
		writer.Flush()
	}
	inW.Close()
}
