package main

import (
	"fmt"
	"sync"
	"time"
)

// prefer using formal arguments for the channels you pass to go-routines
// instead of accessing channels in global scope. You can get more compiler
// checking this way, and better modularity too.
//
// avoid both reading and writing on the same channel in a particular go-routine
// (including the 'main' one). Otherwise, deadlock is a much greater risk.
//
// a feature of Go channels: it is possible to have multiple writers sharing
// one channel; Go will interleave the messages automatically.
// The same applies for one writer and multiple readers on one channel,
// as seen in the second example here:
//
// It is generally a good principle to view buffering as a performance enhancer
// only. If your program does not deadlock without buffers, it won't deadlock
// with buffers either (but the converse is not always true). So, as another
// rule of thumb, start without buffering then add it later as needed.
func main() {

	c := make(chan string)
	var wg sync.WaitGroup

	j := 0
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int, c <-chan string) {
			defer wg.Done()
			msg := <-c
			fmt.Printf("got %q in goroutine %d\n", msg, i)

		}(j, c)
		j++
	}
	for i := 0; i < 5; i++ {
		c <- "original"
	}

	wg.Wait()

	// c := make(chan string)
	//
	// for i := 1; i <= 5; i++ {
	// 	go func(i int, co chan<- string) {
	// 		for j := 1; j <= 5; j++ {
	// 			co <- fmt.Sprintf("hi from %d.%d", i, j)
	// 		}
	// 	}(i, c)
	// }
	//
	// for i := 1; i <= 25; i++ {
	// 	fmt.Println(<-c)
	// }
}

// many writers & one reader on a channel:
func manyWriterOneReader() {
	c := make(chan string)

	for i := 1; i <= 5; i++ {
		go func(i int, co chan<- string) {
			for j := 1; j <= 5; j++ {
				co <- fmt.Sprintf("hi from %d.%d", i, j)
			}
		}(i, c)
	}

	for i := 1; i <= 25; i++ {
		fmt.Println(<-c)
	}
}

// one writer and multiple readers on one channel
func oneWriterManyReader() {
	c := make(chan int)
	var w sync.WaitGroup
	w.Add(5)

	for i := 1; i <= 5; i++ {
		go func(i int, ci <-chan int) {
			j := 1
			for v := range ci {
				time.Sleep(time.Millisecond)
				fmt.Printf("%d.%d got %d\n", i, j, v)
				j += 1
			}
			w.Done()
		}(i, c)
	}

	for i := 1; i <= 25; i++ {
		c <- i
	}
	close(c)
	w.Wait()
}
