package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// create aprilsh-[PID].log in tmp directory
	name := joinPath(os.TempDir(), logFileName("aprilsh", os.Getpid(), "log"))
	file, _ := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	fp := strings.TrimPrefix(file.Name(), os.TempDir())
	fmt.Println(file.Name())
	fmt.Println(fp)
}

func joinPath(dir, name string) string {
	if len(dir) > 0 && os.IsPathSeparator(dir[len(dir)-1]) {
		return dir + name
	}
	return dir + string(os.PathSeparator) + name
}

func logFileName(first string, second int, ext string) string {
	return fmt.Sprintf("%s-%d.%s", first, second, ext)
}
