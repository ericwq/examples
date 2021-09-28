package so

import "fmt"

type serverOptions struct {
	numServerWorkers int
}

type ServerOption interface {
	apply(*serverOptions)
}

/*

EmptyServerOption implements ServerOption interface, provided default apply method

*/
type EmptyServerOption struct{}

func (EmptyServerOption) apply(o *serverOptions) {
	o.numServerWorkers = 17
	fmt.Println("the apply method of EmptyServerOption is called.")
}

/*

DoApply function is exported.

*/
func DoApply(s ServerOption) {
	s.apply(new(serverOptions))
}
