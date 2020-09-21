package main

import (
	"testing"
)

type FilterTest struct {
	in  []string
	out string
}

var filtertests = []FilterTest{
	0: {[]string{"aaa", "bb", "c"}, "aaa bb c"},
	1: {[]string{"aaa\\nbb\\nc"}, "aaa\nbb\nc"},
}

func TestFilter(t *testing.T) {
	for i, tt := range filtertests {
		have := filter(tt.in)
		if have != tt.out {
			t.Errorf("#%d filter(%v)\nhave:\n%s,\nwant:\n%s",
				i, tt.in, have, tt.out)
		}
	}
}
