package main

import (
	"bytes"
	"testing"
)

type CmdVertWriteSliceTest struct {
	ss  []string
	w   int
	h   int
	out string
}

var cmdvertwriteslicetests = []CmdVertWriteSliceTest{
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

func TestCmdVertWriteSlice(t *testing.T) {
	for i, tt := range cmdvertwriteslicetests {
		var buf bytes.Buffer
		writeSlice(&buf, tt.ss, tt.w, tt.h)
		out := string(buf.Bytes())
		if out != tt.out {
			t.Errorf("#%d writeSlice(buf, %q, %d, %d) = %q, want: %q",
				i, tt.ss, tt.w, tt.h, out, tt.out)
		}
	}
}
