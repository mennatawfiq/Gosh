package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var Builtins []string

func HandleExit() {
	if len(Args) > 1 {
		if Args[1] == "0" {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	} else {
		PrintOSln("Did you mean 'exit 0'?")
	}
}

func HandleEcho() {
	if len(Args) > 1 {
		PrintOSln(strings.Join(Args[1:], " "))
	} else {
		PrintOSln("")
	}
}

func HandleType() {
	if len(Args) > 1 {
		if Contains(Builtins, Args[1]) {
			PrintOSln(Args[1] + " is a shell builtin")
		} else {
			// Looking for the input command in the paths of the PATH variable
			path, err := exec.LookPath(Args[1])
			if err != nil {
				PrintOSln(Args[1] + ": not found")
				return
			}
			PrintOSln(Args[1] + " is " + path)
		}
	} else {
		PrintOSln("")
	}
}

func HandlePwd() {
	// Abs function returns an absolute representation of input path
	abs, err := filepath.Abs(".")
	if err != nil {
		PrintOSln("errr: " + err.Error())
	}
	PrintOSln(abs)
}

func HandleCd() {
	var err error
	if len(Args) == 1 || Args[1] == "~" {
		home, _ := os.UserHomeDir()
		err = os.Chdir(home)
	} else {
		if Args[1][0] == '/' {
			Args[1] = strings.TrimPrefix(Args[1], "/")
		}
		err = os.Chdir(Args[1])
	}
	if err != nil {
		PrintOSln("cd: " + Args[1] + ": No such file or directory")
	}
}

func BuiltinsHandler(ofile string, appnd, er, out bool, handler func()) {
	if ofile != "" {
		// Redirect output to the file
		file, err := Redir(appnd, ofile)
		if err != nil {
			PrintOSln("Error opening file: " + err.Error())
			return
		}
		defer file.Close()
		if out {
			// Redirect stdout to the file
			oldStdout := os.Stdout
			os.Stdout = file
			handler()
			os.Stdout = oldStdout // Restore stdout
		} else if er {
			// Redirect stderr to the file
			oldStderr := os.Stderr
			os.Stderr = file
			handler()
			os.Stderr = oldStderr // Restore stderr
		}
	} else {
		handler()
	}
}

func ExternalCmdsHandler(ofile string, appnd, er, out bool) {
	// Handle external commands
	command := exec.Command(Args[0], Args[1:]...)
	if ofile != "" {
		// Redirect output to the file
		file, err := Redir(appnd, ofile)
		if err != nil {
			PrintOSln("Error opening file: " + err.Error())
			return
		}
		defer file.Close()

		if out {
			// Redirect stdout to the file
			command.Stdout = file
		} else if er {
			// Redirect stderr to the file
			command.Stderr = file
		}
	} else {
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
	}

	err := command.Run()
	if err != nil && !er {
		PrintOSln(Args[0] + ": command not found")
	}
}
