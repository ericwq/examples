// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"reflect"
)

type action interface {
	readme() string
}

type base struct {
	color int
}

func (b base) readme() string {
	return fmt.Sprintf("base %d", b.color)
}

type red struct {
	base
}

func (act red) readme() string { return fmt.Sprintf("red %d", act.color) }

type green struct {
	base
}

func (act green) readme() string { return fmt.Sprintf("green %d", act.color) }

func main() {
	a := make([]action, 2)
	a[0] = red{base{100}}
	a[1] = green{base{245}}

	b := make([]action, 2)
	b[0] = red{base{100}}
	b[1] = green{base{245}}
	fmt.Printf("a= %s\n", a)
	fmt.Printf("b= %s\n", b)

	// interface value with pointer receiver can't compare with ==
	for i := range a {
		if a[i] != b[i] {
			fmt.Printf("item [%d] is not equal, [%T, %T]\n", i, a[i], b[i])
		}
	}

	// interface value with pointer receiver can be compared with reflect.DeepEqual()
	a[0] = &red{base{100}}
	a[1] = &green{base{245}}
	b[0] = &red{base{100}}
	b[1] = &green{base{245}}
	fmt.Printf("interface value with pointer receiver result =%t\n", reflect.DeepEqual(a, b))
}
