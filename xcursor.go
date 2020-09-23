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
		if err := writePNG(filepath.Join(TEMPDIR, name), img); err != nil {
			return err
		}
		if _, err := conf.WriteString(fmt.Sprintf(
			// size xhot yhot filename
			"%v	%v	%v	%v %v\n", // plane/bits in cur files are x/y hotspots
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
		// try ico
		return convertCur(filename, cleanname, conf)
	}
	for i, c := range cursor.Icons {
		// assume Icons in animated cursors only have one image
		name := fmt.Sprintf(cleanname+"-%v.png", i)
		if err := writePNG(filepath.Join(TEMPDIR, name), c.Images[0]); err != nil {
			return err
		}
	}
	for _, i := range cursor.Seq {
		c := cursor.Icons[i]
		name := fmt.Sprintf(cleanname+"-%v.png", i)
		if _, err := conf.WriteString(fmt.Sprintf(
			// size xhot yhot filename ms-delay
			"%v	%v	%v	%v %v\n", // plane/bits in cur is xy hotspots
			c.Direntries[0].Width, c.Direntries[0].Plane, c.Direntries[0].Bits,
			name, int(cursor.Rate[i]*(100/6)),
		)); err != nil {
			return err
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
		return err
	}

	cmd := exec.Command("xcursorgen", "--prefix", TEMPDIR,
		filepath.Join(TEMPDIR, cleanname+".config"),
		filepath.Join(OUTPUTDIR, cleanname+".xcur"))
	if _, err := cmd.Output(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			os.Stderr.Write(exiterr.Stderr)
		}
		return err
	}
	return nil
}
