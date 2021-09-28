package main

import (
	"fmt"

	"github.com/ericwq/examples/struct_func/so"
)

type serverOptions struct {
	load       int
	numWorkers uint32
}

type serverOption struct {
	so.EmptyServerOption
	apply func(*serverOptions)
}

func main() {

	opt := serverOption{}

	opt.apply = func(o *serverOptions) {
		fmt.Println("apply field function is called.")
	}

	// DoApply is exported, because we can't call the apply method
	so.DoApply(opt)

	// apply field is a function
	opt.apply(new(serverOptions))
}
