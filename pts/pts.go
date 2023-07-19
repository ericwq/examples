package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	// fmt.Printf("The name of os.Stdin is %s\n", os.Stdin.Name())
	// fmt.Printf("The name of os.Stderr is %s\n", os.Stderr.Name())
	// fmt.Printf("The name of os.Stdout is %s\n", os.Stdout.Name())

	if err := getDeviceName(); err != nil {
		fmt.Printf("error:%s\n", err)
	}
}

func getDeviceName() error {
	// run the following command to get the command output
	//
	// % ls /proc/[PID]/fd/ -l
	// total 0
	// lrwx------    1 ide      develop         64 Jul 17 09:09 0 -> /dev/pts/1
	// lrwx------    1 ide      develop         64 Jul 17 09:09 1 -> /dev/pts/1
	// lrwx------    1 ide      develop         64 Jul 19 09:44 10 -> /dev/tty
	// lrwx------    1 ide      develop         64 Jul 17 09:09 2 -> /dev/pts/1
	prog := "ls"
	arg1 := fmt.Sprintf("/proc/%d/fd/", os.Getpid())
	arg := []string{arg1, "-l"}
	out, err := exec.Command(prog, arg...).Output()
	if err != nil {
		return err
	}

	// print the command output
	outString := string(out)
	fmt.Printf("PID=%d, output:\n%s", os.Getpid(), outString)

	// use regexp to parse the result
	re := regexp.MustCompile(" ([0-2]) -> (/[^ ]+)\n")
	matched := re.FindAllStringSubmatch(string(out), 3)
	for i := range matched {
		fmt.Printf("Find %#v\n", matched[i])
	}

	return nil
}
