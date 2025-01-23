package utils

import (
	"os"
)

var Args []string

func RedirParser(ins []string) (ofile string, appnd, out, er bool) {
	for i, in := range ins {
		if in == ">" || in == "1>" || in == ">>" || in == "2>" || in == "2>>" {
			if i+1 < len(ins) {
				ofile = ins[i+1]
				Args = ins[:i]
			} else {
				PrintOSln("arguments not sufficient")
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
		Args = ins
	}
	return
}

func Redir(appnd bool, ofile string) (*os.File, error) {
	var file *os.File
	var err error
	if appnd {
		file, err = os.OpenFile(ofile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	} else {
		file, err = os.Create(ofile)
	}
	return file, err
}
