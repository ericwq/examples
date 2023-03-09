package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"golang.org/x/sys/unix"
)

func main() {
	// get initial window size
	windowSize, err := unix.IoctlGetWinsize(int(os.Stdin.Fd()), unix.TIOCGWINSZ)
	// windowSize, err := pty.GetsizeFull(os.Stdin)
	if err != nil || windowSize.Col == 0 || windowSize.Row == 0 {
		// Fill in sensible defaults. */
		// They will be overwritten by client on first connection.
		windowSize.Col = 80
		windowSize.Row = 24
	}

	cmd := exec.Command("sleep", "3")

	// pty.StartWithSize(cmd, convertWinsize(windowSize))

	fmt.Printf("#proc2  cmd: %s\n", cmd)
	fmt.Printf("#proc2 size: %v\n", windowSize)

	// var n [100000]byte
	//
	// for _, v := range n {
	// 	fmt.Printf(" 0x%02x", v)
	// }
	// fmt.Printf("\n")

	// real    0m 0.36s
	// user    0m 0.19s
	// sys     0m 0.32s

	var n [100000]byte

	bufStdout := bufio.NewWriter(os.Stdout)
	defer bufStdout.Flush()

	for _, v := range n {
		fmt.Fprintf(bufStdout, " 0x%02x", v)
	}
	fmt.Fprintf(bufStdout, "\n")
}

func convertWinsize(windowSize *unix.Winsize) *pty.Winsize {
	var sz pty.Winsize
	sz.Cols = windowSize.Col
	sz.Rows = windowSize.Row
	sz.X = windowSize.Xpixel
	sz.Y = windowSize.Ypixel

	return &sz
}
