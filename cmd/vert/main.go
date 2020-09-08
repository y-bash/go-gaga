// TODO should have package comment
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/y-bash/go-gaga"
	"os"
	"strings"
)

var version = "v0.0.0" // set value by go build -ldflags

func max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}
	return n2
}

func min(n1, n2 int) int {
	if n1 > n2 {
		return n2
	}
	return n1
}

func readStdin() (out string, cols, rows int) {
	var builder strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		slen := len([]rune(s))
		cols = max(cols, slen)
		rows++
		builder.WriteString(s)
		builder.WriteString("\n")
	}
	return builder.String(), cols, rows
}

func main() {
	var v bool
	var w, h int
	flag.BoolVar(&v, "v", false, "show version")
	flag.IntVar(&w, "w", 40, "width of the output device (positive number)")
	flag.IntVar(&h, "h", 25, "height of the output device (positive number)")
	flag.Parse()

	if v {
		fmt.Println("version:", version)
		return
	}

	if w <= 0 || h <= 0 {
		flag.Usage()
		os.Exit(2)
	}

	if flag.NArg() == 0 {
		// TODO Implement reading of Stdin
	} else {
		// TODO Implement reading of files flag.Arg(n)
	}
	in, cols, rows := readStdin()
	w = min(w, rows)
	h = min(h, cols)
	ss := gaga.Vert(in, w, h)
	if len(ss) > 0 {
		fmt.Print(ss[0])
		for i := 1; i < len(ss); i++ {
			fmt.Println()
			fmt.Print(ss[i])
		}
	}
}
