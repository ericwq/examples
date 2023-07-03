//go:build !darwin && !freebsd && !netbsd && !openbsd && !windows

package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

const (
	GetTermios = unix.TCGETS
	SetTermios = unix.TCSETS
)

// run % stty -a to check the current flag.
// run % stty -iutf8 to disable the flag
// run % stty iutf8 to enable the flag
// run % go run t.go to change the IUTF8 flag. (enable it)
// run % stty -a to check the current flag.
func main() {
	CheckIUTF8(int(os.Stdin.Fd()))
	SetIUTF8(int(os.Stdin.Fd()))
	CheckIUTF8(int(os.Stdin.Fd()))
}

func SetIUTF8(fd int) error {
	termios, err := unix.IoctlGetTermios(fd, GetTermios)
	if err != nil {
		return err
	}

	termios.Iflag |= unix.IUTF8 // enable the flag.
	unix.IoctlSetTermios(fd, SetTermios, termios)

	return nil
}

func CheckIUTF8(fd int) (bool, error) {
	termios, err := unix.IoctlGetTermios(fd, GetTermios)
	if err != nil {
		return false, err
	}

	fmt.Printf("#CheckIUTF8() raw   termios.Iflag=%016b\n", termios.Iflag)
	fmt.Printf("#CheckIUTF8() raw   unix.IUTF8   =%016b\n", unix.IUTF8)
	// termios.Iflag &^= unix.IUTF8
	p := (unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON)
	p |= unix.IUTF8

	fmt.Printf("#CheckIUTF8() &^    termios.Iflag=%016b\n", termios.Iflag&^unix.IUTF8)
	fmt.Printf("#CheckIUTF8() &^all termios.Iflag=%016b\n", p)
	fmt.Printf("#CheckIUTF8() &     termios.Iflag=%016b, enable=%t\n", termios.Iflag&unix.IUTF8, termios.Iflag&unix.IUTF8 > 0)
	// Input is UTF-8 (since Linux 2.6.4)
	return (termios.Iflag & unix.IUTF8) != 0, nil
}
