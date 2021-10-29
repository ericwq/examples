package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {

	var isChild bool
	env := os.Environ()

	//check env.
	if env != nil && len(env[0]) > 4 {
		if strings.HasPrefix(env[0], "STOP") {
			isChild = true
		}
	}

	var wg sync.WaitGroup
	if !isChild {
		fmt.Printf("#1 isChild is %t\n", isChild)
		inR, inW, _ := os.Pipe()
		outR, outW, _ := os.Pipe()

		//#1 start process.
		os.StartProcess(os.Args[0], os.Args[1:], &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr, inW, outR},
			Env:   []string{"STOP=START"},
		})

		//#1 read from inR.
		wg.Add(1)
		go readFrom(inR, "#1 ", &wg)

		//#1 write to outW.
		writeTo(outW)

	} else {
		fmt.Printf("#2 isChild is %t\n", isChild)

		//#2 open inW. TODO try to use the passed file instead
		inW := os.NewFile(3, "inW")

		//#2 open outR.
		outR := os.NewFile(4, "outR")

		//#2 read from outR.
		wg.Add(1)
		go readFrom(outR, "#2 ", &wg)

		//#2 write to inW.
		writeTo(inW)
	}

	wg.Wait()
}

func readFrom(pfile *os.File, prefix string, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(pfile)
	for i := 0; i < 3; i++ {
		if scanner.Scan() {
			fmt.Println(prefix + "receive: " + scanner.Text())
		} else {
			fmt.Println(prefix + scanner.Err().Error())
		}
	}
	fmt.Println(prefix + "finish")
	// pfile.Close()
}

func writeTo(pfile *os.File) {
	writer := bufio.NewWriter(pfile)
	for i := 0; i < 3; i++ {
		time.Sleep(time.Millisecond * 100)
		writer.WriteString(time.Now().String() + "\n")
		writer.Flush()
	}
	//pfile.Close()
}
