// +build ignore

//
// usage:
//
// go run genucdcsv.go -output ucd.csv
//

package main

import (
	"bytes"
	"flag"
	"github.com/y-bash/go-gaga/gen/lib"
	"io/ioutil"
	"log"
	"os"
)

var filename = flag.String("output", "", "output file name")

func main() {
	flag.Parse()
	var err error
	var buf bytes.Buffer
	err = lib.GenUCD(&buf)
	if err != nil {
		log.Fatal(err)
	}
	if len(*filename) == 0 {
		_, err = os.Stdout.Write(buf.Bytes())
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	err = ioutil.WriteFile(*filename, buf.Bytes(), 0664)
	if err != nil {
		log.Fatal(err)
	}
}
