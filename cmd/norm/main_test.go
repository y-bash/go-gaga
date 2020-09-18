package main

import (
	"bytes"
	"github.com/y-bash/go-gaga"
	"io/ioutil"
	"log"
	"strings"
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
		// Supports Windows environment where git config core.autocrlf = true
		wantS := strings.Replace(string(want), "\r", "", -1)

		var buf bytes.Buffer
		ss, err := readfiles([]string{tt.in})
		if err != nil {
			log.Fatal(err)
		}
		err = normstrs(&buf, ss, tt.flag)
		if err != nil {
			log.Fatal(err)
		}
		have := buf.Bytes()
		haveS := string(have)

		if haveS != wantS {
			t.Errorf("#%d\nin(len=%d):\n%s,\nhave(len=%d):\n%s,\nwant(len=%d):\n%s",
				i, len([]rune(ss[0])), ss[0], len([]rune(haveS)), haveS, len([]rune(wantS)), wantS)
		}
	}
}
