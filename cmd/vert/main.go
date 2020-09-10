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

func read(r io.Reader) string {
	var sb strings.Builder
	sb.Grow(1024)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		sb.WriteString(sc.Text())
		sb.WriteString("\n")
	}
	return sb.String()
}

func writeString(f io.Writer, s string, w, h int) {
	ss := gaga.VertToFitStrings(s, w, h)
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
	var v, h bool
	var width, height int
	flag.BoolVar(&v, "v", false, "show version")
	flag.BoolVar(&h, "h", false, "show help")
	flag.IntVar(&width, "width", 40, "maximum width of output")
	flag.IntVar(&height, "height", 25, "maximum height of output")
	flag.Parse()
	if v {
		fmt.Println("version:", version)
		return
	}
	if h {
		flag.Usage()
		return
	}
	if width <= 0 || height <= 0 {
		flag.Usage()
		os.Exit(2)
	}
	var ss []string
	if flag.NArg() == 0 {
		ss = []string{read(os.Stdin)}
	} else {
		args := flag.Args()
		for _, path := range args {
			f, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			ss = append(ss, read(f))
		}
	}
	writeSlice(os.Stdout, ss, width, height)
}
