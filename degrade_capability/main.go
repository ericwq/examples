package main

import (
	"bytes"
	"fmt"
	"io"
)

type onlyWriter struct {
	io.Writer
}

func main() {
	bytesWriter := bytes.NewBuffer(make([]byte, 256))

	// onlyWriter degrade the capability of bytesWriter.
	// bytesWriter implements several interface.
	ow := onlyWriter{bytesWriter}

	var data []byte = []byte("a good day")

	// ow.Write is the ony method that can be called.
	num, err := ow.Write(data)
	if err == nil {
		fmt.Printf("write %d bytes to destination\n", num)
	} else {
		fmt.Printf("write error: %s\n", err)
	}

}
