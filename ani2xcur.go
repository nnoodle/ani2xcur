package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var TEMPDIR string
var OUTPUTDIR string

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"usage: %s [options] files...\n\nconverts windows cursor files to to xcur files\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&OUTPUTDIR, "o", ".", "output `dir`ectory")
	flag.StringVar(&TEMPDIR, "tmp", os.TempDir(), "temporary `dir`ectory [for debugging]")
}

func ensureDirectory(name string) {
	finfo, err := os.Stat(name)
	if os.IsNotExist(err) {
		if err = os.Mkdir(name, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error when creating %v: %v\n", name, err)
			os.Exit(2)
		}
	} else if !finfo.IsDir() {
		fmt.Fprintln(os.Stderr, "error:", name, "is not a directory")
		os.Exit(20)
	}
}

func main() {
	var err error
	flag.Parse()

	ensureDirectory(OUTPUTDIR)

	if TEMPDIR == os.TempDir() {
		TEMPDIR, err = ioutil.TempDir(os.TempDir(), "ani2xcur-")
		defer os.RemoveAll(TEMPDIR)
		defer os.Remove(TEMPDIR)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error making temporary directory:", err)
			os.Exit(1)
		}
	} else {
		ensureDirectory(TEMPDIR)
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	for _, f := range flag.Args() {
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
