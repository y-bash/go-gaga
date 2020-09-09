package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/y-bash/go-gaga"
	"io"
	"log"
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

func read(r io.Reader, maxcol int) (out string, row, col int) {
	if maxcol <= 0 {
		maxcol = 1
	}
	var builder strings.Builder
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s := scanner.Text()
		l := len([]rune(s))
		col = max(col, min(maxcol, l))
		if l == 0 {
			row++
		} else {
			row += (l + maxcol - 1) / maxcol
		}
		builder.WriteString(s)
		builder.WriteString("\n")
	}
	out = builder.String()
	return
}

func writeString(f io.Writer, s string, w, h int) {
	ss := gaga.Vert(s, w, h)
	if len(ss) > 0 {
		fmt.Fprint(f, ss[0])
		for i := 1; i < len(ss); i++ {
			fmt.Fprintln(f)
			fmt.Fprint(f, ss[i])
		}
	}
}

func writeSlice(f io.Writer, in []string, w, h int) {
	if len(in) > 0 {
		writeString(f, in[0], w, h)
		for i := 1; i < len(in); i++ {
			fmt.Fprintln(f)
			writeString(f, in[i], w, h)
		}
	}
}

func main() {
	var ver, help bool
	var maxw, maxh int
	flag.BoolVar(&ver, "v", false, "show version")
	flag.BoolVar(&help, "h", false, "show help")
	flag.IntVar(&maxw, "width", 40, "maximum width of output")
	flag.IntVar(&maxh, "height", 25, "maximum height of output")
	flag.Parse()
	if ver {
		fmt.Println("version:", version)
		return
	}
	if help {
		flag.Usage()
		return
	}
	if maxw <= 0 || maxh <= 0 {
		flag.Usage()
		os.Exit(2)
	}
	var ss []string
	var w, h int
	if flag.NArg() == 0 {
		var s string
		s, w, h = read(os.Stdin, maxh)
		ss = []string{s}
	} else {
		args := flag.Args()
		for _, path := range args {
			f, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			s, r, c := read(f, maxh)
			w = max(w, r)
			h = max(h, c)
			ss = append(ss, s)
		}
	}
	w = min(maxw, max(w, 1))
	h = min(maxh, max(h, 1))
	writeSlice(os.Stdout, ss, w, h)
}
