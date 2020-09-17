package main

import (
	"flag"
	"fmt"
	"strings"
)

var version = "v0.0.0" // set value bygo build -ldflags

func filter(in []string) string {
	out := strings.Join(in, " ")
	return strings.Replace(out, "\\n", "\n", -1)
}

func main() {
	var v, h bool
	flag.BoolVar(&v, "v", false, "show version")
	flag.BoolVar(&h, "h", false, "show help")
	flag.Parse()
	if v {
		fmt.Println("version:", version)
		return
	}
	if h {
		flag.Usage()
		return
	}
	args := flag.Args()
	out := filter(args)
	fmt.Println(out)
}
