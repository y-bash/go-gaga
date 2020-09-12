package gaga

import (
	//	"fmt"
	"strings"
	"testing"
)

type EstimateSizeTest struct {
	s    string
	maxh int
	w    int
	h    int
}

var estimatesizetests = []EstimateSizeTest{
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

func TestEstimateSize(t *testing.T) {
	for i, tt := range estimatesizetests {
		w, h := estimateSize(tt.s, tt.maxh)
		if w != tt.w || h != tt.h {
			t.Errorf("#%d EstimateSize(%q, %d)=(w: %d, h: %d), want(w: %d, h: %d)",
				i, tt.s, tt.maxh, w, h, tt.w, tt.h)
		}
	}
}

type VertFixTest struct {
	w   int
	h   int
	in  []string
	out []string
}

var vertfixtests = []VertFixTest{
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
	26: {1, 1,
		[]string{
			"\r\t\r\ta\r\t\r\t",
			"",
		},
		[]string{
			" a",
			"",
		},
	},
	27: {2, 2,
		[]string{
			"ab",
			"cd",
			"",
		},
		[]string{
			" c a",
			" d b",
			"",
		},
	},
	28: {2, 2,
		[]string{
			"\r\t\r\tab",
			"cd",
			"",
		},
		[]string{
			" c a",
			" d b",
			"",
		},
	},
	29: {2, 2,
		[]string{
			"\r\t\r\ta\r\t\r\tb",
			"cd",
			"",
		},
		[]string{
			" c a",
			" d b",
			"",
		},
	},
	30: {2, 2,
		[]string{
			"\r\t\r\ta\r\t\r\tb\r\t\r\t",
			"cd",
			"",
		},
		[]string{
			" c a",
			" d b",
			"",
		},
	},
	31: {2, 2,
		[]string{
			"ab",
			"\r\t\r\tcd",
			"",
		},
		[]string{
			" c a",
			" d b",
			"",
		},
	},
	32: {2, 2,
		[]string{
			"ab",
			"\r\t\r\tcd\r\t\r\t",
			"",
		},
		[]string{
			" c a",
			" d b",
			"",
		},
	},
	33: {2, 2,
		[]string{
			"ab",
			"\r\t\r\tc\r\t\r\td\r\t\r\t",
			"",
		},
		[]string{
			" c a",
			" d b",
			"",
		},
	},
}

func TestVertFix(t *testing.T) {
	for i, tt := range vertfixtests {
		in := strings.Join(tt.in, "\n")
		exp := strings.Join(tt.out, "\n")
		ss := VertFix(in, tt.w, tt.h)
		got := strings.Join(ss, "\n")
		if got != exp {
			t.Errorf("#%d VertFix(in,%d,%d):\nin=(\n%s\n),\nexpected=(\n%s\n),\ngot=(\n%s\n)",
				i, tt.w, tt.h, in, exp, got)
		}
	}
}

func TestVertFixCatchesOverflow(t *testing.T) {
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
		ss := VertFix("foo,bar,baz", tt.w, tt.h)
		if len(ss) != 0 {
			t.Errorf("#%d expected zero length string slice, got %q", i, ss)
		}
	}
}

type VertShrinkTest struct {
	s   string
	w   int
	h   int
	out []string
}

var vertshrinktests = []VertShrinkTest{
	0:  {"", 0, 0, []string{}},
	1:  {"a", 1, 1, []string{" a\n"}},
	2:  {"a\nbb\nccc\n", 4, 4, []string{" c b a\n c b\n c\n"}},
	3:  {"a\nbb\nccc\n", 4, 3, []string{" c b a\n c b\n c\n"}},
	4:  {"a\nbb\nccc\n", 3, 4, []string{" c b a\n c b\n c\n"}},
	5:  {"a\nbb\nccc\n", 3, 3, []string{" c b a\n c b\n c\n"}},
	6:  {"a\nbb\nccc\n", 3, 2, []string{" c b a\n c b\n", "     c\n\n"}},
	7:  {"a\nbb\nccc\n", 2, 3, []string{" b a\n b\n\n", "   c\n   c\n   c\n"}},
	8:  {"a\nbb\nccc\n", 2, 2, []string{" b a\n b\n", " c c\n   c\n"}},
	9:  {"a\nbb\nccc\n", 2, 1, []string{" b a\n", " c b\n", " c c\n"}},
	10: {"a\nbb\nccc\n", 1, 1, []string{" a\n", " b\n", " b\n", " c\n", " c\n", " c\n"}},
	11: {"a\nbb\nccc\n", 1, 0, []string{}},
	12: {"a\nbb\nccc\n", 0, 1, []string{}},
	13: {"a\nbb\nccc\n", 0, 0, []string{}},
	14: {"a\nbb\nccc\n", 0, -1, []string{}},
	15: {"a\nbb\nccc\n", -1, -1, []string{}},
	16: {"a\nbb\nccc\n", 100, 1000, []string{" c b a\n c b\n c\n"}},
	17: {"1234567890\n\n12345\n", 10, 10, []string{
		" 1   1\n 2   2\n 3   3\n 4   4\n 5   5\n     6\n     7\n     8\n     9\n     0\n"}},
	18: {"1234567890\n\n12345\n", 3, 3, []string{
		" 7 4 1\n 8 5 2\n 9 6 3\n", " 1   0\n 2\n 3\n", "     4\n     5\n\n"}},
}

func TestVertShrink(t *testing.T) {
	for i, tt := range vertshrinktests {
		ss := VertShrink(tt.s, tt.w, tt.h)
		if len(ss) != len(tt.out) {
			t.Errorf("#%d VertShrink(%q, %d, %d),\n\thave:(%q),\n\twant:(%q)",
				i, tt.s, tt.w, tt.h, ss, tt.out)
			continue
		}
		for j, s := range ss {
			if s != tt.out[j] {
				t.Errorf("#%d VertShrink(%q, %d, %d)[%d],\n\thave:%q,\n\twant:%q",
					i, tt.s, tt.w, tt.h, j, s, tt.out[j])
			}
		}
	}
}

func benchmarkVert(b *testing.B, s string) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Vert(s)
	}
	b.StopTimer()
}

func BenchmarkVert100(b *testing.B) {
	s := strings.Repeat("aあbいcう", 100)
	benchmarkVert(b, s)
}

func BenchmarkVert1000(b *testing.B) {
	s := strings.Repeat("aあbいcう", 1000)
	benchmarkVert(b, s)
}

func BenchmarkVert10000(b *testing.B) {
	s := strings.Repeat("aあbいcう", 10000)
	benchmarkVert(b, s)
}

