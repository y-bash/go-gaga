package gaga

import (
	//	"fmt"
	"strings"
	"testing"
)

type VertTest struct {
	w   int
	h   int
	in  []string
	out []string
}

var verttests = []VertTest{
	0: {3, 3,
		[]string{
			"",
		},
		[]string{},
	},
	1: {3, 3,
		[]string{
			"a",
		},
		[]string{
			"     a",
			"",
			"",
			"",
		},
	},
	2: {3, 3,
		[]string{
			"a",
			"",
		},
		[]string{
			"     a",
			"",
			"",
			"",
		},
	},
	3: {3, 3,
		[]string{
			"ab",
		},
		[]string{
			"     a",
			"     b",
			"",
			"",
		},
	},
	4: {3, 3,
		[]string{
			"ab",
			"",
		},
		[]string{
			"     a",
			"     b",
			"",
			"",
		},
	},
	5: {3, 3,
		[]string{
			"abc",
		},
		[]string{
			"     a",
			"     b",
			"     c",
			"",
		},
	},
	6: {3, 3,
		[]string{
			"abc",
			"",
		},
		[]string{
			"     a",
			"     b",
			"     c",
			"",
		},
	},
	7: {3, 3,
		[]string{
			"abcd",
		},
		[]string{
			"   d a",
			"     b",
			"     c",
			"",
		},
	},
	8: {3, 3,
		[]string{
			"abcde",
		},
		[]string{
			"   d a",
			"   e b",
			"     c",
			"",
		},
	},
	9: {3, 3,
		[]string{
			"abcdef",
		},
		[]string{
			"   d a",
			"   e b",
			"   f c",
			"",
		},
	},
	10: {3, 3,
		[]string{
			"abcdefg",
		},
		[]string{
			" g d a",
			"   e b",
			"   f c",
			"",
		},
	},
	11: {3, 3,
		[]string{
			"abcdefgh",
		},
		[]string{
			" g d a",
			" h e b",
			"   f c",
			"",
		},
	},
	12: {3, 3,
		[]string{
			"abcdefgh",
			"",
		},
		[]string{
			" g d a",
			" h e b",
			"   f c",
			"",
		},
	},
	13: {3, 3,
		[]string{
			"abcdefghi",
		},
		[]string{
			" g d a",
			" h e b",
			" i f c",
			"",
		},
	},
	14: {3, 3,
		[]string{
			"abcdefghi",
			"",
		},
		[]string{
			" g d a",
			" h e b",
			" i f c",
			"",
		},
	},
	15: {3, 3,
		[]string{
			"abcdefghij",
		},
		[]string{
			" g d a",
			" h e b",
			" i f c",
			"",
			"     j",
			"",
			"",
			"",
		},
	},
	16: {3, 3,
		[]string{
			"abcdefghij",
			"",
		},
		[]string{
			" g d a",
			" h e b",
			" i f c",
			"",
			"     j",
			"",
			"",
			"",
		},
	},
	17: {3, 3,
		[]string{
			"ab",
			"c",
		},
		[]string{
			"   c a",
			"     b",
			"",
			"",
		},
	},
	18: {3, 3,
		[]string{
			"abc",
			"de",
			"f",
		},
		[]string{
			" f d a",
			"   e b",
			"     c",
			"",
		},
	},
	19: {3, 3,
		[]string{
			"abcd",
			"efg",
			"hi",
			"j",
		},
		[]string{
			" e d a",
			" f   b",
			" g   c",
			"",
			"   j h",
			"     i",
			"",
			"",
		},
	},
	20: {3, 3,
		[]string{
			"abcd",
			"efghijk",
			"lm",
			"n",
		},
		[]string{
			" e d a",
			" f   b",
			" g   c",
			"",
			" l k h",
			" m   i",
			"     j",
			"",
			"     n",
			"",
			"",
			"",
		},
	},
	21: {5, 6,
		[]string{
			"閑さや",
			"岩にしみ入る",
			"蝉の声",
			"",
			"芭蕉",
		},
		[]string{
			"芭  蝉岩閑",
			"蕉  のにさ",
			"    声しや",
			"      み",
			"      入",
			"      る",
			"",
		},
	},
	22: {4, 5,
		[]string{
			"好きなもの",
			"いちご",
			"珈琲",
			"花美人",
			"懐手して",
			"宇宙見物",
			"",
			"寺田寅彦",
		},
		[]string{
			"花珈い好",
			"美琲ちき",
			"人  ごな",
			"      も",
			"      の",
			"",
			"寺  宇懐",
			"田  宙手",
			"寅  見し",
			"彦  物て",
			"",
			"",
		},
	},
	23: {6, 7,
		[]string{
			"あいうえお",
			"abcdefg",
			"かきくけこ",
			"さしすせそ",
			"hijklmn",
			"opq   \n",
		},
		[]string{
			" o hさか aあ",
			" p iしき bい",
			" q jすく cう",
			"   kせけ dえ",
			"   lそこ eお",
			"   m     f",
			"   n     g",
			"",
		},
	},
	24: {1, 1,
		[]string{
			"a",
			"",
		},
		[]string{
			" a",
			"",
		},
	},
	25: {1, 1,
		[]string{
			"\ra",
			"",
		},
		[]string{
			" a",
			"",
		},
	},
}

func TestVert(t *testing.T) {
	for i, tt := range verttests {
		in := strings.Join(tt.in, "\n")
		exp := strings.Join(tt.out, "\n")
		ss := Vert(in, tt.w, tt.h)
		got := strings.Join(ss, "\n")
		if got != exp {
			t.Errorf("#%d Vert(in,%d,%d):\nin=(\n%s\n),\nexpected=(\n%s\n),\ngot=(\n%s\n)",
				i, tt.w, tt.h, in, exp, got)
		}
	}
}

func TestVertCatchesOverflow(t *testing.T) {
	tests := [...]struct {
		w int
		h int
	}{
		0: {0, 1},
		1: {1, 0},
		2: {-1, 1},
		3: {1, -1},
		4: {-2147483647, 1},
	}

	for i, tt := range tests {
		ss := Vert("foo,bar,baz", tt.w, tt.h)
		if len(ss) != 0 {
			t.Errorf("#%d expected zero length string slice, got %q", i, ss)
		}
	}
}

func benchmarkVert(b *testing.B, s string, w, h int) {
	for i := 0; i < b.N; i++ {
		Vert(s, w, h)
	}
}

func BenchmarkVert100(b *testing.B) {
	s := strings.Repeat("aあbいcう", 100)
	benchmarkVert(b, s, 40, 25)
}

func BenchmarkVert1000(b *testing.B) {
	s := strings.Repeat("aあbいcう", 1000)
	benchmarkVert(b, s, 40, 25)
}

func BenchmarkVert10000(b *testing.B) {
	s := strings.Repeat("aあbいcう", 10000)
	benchmarkVert(b, s, 40, 25)
}
