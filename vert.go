package gaga

import (
	"github.com/mattn/go-runewidth"
	"strings"
)

func max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}
	return n2
}

func min(n1, n2 int) int {
	if n1 > n2 {
		return n2
	}
	return n1
}

func wordwrap(in string, w int) (out [][]rune) {
	rs := []rune(in)
	out = make([][]rune, 0, len(rs) * 2 / w)
	for i := 0; i < len(rs); {
		if rs[i] != '\n' && runewidth.RuneWidth(rs[i]) <= 0 {
			rs = append(rs[:i], rs[i+1:]...)
			continue
		}
		i++
	}
	width := 0
	start := 0
	for i := 0; i < len(rs); i++ {
		if rs[i] == '\n' {
			out = append(out, rs[start:i])
			start = i + 1
			width = 0
			continue
		}
		if width >= w {
			out = append(out, rs[start:i])
			start = i
			width = 1
			continue
		}
		width++
	}
	if width > 0 {
		out = append(out, rs[start:len(rs)])
	}
	return
}

func vert(in [][]rune, w, h int) (out []string) {
	if w <= 0 || h <= 0 {
		return
	}
	rs := make([]rune, 0, w*h+h)
	n := (len(in) + w - 1) / w
	for i := 0; i < n; i++ {
		tail := i * w
		head := tail + w - 1
		for j := 0; j < h; j++ {
			for k := head; k >= tail; k-- {
				if k >= len(in) || j >= len(in[k]) {
					rs = append(rs, ' ', ' ')
					continue
				}
				if runewidth.RuneWidth(in[k][j]) == 1 {
					rs = append(rs, ' ')
				}
				rs = append(rs, in[k][j])
			}
			k := len(rs) - 1
			for ; k >= 0 && rs[k] == ' '; k-- {
			}
			rs = rs[:k+1]
			rs = append(rs, '\n')
		}
		out = append(out, string(rs))
		rs = rs[:0]
	}
	return
}

func estimateSize(s string, maxh int) (w, h int) {
	if maxh <= 0 {
		maxh = 1
	}
	height := 0
	rs := []rune(s)
	for _, r := range rs {
		if r == '\n' {
			w++
			h = max(h, height)
			height = 0
			continue
		}
		if runewidth.RuneWidth(r) <= 0 {
			continue
		}
		if height >= maxh {
			w++
			h = maxh
			height = 1
			continue
		}
		height++
	}
	if height > 0 {
		w++
		h = max(h, height)
	}
	return
}

// VertFix returns out that is a vertical conversion of the in.
// If the converted string fits in a matrix of size w and h,
// out is a string slice with one element, if not, out is a
// string slice with multiple elements.
// If in contains half-width or narrow-width characters, space
// is added to the left of it.
// TODO renew comments
func VertFix(in string, w int, h int) []string {
	if w <= 0 || h <= 0 {
		return []string{}
	}
	rss := wordwrap(in, h)
	return vert(rss, w, h)
}

func VertShrink(in string, w, h int) []string {
	ew, eh := estimateSize(in, h)
	w = min(w, ew)
	h = min(h, eh)
	return VertFix(in, w, h)
}

func Vert(s string) string {
	ss := VertShrink(s, 40, 25)
	return strings.Join(ss, "\n")
}
