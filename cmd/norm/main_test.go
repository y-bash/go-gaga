package main

import (
	"bytes"
	"github.com/y-bash/go-gaga"
	"io/ioutil"
	"log"
	"testing"
)

type CmdNormReadWriteTest struct {
	in   string
	out  string
	flag gaga.NormFlag
}

var cmdnormreadwritetests = []CmdNormReadWriteTest{
	0: {"testdata/norm_in01.txt", "testdata/norm_out01.txt", gaga.HiraganaToKatakana},
	1: {"testdata/norm_in02.txt", "testdata/norm_out02.txt", gaga.LatinToWide | gaga.AlphaToUpper},
}

func TestCmdNormReadWrite(t *testing.T) {
	for i, tt := range cmdnormreadwritetests {
		want, err := ioutil.ReadFile(tt.out)
		if err != nil {
			log.Fatal(err)
		}
		var have bytes.Buffer
		ss, err := readfiles([]string{tt.in})
		if err != nil {
			log.Fatal(err)
		}
		err = normstrs(&have, ss, tt.flag)
		if err != nil {
			log.Fatal(err)
		}
		if string(have.Bytes()) != string(want) {
			t.Errorf("#%d\nin:\n%s,\nhave:\n%s,\nwant:\n%s",
				i, ss[0], have.Bytes(), want)
		}
	}
}
