package main

import (
	"bufio"
	"fmt"
	"os"
)

const redColor = "\033[1;31m%s\033[0m"

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(redColor,"Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}
