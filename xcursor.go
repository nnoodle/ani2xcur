package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/nnoodle/ani2xcur/ico"
)

func writePNG(name string, img image.Image) error {
	out, err := os.Create(name)
	defer out.Close()
	if err != nil {
		return err
	}
	return png.Encode(out, img)
}

func cleanFilename(fn string) string {
	return strings.ReplaceAll(strings.TrimSuffix(filepath.Base(fn), path.Ext(fn)), " ", "-")
}

func convertCur(filename string, cleanname string, conf io.StringWriter) error {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return err
	}

	ico, err := ico.DecodeIcons(file)
	if err != nil {
		return err
	}

	for j, img := range ico.Images {
		name := fmt.Sprintf(cleanname+"-%v.png", j)
		writePNG(filepath.Join(TEMPDIR, name), img)
		if _, err := conf.WriteString(fmt.Sprintf(
			// size xhot yhot filename
			"%v	%v	%v	%v %v\n", // plane/bits in cur is xy hotspots
			ico.Direntries[j].Width, ico.Direntries[j].Plane, ico.Direntries[j].Bits, name,
		)); err != nil {
			return err
		}
	}
	return nil
}

func convertAni(filename string, cleanname string, conf io.StringWriter) error {
	cursor, err := readRiff(filename)
	if err != nil {
		return err
	}

	for i, c := range cursor.Icons {
		for j, img := range c.Images {
			name := fmt.Sprintf(cleanname+"-%v-%v.png", i, j)
			writePNG(filepath.Join(TEMPDIR, name), img)
			if _, err := conf.WriteString(fmt.Sprintf(
				// size xhot yhot filename ms-delay
				"%v	%v	%v	%v %v\n", // plane/bits in cur is xy hotspots
				c.Direntries[j].Width, c.Direntries[j].Plane, c.Direntries[j].Bits,
				name, int(cursor.Header.JifRate*10/6*10),
			)); err != nil {
				return err
			}
		}
	}
	return nil
}

func convertFile(filename string) error {
	cleanname := cleanFilename(filename)
	conf, err := os.Create(filepath.Join(TEMPDIR, cleanname+".config"))
	defer conf.Close()
	if err != nil {
		return err
	}

	if err := convertAni(filename, cleanname, conf); err != nil {
		// try ico
		if err := convertCur(filename, cleanname, conf); err != nil {
			return err
		}
	}

	return exec.Command("xcursorgen", "--prefix", TEMPDIR,
		filepath.Join(TEMPDIR, cleanname+".config"),
		filepath.Join(OUTPUTDIR, cleanname+".xcur")).Run()
}
