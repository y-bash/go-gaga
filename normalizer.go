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
	// Example: 'Ａ' => 'A'
	AlphaToNarrow

	// AlphaToWide converts all the half-width Latin letters to their full-width.
	// Example: 'A' => 'Ａ'
	AlphaToWide

	// AlphaToUpper converts all the lower case Latin letters to their upper case.
	// Examples: 'a' => 'A',  'ａ' => 'Ａ'
	AlphaToUpper

	// AlphaToLower converts all the upper case Latin letters to their lower case.
	// Examples: 'A' => 'a',  'Ａ' => 'ａ'
	AlphaToLower

	// DigitToNarrow converts all the full-width Latin digits to their half-width.
	// Example: '１' => '1'
	DigitToNarrow

	// DigitToWide converts all the half-width Latin digits to their full-width.
	// Example: '1' => '１'
	DigitToWide

	// SymbolToNarrow converts all the full-width Latin symbols to their half-width.
	// Example: '？' => '?'
	SymbolToNarrow

	// SymbolToWide converts all the half-width Latin symbols to their full-width.
	// Example: '?' => '？'
	SymbolToWide

	// HiraganaToNarrow converts the full-width Hiragana letters to
	// their half-width Katakana as much as possible.
	// Example: 'あ' => 'ｱ'
	HiraganaToNarrow

	// HiraganaToKatakana converts the full-width Hiragana letters to
	// their full-width Katakana as much as possible.
	// Example: 'あ' => 'ア'
	HiraganaToKatakana

	// KatakanaToNarrow converts the full-width Katakana letters to
	// their half-width Katakana as much as possible.
	// Example: 'ア' => 'ｱ'
	KatakanaToNarrow

	// KatakanaToWide converts all the half-width Katakana letters to
	// their full-width Katakana.
	// Example: 'ｱ' => 'ア'
	KatakanaToWide

	// KatakanaToHiragana converts the half-width or full-width Katakana
	// letters to their full-width Hiragana as much as possible.
	// Examples: 'ア' => 'あ',  'ｱ' => 'あ'
	KatakanaToHiragana

	// KanaSymToNarrow converts the full-width Hiragana-Katakana symbols
	// to their half-width as much as possible.
	// Example: '、' => '､'
	KanaSymToNarrow

	// KanaSymToWide converts all the half-width Hiragana-Katakana symbols
	// to their full-width.
	// Example: '､' => '、'
	KanaSymToWide

	// VoicedKanaToTraditional combines voiced or semi-voiced sound marks behind
	// Hiragana-Katakana in a traditional style.
	// TODO Voiced character, Semi-voiced character
	// Examples:
	//	"か゛" => "が",  "ｶ゛"  => "ｶﾞ",   "は゜" => "ぱ"
	//	"ヰ゛" => "ヸ",  "ゐ゛" => "ゐ゛"
	VoicedKanaToTraditional

	// VoicedKanaToCombining combines voiced or semi-voiced sound marks behind
	// Hiragana-Katakana in a Unicode combining style.
	// TODO Voiced character, Semi-voiced character
	// Examples:
	//	"が" => "か\u3099",  "か゛" => "か\u3099",  "ｶ゛"  => "ｶ\u3099",
	//	"ぱ" => "は\u309A",  "ヰ゛" => "ヰ\u3099",  "ゐ゛" => "ゐ\u3099"
	VoicedKanaToCombining

	// TODO comment Isolated
	IsolatedVsmToNarrow

	// TODO comment
	IsolatedVsmToWide

	// TODO comment
	IsolatedVsmToCombining

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
	//          | Voiced/Semi-voiced Kana Letter | Traditional combining
	//          | Isolated Wide VSM/SVSM         | Narrow VSM/SVSM
	// ---------+--------------------------------+----------------------
	// Example  | "あイ、が゛"                   | "ｱｲ､ｶﾞﾞ"
	//
	KanaToNarrow = HiraganaToNarrow | KatakanaToNarrow | KanaSymToNarrow |
		IsolatedVsmToNarrow | VoicedKanaToTraditional

	// KanaToWide is a combination of normalization flags for converting
	// all the half-width Hiragana-Katakana characters to their full-width.
	//
	//          | CHARACTER                      | CONVERT TO
	// ---------+--------------------------------+----------------------
	//          | Narrow Katakana                | Wide Katakana
	// Category | Narrow Kana Symbol             | Wide Kana Symbol
	//          | Voiced/Semi-voiced Kana Letter | Traditional combining
	//          | Isolated Narrow VSM/SVSM       | Wide VSM/SVSM
	// ---------+--------------------------------+----------------------
	// Example  | "ｱ､ｶﾞﾞ"                        | "ア、ガ゛"
	//
	KanaToWide = KatakanaToWide | KanaSymToWide | IsolatedVsmToWide |
		VoicedKanaToTraditional

	//
	//          | CHARACTER                      | CONVERT TO
	// ---------+--------------------------------+----------------------
	//          | Hiragana                       | Wide Katakana
	// Category | Narrow Katakana                | Wide Katakana
	//          | Narrow Kana Symbol             | Wide Kana Symbol
	//          | Voiced/Semi-voiced Kana Letter | Traditional combining
	//          | Isolated Narrow VSM/SVSM       | Wide VSM/SVSM
	// ---------+--------------------------------+----------------------
	// Example  | "あｲ､ｶﾞﾞ"                       | "アイ、ガ゛"
	//
	KanaToWideKatakana = KatakanaToWide | HiraganaToKatakana | KanaSymToWide |
		IsolatedVsmToWide | VoicedKanaToTraditional

	//
	//          | CHARACTER                      | CONVERT TO
	// ---------+--------------------------------+----------------------
	//          | Hiragana                       | Narrow Katakana
	// Category | Wide Katakana                  | Narrow Katakana
	//          | Wide Kana Symbol               | Narrow Kana Symbol
	//          | Voiced/Semi-voiced Kana Letter | Traditional combining
	//          | Isolated Wide VSM/SVSM         | Narrow VSM/SVSM
	// ---------+--------------------------------+----------------------
	// Example  | "あイ、が゛"                   | "ｱｲ､ｶﾞﾞ"
	//
	KanaToNarrowKatakana = KatakanaToNarrow | HiraganaToNarrow |
		KanaSymToNarrow | IsolatedVsmToNarrow | VoicedKanaToTraditional

	//
	//          | CHARACTER                      | CONVERT TO
	// ---------+--------------------------------+----------------------
	//          | Wide Katakana                  | Hiragana
	// Category | Narrow Katakana                | Hiragana
	//          | Narrow Kana Symbol             | Wide Kana Symbol
	//          | Voiced/Semi-voiced Kana Letter | Traditional combining
	//          | Isolated Narrow VSM/SVSM       | Wide VSM/SVSM
	// ---------+--------------------------------+----------------------
	// Example  | "アｲ､ガ゛"                     | "あい、が゛"
	//
	KanaToHiragana = KatakanaToHiragana | KanaSymToWide |
		IsolatedVsmToWide | VoicedKanaToTraditional
)

func (f NormFlag) has(f2 NormFlag) bool { return f&f2 != 0 }

var normflagNames = map[NormFlag]string{
	AlphaToNarrow:           "AlphaToNarrow",
	AlphaToWide:             "AlphaToWide",
	AlphaToUpper:            "AlphaToUpper",
	AlphaToLower:            "AlphaToLower",
	DigitToNarrow:           "DigitToNarrow",
	DigitToWide:             "DigitToWide",
	SymbolToNarrow:          "SymbolToNarrow",
	SymbolToWide:            "SymbolToWide",
	HiraganaToNarrow:        "HiraganaToNarrow",
	HiraganaToKatakana:      "HiraganaToKatakana",
	KatakanaToNarrow:        "KatakanaToNarrow",
	KatakanaToWide:          "KatakanaToWide",
	KatakanaToHiragana:      "KatakanaToHiragana",
	KanaSymToNarrow:         "KanaSymToNarrow",
	KanaSymToWide:           "KanaSymToWide",
	VoicedKanaToTraditional: "VoicedKanaToTraditional",
	VoicedKanaToCombining:   "VoicedKanaToCombining",
	IsolatedVsmToNarrow:     "IsolatedVsmToNarrow",
	IsolatedVsmToWide:       "IsolatedVsmToWide",
	IsolatedVsmToCombining:  "IsolatedVsmToCombining",
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
	VoicedKanaToTraditional | VoicedKanaToCombining,
	IsolatedVsmToNarrow | IsolatedVsmToWide,
	IsolatedVsmToNarrow | IsolatedVsmToCombining,
	IsolatedVsmToWide | IsolatedVsmToCombining,
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

func (n *Normalizer) maybeCombineVsm(r1, r2 rune) ([]rune, bool) {
	isVsm := isVoicedSoundMark(r2)
	isSvsm := isSemivoicedSoundMark(r2)
	if !isVsm && !isSvsm {
		return n.NormalizeRune(r1), false
	}

	c, ok := getUnichar(r1)
	if !ok {
		return []rune{r1}, false
	}

	if c.category != ctKanaLetter {
		return n.NormalizeRune(r1), false
	}

	if c.voicing == vcVoiced || c.voicing == vcSemivoiced {
		return n.NormalizeRune(r1), false
	}

	nr := n.NormalizeRune(r1)[0]
	cc, ok := getUnichar(nr)
	switch {
	case isVsm:
		switch {
		case n.flag.has(VoicedKanaToTraditional):
			return cc.toTraditionalVoiced(), true
		case n.flag.has(VoicedKanaToCombining):
			return cc.toCombiningVoiced(), true
		default:
			vsm := n.NormalizeRune(r2)[0]
			return []rune{cc.codepoint, vsm}, true
		}
	case isSvsm:
		switch {
		case n.flag.has(VoicedKanaToTraditional):
			return cc.toTraditionalSemivoiced(), true
		case n.flag.has(VoicedKanaToCombining):
			return cc.toCombiningSemivoiced(), true
		default:
			svsm := n.NormalizeRune(r2)[0]
			return []rune{cc.codepoint, svsm}, true
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
func (n *Normalizer) NormalizeRune(r rune) []rune {
	// TEST_Fc68JR9i knows about the number of elements in
	// the return value of this function
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
			panic("unreachable")
		}

		switch c.voicing {
		case vcUndefined, vcUnvoiced:
			return []rune{cc.codepoint}

		case vcVoiced:
			switch {
			case n.flag.has(VoicedKanaToTraditional):
				return cc.toTraditionalVoiced()
			case n.flag.has(VoicedKanaToCombining):
				return cc.toCombiningVoiced()
			default:
				return cc.toTraditionalVoiced() // fix for TEST_L7tADs2z.
			}

		case vcSemivoiced:
			switch {
			case n.flag.has(VoicedKanaToTraditional):
				return cc.toTraditionalSemivoiced()
			case n.flag.has(VoicedKanaToCombining):
				return cc.toCombiningSemivoiced()
			default:
				return cc.toTraditionalSemivoiced() // fix for TEST_K6t8hQYp
			}

		default:
			// TEST_R8jrnbCz knows that the program never passes here
			panic("unreachable")
		}

	case ctKanaSymbol:
		switch {
		case n.flag.has(KanaSymToNarrow):
			return []rune{c.toNarrow()}
		case n.flag.has(KanaSymToWide):
			return []rune{c.toWide()}
		default:
			return []rune{c.codepoint}
		}

	case ctKanaVsm:
		switch {
		case n.flag.has(IsolatedVsmToNarrow):
			return []rune{c.toNarrow()}
		case n.flag.has(IsolatedVsmToWide):
			return []rune{c.toTraditionalMarkUnichar().toWide()}
		case n.flag.has(IsolatedVsmToCombining):
			return []rune{c.toCombiningMark()}
		default:
			return []rune{c.codepoint}
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
	var nrs []rune
	var sb strings.Builder
	sb.Grow(len(rs) * 2)
	for i := 0; i < len(rs); i++ {
		if i < len(rs)-1 {
			var ok bool
			if nrs, ok = n.maybeCombineVsm(rs[i], rs[i+1]); ok {
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
