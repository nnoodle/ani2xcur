package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var TEMPDIR string
var OUTPUTDIR string = "."

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
		os.Exit(22)
	case 2:
		input = []string{os.Args[1]}
	default:
		l := len(os.Args)
		input = os.Args[1 : l-1]
		OUTPUTDIR = os.Args[l-1]

		finfo, err := os.Stat(OUTPUTDIR)
		if os.IsNotExist(err) {
			if err := os.Mkdir(OUTPUTDIR, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "error when creating %v: %v\n", OUTPUTDIR, err)
				os.Exit(2)
			}
		} else if !finfo.IsDir() {
			fmt.Fprintln(os.Stderr, "error:", OUTPUTDIR, "is not a directory")
			os.Exit(20)
		}
	}

	TEMPDIR, err = ioutil.TempDir(os.TempDir(), "ani2xcur-")
	defer os.RemoveAll(TEMPDIR)
	defer os.Remove(TEMPDIR)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error making temporary directory:", err)
		os.Exit(1)
	}

	for _, f := range input {
		ext := filepath.Ext(f)
		switch ext {
		case ".ico":
			fallthrough
		case ".cur":
			fallthrough
		case ".ani":
			if err := convertFile(f); err != nil {
				fmt.Println("error when processing", f, err)
			}
		default:
			fmt.Println("ignoring file with invalid extension:", f)
		}
	}
}
