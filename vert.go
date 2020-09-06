package gaga

import (
	"github.com/mattn/go-runewidth"
	"strings"
)

type matrix [][]string

func makeMatrix(w int, h int) matrix {
	ss := make([][]string, h)
	for i := 0; i < h; i++ {
		ss[i] = make([]string, w)
	}
	return matrix(ss)
}

func (m matrix) init() {
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			m[i][j] = "  "
		}
	}
}

func (m matrix) String() string {
	var sb strings.Builder
	for _, row := range m {
		for _, cell := range row {
			sb.WriteString(cell)
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (m matrix) vert(s string, start int) (string, int) {
	m.init()
	rs := []rune(s)
	x := len(m[0]) - 1
	y := 0
	var i int
	for i = start; i < len(rs); i++ {
		switch runewidth.RuneWidth(rs[i]) {
		case 0:
			if rs[i] != '\n' {
				continue
			}
		case 1:
			m[y][x] = string([]rune{' ', rs[i]})
		case 2:
			m[y][x] = string([]rune{rs[i]})
		default:
			panic("vert: RuneWidth returned unexpected value")
		}
		if rs[i] == '\n' || y >= len(m)-1 {
			y = 0
			x--
			if rs[i] != '\n' && i < len(rs)-1 && rs[i+1] == '\n' {
				i++
			}
			if x < 0 {
				return m.String(), i + 1
			}
		} else {
			y++
		}
	}
	return m.String(), i
}

func Vert(in string, w int, h int) (out []string) {
	if w <= 0 || h <= 0 {
		panic("vert: non-positive Vert size")
	}
	inlen := len([]rune(in))
	m := makeMatrix(w, h)
	for pos := 0; pos < inlen; {
		var s string
		s, pos = m.vert(in, pos)
		out = append(out, s)
	}
	return out
}
