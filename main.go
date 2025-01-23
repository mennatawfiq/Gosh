package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var args, builtins []string

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

func input(reader *bufio.Reader) []string {
	fmt.Fprint(os.Stdout, "\033[1m\033[36m$ ")
	fmt.Fprint(os.Stdout, "\033[1m\033[33m")

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprint(os.Stderr, "ERROR: ", err)
		return []string{}
	}

	// reset the formatting
	fmt.Fprint(os.Stdout, "\033[0m")

	return strings.Split(strings.TrimSpace(input), " ")
}

func redirParser(ins []string) (ofile string, appnd, out, er bool) {
	for i, in := range ins {
		if in == ">" || in == "1>" || in == ">>" || in == "2>" || in == "2>>" {
			if i+1 < len(ins) {
				ofile = ins[i+1]
				args = ins[:i]
			} else {
				printOSln("arguments not sufficient")
				continue
			}
			if in == ">>" || in == "2>>" {
				appnd = true
			}
			if in == ">" || in == "1>" || in == ">>" {
				out = true
			} else if in == "2>" || in == "2>>" {
				er = true
			}
			break
		}
	}
	if ofile == "" {
		args = ins
	}
	return
}

func redir(appnd bool, ofile string) (*os.File, error) {
	var file *os.File
	var err error
	if appnd {
		file, err = os.OpenFile(ofile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	} else {
		file, err = os.Create(ofile)
	}
	return file, err
}

func handleExit() {
	if len(args) > 1 {
		if args[1] == "0" {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	} else {
		printOSln("Did you mean 'exit 0'?")
	}
}

func handleEcho() {
	if len(args) > 1 {
		printOSln(strings.Join(args[1:], " "))
	} else {
		printOSln("")
	}
}

func handleType() {
	if len(args) > 1 {
		if contains(builtins, args[1]) {
			printOSln(args[1] + " is a shell builtin")
		} else {
			// looking for the input command in the paths of the PATH variable
			path, err := exec.LookPath(args[1])
			if err != nil {
				printOSln(args[1] + ": not found")
				return
			}
			printOSln(args[1] + " is " + path)
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

func handleCd() {
	var err error
	if len(args) == 1 || args[1] == "~" {
		home, _ := os.UserHomeDir()
		err = os.Chdir(home)
	} else {
		if args[1][0] == '/' {
			args[1] = strings.TrimPrefix(args[1], "/")
		}
		err = os.Chdir(args[1])
	}
	if err != nil {
		printOSln("cd: " + args[1] + ": No such file or directory")
	}
}

func builtinsRedirHandle(ofile string, appnd, er, out bool, handler func()) {
	if ofile != "" {
		// redirect output to the file
		file, err := redir(appnd, ofile)
		if err != nil {
			printOSln("Error opening file: " + err.Error())
			return
		}
		defer file.Close()
		if out {
			// redirect stdout to the file
			oldStdout := os.Stdout
			os.Stdout = file
			handler()
			os.Stdout = oldStdout // restore stdout
		} else if er {
			// redirect stderr to the file
			oldStderr := os.Stderr
			os.Stderr = file
			handler()
			os.Stderr = oldStderr // restore stderr
		}
	} else {
		handler()
	}
}

func externalCmdsHandler(ofile string, appnd, er, out bool) {
	// handle external commands
	command := exec.Command(args[0], args[1:]...)
	if ofile != "" {
		// redirect output to the file
		file, err := redir(appnd, ofile)
		if err != nil {
			printOSln("Error opening file: " + err.Error())
			return
		}
		defer file.Close()

		if out {
			// redirect stdout to the file
			command.Stdout = file
		} else if er {
			// redirect stderr to the file
			command.Stderr = file
		}
	} else {
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
	}

	err := command.Run()
	if err != nil && !er {
		printOSln(args[0] + ": command not found")
	}
}

func main() {
	builtins = []string{"type", "echo", "exit", "pwd", "cd"}
	reader := bufio.NewReader(os.Stdin)

	for {
		ins := input(reader)
		if len(ins) == 0 {
			printOSln("Error with IO\nProgram exit..")
			return
		}

		ofile, appnd, out, er := redirParser(ins)

		switch args[0] {
		case "exit":
			handleExit()
		case "echo":
			builtinsRedirHandle(ofile, appnd, er, out, handleEcho)
		case "type":
			builtinsRedirHandle(ofile, appnd, er, out, handleType)
		case "pwd":
			builtinsRedirHandle(ofile, appnd, er, out, handlePwd)
		case "cd":
			builtinsRedirHandle(ofile, appnd, er, out, handleCd)
		default:
			externalCmdsHandler(ofile, appnd, er, out)
		}
	}
}
