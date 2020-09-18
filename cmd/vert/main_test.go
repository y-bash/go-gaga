package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

type CmdVertStrsTest struct {
	ss  []string
	w   int
	h   int
	out string
}

var cmdvertvertstrstests = []CmdVertStrsTest{
	0: {[]string{"a\nbb\ncc\n"}, 3, 2,
		" c b a\n" +
			" c b\n"},
	1: {[]string{"a\nbb\nccc\n"}, 3, 3,
		" c b a\n" +
			" c b\n" +
			" c\n"},
	2: {[]string{"a\nbb\ncc\na\nbb\ncc\n"}, 3, 2,
		" c b a\n" +
			" c b\n" +
			"\n" +
			" c b a\n" +
			" c b\n"},
	3: {[]string{"a\nbb\nccc\na\nbb\ncc\n"}, 3, 3,
		" c b a\n" +
			" c b\n" +
			" c\n" +
			"\n" +
			" c b a\n" +
			" c b\n" +
			"\n"},
	4: {[]string{"a\nbb\ncc\na\nbb\nccc\n"}, 3, 3,
		" c b a\n" +
			" c b\n" +
			"\n" +
			"\n" +
			" c b a\n" +
			" c b\n" +
			" c\n"},
	5: {[]string{"a\nbb\ncc\na\nbb\nccc\n", "a\nbb\ncc\na\nbb\nccc\n"}, 3, 3,
		" c b a\n" +
			" c b\n" +
			"\n" +
			"\n" +
			" c b a\n" +
			" c b\n" +
			" c\n" +
			"\n" +
			" c b a\n" +
			" c b\n" +
			"\n" +
			"\n" +
			" c b a\n" +
			" c b\n" +
			" c\n"},
}

func TestCmdVertStrs(t *testing.T) {
	for i, tt := range cmdvertvertstrstests {
		var buf bytes.Buffer
		vertstrs(&buf, tt.ss, tt.w, tt.h)
		out := string(buf.Bytes())
		if out != tt.out {
			t.Errorf("#%d vertstrs(buf, %q, %d, %d) = %q, want: %q",
				i, tt.ss, tt.w, tt.h, out, tt.out)
		}
	}
}

type CmdVertReadWrieTest struct {
	in  string
	out string
}

var cmdvertreadwritetests = []CmdVertReadWrieTest{
	0: {"testdata/vert_in01.txt", "testdata/vert_out01.txt"},
	1: {"testdata/vert_in02.txt", "testdata/vert_out02.txt"},
	2: {"testdata/vert_in03.txt", "testdata/vert_out03.txt"},
	3: {"testdata/vert_in04.txt", "testdata/vert_out04.txt"},
}

func TestCmdVertReadWrite(t *testing.T) {
	for i, tt := range cmdvertreadwritetests {
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
		vertstrs(&buf, ss, 40, 25)
		have := buf.Bytes()
		haveS := string(have)

		if haveS != wantS {
			t.Errorf("#%d\nin(len=%d):\n%s,\nhave(len=%d):\n%s,\nwant(len=%d):\n%s",
				i, len([]rune(ss[0])), ss[0], len([]rune(haveS)), haveS, len([]rune(wantS)), wantS)
		}
	}
}
