//go:build darwin || freebsd || openbsd || netbsd
// +build darwin freebsd openbsd netbsd

package main

import (
	// "bufio"
	"bufio"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

const redColor = "\033[1;31m%s\033[0m"

func checkIUTF8(fd int) (bool, error) {
	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return false, err
	}

	// Input is UTF-8 (since Linux 2.6.4)
	return (termios.Iflag & unix.IUTF8) != 0, nil
}

func setIUTF8(fd int) error {
	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return err
	}

	termios.Iflag |= unix.IUTF8

	if err := unix.IoctlSetTermios(fd, unix.TCSETS, termios); err != nil {
		return err
	}
	return nil
}

func main() {
	// stdin should return false for first calling
	flag, err := checkIUTF8(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("#checkIUTF8 stdin report error: %s\n", err)
	} else {
		fmt.Printf("#checkIUTF8 stadin report %t\n", flag)
	}

	// set IUTF8 for stdin
	err = setIUTF8(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("#setIUTF8 should report nil, got %s\n", err)
	} else {
		fmt.Printf("#setIUTF8 stadin done\n")
	}

	// stdin should return true after setIUTF8
	flag, err = checkIUTF8(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("#checkIUTF8 stdin report error: %s\n", err)
	} else {
		fmt.Printf("#checkIUTF8 stadin report %t\n", flag)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(redColor, "Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}
