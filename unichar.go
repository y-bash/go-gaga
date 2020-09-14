package gaga

import (
	"fmt"
)

//go:generate go run gen/gentables.go -output unichar_tables.go

// A voicing modifier (voiced or semi-voiced sound mark)
type vom rune

const (
	vmNone        vom = 0
	vmVsmNonspace vom = 0x3099 // Combining voiced sound mark
	vmSsmNonspace vom = 0x309A // Combining semi-voiced sound mark
	vmVsmWide     vom = 0x309B // Wide voiced sournd mark
	vmSsmWide     vom = 0x309C // Wide semi-voiced sound mark
	vmVsmNarrow   vom = 0xFF9E // Narrow voiced sound mark
	vmSsmNarrow   vom = 0xFF9F // Narrow semi-voiced sound mark
)

func (m vom) isVsm() bool {
	switch m {
	case vmVsmNonspace, vmVsmWide, vmVsmNarrow:
		return true
	default:
		return false
	}
}

func (m vom) isSsm() bool {
	switch m {
	case vmSsmNonspace, vmSsmWide, vmSsmNarrow:
		return true
	default:
		return false
	}
}

func (m vom) isVom() bool {
	return m.isVsm() || m.isSsm()
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
		panic(fmt.Sprintf("Unexpectedly %#U was not found in the table", c.codepoint))
	}
	return c
}

func (c *unichar) getCompatCaseC() *unichar {
	// TEST_fm8XjZTB ensured that all compatCases are in the tables.
	return getUnicharForSure(c.compatCase)
}

func (c *unichar) getCompatWidthC() *unichar {
	// TEST_T3bc4Nh7 ensured that all compatWidth are in the table.
	return getUnicharForSure(c.compatWidth)
}

func (c *unichar) getCompatVoicedC() *unichar {
	// TEST_Cu8iKMxF ensured that all compatVoiced are in the tables.
	return getUnicharForSure(c.compatVoiced)
}

func (c *unichar) getCompatSemivoicedC() *unichar {
	// TEST_rW4UiNHC ensured that all compatSemivoiced are in the tables.
	return getUnicharForSure(c.compatSemivoiced)
}

func (c *unichar) existsCompatVoiced() bool {
	return c.codepoint != c.compatVoiced
}

func (c *unichar) existsCompatSemivoiced() bool {
	return c.codepoint != c.compatSemivoiced
}

func (c *unichar) toUpperR() rune {
	if c.charCase != ccLower {
		return c.codepoint
	}
	return c.compatCase
}

func (c *unichar) toLowerR() rune {
	if c.charCase != ccUpper {
		return c.codepoint
	}
	return c.compatCase
}

func (c *unichar) toHiraganaC() *unichar {
	if c.charCase != ccKatakana {
		return c
	}
	return c.getCompatCaseC()
}

func (c *unichar) toKatakanaC() *unichar {
	if c.charCase != ccHiragana {
		return c
	}
	return c.getCompatCaseC()
}

func (c *unichar) toWideR() rune {
	if c.charWidth != cwNarrow {
		return c.codepoint
	}
	return c.compatWidth
}

func (c *unichar) toWideC() *unichar {
	if c.charWidth != cwNarrow {
		return c
	}
	return c.getCompatWidthC()
}

func (c *unichar) toNarrowR() rune {
	if c.charWidth != cwWide {
		return c.codepoint
	}
	return c.compatWidth
}

func (c *unichar) toNarrowC() *unichar {
	if c.charWidth != cwWide {
		return c
	}
	return c.getCompatWidthC()
}

// for KanaVom
func (c *unichar) toLegacyC() *unichar {
	if c.charCase != ccCombining {
		return c
	}
	return c.getCompatCaseC()
}

// for KanaVom
func (c *unichar) toCombiningR() rune {
	if c.charCase != ccLegacy {
		return c.codepoint
	}
	return c.compatCase
}

// for Hiragana-Katakana letters.
// TEST_Vs4Ad89Z knows that this function returns a rune array with
// 1 or 2 elements and no other number of elements.
// TODO comment, func name
func (c *unichar) composeVoiced() (rune, vom) {
	switch c.voicing {
	case vcVoiced:
		return c.codepoint, vmNone
	case vcSemivoiced:
		// TEST_fW6auXUi knows that every semi-voiced character has
		// a corresponding unvoiced character, and that unvoiced
		// character has a corresponding voiced character.
		return c.getCompatSemivoicedC().compatVoiced, vmNone
	case vcUnvoiced:
		// TEST_Jt3UaWwr knows that every unvoiced character has a
		// corresponding voiced character.
		return c.compatVoiced, vmNone
	case vcUndefined:
		switch c.charWidth {
		case cwNarrow:
			return c.codepoint, vmVsmNarrow
		case cwWide:
			return c.codepoint, vmVsmWide
		case cwUndefined:
			// These characters (U+3040, U+3097, U+3098, U+FF00) are not in the UCD.
			return c.codepoint, vmNone
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
func (c *unichar) composeSemivoiced() (rune, vom) {
	switch c.voicing {
	case vcSemivoiced:
		return c.codepoint, vmNone
	case vcVoiced:
		unvoiced := c.getCompatVoicedC()
		if unvoiced.existsCompatSemivoiced() {
			return unvoiced.compatSemivoiced, vmNone
		}
		switch c.charWidth {
		case cwNarrow:
			return c.compatVoiced, vmSsmNarrow
		case cwWide:
			return c.compatVoiced, vmSsmWide
		default:
			// TEST_T2eKd76G knows that the program never passes here
			panic("unreachable")
		}
	case vcUnvoiced:
		if c.existsCompatSemivoiced() {
			return c.compatSemivoiced, vmNone
		}
		switch c.charWidth {
		case cwNarrow:
			return c.codepoint, vmSsmNarrow
		case cwWide:
			return c.codepoint, vmSsmWide
		default:
			// TEST_Mw87qjkF knows that the program never passes here
			panic("unreachable")
		}
	case vcUndefined:
		switch c.charWidth {
		case cwNarrow:
			return c.codepoint, vmSsmNarrow
		case cwWide:
			return c.codepoint, vmSsmWide
		case cwUndefined:
			// These characters (U+3040, U+3097, U+3098, U+FF00) are not in the UCD.
			return c.codepoint, vmNone
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
func (c *unichar) decomposeVoiced() (rune, vom) {
	switch c.voicing {
	case vcUnvoiced, vcUndefined:
		return c.codepoint, vmVsmNonspace
	case vcVoiced:
		return c.compatVoiced, vmVsmNonspace
	case vcSemivoiced:
		return c.compatSemivoiced, vmVsmNonspace
	default:
		// TEST_R8jrnbCz knows that the program never passes here
		panic("unreachable")
	}
}

// for Hiragana-Katakana letters.
// TEST_Pp9gBVj2 knows that this function returns a rune array with
// 1 or 2 elements and no other number of elements.
func (c *unichar) decomposeSemivoiced() (rune, vom) {
	switch c.voicing {
	case vcUnvoiced, vcUndefined:
		return c.codepoint, vmSsmNonspace
	case vcVoiced:
		return c.compatVoiced, vmSsmNonspace
	case vcSemivoiced:
		return c.compatSemivoiced, vmSsmNonspace
	default:
		// TEST_R8jrnbCz knows that the program never passes here
		panic("unreachable")
	}
}
