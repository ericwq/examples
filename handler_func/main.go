package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// fakeWriter implements http.ResponseWriter interface
type fakeWriter struct {
	io.Writer
}

func (fakeWriter) Header() http.Header { panic("should not be called") }
func (fakeWriter) WriteHeader(int)     { panic("should not be called") }

// http.HandlerFunc will convert a function into a type which implements http.Handler interface
//
// TODO  Please check the definition of http.HandlerFunc
//
func main() {

	// printHandler is a value of type http.HandlerFunc
	printHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Println("ServeHTTP is called, so this function is called.")
	})

	// create a Writer for the fakeWriter
	bytesWriter := bytes.NewBuffer(make([]byte, 256))

	fmt.Println("Let's call the ServeHTTP method of http.HandlerFunc.")
	printHandler.ServeHTTP(fakeWriter{bytesWriter}, new(http.Request))
}
