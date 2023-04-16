package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

	shellNotWorking()
	cmdStartWait()
	shellStart()
}

func cmdStartWait() {
	fmt.Printf("\ncmdStartWait:\n")
	cmd := exec.Command("sleep", "5000")
	// cmd := exec.Command("sh", "-sh")
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(time.Now().UnixMilli(), "Process started ", cmd, " PID:", cmd.Process.Pid)
	time.AfterFunc(40*time.Millisecond, func() {
		fmt.Println(time.Now().UnixMilli(), "-- Process state: ", cmd.ProcessState)
		err = cmd.Process.Kill()
		if err != nil {
			fmt.Println("-- errors:", err)
		}
		fmt.Println(time.Now().UnixMilli(), "-- Process killed with PID:", cmd.Process.Pid)
	})

	fmt.Println(time.Now().UnixMilli(), "cmd Wait.")
	cmd.Wait()
	fmt.Println(time.Now().UnixMilli(), "cmd Wait finished: ", cmd.Process.Pid)
}

func shellNotWorking() {
	fmt.Printf("\nshellNotWorking:\n")
	cmd := exec.Command("sh", "-sh")
	// cmd := exec.Command("sleep", "5000")
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("cmd start ", cmd, " PID:", cmd.Process.Pid)
	fmt.Println("cmd state: ", cmd.ProcessState)

	timer1 := time.NewTimer(time.Duration(200) * time.Millisecond)
	go func() {
		<-timer1.C
		fmt.Println("  Process state:", cmd.ProcessState)
		fmt.Println("  Process killing")
		err = cmd.Process.Kill()
		if err != nil {
			fmt.Println("  errors:", err)
		}
		fmt.Println("  Process killed with PID:", cmd.Process.Pid)
	}()

	fmt.Println("cmd wait.")
	cmd.Wait()
	fmt.Println("cmd wait finished.", cmd.Process.Pid)
}

// extract shell name from shellPath and prepend '-' to the returned shell name
func getShellNameFrom(shellPath string) (shellName string) {
	shellSplash := strings.LastIndex(shellPath, "/")
	if shellSplash == -1 {
		shellName = shellPath
	} else {
		shellName = shellPath[shellSplash+1:]
	}

	// prepend '-' to make login shell
	shellName = "-" + shellName

	return
}

// get current user home directory
func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return home
}

// http://technosophos.com/2014/07/11/start-an-interactive-shell-from-within-go.html
// it turns out to be that exec.Cmd can't be used to start shell.
func shellStart() {
	fmt.Printf("\nshellStart:\n")
	// Get the current user.
	// me, err := user.Current()
	// if err != nil {
	// 	panic(err)
	// }

	// Get the current working directory.
	cwd := getHomeDir()

	// Set an environment variable.
	os.Setenv("SOME_VAR", "1")

	// Transfer stdin, stdout, and stderr to the new process
	// and also set target directory for the shell to start in.
	pa := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir:   cwd,
	}

	// Start up a new shell.
	// Note that we supply "login" twice.
	// -fpl means "don't prompt for PW and pass through environment."
	fmt.Println("Starting a new interactive shell")

	// prepare shell path and parameters
	shellPath := os.Getenv("SHELL")
	argv := getShellNameFrom(shellPath)
	fmt.Println("shell path: ", shellPath, " argv:", argv)

	proc, err := os.StartProcess(shellPath, []string{argv}, &pa)
	if err != nil {
		panic(err)
	}

	// kill the shell process
	time.AfterFunc(10*time.Millisecond, func() {
		fmt.Println("")
		fmt.Println("-- kill the process with PID:", proc.Pid)
		err = proc.Kill()
		if err != nil {
			fmt.Println("-- errors:", err)
		}
		fmt.Println("-- Process killed with PID:", proc.Pid)
	})

	// Wait until user exits the shell
	state, err := proc.Wait()
	if err != nil {
		panic(err)
	}

	// Keep on keepin' on.
	fmt.Printf("Exited shell: %s\n", state.String())
}
