package lib

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"strconv"
	"strings"
)

// Handling range of code points
type blockRange struct {
	first rune
	last  rune
	name  string
}

// Latin (excepted control characters)
var basicLatinBlock = blockRange{0x0020, 0x007F, "latin"}

// CJK symbols 3000-303F, Hiragana 3040-309F, Katakana 30A0-30FF
var jpKanaBlock = blockRange{0x3000, 0x30FF, "kana"}

// Katakana Phonetic Extensions 31F0-31FF
var jpKanaExtBlock = blockRange{0x31F0, 0x31FF, "kanaExt"}

// Full width latin letter, Half width kana (excepted FFA0-FFEF)
var widthFormBlock = blockRange{0xFF00, 0xFF9F, "width"}

var blockRanges = []blockRange{
	basicLatinBlock,
	jpKanaBlock,
	jpKanaExtBlock,
	widthFormBlock,
}

// Constant values structure
type constValues struct {
	title string
	list  []string
}

func (cv *constValues) name(u uint8) string {
	i := int(u)
	if i < 0 || i >= len(cv.list) {
		return fmt.Sprintf("Unknown %s: %d", cv.title, i)
	}
	return cv.list[i]
}

// Category values
const (
	ctUndefined = iota
	ctLatinLetter
	ctLatinDigit
	ctLatinSymbol
	ctKanaLetter
	ctKanaSymbol
	ctKanaVsm
	ctMax
)

var categories = constValues{
	title: "Category values",
	list: []string{
		"ctUndefined",
		"ctLatinLetter",
		"ctLatinDigit",
		"ctLatinSymbol",
		"ctKanaLetter",
		"ctKanaSymbol",
		"ctKanaVsm",
		"ctMax",
	},
}

// Character case values
const (
	ccUndefined = iota
	ccUpper
	ccLower
	ccHiragana
	ccKatakana
	ccTraditional
	ccNonspace
	ccMax
)

var charCases = constValues{
	title: "Character case values",
	list: []string{
		"ccUndefined",
		"ccUpper",
		"ccLower",
		"ccHiragana",
		"ccKatakana",
		"ccTraditional",
		"ccNonspace",
		"ccMax",
	},
}

// Character width values
const (
	cwUndefined = iota
	cwNarrow
	cwWide
	cwMax
)

var charWidths = constValues{
	title: "Character width values",
	list: []string{
		"cwUndefined",
		"cwNarrow",
		"cwWide",
		"cwMax",
	},
}

// Voicing values
const (
	vcUndefined = iota
	vcUnvoiced
	vcVoiced
	vcSemivoiced
	vcMax
)

var voicings = constValues{
	title: "Voicing values",
	list: []string{
		"vcUndefined",
		"vcUnvoiced",
		"vcVoiced",
		"vcSemivoiced",
		"vcMax",
	},
}

type charex struct {
	codepoint   rune   // Unicode code point value
	blk         string // Block in UCD
	na          string // Name  in UCD
	age         string // Age   in UCD
	gc          string // General category in UCD
	category    uint8  // ctUndefined/ctLatinLetter/ctLatinDigit/ctLatinSymbol/ctKanaLetter/ctKanaSymbol
	charCase    uint8  // ccUndefined/ccUpper/ccLower/ccHiragana/ccKatakana
	charWidth   uint8  // cwUndefined/cwNarrow/cwWide
	voicing     uint8  // vcUndefined/vcUnvoiced/vcVoiced/vcSemivoiced
	compatCase  rune   // Charcase compatible character
	compatWidth rune   // Width compatible character
	compatVs    rune   // Voiced sound compatible character
	compatSvs   rune   // Semi-voiced sound compatible character
}

type ucdex map[rune]*charex

type charmap struct {
	lo   rune
	hi   rune
	diff rune
}

var hiragana2katakana = []charmap{
	{'ぁ', 'ゖ', 'ア' - 'あ'}, // 3041[ぁ]-3096[ゖ] -> 30A1[ァ]-30F6 [ヶ]
	{'ゝ', 'ゞ', 'ア' - 'あ'}, // 309D[ゝ]-309E[ゞ] -> 30FD[ヽ]-30FE [ヾ]
}

// Additional ref directive
const (
	arCompatCase = iota
	arCompatWidth
	arCompatVs
	arCompatSvs
)

type additionalRef struct {
	codepoint rune
	ref       int
	data      rune
}

var additionalRefList = []additionalRef{
	{'ヮ', arCompatWidth, 'ﾜ'},      // 30EE [ヮ] -> FF9C [ﾜ]
	{'ヰ', arCompatWidth, 'ｲ'},      // 30F0 [ヰ] -> FF72 [ｲ]
	{'ヱ', arCompatWidth, 'ｴ'},      // 30F1 [ヱ] -> FF74 [ｴ]
	{'ヵ', arCompatWidth, 'ｶ'},      // 30F5 [ヵ] -> FF76 [ｶ]
	{'ヶ', arCompatWidth, 'ｹ'},      // 30F6 [ヶ] -> FF79 [ｹ]
	{'ゎ', arCompatWidth, 'ﾜ'},      // 308E [ゎ] -> FF9C [ﾜ]
	{'ゐ', arCompatWidth, 'ｲ'},      // 3090 [ゐ] -> FF72 [ｲ]
	{'ゑ', arCompatWidth, 'ｴ'},      // 3091 [ゑ] -> FF74 [ｴ]
	{'ゕ', arCompatWidth, 'ｶ'},      // 3095 [ゕ] -> FF76 [ｶ]
	{'ゖ', arCompatWidth, 'ｹ'},      // 3096 [ゖ] -> FF79 [ｹ]
	{'\u3099', arCompatWidth, 'ﾞ'}, // 3099 [゙ ] -> FF9E [ﾞ] // fix for TEST_N9x6dneg
	{'\u309A', arCompatWidth, 'ﾟ'}, // 309A [゚ ] -> FF9F [ﾟ]
	{'゛', arCompatWidth, 'ﾞ'},      // 309B [゙゛] -> FF9E [ﾞ] // fix for TEST_N9x6dneg
	{'゜', arCompatWidth, 'ﾟ'},      // 309C [゚゜] -> FF9F [ﾟ]
	{'ﾞ', arCompatWidth, '゛'},      // FF9E [ﾞ] -> 309B [゛]
	{'ﾟ', arCompatWidth, '゜'},      // FF9F [ﾟ] -> 309C [゜]
	{'\u3099', arCompatCase, '゛'},  // 3099 [ ] -> 309B [゛]
	{'\u309A', arCompatCase, '゜'},  // 309A [ ] -> 309C [゜]
	{'゛', arCompatCase, '\u3099'},  // 309B [゛] -> 3099 [ ]
	{'゜', arCompatCase, '\u309A'},  // 309C [゜ﾟ] -> 309A [ ]
	{'ﾞ', arCompatCase, '\u3099'},  // FF9E [ﾞ] -> 3099 [ ]
	{'ﾟ', arCompatCase, '\u309A'},  // FF9F [ﾟ] -> 309A [ ]
	{'ㇰ', arCompatCase, 'く'},       // 31F0 [ㇰ] -> 304F [く]
	{'ㇱ', arCompatCase, 'し'},       // 31F1 [ㇱ] -> 3057 [し]
	{'ㇲ', arCompatCase, 'す'},       // 31F2 [ㇲ] -> 3059 [す]
	{'ㇳ', arCompatCase, 'と'},       // 31F3 [ㇳ] -> 3068 [と]
	{'ㇴ', arCompatCase, 'ぬ'},       // 31F4 [ㇴ] -> 306C [ぬ]
	{'ㇵ', arCompatCase, 'は'},       // 31F5 [ㇵ] -> 306F [は]
	{'ㇶ', arCompatCase, 'ひ'},       // 31F6 [ㇶ] -> 3072 [ひ]
	{'ㇷ', arCompatCase, 'ふ'},       // 31F7 [ㇷ] -> 3075 [ふ]
	{'ㇸ', arCompatCase, 'へ'},       // 31F8 [ㇸ] -> 3078 [へ]
	{'ㇹ', arCompatCase, 'ほ'},       // 31F9 [ㇹ] -> 307B [ほ]
	{'ㇺ', arCompatCase, 'む'},       // 31FA [ㇺ] -> 3080 [む]
	{'ㇻ', arCompatCase, 'ら'},       // 31FB [ㇻ] -> 3089 [ら]
	{'ㇼ', arCompatCase, 'り'},       // 31FC [ㇼ] -> 308A [り]
	{'ㇽ', arCompatCase, 'る'},       // 31FD [ㇽ] -> 308B [る]
	{'ㇾ', arCompatCase, 'れ'},       // 31FE [ㇾ] -> 308C [れ]
	{'ㇿ', arCompatCase, 'ろ'},       // 31FF [ㇿ] -> 308D [ろ]
	{'ㇰ', arCompatWidth, 'ｸ'},      // 31F0 [ㇰ] -> FF78 [ｸ]
	{'ㇱ', arCompatWidth, 'ｼ'},      // 31F1 [ㇱ] -> FF7C [ｼ]
	{'ㇲ', arCompatWidth, 'ｽ'},      // 31F2 [ㇲ] -> FF7D [ｽ]
	{'ㇳ', arCompatWidth, 'ﾄ'},      // 31F3 [ㇳ] -> FF84 [ﾄ]
	{'ㇴ', arCompatWidth, 'ﾇ'},      // 31F4 [ㇴ] -> FF87 [ﾇ]
	{'ㇵ', arCompatWidth, 'ﾊ'},      // 31F5 [ㇵ] -> FF8A [ﾊ]
	{'ㇶ', arCompatWidth, 'ﾋ'},      // 31F6 [ㇶ] -> FF8B [ﾋ]
	{'ㇷ', arCompatWidth, 'ﾌ'},      // 31F7 [ㇷ] -> FF8C [ﾌ]
	{'ㇸ', arCompatWidth, 'ﾍ'},      // 31F8 [ㇸ] -> FF8D [ﾍ]
	{'ㇹ', arCompatWidth, 'ﾎ'},      // 31F9 [ㇹ] -> FF8E [ﾎ]
	{'ㇺ', arCompatWidth, 'ﾑ'},      // 31FA [ㇺ] -> FF91 [ﾑ]
	{'ㇻ', arCompatWidth, 'ﾗ'},      // 31FB [ㇻ] -> FF97 [ﾗ]
	{'ㇼ', arCompatWidth, 'ﾘ'},      // 31FC [ㇼ] -> FF98 [ﾘ]
	{'ㇽ', arCompatWidth, 'ﾙ'},      // 31FD [ㇽ] -> FF99 [ﾙ]
	{'ㇾ', arCompatWidth, 'ﾚ'},      // 31FE [ㇾ] -> FF9A [ﾚ]
	{'ㇿ', arCompatWidth, 'ﾛ'},      // 31FF [ㇿ] -> FF9B [ﾛ]
}

// Additional attr directive
const (
	aaCategory = iota
	aaCharCase
	aaCharWidth
	aaVoicing
)

type additionalAttr struct {
	codepoint rune
	attr      int
	data      uint8
}

var additionalAttrList = []additionalAttr{
	{'　', aaCategory, ctLatinSymbol},   // 0x3000 [　] ctKanaSymbol -> ctLatinSymbol
	{'\u3099', aaCategory, ctKanaVsm},  // 0x3099 [ ◌゙] ctKanaLetter -> ctKanaVsm
	{'\u309A', aaCategory, ctKanaVsm},  // 0x309A [゚゚ ゚] ctKanaLetter -> ctKanaVsm
	{'゛', aaCategory, ctKanaVsm},       // 0x309B [゛] ctKanaLetter -> ctKanaVsm
	{'゜', aaCategory, ctKanaVsm},       // 0x309C [゜] ctKanaLetter -> ctKanaVsm
	{'ﾞ', aaCategory, ctKanaVsm},       // FF9E [ﾞ] ctKanaLetter -> ctKanaVsm
	{'ﾟ', aaCategory, ctKanaVsm},       // FF9F [ﾟ] ctKanaLetter -> ctKanaVsm
	{'\u3099', aaCharCase, ccNonspace}, // 0x3099 [ ◌゙] ccHiragana -> ccTraditional
	{'\u309A', aaCharCase, ccNonspace}, // 0x309A [゚゚ ゚] ccHiragana -> ccTraditional
	{'゛', aaCharCase, ccTraditional},   // 0x309B [゛] ccHiragana -> ccTraditional
	{'゜', aaCharCase, ccTraditional},   // 0x309C [゜] ccHiragana -> ccTraditional
	{'ﾞ', aaCharCase, ccTraditional},   // FF9E [ﾞ] ccKatakana -> ccTraditional
	{'ﾟ', aaCharCase, ccTraditional},   // FF9F [ﾟ] ccKatakana -> ccTraditional
}

// Get the multiple codepoints from Cp(UCD)
func multiRunesFromCp(codepoints string) ([]rune, error) {
	if codepoints == "" || codepoints == "#" {
		return []rune{}, nil
	}
	rs := []rune{}
	ss := strings.Split(codepoints, " ")
	for _, codepoint := range ss {
		n, err := strconv.ParseInt(codepoint, 16, 32)
		if err != nil {
			return []rune{}, err
		}
		rs = append(rs, rune(n))
	}
	return rs, nil
}

// Get the single codepoint from Cp(UCD)
func singleRuneFromCp(codepoints string) (r rune, err error) {
	rs, err := multiRunesFromCp(codepoints)
	if err != nil {
		return rune(0), err
	}
	if len(rs) != 1 {
		return rune(0), fmt.Errorf("want: single codepoint, have: %s", codepoints)
	}
	return rs[0], nil
}

func isTargetRune(codepoint rune) bool {
	for _, b := range blockRanges {
		if b.first <= codepoint && codepoint <= b.last {
			return true
		}
	}
	return false
}

func isTargetCp(codepoints string) (bool, error) {
	rs, err := multiRunesFromCp(codepoints)
	if err != nil {
		return false, err
	}
	if len(rs) <= 0 {
		return false, nil
	}
	return isTargetRune(rs[0]), nil
}

func char2category(char *Char) (uint8, error) {
	switch char.Blk {
	case "ASCII":
		switch char.Gc {
		case "Lu", "Ll":
			return ctLatinLetter, nil
		case "Nd":
			return ctLatinDigit, nil
		default:
			return ctLatinSymbol, nil
		}

	case "CJK_Symbols":
		return ctKanaSymbol, nil

	case "Hiragana":
		return ctKanaLetter, nil

	case "Katakana":
		switch char.Gc {
		case "Lo", "Lm":
			return ctKanaLetter, nil
		default:
			return ctKanaSymbol, nil
		}

	case "Half_And_Full_Forms":
		switch char.Dt {
		case "wide":
			switch char.Gc {
			case "Lu", "Ll":
				return ctLatinLetter, nil
			case "Nd":
				return ctLatinDigit, nil
			default:
				return ctLatinSymbol, nil
			}

		case "nar":
			switch char.Gc {
			case "Lo", "Lm":
				return ctKanaLetter, nil
			default:
				return ctKanaSymbol, nil
			}
		}

	case "Katakana_Ext":
		return ctKanaLetter, nil

	default:
		return 0, fmt.Errorf("unexpected char.Blk: %q", char.Blk)
	}
	return ctUndefined, nil
}

func char2charCase(char *Char) (uint8, error) {
	switch char.Blk {
	case "ASCII":
		switch char.Gc {
		case "Lu":
			return ccUpper, nil
		case "Ll":
			return ccLower, nil
		default:
			return ccUndefined, nil
		}

	case "CJK_Symbols":
		return ccUndefined, nil

	case "Hiragana":
		return ccHiragana, nil

	case "Katakana":
		switch char.Gc {
		case "Lo", "Lm":
			return ccKatakana, nil
		default:
			return ccUndefined, nil
		}

	case "Half_And_Full_Forms":
		switch char.Dt {
		case "wide":
			switch char.Gc {
			case "Lu":
				return ccUpper, nil
			case "Ll":
				return ccLower, nil
			default:
				return ccUndefined, nil
			}

		case "nar":
			switch char.Gc {
			case "Lo", "Lm":
				return ccKatakana, nil
			default:
				return ccUndefined, nil
			}
		}
	case "Katakana_Ext":
		return ccKatakana, nil

	default:
		return 0, fmt.Errorf("unexpected char.Blk: %q", char.Blk)
	}
	return ccUndefined, nil
}

func char2charWidth(char *Char) (charWidth uint8, compatWidth rune, err error) {
	charWidth = cwUndefined
	switch char.Blk {
	case "ASCII":
		charWidth = cwNarrow
	case "CJK_Symbols", "Hiragana", "Katakana", "Katakana_Ext":
		charWidth = cwWide
	case "Half_And_Full_Forms":
		switch char.Dt {
		case "wide":
			charWidth = cwWide
		case "nar":
			charWidth = cwNarrow
		default:
			return 0, rune(0), fmt.Errorf("unexpected char.Dt: %q", char.Dt)
		}
	default:
		return 0, rune(0), fmt.Errorf("unexpected char.Blk: %q", char.Blk)
	}

	compatWidth = rune(0)
	isTarget, err := isTargetCp(char.Dm)
	if err != nil {
		return 0, rune(0), err
	}
	if isTarget && (char.Dt == "wide" || char.Dt == "nar") {
		compatWidth, err = singleRuneFromCp(char.Dm)
		if err != nil {
			return 0, rune(0), err
		}
	}

	return charWidth, compatWidth, nil
}

func char2compatCase(char *Char) (r rune, err error) {
	switch char.Gc {
	case "Lu":
		r, err = singleRuneFromCp(char.Slc)
	case "Ll":
		r, err = singleRuneFromCp(char.Suc)
	default:
		r = rune(0)
	}
	if err != nil {
		return rune(0), err
	}
	return r, nil
}

func char2voicing(char *Char) (voicing uint8, compatVs, compatSvs rune, err error) {
	if char.Blk != "Hiragana" && char.Blk != "Katakana" {
		return vcUndefined, rune(0), rune(0), nil
	}
	if char.Gc != "Lo" {
		return vcUndefined, rune(0), rune(0), nil
	}
	if char.Dt != "can" {
		return vcUndefined, rune(0), rune(0), nil
	}

	rs, err := multiRunesFromCp(char.Dm)
	if err != nil {
		return 0, rune(0), rune(0), err
	}
	if len(rs) != 2 {
		return 0, rune(0), rune(0),
			fmt.Errorf("want: double codepoints, have: %s", char.Dm)
	}
	switch rs[1] {
	case 0x3099:
		return vcVoiced, rs[0], rune(0), nil
	case 0x309A:
		return vcSemivoiced, rune(0), rs[0], nil
	default:
		return 0, rune(0), rune(0),
			fmt.Errorf("want: 2nd charis 3099 or 309A, have: %s", char.Dm)
	}
}

func updateKanaRelation(m ucdex) error {
	for _, charmap := range hiragana2katakana {
		for cpHiragana := charmap.lo; cpHiragana <= charmap.hi; cpHiragana++ {
			ucHiragana, ok := m[cpHiragana]
			if !ok {
				return fmt.Errorf("updateKanaRelation; %#U is not exists in ucdex", cpHiragana)
			}
			if ucHiragana.charCase != ccHiragana {
				return fmt.Errorf("updateKanaRelation; %#U.charCase is %s, want ccHiragana",
					ucHiragana.codepoint, charCases.name(ucHiragana.charCase))
			}

			cpKatakana := cpHiragana + charmap.diff
			ucKatakana, ok := m[cpKatakana]
			if !ok {
				return fmt.Errorf("updateKanaRelation; %#U is not exists in ucdex", cpKatakana)
			}
			if ucKatakana.charCase != ccKatakana {
				return fmt.Errorf("updateKanaRelation; %#U.charCase is %s, want ccKatakana",
					ucKatakana.codepoint, charCases.name(ucKatakana.charCase))
			}

			// The main theme of this function
			ucHiragana.compatCase = ucKatakana.codepoint
			ucKatakana.compatCase = ucHiragana.codepoint
		}
	}
	return nil
}

func updateLatinRelation(c *charex, m ucdex) error {
	// Only full-width latin characters are targeted.
	// However, the following characters are excluded.
	// (because the related half-width characters are out of target range)
	// U+FF5F '｟' FULLWIDTH LEFT WHITE PARENTHESIS
	// U+FF60 '｠' FULLWIDTH RIGHT WHITE PARENTHESIS
	if c.charWidth != cwWide || c.category == ctUndefined ||
		c.category == ctKanaLetter || c.category == ctKanaSymbol || c.category == ctKanaVsm ||
		c.codepoint == '｟' || c.codepoint == '｠' {
		return nil
	}

	wideLatin := c
	if wideLatin.compatWidth == 0 {
		return fmt.Errorf("updateLatinRelation; %#U.compatWidth is rune(0)", wideLatin.codepoint)
	}

	narrowLatin, ok := m[wideLatin.compatWidth]
	if !ok {
		return fmt.Errorf("updateLatinRelation; %#U.compatWidth -> %#U is not exists in ucdex",
			wideLatin.codepoint, wideLatin.compatWidth)
	}

	if narrowLatin.category != ctLatinLetter &&
		narrowLatin.category != ctLatinDigit &&
		narrowLatin.category != ctLatinSymbol {
		return fmt.Errorf("updateLatinRelation; %#U.compatWidth -> %#U.category is %d, want ctLatinXxx",
			wideLatin.codepoint, narrowLatin.codepoint, narrowLatin.category)
	}

	if narrowLatin.charWidth != cwNarrow {
		return fmt.Errorf("updateLatinRelation; %#U.compatWidth -> %#U.charWidth is %d, want cwNarrow",
			wideLatin.codepoint, narrowLatin.codepoint, narrowLatin.charWidth)
	}

	if narrowLatin.compatWidth != 0 {
		return fmt.Errorf("updateLatinRelation; %#U.compatWidth -> %#U.compatWidth is %#U, want 0",
			wideLatin.codepoint, narrowLatin.codepoint, narrowLatin.compatWidth)
	}

	// The main theme of this function
	narrowLatin.compatWidth = wideLatin.codepoint

	return nil
}

func updateKanaLetterRelation(c *charex, m ucdex) error {
	// Only half-width katakana characters are targeted.
	// However, the following special katakana leters are excluded.
	// U+FF70 'ｰ' HALFWIDTH KATAKANA-HIRAGANA PROLONGED SOUND MARK
	// U+FF9E 'ﾞ' HALFWIDTH KATAKANA VOICED SOUND MARK
	// U+FF9F 'ﾟ' HALFWIDTH KATAKANA SEMI-VOICED SOUND MARK
	if c.blk != "Half_And_Full_Forms" ||
		c.charWidth != cwNarrow || c.category == ctKanaSymbol ||
		c.codepoint == 'ｰ' || c.codepoint == 'ﾞ' || c.codepoint == 'ﾟ' {
		return nil
	}

	narrowKatakana := c
	if narrowKatakana.compatCase != 0 {
		return fmt.Errorf("updateKanaLetterRelation; %#U.compatCase is not 0", narrowKatakana.codepoint)
	}
	if narrowKatakana.compatWidth == 0 {
		return fmt.Errorf("updateKanaLetterRelation; %#U.compatWidth is 0", narrowKatakana.codepoint)
	}

	wideKatakana, ok := m[narrowKatakana.compatWidth]
	if !ok {
		return fmt.Errorf("updateKanaLetterRelation; %#U.compatWidth -> %#U is not exists in ucdex",
			narrowKatakana.codepoint, narrowKatakana.compatWidth)
	}
	if wideKatakana.category != ctKanaLetter {
		return fmt.Errorf("updateKanaLetterRelation; %#U.compatWidth -> %#U.category is %s, want ctKanaLetter",
			narrowKatakana.codepoint, wideKatakana.codepoint, categories.name(wideKatakana.category))
	}
	if wideKatakana.charCase != ccKatakana {
		return fmt.Errorf("updateKanaLetterRelation; %#U.compatWidth -> %#U.charCase is %s, want ccKatakana",
			narrowKatakana.codepoint, wideKatakana.codepoint, charCases.name(wideKatakana.charCase))
	}
	if wideKatakana.charWidth != cwWide {
		return fmt.Errorf("updateKanaLetterRelation; %#U.compatWidth -> %#U.charWidth is %s, want cwWide",
			narrowKatakana.codepoint, wideKatakana.codepoint, charWidths.name(wideKatakana.charWidth))
	}
	if wideKatakana.compatWidth != 0 {
		return fmt.Errorf("updateKanaLetterRelation; %#U.compatWidth -> %#U.compatWidth is not 0",
			narrowKatakana.codepoint, wideKatakana.codepoint)
	}
	if wideKatakana.compatCase == 0 {
		return fmt.Errorf("updateKanaLetterRelation; %#U.compatWidth -> %#U.compatCase is 0",
			narrowKatakana.codepoint, wideKatakana.codepoint)
	}

	wideHiragana, ok := m[wideKatakana.compatCase]
	if !ok {
		return fmt.Errorf("updateKanaLetterRelation; %#U.compatCase -> %#U is not exists in ucdex",
			narrowKatakana.codepoint, narrowKatakana.compatCase)
	}
	if wideHiragana.category != ctKanaLetter {
		return fmt.Errorf("updateKanaLetterRelation; %#U.compatWidth -> %#U.compatCase -> %#U.category is %s, want ctKanaLetter",
			narrowKatakana.codepoint, wideKatakana.codepoint,
			wideHiragana.codepoint, categories.name(wideHiragana.category))
	}
	if wideHiragana.charCase != ccHiragana {
		return fmt.Errorf("updateKanaLetterRelation; %#U.compatWidth -> %#U.compatCase -> %#U.charCase is %s, want ccHiragana",
			narrowKatakana.codepoint, wideKatakana.codepoint,
			wideHiragana.codepoint, charCases.name(wideHiragana.charCase))
	}
	if wideHiragana.charWidth != cwWide {
		return fmt.Errorf("updateKanaLetterRelation; %#U.compatWidth -> %#U.compatCase -> %#U.charWidth is %s, want cwWide",
			narrowKatakana.codepoint, wideKatakana.codepoint,
			wideHiragana.codepoint, charWidths.name(wideHiragana.charWidth))
	}
	if wideHiragana.compatWidth != 0 {
		return fmt.Errorf("updateKanaLetterRelation; %#U.compatWidth -> %#U.compatCase -> %#U.compatWidth -> %#U is not 0",
			narrowKatakana.codepoint, wideKatakana.codepoint, wideHiragana.codepoint, wideHiragana.compatWidth)
	}

	// The main theme of this function
	wideKatakana.compatWidth = narrowKatakana.codepoint
	wideHiragana.compatWidth = narrowKatakana.codepoint
	narrowKatakana.compatCase = wideHiragana.codepoint

	return nil
}

func updateKanaSymbolRelation(c *charex, m ucdex) error {
	// Only half-width kana symbol characters are targeted.
	if c.blk != "Half_And_Full_Forms" || c.charWidth != cwNarrow || c.category != ctKanaSymbol {
		return nil
	}

	narrowKanaSymbol := c
	if narrowKanaSymbol.compatWidth == 0 {
		return fmt.Errorf("updateKanaSymbolRelation; %#U.compatWidth is 0", narrowKanaSymbol.codepoint)
	}

	wideKanaSymbol, ok := m[narrowKanaSymbol.compatWidth]
	if !ok {
		return fmt.Errorf("updateKanaSymbolRelation; %#U.compatWidth -> %#U is not exists in ucdex",
			narrowKanaSymbol.codepoint, narrowKanaSymbol.compatWidth)
	}
	if wideKanaSymbol.category != ctKanaSymbol {
		return fmt.Errorf("updateKanaSymbolRelation; %#U.compatWidth -> %#U.category is %s, want ctKanaSymbol",
			narrowKanaSymbol.codepoint, wideKanaSymbol.codepoint, categories.name(wideKanaSymbol.category))
	}
	if wideKanaSymbol.charWidth != cwWide {
		return fmt.Errorf("updateKanaSymbolRelation; %#U.compatWidth -> %#U.charWidth is %s, want cwWide",
			narrowKanaSymbol.codepoint, wideKanaSymbol.codepoint, charWidths.name(wideKanaSymbol.charWidth))
	}
	if wideKanaSymbol.compatWidth != 0 {
		return fmt.Errorf("updateKanaSymbolRelation; %#U.compatWidth -> %#U.compatWidth is not 0",
			narrowKanaSymbol.codepoint, wideKanaSymbol.codepoint)
	}

	// The main theme of this function
	wideKanaSymbol.compatWidth = narrowKanaSymbol.codepoint

	return nil
}

func updateVoicingRelation(c *charex, m ucdex) error {
	// Only voiced sound characters are targeted.
	if c.voicing != vcVoiced {
		return nil
	}

	voiced := c

	if voiced.compatVs == 0 {
		return fmt.Errorf("updateVoicingRelation; %#U.compatVs is not 0", voiced.codepoint)
	}

	unvoiced, ok := m[voiced.compatVs]
	if !ok {
		return fmt.Errorf("updateVoicingRelation; %#U.compatVs -> %#U is not exists in ucdex",
			voiced.codepoint, voiced.compatVs)
	}
	if unvoiced.voicing != vcUndefined && unvoiced.voicing != vcUnvoiced {
		return fmt.Errorf("updateVoicingRelation; %#U.compatVs -> %#U.voiceing = %s, want %s or %s",
			voiced.codepoint, unvoiced.codepoint, voicings.name(unvoiced.voicing),
			voicings.name(vcUndefined), voicings.name(vcUnvoiced))
	}
	if unvoiced.compatVs != 0 {
		return fmt.Errorf("updateVoicingRelation; %#U.compatVs -> %#U.compatVs -> %#U is not 0",
			voiced.compatVs, unvoiced.codepoint, unvoiced.compatVs)
	}

	// The main theme of this function
	unvoiced.voicing = vcUnvoiced
	unvoiced.compatVs = voiced.codepoint

	return nil
}

func updateSemivoicingRelation(c *charex, m ucdex) error {
	// Only semi-voiced sound characters are targeted.
	if c.voicing != vcSemivoiced {
		return nil
	}

	semivoiced := c

	if semivoiced.compatSvs == 0 {
		return fmt.Errorf("updateSemivoicingRelation; %#U.compatSvs is not 0", semivoiced.codepoint)
	}

	unvoiced, ok := m[semivoiced.compatSvs]
	if !ok {
		return fmt.Errorf("updateSemivoicingRelation; %#U.compatSvs -> %#U is not exists in ucdex",
			semivoiced.codepoint, semivoiced.compatSvs)
	}
	if unvoiced.voicing != vcUndefined && unvoiced.voicing != vcUnvoiced {
		return fmt.Errorf("updateSemivoicingRelation; %#U.compatSvs -> %#U.voiceing = %s, want %s or %s",
			semivoiced.codepoint, unvoiced.codepoint, voicings.name(unvoiced.voicing),
			voicings.name(vcUndefined), voicings.name(vcUnvoiced))
	}
	if unvoiced.compatSvs != 0 {
		return fmt.Errorf("updateSemivoicingRelation; %#U.compatSs -> %#U.compatSvs -> %#U is not 0",
			semivoiced.compatSvs, unvoiced.codepoint, unvoiced.compatSvs)
	}

	// The main theme of this function
	unvoiced.voicing = vcUnvoiced
	unvoiced.compatSvs = semivoiced.codepoint

	return nil
}

func updateVoicingWidthRelation(c *charex, m ucdex) error {
	// Only voiced or semi-voiced sound characters are targeted.
	if c.voicing != vcVoiced && c.voicing != vcSemivoiced {
		return nil
	}

	voiced := c

	if voiced.compatWidth != 0 {
		return fmt.Errorf("updateVoicingWidthRelation; %#U.compatWidth is not 0", voiced.codepoint)
	}
	if voiced.compatVs == 0 && voiced.compatSvs == 0 {
		return fmt.Errorf("updateVoicingWidthRelation; %#U.compatVs and compatSvs are 0", voiced.codepoint)
	}
	if voiced.compatVs != 0 && voiced.compatSvs != 0 {
		return fmt.Errorf("updateVoicingWidthRelation; %#U.compatVs is %#U and compatSvs is %#U, want Either one is 0",
			voiced.codepoint, voiced.compatVs, voiced.compatSvs)
	}

	var unvoiced *charex
	var ok bool
	if voiced.compatVs != 0 {
		unvoiced, ok = m[voiced.compatVs]
		if !ok {
			return fmt.Errorf("updateVoicingWidthRelation; %#U.compatVs -> %#U is not exists in ucdex",
				voiced.codepoint, voiced.compatVs)
		}
		if unvoiced.voicing != vcUnvoiced {
			return fmt.Errorf("updateVoicingWidthRelation; %#U.compatVs -> %#U.voicing is %s, want vcUnvoiced",
				voiced.codepoint, unvoiced.codepoint, voicings.name(unvoiced.voicing))
		}
		if unvoiced.compatWidth == 0 {
			// The following characters have no related half-width characters.
			// U+30F1 'ヱ', U+30F0 'ヰ'
			// log.Printf("INFO: %#U have no related half-width character", unvoiced.codepoint)
		}
	} else {
		unvoiced, ok = m[voiced.compatSvs]
		if !ok {
			return fmt.Errorf("updateVoicingWidthRelation; %#U.compatSvs -> %#U is not exists in ucdex",
				voiced.codepoint, voiced.compatSvs)
		}
		if unvoiced.voicing != vcUnvoiced {
			return fmt.Errorf("updateVoicingWidthRelation; %#U.compatSvs -> %#U.voicing is %s, want vcUnvoiced",
				voiced.codepoint, unvoiced.codepoint, voicings.name(unvoiced.voicing))
		}
		if unvoiced.compatWidth == 0 {
			return fmt.Errorf("updateVoicingWidthRelation; %#U.compatSvs -> %#U.compatWidth is 0",
				voiced.codepoint, unvoiced.codepoint)
		}
	}

	// The main theme of this function
	voiced.compatWidth = unvoiced.compatWidth
	if voiced.compatCase == 0 {
		// The following characters have no related hiragana characters.
		// U+30F7 'ヷ', U+30F8 'ヸ', U+30F9 'ヹ', U+30FA 'ヺ'
		// log.Printf("INFO: %#U have no related hiragana character", voiced.codepoint)
		voiced.compatCase = unvoiced.compatCase
	}

	return nil
}

func updateAdditionalRefList(m ucdex) error {
	for _, ad := range additionalRefList {
		c, ok := m[ad.codepoint]
		if !ok {
			return fmt.Errorf("updateAdditionalRefList; %#U is not exists in ucdex", ad.codepoint)
		}
		switch ad.ref {
		case arCompatCase:
			c.compatCase = ad.data
		case arCompatWidth:
			c.compatWidth = ad.data
		case arCompatVs:
			c.compatVs = ad.data
		case arCompatSvs:
			c.compatSvs = ad.data
		default:
			return fmt.Errorf("updateAdditionalRefList; Unexpected ad.ref: %d", ad.ref)
		}
	}
	return nil
}

func updateAdditionalAttrList(m ucdex) error {
	for _, ad := range additionalAttrList {
		c, ok := m[ad.codepoint]
		if !ok {
			return fmt.Errorf("updateAdditionalAttrList; %#U is not exists in ucdex", ad.codepoint)
		}
		switch ad.attr {
		case aaCategory:
			c.category = ad.data
		case aaCharCase:
			c.charCase = ad.data
		case aaCharWidth:
			c.charWidth = ad.data
		case aaVoicing:
			c.voicing = ad.data
		default:
			return fmt.Errorf("updateAdditionalAttrList; Unexpected ad.attr: %d", ad.attr)
		}
	}
	return nil
}

func createUCDEX(ucd *UCD) (ucdex, error) {
	// Create UCDEX from UCD
	m := make(ucdex, len(ucd.Chars))
	for _, char := range ucd.Chars {
		isTarget, err := isTargetCp(char.Cp)
		if err != nil {
			return nil, err
		}
		if !isTarget {
			continue
		}

		codepoints, err := multiRunesFromCp(char.Cp)
		if err != nil {
			return nil, err
		}
		if len(codepoints) == 0 {
			continue
		}

		codepoint := codepoints[0]
		category, err := char2category(&char)
		if err != nil {
			return nil, err
		}
		charCase, err := char2charCase(&char)
		if err != nil {
			return nil, err
		}
		compatCase, err := char2compatCase(&char)
		if err != nil {
			return nil, err
		}
		charWidth, compatWidth, err := char2charWidth(&char)
		if err != nil {
			return nil, err
		}
		voicing, compatVs, compatSvs, err := char2voicing(&char)
		if err != nil {
			return nil, err
		}

		m[codepoint] = &charex{
			codepoint:   codepoint,
			blk:         char.Blk,
			na:          char.Na,
			age:         char.Age,
			gc:          char.Gc,
			category:    category,
			charCase:    charCase,
			charWidth:   charWidth,
			compatCase:  compatCase,
			voicing:     voicing,
			compatWidth: compatWidth,
			compatVs:    compatVs,
			compatSvs:   compatSvs,
		}
	}

	// Update UCDEX
	var err error
	if err = updateKanaRelation(m); err != nil {
		return nil, err
	}
	if err = updateAdditionalRefList(m); err != nil {
		return nil, err
	}
	if err = updateAdditionalAttrList(m); err != nil {
		return nil, err
	}
	for _, unichar := range m {
		if err = updateLatinRelation(unichar, m); err != nil {
			return nil, err
		}
		if err = updateKanaLetterRelation(unichar, m); err != nil {
			return nil, err
		}
		if err = updateKanaSymbolRelation(unichar, m); err != nil {
			return nil, err
		}
		if err = updateVoicingRelation(unichar, m); err != nil {
			return nil, err
		}
		if err = updateSemivoicingRelation(unichar, m); err != nil {
			return nil, err
		}
	}
	for _, unichar := range m {
		if err = updateVoicingWidthRelation(unichar, m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func formatRune(r rune) string {
	if r <= 0 {
		return ""
	}
	return fmt.Sprintf("U+%04X", r)
}

func escapeChar(r rune) string {
	switch {
	case r == 0:
		return ""
	case 0 < r && r < 0x20:
		return fmt.Sprintf("^%s", string([]rune{'@' + r}))
	case r == '"':
		return "(QUOT)"
	case r == ',':
		return "(COMMA)"
	default:
		return string([]rune{r})
	}
}

func writeUCDEX(f io.Writer, m ucdex) {
	fmt.Fprint(f, "Blk,Na,Age,Gc,codepoint,ch,category,char_case,compat_case,ch,")
	fmt.Fprint(f, "char_width,compat_width,ch,voicing,compat_vs,ch,compat_svs,ch\n")

	for _, b := range blockRanges {
		for i := b.first; i <= b.last; i++ {
			if c, ok := m[i]; !ok {
				fmt.Fprintf(f, ",(not present in the ucd),,,U+%04X,,ctUndefined,"+
					"ccUndefined,,,cwUndefined,,,vcUndefined,,,,\n", i)
			} else {
				fmt.Fprintf(f, "%s", c.blk)
				fmt.Fprintf(f, ",%s", c.na)
				fmt.Fprintf(f, ",%s", c.age)
				fmt.Fprintf(f, ",%s", c.gc)
				fmt.Fprintf(f, ",%s", formatRune(c.codepoint))
				fmt.Fprintf(f, ",%s", escapeChar(c.codepoint))
				fmt.Fprintf(f, ",%s", categories.name(c.category))
				fmt.Fprintf(f, ",%s", charCases.name(c.charCase))
				fmt.Fprintf(f, ",%s", formatRune(c.compatCase))
				fmt.Fprintf(f, ",%s", escapeChar(c.compatCase))
				fmt.Fprintf(f, ",%s", charWidths.name(c.charWidth))
				fmt.Fprintf(f, ",%s", formatRune(c.compatWidth))
				fmt.Fprintf(f, ",%s", escapeChar(c.compatWidth))
				fmt.Fprintf(f, ",%s", voicings.name(c.voicing))
				fmt.Fprintf(f, ",%s", formatRune(c.compatVs))
				fmt.Fprintf(f, ",%s", escapeChar(c.compatVs))
				fmt.Fprintf(f, ",%s", formatRune(c.compatSvs))
				fmt.Fprintf(f, ",%s", escapeChar(c.compatSvs))
				fmt.Fprintln(f, "")
			}
		}
	}
}

// GenUCDEX generates the UCDEX in CSV
func GenUCDEX(f io.Writer) error {
	ucd, err := readUCD()
	if err != nil {
		return err
	}

	ucdex, err := createUCDEX(ucd)
	if err != nil {
		return err
	}

	writeUCDEX(f, ucdex)

	return nil
}

func printConstList(f io.Writer, c constValues) {
	fmt.Fprintf(f, "// %s\n", c.title)
	fmt.Fprintln(f, "const (")
	fmt.Fprintf(f, "\t%s = iota\n", c.list[0])
	for i := 1; i < len(c.list); i++ {
		fmt.Fprintf(f, "\t%s\n", c.list[i])
	}
	fmt.Fprint(f, ")\n\n")
}

// TODO consider implementation in uint8 (category, charCase, charWidth, voicing)
func printTypes(f io.Writer) {
	text := `
type unichar struct {
	codepoint   rune  // Unicode code point value
	category    uint8 // ctUndefined/ctLatinLetter/ctLatinDigit/ctLatinSymbol/ctKanaLetter/ctKanaSymbol
	charCase    uint8 // ccUndefined/ccUpper/ccLower/ccHiragana/ccKatakana
	charWidth   uint8 // cwUndefined/cwNarrow/cwWide
	voicing     uint8 // vcUndefined/vcUnvoiced/vcVoiced/vcSemivoiced
	compatCase  rune  // Charcase compatible character (Upper-Lower, Hiragana-Katakana)
	compatWidth rune  // Width compatible character (Narrow-Wide)
	compatVs    rune  // Voiced sound compatible character (Unvoiced-Voiced)
	compatSvs   rune  // Semi-voiced sound compatible character (Unvoiced-Semivoiced)
}

type unichars []unichar
`
	fmt.Fprintln(f, text)
}

func nonZeroOrElse(r rune, alt rune) string {
	if r == 0 {
		return fmt.Sprintf("%q", alt)
	}
	return fmt.Sprintf("%q", r)
}

func generate(f io.Writer, m ucdex, genname string) {
	fmt.Fprintf(f, "// Code generated by %s; DO NOT EDIT.\n", genname)
	fmt.Fprintf(f, "// Based on information from %s\n\n", ucdURL)
	fmt.Fprint(f, "package gaga\n\n")

	printConstList(f, categories)
	printConstList(f, charCases)
	printConstList(f, charWidths)
	printConstList(f, voicings)

	printTypes(f)

	for _, b := range blockRanges {
		fmt.Fprintf(f, "var %sFirst rune = 0x%04X\n", b.name, b.first)
		fmt.Fprintf(f, "var %sLast  rune = 0x%04X\n", b.name, b.last)
		fmt.Fprintf(f, "var %sTable = unichars {\n", b.name)
		for i := b.first; i <= b.last; i++ {
			fmt.Fprintf(f, "\t{")
			if c, ok := m[i]; !ok {
				fmt.Fprintf(f, "0x%04X", i)
				fmt.Fprint(f, ",ctUndefined, ccUndefined, cwUndefined, vcUndefined")
				fmt.Fprintf(f, ",%q", i)
				fmt.Fprintf(f, ",%q", i)
				fmt.Fprintf(f, ",%q", i)
				fmt.Fprintf(f, ",%q", i)
			} else {
				fmt.Fprintf(f, "0x%04X", c.codepoint)
				fmt.Fprintf(f, ",%s", categories.name(c.category))
				fmt.Fprintf(f, ",%s", charCases.name(c.charCase))
				fmt.Fprintf(f, ",%s", charWidths.name(c.charWidth))
				fmt.Fprintf(f, ",%s", voicings.name(c.voicing))
				fmt.Fprintf(f, ",%s", nonZeroOrElse(c.compatCase, c.codepoint))
				fmt.Fprintf(f, ",%s", nonZeroOrElse(c.compatWidth, c.codepoint))
				fmt.Fprintf(f, ",%s", nonZeroOrElse(c.compatVs, c.codepoint))
				fmt.Fprintf(f, ",%s", nonZeroOrElse(c.compatSvs, c.codepoint))
			}
			fmt.Fprintf(f, "}, // 0x%04X %s\n", i, string([]rune{i}))

		}
		fmt.Fprint(f, "}\n\n")
	}
}

// Generate generates the array of UCDEX (Go source code)
func Generate(f io.Writer, genname string) error {
	ucd, err := readUCD()
	if err != nil {
		return err
	}

	ucdex, err := createUCDEX(ucd)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	generate(&buf, ucdex, genname)
	out, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	_, err = f.Write(out)
	if err != nil {
		return err
	}

	return nil
}
