package main

import (
	"bytes"
	"strings"
	"testing"
)

type CmdVertReadTest struct {
	s      string
	maxcol int
	row    int
	col    int
}

var cmdvertreadtests = []CmdVertReadTest{
	0:  {"a\nbb\nccc\n", 1, 6, 1},
	1:  {"a\nbb\nccc\n", 2, 4, 2},
	2:  {"a\nbb\nccc\n", 3, 3, 3},
	3:  {"a\nbb\nccc\n", 4, 3, 3},
	4:  {"a\nbb\n1122334\n", 2, 6, 2},
	5:  {"a\nbb\n11223344\n", 2, 6, 2},
	6:  {"a\nbb\n112233445\n", 2, 7, 2},
	7:  {"a\n112\nccc\n", 2, 5, 2},
	8:  {"a\n1122\nccc\n", 2, 5, 2},
	9:  {"a\n11223\nccc\n", 2, 6, 2},
	10: {"11\nbb\nccc\n", 2, 4, 2},
	11: {"112\nbb\nccc\n", 2, 5, 2},
	12: {"1122\nbb\nccc\n", 2, 5, 2},
	13: {"11223\nbb\nccc\n", 2, 6, 2},
	14: {"ａ\nｂb\nｃcc\n", 1, 6, 1},
	15: {"a\nbｂ\ncｃc\n", 2, 4, 2},
	16: {"a\nbb\nccｃ\n", 3, 3, 3},
	17: {"ａ\nｂｂ\nｃｃｃ\n", 4, 3, 3},
	18: {"", 2, 0, 0},
	19: {"\n", 2, 1, 0},
	20: {"\n\n", 2, 2, 0},
	21: {"\n\n\n", 2, 3, 0},
	22: {"a", 2, 1, 1},
	23: {"a\n", 2, 1, 1},
	24: {"a\n\n", 2, 2, 1},
	25: {"a\n\n\n", 2, 3, 1},
	26: {"a\na", 2, 2, 1},
	27: {"a\nab", 2, 2, 2},
	28: {"a\n\na", 2, 3, 1},
	29: {"a\n\nab", 2, 3, 2},
	30: {"a\n\nabc", 2, 4, 2},
	31: {"a\n\n\na", 2, 4, 1},
	32: {"a\n\n\nab", 2, 4, 2},
	33: {"a\n\n\nabc", 2, 5, 2},
	34: {"a\n\n\nabc", -1, 6, 1},
	35: {"a\n\n\nabc", 0, 6, 1},
	36: {"a\n\n\nabc", 1, 6, 1},
}

func TestCmdVertRead(t *testing.T) {
	for i, tt := range cmdvertreadtests {
		r := strings.NewReader(tt.s)
		out, row, col := read(r, tt.maxcol)
		haves := strings.TrimRight(out, "\n")
		wants := strings.TrimRight(tt.s, "\n")
		if haves != wants || row != tt.row || col != tt.col {
			t.Errorf("#%d read(%q, %d)\n"+
				"\thave: row=%2d, col=%2d, out=%q\n\twant: row=%2d, col=%2d, out=%q",
				i, wants, tt.maxcol, row, col, haves, tt.row, tt.col, tt.s)
		}
	}
}

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
