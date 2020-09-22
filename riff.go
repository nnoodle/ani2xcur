package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"

	"github.com/campoy/riff"
	"github.com/nnoodle/ani2xcur/ico"
)

const (
	ANIFLAG_CUR = 0x01
	ANIFLAG_SEQ = 0x02
)

var (
	RIFF_LIST = riff.NewID("LIST")
	RIFF_ANIH = riff.NewID("anih")
	RIFF_ICON = riff.NewID("icon")
	RIFF_RATE = riff.NewID("rate")
	RIFF_SEQ  = riff.NewID("seq ")
)

// dword = uint32
// https://web.archive.org/web/20130530192915/http://oreilly.com/www/centers/gff/formats/micriff
type ANIHeader struct {
	HeaderSize          uint32 // Num bytes in AniHeader (36 bytes)
	NumFrames           uint32 // Number of unique Icons in this cursor
	NumSteps            uint32 // Number of Blits before the animation cycles
	Width, Height       uint32 // reserved, must be zero.
	BitCount, NumPlanes uint32 // reserved, must be zero.
	JifRate             uint32 // Default Jiffies (1/60th of a second) if rate chunk not present.
	Flags               uint32 // Animation Flag (see AF_ constants)
}

type ANICursor struct {
	Header ANIHeader
	// Iname  []byte
	// Iart   []byte
	Rate  []uint32
	Seq   []uint32
	Icons []ico.Icon
}

func readAnih(r io.Reader) (interface{}, error) {
	var h ANIHeader
	err := binary.Read(r, binary.LittleEndian, &h)
	return h, err
}

func readIcon(r io.Reader) (interface{}, error) {
	return ico.DecodeIcons(r)
}

func readRiff(filename string) (ANICursor, error) {
	ani, err := os.Open(filename)
	defer ani.Close()
	if err != nil {
		return ANICursor{}, err
	}

	riffreader := riff.NewDecoder(ani)
	riffreader.Map(RIFF_ANIH, riff.DecoderFunc(readAnih))
	riffreader.Map(RIFF_ICON, riff.DecoderFunc(readIcon))
	root, err := riffreader.Decode()
	if err != nil {
		return ANICursor{}, err
	}
	return parseRiff(root), nil
}

func parseRiff(root *riff.Chunk) (cursor ANICursor) {
	chunks := make(map[riff.ID]*riff.Chunk)
	for _, chunk := range root.Chunks {
		chunks[chunk.ID] = chunk
	}

	cursor.Header = chunks[RIFF_ANIH].Content.(ANIHeader)
	cursor.Icons = make([]ico.Icon, len(chunks[RIFF_LIST].Chunks))
	for idx, chunk := range chunks[RIFF_LIST].Chunks {
		cursor.Icons[idx] = chunk.Content.(ico.Icon)
	}

	if chunks[RIFF_RATE] != nil {
		cursor.Rate = make([]uint32, cursor.Header.NumSteps)
		r := bytes.NewReader(chunks[RIFF_RATE].Data)
		binary.Read(r, binary.LittleEndian, &cursor.Rate)
	}

	if chunks[RIFF_SEQ] != nil {
		cursor.Seq = make([]uint32, cursor.Header.NumSteps)
		r := bytes.NewReader(chunks[RIFF_SEQ].Data)
		binary.Read(r, binary.LittleEndian, &cursor.Seq)
	}

	return cursor
}
