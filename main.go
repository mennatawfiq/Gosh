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
		if ins[1][0] == '/' {
			ins[1] = strings.TrimPrefix(ins[1], "/")
		}
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
		// print the shell prompt to the terminal with formatting
		fmt.Fprint(os.Stdout, "\033[1m\033[36m$ ")
		fmt.Fprint(os.Stdout, "\033[1m\033[33m")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprint(os.Stderr, "ERROR: ", err)
			return
		}

		// reset the formatting
		fmt.Fprint(os.Stdout, "\033[0m")

		ins := strings.Split(strings.TrimSpace(input), " ")

		var ofile string
		var args []string
		var appnd, out, er bool
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

		switch args[0] {
		case "exit":
			handleExit(args)
		case "echo":
			if ofile != "" {
				// redirect output to the file
				file, err := redir(appnd, ofile)
				if err != nil {
					printOSln("Error opening file: " + err.Error())
					continue
				}
				defer file.Close()
				if out {
					// redirect stdout to the file
					oldStdout := os.Stdout
					os.Stdout = file
					handleEcho(args)
					os.Stdout = oldStdout // restore stdout
				} else if er {
					// redirect stderr to the file
					oldStderr := os.Stderr
					os.Stderr = file
					handleEcho(args)
					os.Stderr = oldStderr // restore stderr
				}

			} else {
				handleEcho(args)
			}
		case "type":
			handleType(args, builtins)
		case "pwd":
			handlePwd()
		case "cd":
			handleCd(args)
		default:
			// handle external commands
			command := exec.Command(args[0], args[1:]...)
			if ofile != "" {
				// redirect output to the file
				file, err := redir(appnd, ofile)
				if err != nil {
					printOSln("Error opening file: " + err.Error())
					continue
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
			if err != nil {
				printOSln(args[0] + ": command not found")
			}
		}
	}
}
