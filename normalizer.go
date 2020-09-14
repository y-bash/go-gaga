package gaga

import (
	"strings"
)

// Normalizer normalizes the input provided and returns
// the normalized string.
type Normalizer struct {
	flag NormFlag
}

func (n *Normalizer) maybeComposeVom(r1, r2 rune) (rune, vom, bool) {
	if !vom(r2).isVom() {
		r, m := n.NormalizeRune(r1)
		return r, m, false
	}

	c1, ok := findUnichar(r1)
	if !ok {
		return r1, vmNone, false
	}

	if c1.category != ctKanaLetter ||
		c1.voicing == vcVoiced || c1.voicing == vcSemivoiced {
		r, m := n.NormalizeRune(r1)
		return r, m, false
	}

	// TEST_nD7FwQUW knows that NormalizeRune() will definitely return
	// a rune and the vmNone.
	nr1, _ := n.NormalizeRune(r1)

	// TEST_G9amUMTr knows that findUnichar() definitely return a rune
	// and the ok value.
	nc1, _ := findUnichar(nr1)

	switch {
	case vom(r2).isVsm():
		switch {
		case n.flag.has(ComposeVom):
			r, m := nc1.composeVoiced()
			return r, m, true
		case n.flag.has(DecomposeVom):
			r, m := nc1.decomposeVoiced()
			return r, m, true
		default:
			vsm, _ := n.NormalizeRune(r2)
			return nc1.codepoint, vom(vsm), true
		}
	case vom(r2).isSsm():
		switch {
		case n.flag.has(ComposeVom):
			r, m := nc1.composeSemivoiced()
			return r, m, true
		case n.flag.has(DecomposeVom):
			r, m := nc1.decomposeSemivoiced()
			return r, m, true
		default:
			svsm, _ := n.NormalizeRune(r2)
			return nc1.codepoint, vom(svsm), true
		}
	}
	panic("unreachable")
}

// NewNormalizer creates a new Normalizer with specified flag
// (LatinToNarrow etc.). If successful, methods on the returned
// Normalizer can be used for normalization.
func NewNormalizer(flag NormFlag) (*Normalizer, error) {
	err := flag.validate()
	if err != nil {
		return nil, err
	}
	n := Normalizer{flag}
	return &n, nil
}

// SetFlag changes the normalization mode with
// the newly specified flag.
func (n *Normalizer) SetFlag(flag NormFlag) error {
	err := flag.validate()
	if err != nil {
		return err
	}
	n.flag = flag
	return nil
}

// NormalizeRune normalizes the rune according to the current
// normalization mode. Depending on the mode, the voiced or
// semi-voiced sound mark may be separated, so it may return
// multiple runes. but, this function allways returns a rune
// array with 1 or 2 elements, and never returns an array with
// any other number of elements.
func (n *Normalizer) NormalizeRune(r rune) (rune, vom) {
	// TEST_Fc68JR9i knows about the number of elements in
	// the return value of this function
	c, ok := findUnichar(r)
	if !ok {
		return r, vmNone
	}

	switch c.category {
	case ctUndefined:
		return c.codepoint, vmNone

	case ctLatinLetter:
		switch {
		case n.flag.has(AlphaToNarrow):
			c = c.toNarrowC()
		case n.flag.has(AlphaToWide):
			c = c.toWideC()
		}

		switch {
		case n.flag.has(AlphaToUpper):
			return c.toUpperR(), vmNone
		case n.flag.has(AlphaToLower):
			return c.toLowerR(), vmNone
		default:
			return c.codepoint, vmNone
		}

	case ctLatinDigit:
		switch {
		case n.flag.has(DigitToNarrow):
			return c.toNarrowR(), vmNone
		case n.flag.has(DigitToWide):
			return c.toWideR(), vmNone
		default:
			return c.codepoint, vmNone
		}

	case ctLatinSymbol:
		switch {
		case n.flag.has(SymbolToNarrow):
			return c.toNarrowR(), vmNone
		case n.flag.has(SymbolToWide):
			return c.toWideR(), vmNone
		default:
			return c.codepoint, vmNone
		}

	case ctKanaLetter:
		var cc *unichar
		switch c.charCase {
		case ccHiragana:
			switch {
			case n.flag.has(HiraganaToNarrow):
				cc = c.toNarrowC()
			case n.flag.has(HiraganaToKatakana):
				cc = c.toKatakanaC()
			default:
				cc = c
			}
		case ccKatakana:
			switch {
			case n.flag.has(KatakanaToNarrow):
				cc = c.toNarrowC()
			case n.flag.has(KatakanaToHiragana):
				cc = c.toHiraganaC()
			case n.flag.has(KatakanaToWide):
				cc = c.toWideC()
			default:
				cc = c
			}
		default:
			// TEST_gT8YJdBc knows that the program never passes here
			panic("unreachable")
		}

		switch c.voicing {
		case vcUndefined, vcUnvoiced:
			return cc.codepoint, vmNone

		case vcVoiced:
			switch {
			case n.flag.has(ComposeVom):
				return cc.composeVoiced()
			case n.flag.has(DecomposeVom):
				return cc.decomposeVoiced()
			default:
				return cc.composeVoiced() // fix for TEST_L7tADs2z.
			}

		case vcSemivoiced:
			switch {
			case n.flag.has(ComposeVom):
				return cc.composeSemivoiced()
			case n.flag.has(DecomposeVom):
				return cc.decomposeSemivoiced()
			default:
				return cc.composeSemivoiced() // fix for TEST_K6t8hQYp
			}

		default:
			// TEST_R8jrnbCz knows that the program never passes here
			panic("unreachable")
		}

	case ctKanaSymbol:
		switch {
		case n.flag.has(KanaSymbolToNarrow):
			return c.toNarrowR(), vmNone
		case n.flag.has(KanaSymbolToWide):
			return c.toWideR(), vmNone
		default:
			return c.codepoint, vmNone
		}

	case ctKanaVom:
		switch {
		case n.flag.has(IsolatedKanaVomToNarrow):
			return c.toNarrowR(), vmNone
		case n.flag.has(IsolatedKanaVomToWide):
			return c.toLegacyC().toWideR(), vmNone
		case n.flag.has(IsolatedKanaVomToNonspace):
			return c.toCombiningR(), vmNone
		default:
			return c.codepoint, vmNone
		}

	default:
		// TEST_P8w4qtsm knows that the program never passes here
		panic("unreachable")
	}
}

// Normalize normalizes the string according to the current
// normalization mode.
func (n *Normalizer) Normalize(s string) string {
	rs := []rune(s)
	var sb strings.Builder
	sb.Grow(len(rs) * 2)
	var r rune
	var m vom
	for i := 0; i < len(rs); i++ {
		if i < len(rs)-1 {
			var ok bool
			if r, m, ok = n.maybeComposeVom(rs[i], rs[i+1]); ok {
				i++
			}
		} else {
			r, m = n.NormalizeRune(rs[i])
		}
		sb.WriteRune(r)
		if m.isVom() {
			sb.WriteRune(rune(m))
		}
	}

	return sb.String()
}
