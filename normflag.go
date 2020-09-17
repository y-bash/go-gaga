package gaga

import (
	"fmt"
	"strings"
)

// NormFlag is the normalization rule used by Normalizer.
type NormFlag int

// Constants to identify various normalization flags.
const (
	// normflagUndefined indicates that the normalization flag is undefined.
	// sequence of 0, 1, 2, 4, 8, etc...
	normflagUndefined NormFlag = (1 << iota) / 2

	// AlphaToNarrow converts all the full-width Latin letters to
	// their half-width.
	// Example: [Ａ] =>[A]
	AlphaToNarrow

	// AlphaToWide converts all the half-width Latin letters to
	// their full-width.
	// Example: [A] => [Ａ]
	AlphaToWide

	// AlphaToUpper converts all the lower case Latin letters to
	// their upper case.
	// Examples: [a] => [A],  [ａ] => [Ａ]
	AlphaToUpper

	// AlphaToLower converts all the upper case Latin letters to
	// their lower case.
	// Examples: [A] => [a],  [Ａ] => [ａ]
	AlphaToLower

	// DigitToNarrow converts all the full-width Latin digits to
	// their half-width.
	// Example: [１] => [1]
	DigitToNarrow

	// DigitToWide converts all the half-width Latin digits to
	// their full-width.
	// Example: [1] => [１]
	DigitToWide

	// SymbolToNarrow converts all the full-width Latin symbols to
	// their half-width.
	// Example: [？] => [?]
	SymbolToNarrow

	// SymbolToWide converts all the half-width Latin symbols to
	// their full-width.
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

	// KanaSymbolToNarrow converts the full-width Hiragana-Katakana
	// symbols to their half-width as much as possible.
	// Example: [、] => [､]
	KanaSymbolToNarrow

	// KanaSymbolToWide converts all the half-width Katakana symbols
	// to their full-width.
	// Example: [､] => [、]
	KanaSymbolToWide

	// ComposeVom composes the voiced or semi-voiced sound letters in
	// the most conventional way.
	// Examples:
	//  [が]     => [が],  [か][゛] => [が],    [か][\u3099] => [が],
	//  [か][ﾞ]  => [が],  [ｶ][゛]  => [ｶ][ﾞ],  [ｶ][ﾞ]       => [ｶ][ﾞ],
	//  [は][゜] => [ぱ],  [ヰ][゛] => [ヸ],    [ゐ][゛]     => [ゐ][゛]
	ComposeVom

	// DecomposeVom decomposes the voiced or semi-voiced sound letters
	// in a way similar to the Unicode canonical decomposition mappings.
	// Examples:
	//  [が]         => [か][\u3099],  [か][゛] => [か][\u3099],
	//  [か][\u3099] => [か][\u3099],  [か][ﾞ]  => [か][\u3099],
	//  [ｶ][゛]      => [ｶ][\u3099],   [ｶ][ﾞ]   => [ｶ][\u3099],
	//  [ぱ]         => [は][\u309A],  [ヰ][゛] => [ヰ][\u3099],
	//  [ゐ][゛]     => [ゐ][\u3099]
	DecomposeVom

	// IsolatedVomToNarrow converts an isolated voicing modifier
	// which was not combined into a base letter into a half-width
	// voiced or semi-voiced sound letter.
	// Examples:
	//  [゛] => [ﾞ],  [\u3099] => [ﾞ],  [゜] => [ﾟ],  [\u309A] => [ﾟ]
	IsolatedVomToNarrow

	// IsolatedVomToWide converts an isolated voicing modifier
	// which was not combined into a base letter into a full-width
	// voiced or semi-voiced sound letter.
	// Examples:
	//  [\u3099] => [゛],  [ﾞ] => [゛],  [\u309A] => [゜],  [ﾟ] => [゜]
	IsolatedVomToWide

	// IsolatedVomToCombining converts an isolated voicing
	// modifier which was not combined into a base letter into a
	// combining voiced or semi-voiced sound letter.
	// Examples:
	//  [゛] => [\u3099],  [ﾞ] => [\u3099],  [゜] = [\u309A],  [ﾟ] => [\u309A]
	IsolatedVomToNonspace

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
	//          | CHARACTER                       | CONVERT TO
	// ---------+---------------------------------+-------------------
	//          | Hiaragana                       | Narrow Katakana
	// Category | Wide Katakana                   | Narrow Katakana
	//          | Wide Kana Symbol                | Narrow Kana Symbol
	//          | Voiced/Semi-voiced Kana Letter  | Legacy composed
	//          | Isolated Voicing Modifier (VOM) | Narrow VOM
	// ---------+---------------------------------+-------------------
	// Example  | "あイ、が゛"                    | "ｱｲ､ｶﾞﾞ"
	//
	KanaToNarrow = HiraganaToNarrow | KatakanaToNarrow | KanaSymbolToNarrow |
		IsolatedVomToNarrow | ComposeVom

	// KanaToWide is a combination of normalization flags for converting
	// all the half-width Katakana characters to their full-width.
	//
	//          | CHARACTER                       | CONVERT TO
	// ---------+---------------------------------+-----------------
	//          | Narrow Katakana                 | Wide Katakana
	// Category | Narrow Kana Symbol              | Wide Kana Symbol
	//          | Voiced/Semi-voiced Kana Letter  | Legacy composed
	//          | Isolated Voicing Modifier (VOM) | Wide VOM
	// ---------+---------------------------------+-----------------
	// Example  | "ｱ､ｶﾞﾞ"                         | "ア、ガ゛"
	//
	KanaToWide = KatakanaToWide | KanaSymbolToWide | IsolatedVomToWide |
		ComposeVom

	// KanaToWideKatakana is a combination of normalization flags for
	// converting all the half-width Katakana characters to their full-width,
	// and the Hiragana characters to their full-width Katakana as much as
	// possible..
	//
	//          | CHARACTER                       | CONVERT TO
	// ---------+---------------------------------+-----------------
	//          | Hiragana                        | Wide Katakana
	// Category | Narrow Katakana                 | Wide Katakana
	//          | Narrow Kana Symbol              | Wide Kana Symbol
	//          | Voiced/Semi-voiced Kana Letter  | Legacy composed
	//          | Isolated Voicing Modifier (VOM) | Wide VOM
	// ---------+---------------------------------+-----------------
	// Example  | "あｲ､ｶﾞﾞ"                       | "アイ、ガ゛"
	//
	KanaToWideKatakana = KatakanaToWide | HiraganaToKatakana | KanaSymbolToWide |
		IsolatedVomToWide | ComposeVom

	// KanaToNarrowKatakana is a combination of normalization flags for
	// converting the full-width Katakana characters to their half-width,
	// and the Hiragana characters to their half-width Katakana as much as
	// possible.
	//
	//          | CHARACTER                       | CONVERT TO
	// ---------+---------------------------------+-------------------
	//          | Hiragana                        | Narrow Katakana
	// Category | Wide Katakana                   | Narrow Katakana
	//          | Wide Kana Symbol                | Narrow Kana Symbol
	//          | Voiced/Semi-voiced Kana Letter  | Legacy composed
	//          | Isolated Voicing Modifier (VOM) | Narrow VOM
	// ---------+---------------------------------+-------------------
	// Example  | "あイ、が゛"                    | "ｱｲ､ｶﾞﾞ"
	//
	KanaToNarrowKatakana = KatakanaToNarrow | HiraganaToNarrow |
		KanaSymbolToNarrow | IsolatedVomToNarrow | ComposeVom

	// KanaToHiragana is a combination of normalization flags for
	// converting the full-width Katakana characters to their Hiragana
	// as much as possible, and all the half-width Katakana characters
	// to their Hiragana.
	//
	//          | CHARACTER                       | CONVERT TO
	// ---------+---------------------------------+----------------------
	//          | Wide Katakana                   | Hiragana
	// Category | Narrow Katakana                 | Hiragana
	//          | Narrow Kana Symbol              | Wide Kana Symbol
	//          | Voiced/Semi-voiced Kana Letter  | Legacy composed
	//          | Isolated Voicing Modifier (VOM) | Wide VOM
	// ---------+---------------------------------+----------------------
	// Example  | "アｲ､ガ゛"                      | "あい、が゛"
	//
	KanaToHiragana = KatakanaToHiragana | KanaSymbolToWide |
		IsolatedVomToWide | ComposeVom

	// Fold is a combination of normalization flags for converting
	// the Latin characters and the Hiragana-Katakana characters to
	// their canonical width.
	//
	//          | CHARACTER                       | CONVERT TO
	// ---------+---------------------------------+-----------------
	//          | Wide Alphabet                   | Narrow Alphabet
	//          | Wide Digit                      | Narrow Digit
	//          | Wide Symbol                     | Narrow Symbol
	// Category | Narrow Katakana                 | Wide Katakana
	//          | Narrow Kana Symbol              | Wide Kana Symbol
	//          | Voiced/Semi-voiced Kana Letter  | Legacy composed
	//          | Isolated Voicing Modifier (VOM) | Wide VOM
	// ---------+---------------------------------+-----------------
	// Example  | "Ａ１？ｱ､ｶﾞﾞ"                   | "A1?ア、ガ゛"
	//
	Fold = LatinToNarrow | KanaToWide
)

var normflagMap = map[NormFlag]string{
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
	IsolatedVomToNarrow:   "IsolatedVomToNarrow",
	IsolatedVomToWide:     "IsolatedVomToWide",
	IsolatedVomToNonspace: "IsolatedVomToNonspace",
}

var combflagList = []struct {
	flag NormFlag
	name string
}{
	{LatinToNarrow, "LatinToNarrow"},
	{LatinToWide, "LatinToWide"},
	{KanaToNarrow, "KanaToNarrow"},
	{KanaToWide, "KanaToWide"},
	{KanaToWideKatakana, "KanaToWideKatakana"},
	{KanaToNarrowKatakana, "KanaToNarrowKatakana"},
	{KanaToHiragana, "KanaToHiragana"},
	{Fold, "Fold"},
}

var normflagRevMap = func() map[string]NormFlag {
	l := len(normflagMap) + len(combflagList)
	revmap := make(map[string]NormFlag, l)
	for k, v := range normflagMap {
		revmap[v] = k
	}
	for _, combflag := range combflagList {
		revmap[combflag.name] = combflag.flag
	}
	return revmap

}()

// invalid combination of normalization flags.
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
	IsolatedVomToNarrow | IsolatedVomToWide,
	IsolatedVomToNarrow | IsolatedVomToNonspace,
	IsolatedVomToWide | IsolatedVomToNonspace,
}

func (f NormFlag) has(f2 NormFlag) bool { return f&f2 != 0 }

// String returns the name of a flag
func (f NormFlag) String() string {
	var ss []string
	for f2 := NormFlag(1); f2 < normflagMax; f2 <<= 1 {
		if f.has(f2) {
			ss = append(ss, normflagMap[f2])
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

//ParseNormflag returns a flags parsed names
func ParseNormFlag(names string) (flags NormFlag, err error) {
	ss := strings.Split(names, "|")
	for _, s := range ss {
		name := strings.Trim(s, " ()")
		if len(name) <= 0 {
			continue
		}
		flag, ok := normflagRevMap[name]
		if !ok {
			return flags, fmt.Errorf("invalid normalization flag: %s", name)
		}
		flags |= flag
	}
	if err = flags.validate(); err != nil {
		return flags, err
	}
	return flags, nil
}

func (f NormFlag) validate() error {
	if f <= normflagUndefined || f >= normflagMax {
		return fmt.Errorf("invalid normalization flag: %s", f)
	}
	for _, invalid := range invalidFlagsList {
		if f&invalid == invalid {
			return fmt.Errorf(
				"invalid normalization flag: %s, invalid combination: %s",
				f, invalid)
		}
	}
	return nil
}
