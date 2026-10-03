package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/term"
)

func main() {
	width := flag.Int("width", 0, "column count; default is the terminal width")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: imgascii [--width N] <file-or-url>\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	cols := *width
	if cols == 0 {
		var err error
		cols, err = termCols()
		if err != nil {
			fmt.Fprintf(os.Stderr, "imgascii: %v\n", err)
			os.Exit(1)
		}
	}

	img, err := load(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "imgascii: %v\n", err)
		os.Exit(1)
	}
	if err := render(os.Stdout, img, cols); err != nil {
		fmt.Fprintf(os.Stderr, "imgascii: %v\n", err)
		os.Exit(1)
	}
}

func termCols() (int, error) {
	fd := int(os.Stdout.Fd())
	if !term.IsTerminal(fd) {
		return 0, fmt.Errorf("stdout is not a terminal; pass --width")
	}
	w, _, err := term.GetSize(fd)
	if err != nil || w < 1 {
		return 0, fmt.Errorf("stdout is not a terminal; pass --width")
	}
	return w, nil
}
