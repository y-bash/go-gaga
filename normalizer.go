package gaga

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
	// Example: [Ａ] =>[A]
	AlphaToNarrow

	// AlphaToWide converts all the half-width Latin letters to their full-width.
	// Example: [A] => [Ａ]
	AlphaToWide

	// AlphaToUpper converts all the lower case Latin letters to their upper case.
	// Examples: [a] => [A],  [ａ] => [Ａ]
	AlphaToUpper

	// AlphaToLower converts all the upper case Latin letters to their lower case.
	// Examples: [A] => [a],  [Ａ] => [ａ]
	AlphaToLower

	// DigitToNarrow converts all the full-width Latin digits to their half-width.
	// Example: [１] => [1]
	DigitToNarrow

	// DigitToWide converts all the half-width Latin digits to their full-width.
	// Example: [1] => [１]
	DigitToWide

	// SymbolToNarrow converts all the full-width Latin symbols to their half-width.
	// Example: [？] => [?]
	SymbolToNarrow

	// SymbolToWide converts all the half-width Latin symbols to their full-width.
	// Example: [?] => [？]
	SymbolToWide

	// HiraganaToNarrow converts the full-width Hiragana letters to
	// their half-width Katakana as much as possible.
	// Example: [あ] => [ｱ]
	HiraganaToNarrow

	// HiraganaToKatakana converts the full-width Hiragana letters to
	// their full-width Katakana as much as possible.
	// Example: [あ] => [ア]
	HiraganaToKatakana

	// KatakanaToNarrow converts the full-width Katakana letters to
	// their half-width Katakana as much as possible.
	// Example: [ア] => [ｱ]
	KatakanaToNarrow

	// KatakanaToWide converts all the half-width Katakana letters to
	// their full-width Katakana.
	// Example: [ｱ] => [ア]
	KatakanaToWide

	// KatakanaToHiragana converts the half-width or full-width Katakana
	// letters to their full-width Hiragana as much as possible.
	// Examples: [ア] => [あ],  [ｱ] => [あ]
	KatakanaToHiragana

	// KanaSymbolToNarrow converts the full-width Hiragana-Katakana symbols
	// to their half-width as much as possible.
	// Example: [、] => [､]
	KanaSymbolToNarrow

	// KanaSymbolToWide converts all the half-width Hiragana-Katakana symbols
	// to their full-width.
	// Example: [､] => [、]
	KanaSymbolToWide

	// ComposeVom combines voiced or semi-voiced sound letter
	// in the traditional JIS X 0201 or JIS X 0208 sytle.
	// marks behind
	// Hiragana-Katakana in a traditional style.
	// JIS X 0201 and JIS X 0208
	// TODO Voiced character, Semi-voiced character, CombinedModmark, Traditional combining
	// Examples:
	//	[か][゛] => [が],    [ｶ][゛]  => [ｶ][ﾞ],     [は][゜] => [ぱ]
	//	[ヰ][゛] => [ヸ],    [ゐ][゛] => [ゐ][゛]
	ComposeVom

	// DecomposeVom combines voiced or semi-voiced sound marks behind
	// Hiragana-Katakana in a Unicode combining style.
	// TODO Voiced character, Semi-voiced character, Unicode Canonical combining
	// Examples:
	//	[が] => [か][\u3099],  [か][゛] => [か][\u3099],  [ｶ][゛]  => [ｶ][\u3099],
	//	[ぱ] => [は][\u309A],  [ヰ][゛] => [ヰ][\u3099],  [ゐ][゛] => [ゐ][\u3099]
	DecomposeVom

	// TODO comment Isolated
	IsolatedKanaVomToNarrow

	// TODO comment
	IsolatedKanaVomToWide

	// TODO comment
	IsolatedKanaVomToNonspace

	normflagMax
)

// Combination of normalization flags
const (
	// LatinToNarrow is a combination of normalization flags for converting
	// all the full-width Latin characters to their half-width.
	//
	//          | CHARACTER     | CONVERT TO
	// ---------+---------------+----------------
	//          | Wide Alphabet | Narrow Alphabet
	// Category | Wide Digit    | Narrow Digit
	//          | Wide Symbol   | Narrow Symbol
	// ---------+---------------+----------------
	// Example  | "Ａ１？"      | "A1?"
	//
	LatinToNarrow = AlphaToNarrow | DigitToNarrow | SymbolToNarrow

	// LatinToWide is a combination of normalization flags for converting
	// all the half-width Latin characters to their full-width.
	//
	//          | CHARACTER       | CONVERT TO
	// ---------+-----------------+--------------
	//          | Narrow Alphabet | Wide Alphabet
	// Category | Narrow Digit    | Wide Digit
	//          | Narrow Symbol   | Wide Symbol
	// ---------+-----------------+--------------
	// Example  | "A1?"           | "Ａ１？"
	//
	LatinToWide = AlphaToWide | DigitToWide | SymbolToWide

	// KanaToNarrow is a combination of normalization flags for converting
	// the full-width Hiragana-Katakana characters to their half-width as
	// much as possible.
	//
	//          | CHARACTER                      | CONVERT TO
	// ---------+--------------------------------+----------------------
	//          | Hiaragana                      | Narrow Katakana
	// Category | Wide Katakana                  | Narrow Katakana
	//          | Wide Kana Symbol               | Narrow Kana Symbol
	//          | Voiced/Semi-voiced Kana Letter | Traditional voicing
	//          | Isolated Wide VSM/SVSM         | Narrow VSM/SVSM
	// ---------+--------------------------------+----------------------
	// Example  | "あイ、が゛"                   | "ｱｲ､ｶﾞﾞ"
	//
	KanaToNarrow = HiraganaToNarrow | KatakanaToNarrow | KanaSymbolToNarrow |
		IsolatedKanaVomToNarrow | ComposeVom

	// KanaToWide is a combination of normalization flags for converting
	// all the half-width Hiragana-Katakana characters to their full-width.
	//
	//          | CHARACTER                      | CONVERT TO
	// ---------+--------------------------------+----------------------
	//          | Narrow Katakana                | Wide Katakana
	// Category | Narrow Kana Symbol             | Wide Kana Symbol
	//          | Voiced/Semi-voiced Kana Letter | Traditional voicing
	//          | Isolated Narrow VSM/SVSM       | Wide VSM/SVSM
	// ---------+--------------------------------+----------------------
	// Example  | "ｱ､ｶﾞﾞ"                        | "ア、ガ゛"
	//
	KanaToWide = KatakanaToWide | KanaSymbolToWide | IsolatedKanaVomToWide |
		ComposeVom

	//
	//          | CHARACTER                      | CONVERT TO
	// ---------+--------------------------------+----------------------
	//          | Hiragana                       | Wide Katakana
	// Category | Narrow Katakana                | Wide Katakana
	//          | Narrow Kana Symbol             | Wide Kana Symbol
	//          | Voiced/Semi-voiced Kana Letter | Traditional voicing
	//          | Isolated Narrow VSM/SVSM       | Wide VSM/SVSM
	// ---------+--------------------------------+----------------------
	// Example  | "あｲ､ｶﾞﾞ"                       | "アイ、ガ゛"
	//
	KanaToWideKatakana = KatakanaToWide | HiraganaToKatakana | KanaSymbolToWide |
		IsolatedKanaVomToWide | ComposeVom

	//
	//          | CHARACTER                      | CONVERT TO
	// ---------+--------------------------------+----------------------
	//          | Hiragana                       | Narrow Katakana
	// Category | Wide Katakana                  | Narrow Katakana
	//          | Wide Kana Symbol               | Narrow Kana Symbol
	//          | Voiced/Semi-voiced Kana Letter | Traditional voicing
	//          | Isolated Wide VSM/SVSM         | Narrow VSM/SVSM
	// ---------+--------------------------------+----------------------
	// Example  | "あイ、が゛"                   | "ｱｲ､ｶﾞﾞ"
	//
	KanaToNarrowKatakana = KatakanaToNarrow | HiraganaToNarrow |
		KanaSymbolToNarrow | IsolatedKanaVomToNarrow | ComposeVom

	//
	//          | CHARACTER                      | CONVERT TO
	// ---------+--------------------------------+----------------------
	//          | Wide Katakana                  | Hiragana
	// Category | Narrow Katakana                | Hiragana
	//          | Narrow Kana Symbol             | Wide Kana Symbol
	//          | Voiced/Semi-voiced Kana Letter | Traditional voicing
	//          | Isolated Narrow VSM/SVSM       | Wide VSM/SVSM
	// ---------+--------------------------------+----------------------
	// Example  | "アｲ､ガ゛"                     | "あい、が゛"
	//
	KanaToHiragana = KatakanaToHiragana | KanaSymbolToWide |
		IsolatedKanaVomToWide | ComposeVom
)

func (f NormFlag) has(f2 NormFlag) bool { return f&f2 != 0 }

var normflagNames = map[NormFlag]string{
	AlphaToNarrow:             "AlphaToNarrow",
	AlphaToWide:               "AlphaToWide",
	AlphaToUpper:              "AlphaToUpper",
	AlphaToLower:              "AlphaToLower",
	DigitToNarrow:             "DigitToNarrow",
	DigitToWide:               "DigitToWide",
	SymbolToNarrow:            "SymbolToNarrow",
	SymbolToWide:              "SymbolToWide",
	HiraganaToNarrow:          "HiraganaToNarrow",
	HiraganaToKatakana:        "HiraganaToKatakana",
	KatakanaToNarrow:          "KatakanaToNarrow",
	KatakanaToWide:            "KatakanaToWide",
	KatakanaToHiragana:        "KatakanaToHiragana",
	KanaSymbolToNarrow:        "KanaSymbolToNarrow",
	KanaSymbolToWide:          "KanaSymbolToWide",
	ComposeVom:                "ComposeVom",
	DecomposeVom:              "DecomposeVom",
	IsolatedKanaVomToNarrow:   "IsolatedKanaVomToNarrow",
	IsolatedKanaVomToWide:     "IsolatedKanaVomToWide",
	IsolatedKanaVomToNonspace: "IsolatedKanaVomToNonspace",
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
	KanaSymbolToNarrow | KanaSymbolToWide,
	ComposeVom | DecomposeVom,
	IsolatedKanaVomToNarrow | IsolatedKanaVomToWide,
	IsolatedKanaVomToNarrow | IsolatedKanaVomToNonspace,
	IsolatedKanaVomToWide | IsolatedKanaVomToNonspace,
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

func (n *Normalizer) maybeComposeVom(r1, r2 rune) (rune, vom, bool) {
	if !vom(r2).isVom() {
		r, m := n.NormalizeRune(r1)
		return r, m, false
	}

	c1, ok := getUnichar(r1)
	if !ok {
		return r1, vmNone, false
	}

	if c1.category != ctKanaLetter || c1.voicing == vcVoiced || c1.voicing == vcSemivoiced {
		r, m := n.NormalizeRune(r1)
		return r, m, false
	}

	nr1, _ := n.NormalizeRune(r1) // // TEST_nD7FwQUW knows that NormalizeRune() definitely return the vmNone
	nc1, _ := getUnichar(nr1)     // TEST_G9amUMTr knows that getUnichar() definitely return the ok value
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

// NewNormalizer creates a new Normalizer with specified flag (LatinToNarrow etc.).
// If successful, methods on the returned Normalizer can be used for normalization.
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

// NormalizeRune normalizes the rune according to the current normalization
// mode. Depending on the mode, the voiced or semi-voiced sound mark may be
// separated, so it may return multiple runes. but, this function allways
// returns a rune array with 1 or 2 elements, and never returns an array with
// any other number of elements.
func (n *Normalizer) NormalizeRune(r rune) (rune, vom) {
	// TEST_Fc68JR9i knows about the number of elements in
	// the return value of this function
	c, ok := getUnichar(r)
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
