package gaga

import (
	"bufio"
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

func ss2s(ss []string) string {
	if len(ss) <= 0 {
		return ""
	}
	var size int
	for _, s := range ss {
		size += len([]rune(s)) + 1
	}
	var sb strings.Builder
	sb.Grow(size)
	sb.WriteString(ss[0])
	for i := 1; i < len(ss); i++ {
		sb.WriteRune('\n')
		sb.WriteString(ss[i])
	}
	return sb.String()
}

func wordwrap(in string, w int) (out [][]rune) {
	r := strings.NewReader(in)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		rs := []rune(sc.Text())
		for i := 0; i < len(rs); {
			if runewidth.RuneWidth(rs[i]) > 0 {
				i++
			} else {
				rs = append(rs[:i], rs[i+1:]...)
			}
		}
		if len(rs) > 0 {
			for i := 0; i < len(rs); i += w {
				e := min(i+w, len(rs))
				out = append(out, rs[i:e])
			}
		} else {
			out = append(out, []rune{})
		}
	}
	return
}

func vert(in [][]rune, w, h int) (out []string) {
	if w <= 0 || h <= 0 {
		return
	}
	n := (len(in) + w - 1) / w
	for i := 0; i < n; i++ {
		tail := i * w
		head := tail + w - 1
		var rsb strings.Builder
		rsb.Grow(w * h * 4)
		for j := 0; j < h; j++ {
			var csb strings.Builder
			csb.Grow(h * 4)
			for k := head; k >= tail; k-- {
				if k >= len(in) || j >= len(in[k]) {
					csb.WriteString("  ")
					continue
				}
				if runewidth.RuneWidth(in[k][j]) == 1 {
					csb.WriteRune(' ')
				}
				csb.WriteRune(in[k][j])
			}
			s := strings.TrimRight(csb.String(), " ")
			rsb.WriteString(s)
			rsb.WriteRune('\n')
		}
		out = append(out, rsb.String())
	}
	return
}

func ScanNowrap(s string) (w, h int) {
	r := strings.NewReader(s)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s := scanner.Text()
		rs := []rune(s)
		l := 0
		for _, r := range rs {
			if runewidth.RuneWidth(r) > 0 {
				l++
			}
		}

		h = max(h, l)
		w++
	}
	return
}

func ScanWordwrap(s string, maxh int) (w, h int) {
	if maxh <= 0 {
		maxh = 1
	}
	r := strings.NewReader(s)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s := scanner.Text()
		rs := []rune(s)
		l := 0
		for _, r := range rs {
			if runewidth.RuneWidth(r) > 0 {
				l++
			}
		}
		h = max(h, min(maxh, l))
		if l == 0 {
			w++
		} else {
			w += (l + maxh - 1) / maxh
		}
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
func VertFixStrings(in string, w int, h int) []string {
	if w <= 0 || h <= 0 {
		return []string{}
	}
	rss := wordwrap(in, h)
	return vert(rss, w, h)
}

func VertToFitStrings(in string, w, h int) []string {
	sw, sh := ScanWordwrap(in, h)
	w = min(w, sw)
	h = min(h, sh)
	return VertFixStrings(in, w, h)
}

func VertFix(s string, w, h int) string {
	ss := VertFixStrings(s, w, h)
	return ss2s(ss)
}

func VertToFit(s string, w, h int) string {
	ss := VertToFitStrings(s, w, h)
	return ss2s(ss)
}

func VertNowrap(s string) string {
	w, h := ScanNowrap(s)
	ss := VertFixStrings(s, w, h)
	return ss2s(ss)
}

func Vert(s string) string {
	return VertToFit(s, 40, 25)
}
