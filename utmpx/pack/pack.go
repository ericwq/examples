package pack

/*
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
)

func main() {
	unpack := C.unpacked{}
	pack := C.packed{}

	fmt.Println("Printing the structure of the unpacked struct")
	spew.Dump(unpack)

	fmt.Println("Printing the structure of the packed struct")
	spew.Dump(pack)
}
