package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func printOSln(out string) {
	fmt.Fprintln(os.Stdout, out)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func handleExit(ins []string) {
	if len(ins) > 1 {
		if ins[1] == "0" {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	} else {
		printOSln("Did you mean 'exit 0'?")
	}
}

func handleEcho(ins []string) {
	if len(ins) > 1 {
		printOSln(strings.Join(ins[1:], " "))
	} else {
		printOSln("")
	}
}

func handleType(ins []string, builtins []string) {
	if len(ins) > 1 {
		if contains(builtins, ins[1]) {
			printOSln(ins[1] + " is a shell builtin")
		} else {
			// looking for the input command in the paths of the PATH variable
			path, err := exec.LookPath(ins[1])
			if err != nil {
				printOSln(ins[1] + ": not found")
				return
			}
			printOSln(ins[1] + " is " + path)
		}
	} else {
		printOSln("")
	}
}

func handlePwd() {
	// Abs function returns an absolute representation of input path
	abs, err := filepath.Abs(".")
	if err != nil {
		printOSln("errr: " + err.Error())
	}
	printOSln(abs)
}

func handleCd(ins []string) {
	var err error
	if len(ins) == 1 || ins[1] == "~" {
		home, _ := os.UserHomeDir()
		err = os.Chdir(home)
	} else {
		// if path[0] == '/' {
		// 	path = strings.TrimPrefix(path, "/")
		// }
		err = os.Chdir(ins[1])
	}
	if err != nil {
		printOSln("cd: " + ins[1] + ": No such file or directory")
	}
}

func main() {
	builtins := []string{"type", "echo", "exit", "pwd", "cd"}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, "\033[1m\033[36m$ ")

		// formatting the command to bold and yellow
		fmt.Fprint(os.Stdout, "\033[1m\033[33m")

		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprint(os.Stderr, "ERROR: ", err)
			return
		}

		// reset the formatting
		fmt.Fprint(os.Stdout, "\033[0m")

		// to remove the \n
		input = strings.TrimSpace(input)

		ins := strings.Split(input, " ")

		switch ins[0] {
		case "exit":
			handleExit(ins)
		case "echo":
			handleEcho(ins)
		case "type":
			handleType(ins, builtins)
		case "pwd":
			handlePwd()
		case "cd":
			handleCd(ins)
		default:
			/* handling the external commands */
			// prepare the command to be runnable
			command := exec.Command(ins[0], ins[1:]...)
			// assign the standard output and error to the command
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			err := command.Run()
			//if error, then the command does not exist
			if err != nil {
				printOSln(ins[0] + ": command not found")
			}
		}
	}
}
