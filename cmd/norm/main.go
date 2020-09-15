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

func read(f io.Reader) string {
	var sb strings.Builder
	sb.Grow(1024)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		sb.WriteString(sc.Text())
		sb.WriteString("\n")
	}
	return sb.String()
}

func readfiles(paths []string) (out []string, err error) {
	if len(paths) == 0 {
		out = []string{read(os.Stdin)}
		return
	}
	for _, path := range paths {
		var f *os.File
		f, err = os.Open(path)
		if err != nil {
			return
		}
		defer f.Close()
		out = append(out, read(f))
	}
	return
}

func normstrs(f io.Writer, in []string, flag gaga.NormFlag) error {
	if len(in) <= 0 {
		return nil
	}
	n, err := gaga.Norm(flag)
	if err != nil {
		return err
	}
	fmt.Fprint(f, n.String(in[0]))
	for i := 1; i < len(in); i++ {
		fmt.Fprintln(f)
		fmt.Fprint(f, n.String(in[i]))
	}
	return nil
}

func main() {
	var v, h bool
	var normflag string
	flag.BoolVar(&v, "v", false, "show version")
	flag.BoolVar(&h, "h", false, "show help")
	flag.StringVar(&normflag, "flag", "Fold", "normalization flag")
	flag.Parse()
	if v {
		fmt.Println("version:", version)
		return
	}
	if h {
		flag.Usage()
		return
	}
	nf, err := gaga.ParseNormFlag(normflag)
	if err != nil {
		log.Fatal(err)
	}
	var ss []string
	ss, err = readfiles(flag.Args())
	if err != nil {
		log.Fatal(err)
	}
	err = normstrs(os.Stdout, ss, nf)
	if err != nil {
		log.Fatal(err)
	}
}
