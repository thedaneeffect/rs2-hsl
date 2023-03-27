package main

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"
)

//go:embed palette.bin
var palette []byte

func main() {
	// catchall
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("usage: hsl [0, 65535]")
		}
	}()

	// parse arg
	var index uint64
	index, err := strconv.ParseUint(os.Args[1], 0, 64)
	if err != nil {
		panic(err)
	}

	// unpack rs2 hsl
	h := (index >> 10)
	s := (index >> 7) & 7
	l := index & 127

	// scale up
	hf := float64(h*360) / 64.0
	sf := float64(s*100) / 7.0
	lf := float64(l*100) / 127.0

	// move to pixel space
	index *= 3
	r := palette[index+0]
	g := palette[index+1]
	b := palette[index+2]

	// decide text color
	var tc byte
	if l < 63 {
		tc = 255
	}

	start := setcolor(48, r, g, b) + setcolor(38, tc, tc, tc)
	end := "\033[0m\n"
	fmt.Print(start, fmt.Sprintf("rgb %20s", fmt.Sprintf("%d %d %d", r, g, b)), end)
	fmt.Print(start, fmt.Sprintf("hex %20s", fmt.Sprintf("%06X", int(r)<<16|int(g)<<8|int(b))), end)
	fmt.Print(start, fmt.Sprintf("hsl %20s", fmt.Sprintf("%.1f %.1f%% %.1f%%", hf, sf, lf)), end)
}

func setcolor(code, r, g, b byte) string {
	return fmt.Sprintf("\033[%d;2;%d;%d;%dm", code, r, g, b)
}
