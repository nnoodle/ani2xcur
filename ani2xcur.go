package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

var TEMPDIR string
var OUTPUTDIR string = "."

func printFatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}

func main() {
	var (
		input []string
		err   error
	)
	switch len(os.Args) {
	case 0:
		panic("unreachable")
	case 1:
		fmt.Println("usage: ani2xcur input-file [input-file]... [output-directory]")
		os.Exit(1)
	case 2:
		input = []string{os.Args[1]}
	case 3:
	default:
		l := len(os.Args)
		input = os.Args[1 : l-1]
		OUTPUTDIR = os.Args[l-1]
	}

	TEMPDIR, err = ioutil.TempDir(os.TempDir(), "ani2xcur-")
	if err != nil {
		fmt.Println(TEMPDIR)
		printFatal(err)
	}
	//defer os.RemoveAll(TEMPDIR)

	for _, f := range input {
		if err := convertFile(f); err != nil {
			printFatal(err)
		}
	}
}
