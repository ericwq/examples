package main

/*

#include <utmps/utmpx.h>

typedef struct{
	unsigned char a;
	char b;
	int c;
	unsigned int d;
	char e[10];
}unpacked;

#pragma pack(1)
typedef struct{
	unsigned char a;
	char b;
	int c;
	unsigned int d;
	char e[10];
}packed;

*/
import "C"

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/ericwq/examples/utmpx/base"
)

func main() {
	// unpack := C.unpacked{}
	// pack := C.packed{}
	//
	// fmt.Println("Printing the structure of the unpacked struct")
	// spew.Dump(unpack)
	//
	// fmt.Println("Printing the structure of the packed struct")
	// spew.Dump(pack)
	//
	fmt.Println("Printing the structure of the TimeVal struct")
	spew.Dump(base.TimeVal{})

	fmt.Println("Printing the structure of the C.timeval struct")
	spew.Dump(C.struct_timeval{})

	fmt.Println("Printing the structure of the ExitStatus struct")
	spew.Dump(base.ExitStatus{})

	fmt.Println("Printing the structure of the C.exit_status struct")
	spew.Dump(C.struct_exit_status{})

	// fmt.Println("Printing the structure of the Utmpx struct")
	// spew.Dump(base.Utmpx{})
}
