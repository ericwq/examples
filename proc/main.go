package main

import (
	"fmt"
	"os/exec"
	// "os/signal"
	"syscall"
	"time"
)

func terminate(pid int) {
	for {
		<-time.After(1 * time.Second)
		fmt.Println("trying to kill", pid)
		syscall.Kill(pid, syscall.SIGTERM)
	}
}

func main() {
	// signal.Ignore(syscall.SIGTERM)

	/*
		https://github.com/golang/go/issues/20479

		POSIX specifies that for execve signals whose handlers are either SIG_IGN
		or SIG_DFL are left unchanged: that is, an ignored signal is still ignored
		after an execve. Calling signal.Ignore sets the handler to SIG_IGN, so what
		are you seeing is that the execution is preserving the ignored state of the
		signal. Go doesn't have to follow POSIX, of course, but you are basically
		suggesting that os.StartProcess should override ignored signals and set them
		back to the default. We could do that, but is it the right thing to do?

		There is a clear utility to being able to ignore signals across execve:
		that's how the nohup utility works. Today, we could write the nohup utility
		in Go. With your proposed change, we would not be able to. It's easy to catch
		but ignore a signal in Go: just call signal.Notify(make(chan os.Signal)).
		So your proposed change would make Go strictly less useful than it is today.
		That does not seem wise.

		One possibility would be to add something to syscall.SysProcAttr that lets
		you specify signal dispositions (either ignored or default). Then you could
		easily control signals. But given that it is fairly simple to do that anyhow,
		I'm not sure how useful that would be.
	*/

	cmd := exec.Command("sleep", "3")

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	fmt.Println("PID:", cmd.Process.Pid)
	go terminate(cmd.Process.Pid)
	if state, err := cmd.Process.Wait(); err != nil {
		panic(err)
	} else {
		fmt.Println("Success:", state.Success(), state)
	}
}
