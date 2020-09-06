// Package gaga implements Japanese string conversion.
//
// Usage
//
// TODO update package comments
// Define a normalizer using NewNormalizer() with the normalization flag.
//
// This declares a normalizer, that converts Latin characters to
// half-width and Kana characters to full-width.
// 	import "github.com/y-bash/go-gaga"
// 	n:= gaga.NewNormalizer(gaga.LatinToNarrow | gaga.KanaToWide)
//
// After normalizer is defined, call
// 	s := n.Normalize("ＡＢＣｱｲｳ")
// to normalize the string using the normalization flags.
//
// Then the string is converted to
// 	fmt.Printf("%q", s) // Stdout: "ABCアイウ"
package gaga

// TODO support for Katakana_Extの対応

import (
	"fmt"
	//	"github.com/mattn/go-runewidth"
	"strings"
)

// NormFlag is the normalization rule used by Normalizer.
type NormFlag int

// Constants to identify various normalization flags.
const (
	// normflagUndefined indicates that the normalization flag is undefined.
	normflagUndefined = (1 << iota) / 2 // Sequence of 0, 1, 2, 4, 8, etc...

	// AlphaToNarrow converts all the full-width Latin letters to their half-width.
	AlphaToNarrow

	// AlphaToWide converts all the half-width Latin letters to their full-width.
	AlphaToWide

	// AlphaToUpper converts all the lower case Latin letters to their upper case.
	AlphaToUpper

	// AlphaToLower converts all the upper case Latin letters to their lower case.
	AlphaToLower

	// DigitToNarrow converts all the full-width Latin digits to their half-width.
	DigitToNarrow

	// DigitToWide converts all the half-width Latin digits to their full-width.
	DigitToWide

	// SymbolToNarrow converts all the full-width Latin symbols to their half-width.
	SymbolToNarrow

	// SymbolToWide converts all the half-width Latin symbols to their full-width.
	SymbolToWide

	// HiraganaToNarrow converts the full-width Hiragana letters to
	// their half-width Katakana as much as possible.
	HiraganaToNarrow

	// HiraganaToKatakana converts the full-width Hiragana letters to
	// their full-width Katakana as much as possible.
	HiraganaToKatakana

	// KatakanaToNarrow converts the full-width Katakana letters to
	// their half-width Katakana as much as possible.
	KatakanaToNarrow

	// KatakanaToWide converts all the half-width Katakana letters to
	// their full-width Katakana.
	KatakanaToWide

	// KatakanaToHiragana converts the half-width or full-width Katakana
	// letters to their full-width Hiragana as much as possible.
	KatakanaToHiragana

	// KanaSymToNarrow converts the full-width Hiragana-Katakana symbols
	// to their half-width as much as possible.
	KanaSymToNarrow

	// KanaSymToWide converts all the half-width Hiragana-Katakana symbols
	// to their full-width.
	KanaSymToWide

	// KanaVsmToClassical converts the voiced or semi-voiced sound marks of
	// Hiragana-Katakana to the classical style as much as possible.
	KanaVsmToClassical

	// KanaVsmToCombining converts the voiced or semi-voiced sound marks of
	// Hiragana-Katakana to the combining style as much as possible.
	KanaVsmToCombining

	normflagMax
)

// Combination of normalization flags
const (
	// LatinToNarrow is a combination of normalization flags for converting
	// all the full-width Latin characters to their half-width.
	LatinToNarrow = AlphaToNarrow | DigitToNarrow | SymbolToNarrow

	// LatinToWide is a combination of normalization flags for converting
	// all the half-width Latin characters to their full-width.
	LatinToWide = AlphaToWide | DigitToWide | SymbolToWide

	// KanaToNarrow is a combination of normalization flags for converting
	// the full-width Hiragana-Katakana characters to their half-width as
	// much as possible.
	KanaToNarrow = HiraganaToNarrow | KatakanaToNarrow | KanaSymToNarrow | KanaVsmToClassical

	// KanaToWide is a combination of normalization flags for converting
	// all the half-width Hiragana-Katakana characters to their full-width.
	KanaToWide = KatakanaToWide | KanaSymToWide | KanaVsmToClassical
)

func (f NormFlag) has(f2 NormFlag) bool { return f&f2 != 0 }

var normflagNames = map[NormFlag]string{
	AlphaToNarrow:      "AlphaToNarrow",
	AlphaToWide:        "AlphaToWide",
	AlphaToUpper:       "AlphaToUpper",
	AlphaToLower:       "AlphaToLower",
	DigitToNarrow:      "DigitToNarrow",
	DigitToWide:        "DigitToWide",
	SymbolToNarrow:     "SymbolToNarrow",
	SymbolToWide:       "SymbolToWide",
	HiraganaToNarrow:   "HiraganaToNarrow",
	HiraganaToKatakana: "HiraganaToKatakana",
	KatakanaToNarrow:   "KatakanaToNarrow",
	KatakanaToWide:     "KatakanaToWide",
	KatakanaToHiragana: "KatakanaToHiragana",
	KanaSymToNarrow:    "KanaSymToNarrow",
	KanaSymToWide:      "KanaSymToWide",
	KanaVsmToClassical: "KanaVsmToClassical",
	KanaVsmToCombining: "KanaVsmToCombining",
}

func (f NormFlag) String() string {
	var ss []string
	for f2 := NormFlag(1); f2 < normflagMax; f2 <<= 1 {
		if f.has(f2) {
			ss = append(ss, normflagNames[f2])
		}
	}
	switch len(ss) {
	case 0:
		return "<undefined>"
	case 1:
		return ss[0]
	default:
		return "(" + strings.Join(ss, " | ") + ")"
	}
}

// Normalizer normalizes the input provided and returns the normalized string.
type Normalizer struct {
	flag NormFlag
}

// invalid combination of normalization flags
var invalidFlagsList = []NormFlag{
	AlphaToUpper | AlphaToLower,
	AlphaToNarrow | AlphaToWide,
	DigitToNarrow | DigitToWide,
	SymbolToNarrow | SymbolToWide,
	KatakanaToWide | KatakanaToNarrow,
	KatakanaToWide | KatakanaToHiragana,
	KatakanaToNarrow | KatakanaToHiragana,
	HiraganaToNarrow | HiraganaToKatakana,
	KanaSymToNarrow | KanaSymToWide,
	KanaVsmToClassical | KanaVsmToCombining,
}

func validateNormFlag(flag NormFlag) error {
	if flag <= normflagUndefined || flag >= normflagMax {
		return fmt.Errorf("invalid normalization flag value: %d", flag)
	}
	for _, invalid := range invalidFlagsList {
		if flag&invalid == invalid {
			return fmt.Errorf("invalid normalization flag: %d, invalid combination: %d", flag, invalid)
		}
	}
	return nil
}

func (n *Normalizer) maybeJoinClassicalVsm(r1, r2 rune) ([]rune, bool) {
	isVsm := isVoicedSoundMark(r2)
	isSvsm := isSemivoicedSoundMark(r2)
	if !isVsm && !isSvsm {
		return n.NormalizeRune(r1), false
	}

	c, ok := getUnichar(r1)
	if !ok {
		return []rune{r1}, false
	}

	if c.category != ctLatinLetter && c.category != ctLatinDigit && c.category != ctKanaLetter {
		return n.NormalizeRune(r1), false
	}

	if c.voicing == vcVoiced || c.voicing == vcSemivoiced {
		return n.NormalizeRune(r1), false
	}

	rs := n.NormalizeRune(r1)
	if len(rs) != 1 {
		neverBeCalled()
	}
	cc, ok := getUnichar(rs[0])

	if isVsm {
		rs = cc.toClassicalVoiced()
	} else {
		rs = cc.toClassicalSemivoiced()
	}
	return rs, true
}

func (n *Normalizer) normalizeDontcareVsm(s string) string {
	if len(s) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.Grow(len(s) * 2)
	rs := []rune(s)
	for _, r := range rs {
		nrs := n.NormalizeRune(r)
		for _, nr := range nrs {
			sb.WriteRune(nr)
		}
	}
	return sb.String()
}

func (n *Normalizer) normalizeCombiningVsm(s string) string {
	if len(s) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.Grow(len(s) * 2)
	rs := []rune(s)
	for _, r := range rs {
		var nrs []rune
		switch {
		case isVoicedSoundMark(r):
			nrs = []rune{vsmCombining}
		case isSemivoicedSoundMark(r):
			nrs = []rune{svsmCombining}
		default:
			nrs = n.NormalizeRune(r)
		}
		for _, nr := range nrs {
			sb.WriteRune(nr)
		}
	}
	return sb.String()
}

func (n *Normalizer) normalizeClassicalVsm(s string) string {
	rs := []rune(s)
	var nrs []rune
	var sb strings.Builder
	sb.Grow(len(rs) * 2)
	for i := 0; i < len(rs); i++ {
		if i < len(rs)-1 {
			var ok bool
			if nrs, ok = n.maybeJoinClassicalVsm(rs[i], rs[i+1]); ok {
				i++
			}
		} else {
			nrs = n.NormalizeRune(rs[i])
		}
		for _, nr := range nrs {
			sb.WriteRune(nr)
		}
	}

	return sb.String()
}

// NewNormalizer creates a new Normalizer with specified flag
// (LatinToNarrow etc.). If successful, methods on the returned
// Normalizer can be used for normalization.
func NewNormalizer(flag NormFlag) (*Normalizer, error) {
	err := validateNormFlag(flag)
	if err != nil {
		return nil, err
	}
	n := Normalizer{flag}
	return &n, nil
}

// SetFlag changes the normalization mode with the newly specified flag.
func (n *Normalizer) SetFlag(flag NormFlag) error {
	err := validateNormFlag(flag)
	if err != nil {
		return err
	}
	n.flag = flag
	return nil
}

// NormalizeRune normalizes the rune according to the current
// normalization mode. Depending on the mode, the voiced or
// Semi-voiced sound mark may be separated, so it may return
// multiple runes.
func (n *Normalizer) NormalizeRune(r rune) []rune {
	c, ok := getUnichar(r)
	if !ok {
		return []rune{r}
	}

	switch c.category {
	case ctUndefined:
		return []rune{c.codepoint}

	case ctLatinLetter:
		switch {
		case n.flag.has(AlphaToNarrow):
			c = c.toNarrowUnichar()
		case n.flag.has(AlphaToWide):
			c = c.toWideUnichar()
		}

		switch {
		case n.flag.has(AlphaToUpper):
			return []rune{c.toUpper()}
		case n.flag.has(AlphaToLower):
			return []rune{c.toLower()}
		default:
			return []rune{c.codepoint}
		}

	case ctLatinDigit:
		switch {
		case n.flag.has(DigitToNarrow):
			return []rune{c.toNarrow()}
		case n.flag.has(DigitToWide):
			return []rune{c.toWide()}
		default:
			return []rune{c.codepoint}
		}

	case ctLatinSymbol:
		switch {
		case n.flag.has(SymbolToNarrow):
			return []rune{c.toNarrow()}
		case n.flag.has(SymbolToWide):
			return []rune{c.toWide()}
		default:
			return []rune{c.codepoint}
		}

	case ctKanaLetter:
		var cc *unichar
		switch c.charCase {
		case ccHiragana:
			switch {
			case n.flag.has(HiraganaToNarrow):
				cc = c.toNarrowUnichar()
			case n.flag.has(HiraganaToKatakana):
				cc = c.toKatakanaUnichar()
			default:
				cc = c
			}
		case ccKatakana:
			switch {
			case n.flag.has(KatakanaToNarrow):
				cc = c.toNarrowUnichar()
			case n.flag.has(KatakanaToHiragana):
				cc = c.toHiraganaUnichar()
			case n.flag.has(KatakanaToWide):
				cc = c.toWideUnichar()
			default:
				cc = c
			}
		default:
			// TEST_gT8YJdBc knows that the program never passes here
			return neverBeCalled()
		}

		switch c.voicing {
		case vcUndefined, vcUnvoiced:
			return []rune{cc.codepoint}

		case vcVoiced:
			switch {
			case n.flag.has(KanaVsmToClassical):
				return cc.toClassicalVoiced()
			case n.flag.has(KanaVsmToCombining):
				return cc.toCombiningVoiced()
			default:
				return cc.toClassicalVoiced() // fix for TEST_L7tADs2z.
			}

		case vcSemivoiced:
			switch {
			case n.flag.has(KanaVsmToClassical):
				return cc.toClassicalSemivoiced()
			case n.flag.has(KanaVsmToCombining):
				return cc.toCombiningSemivoiced()
			default:
				return cc.toClassicalSemivoiced() // fix for TEST_K6t8hQYp
			}

		default:
			// TEST_R8jrnbCz knows that the program never passes here
			return neverBeCalled()
		}

	case ctKanaSymbol:
		switch {
		case n.flag.has(KanaSymToNarrow):
			return []rune{c.toNarrow()}
		case n.flag.has(KanaSymToWide):
			// fix for TEST_A9fCxUi6
			// In this case, it is hard-coded, because it cannot be represented by unichar_table.
			if c.codepoint == vsmCombining {
				return []rune{vsmWide}
			} else if c.codepoint == svsmCombining {
				return []rune{svsmWide}
			} else {
				return []rune{c.toWide()}
			}
		default:
			return []rune{c.codepoint}
		}

	default:
		// TEST_P8w4qtsm knows that the program never passes here
		return neverBeCalled()
	}
}

// Normalize normalizes the string according to the current
// normalization mode.
func (n *Normalizer) Normalize(s string) string {
	switch {
	case n.flag.has(KanaVsmToClassical):
		return n.normalizeClassicalVsm(s)
	case n.flag.has(KanaVsmToCombining):
		return n.normalizeCombiningVsm(s)
	default:
		return n.normalizeDontcareVsm(s)
	}
}
