package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		var a chan string
		fmt.Println("A send to a nil channel blocks forever.")
		a <- "let's get started" // deadlock
	}()

	go func() {
		var b chan struct{}
		fmt.Println("A receive from a nil channel blocks forever")
		fmt.Println(<-b) // deadlock
		fmt.Println("A receive from a nil channel blocks forever")
	}()

	go func() {
		fmt.Println("A receive from a closed channel returns the zero value immediately")
		c := make(chan int, 3)
		c <- 1
		c <- 2
		c <- 3
		close(c)
		for i := 0; i < 4; i++ {
			fmt.Printf("%d: %d \n", i, <-c) // prints 1 2 3 0
		}
	}()

	time.Sleep(time.Second * 3)

	fmt.Println("A send to a closed channel panics")
	var c = make(chan int, 100)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				c <- j
			}
			close(c)
		}()
	}
	for i := range c {
		fmt.Println(i)
	}
}
