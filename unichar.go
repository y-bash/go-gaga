package gaga

import (
	"fmt"
)

//go:generate go run gen/gentables.go -output unichar_tables.go

// Voicde or Semi-voiced sound mark
const (
	vsmCombining  = rune(0x3099)
	vsmWide       = rune(0x309B)
	vsmNarrow     = rune(0xFF9E)
	svsmCombining = rune(0x309A)
	svsmWide      = rune(0x309C)
	svsmNarrow    = rune(0xFF9F)
)

func isVoicedSoundMark(r rune) bool {
	return r == vsmCombining || r == vsmWide || r == vsmNarrow
}

func isSemivoicedSoundMark(r rune) bool {
	return r == svsmCombining || r == svsmWide || r == svsmNarrow
}

func getUnichar(r rune) (c *unichar, ok bool) {
	switch {
	case latinFirst <= r && r <= latinLast:
		return &latinTable[r-latinFirst], true
	case kanaFirst <= r && r <= kanaLast:
		return &kanaTable[r-kanaFirst], true
	case kanaExtFirst <= r && r <= kanaExtLast:
		return &kanaExtTable[r-kanaExtFirst], true
	case widthFirst <= r && r <= widthLast:
		return &widthTable[r-widthFirst], true
	default:
		return nil, false
	}
}

// This function can only be called if r is known to exist in the tables.
func getUnicharForSure(r rune) *unichar {
	c, ok := getUnichar(r)
	if !ok {
		panic(fmt.Sprintf("Unexpectedly %#U.compatCase %#U was not found in the table",
			c.codepoint, c.compatCase))
	}
	return c
}

func (c *unichar) getCompatCaseUnichar() *unichar {
	// TEST_fm8XjZTB ensured that all compatCases are in the tables.
	return getUnicharForSure(c.compatCase)
}

func (c *unichar) getCompatWidthUnichar() *unichar {
	// TEST_T3bc4Nh7 ensured that all compatWidth are in the table.
	return getUnicharForSure(c.compatWidth)
}

func (c *unichar) getCompatVsUnichar() *unichar {
	// TEST_Cu8iKMxF ensured that all compatVs are in the tables.
	return getUnicharForSure(c.compatVs)
}

func (c *unichar) getCompatSvsUnichar() *unichar {
	// TEST_rW4UiNHC ensured that all compatSvs are in the tables.
	return getUnicharForSure(c.compatSvs)
}

func (c *unichar) existsCompatVs() bool {
	return c.codepoint != c.compatVs
}

func (c *unichar) existsCompatSvs() bool {
	return c.codepoint != c.compatSvs
}

func (c *unichar) toUpper() rune {
	if c.charCase != ccLower {
		return c.codepoint
	}
	return c.compatCase
}

func (c *unichar) toLower() rune {
	if c.charCase != ccUpper {
		return c.codepoint
	}
	return c.compatCase
}

func (c *unichar) toHiraganaUnichar() *unichar {
	if c.charCase != ccKatakana {
		return c
	}
	return c.getCompatCaseUnichar()
}

func (c *unichar) toKatakanaUnichar() *unichar {
	if c.charCase != ccHiragana {
		return c
	}
	return c.getCompatCaseUnichar()
}

func (c *unichar) toWide() rune {
	if c.charWidth != cwNarrow {
		return c.codepoint
	}
	return c.compatWidth
}

func (c *unichar) toWideUnichar() *unichar {
	if c.charWidth != cwNarrow {
		return c
	}
	return c.getCompatWidthUnichar()
}

func (c *unichar) toNarrow() rune {
	if c.charWidth != cwWide {
		return c.codepoint
	}
	return c.compatWidth
}

func (c *unichar) toNarrowUnichar() *unichar {
	if c.charWidth != cwWide {
		return c
	}
	return c.getCompatWidthUnichar()
}

// for voiced or semi-voiced sound mark characters.
func (c *unichar) toTraditionalMarkUnichar() *unichar {
	if c.charCase != ccCombining {
		return c
	}
	return c.getCompatCaseUnichar()
}

// for voiced or semi-voiced sound mark characters.
func (c *unichar) toCombiningMark() rune {
	if c.charCase != ccTraditional {
		return c.codepoint
	}
	return c.compatCase
}

// for Hiragana-Katakana letters.
// TEST_Vs4Ad89Z knows that this function returns a rune array with
// 1 or 2 elements and no other number of elements.
func (c *unichar) toTraditionalVoiced() []rune {
	switch c.voicing {
	case vcVoiced:
		return []rune{c.codepoint}
	case vcSemivoiced:
		// TEST_fW6auXUi knows that every semi-voiced character has
		// a corresponding unvoiced character, and that unvoiced
		// character has a corresponding voiced character.
		return []rune{c.getCompatSvsUnichar().compatVs}
	case vcUnvoiced:
		// TEST_Jt3UaWwr knows that every unvoiced character has a
		// corresponding voiced character.
		return []rune{c.compatVs}
	case vcUndefined:
		switch c.charWidth {
		case cwNarrow:
			return []rune{c.codepoint, vsmNarrow}
		case cwWide:
			return []rune{c.codepoint, vsmWide}
		case cwUndefined:
			// These characters (U+3040, U+3097, U+3098, U+FF00) are not in the UCD.
			return []rune{c.codepoint}
		default:
			// TEST_U2mt8xTY knows that the program never passes here
			panic("unreachable")
		}
	default:
		// TEST_R8jrnbCz knows that the program never passes here
		panic("unreachable")
	}
}

// for Hiragana-Katakana letters.
// TEST_s8U59Hzf knows that this function returns a rune array with
// 1 or 2 elements and no other number of elements.
func (c *unichar) toTraditionalSemivoiced() []rune {
	switch c.voicing {
	case vcSemivoiced:
		return []rune{c.codepoint}
	case vcVoiced:
		unvoiced := c.getCompatVsUnichar()
		if unvoiced.existsCompatSvs() {
			return []rune{unvoiced.compatSvs}
		}
		switch c.charWidth {
		case cwNarrow:
			return []rune{c.compatVs, svsmNarrow}
		case cwWide:
			return []rune{c.compatVs, svsmWide}
		default:
			// TEST_T2eKd76G knows that the program never passes here
			panic("unreachable")
		}
	case vcUnvoiced:
		if c.existsCompatSvs() {
			return []rune{c.compatSvs}
		}
		switch c.charWidth {
		case cwNarrow:
			return []rune{c.codepoint, svsmNarrow}
		case cwWide:
			return []rune{c.codepoint, svsmWide}
		default:
			// TEST_Mw87qjkF knows that the program never passes here
			panic("unreachable")
		}
	case vcUndefined:
		switch c.charWidth {
		case cwNarrow:
			return []rune{c.codepoint, svsmNarrow}
		case cwWide:
			return []rune{c.codepoint, svsmWide}
		case cwUndefined:
			// These characters (U+3040, U+3097, U+3098, U+FF00) are not in the UCD.
			return []rune{c.codepoint}
		default:
			// TEST_U2mt8xTY knows that the program never passes here
			panic("unreachable")
		}
	default:
		// TEST_R8jrnbCz knows that the program never passes here
		panic("unreachable")
	}
}

// for Hiragana-Katakana letters.
// TEST_R4gNVpGj knows that this function returns a rune array with
// 1 or 2 elements and no other number of elements.
func (c *unichar) toCombiningVoiced() []rune {
	switch c.voicing {
	case vcUnvoiced, vcUndefined:
		return []rune{c.codepoint, vsmCombining}
	case vcVoiced:
		return []rune{c.compatVs, vsmCombining}
	case vcSemivoiced:
		return []rune{c.compatSvs, vsmCombining}
	default:
		// TEST_R8jrnbCz knows that the program never passes here
		panic("unreachable")
	}
}

// for Hiragana-Katakana letters.
// TEST_Pp9gBVj2 knows that this function returns a rune array with
// 1 or 2 elements and no other number of elements.
func (c *unichar) toCombiningSemivoiced() []rune {
	switch c.voicing {
	case vcUnvoiced, vcUndefined:
		return []rune{c.codepoint, svsmCombining}
	case vcVoiced:
		return []rune{c.compatVs, svsmCombining}
	case vcSemivoiced:
		return []rune{c.compatSvs, svsmCombining}
	default:
		// TEST_R8jrnbCz knows that the program never passes here
		panic("unreachable")
	}
}
