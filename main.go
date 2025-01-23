package main

import (
	"bufio"
	"Gosh/utils"
	"os"
)

var Args, Builtins []string

func main() {
	Builtins = []string{"type", "echo", "exit", "pwd", "cd"}
	reader := bufio.NewReader(os.Stdin)

	for {
		Args = utils.Input(reader)
		if len(Args) == 0 {
			utils.PrintOSln("Error with IO\nProgram exit..")
			return
		}

		ofile, appnd, out, er := utils.RedirParser(Args)

		switch Args[0] {
		case "exit":
			utils.HandleExit()
		case "echo":
			utils.BuiltinsHandler(ofile, appnd, er, out, utils.HandleEcho)
		case "type":
			utils.BuiltinsHandler(ofile, appnd, er, out, utils.HandleType)
		case "pwd":
			utils.BuiltinsHandler(ofile, appnd, er, out, utils.HandlePwd)
		case "cd":
			utils.BuiltinsHandler(ofile, appnd, er, out, utils.HandleCd)
		default:
			utils.ExternalCmdsHandler(ofile, appnd, er, out)
		}
	}
}
