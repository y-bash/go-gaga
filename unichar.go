package gaga

import (
	"fmt"
)

//go:generate go run gen/gentables.go -output unichar_tables.go

// Modifier mark (Voiced or Semi-voiced sound mark)
type modmark rune

const (
	mmNone         modmark = 0
	mmVsmNonspace  modmark = 0x3099
	mmVsmWide      modmark = 0x309B
	mmVsmNarrow    modmark = 0xFF9E
	mmSvsmNonspace modmark = 0x309A
	mmSvsmWide     modmark = 0x309C
	mmSvsmNarrow   modmark = 0xFF9F
)

func (m modmark) isModmark() bool {
	return m != mmNone
}

func isVoicedSoundMark(r rune) bool {
	switch modmark(r) {
	case mmVsmNonspace, mmVsmWide, mmVsmNarrow:
		return true
	default:
		return false
	}
}

func isSemivoicedSoundMark(r rune) bool {
	switch modmark(r) {
	case mmSvsmNonspace, mmSvsmWide, mmSvsmNarrow:
		return true
	default:
		return false
	}
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
	if c.charCase != ccNonspace {
		return c
	}
	return c.getCompatCaseUnichar()
}

// for voiced or semi-voiced sound mark characters.
func (c *unichar) toNonspaceMark() rune {
	if c.charCase != ccTraditional {
		return c.codepoint
	}
	return c.compatCase
}

// for Hiragana-Katakana letters.
// TEST_Vs4Ad89Z knows that this function returns a rune array with
// 1 or 2 elements and no other number of elements.
func (c *unichar) toTraditionalVoiced() (rune, modmark) {
	switch c.voicing {
	case vcVoiced:
		return c.codepoint, mmNone
	case vcSemivoiced:
		// TEST_fW6auXUi knows that every semi-voiced character has
		// a corresponding unvoiced character, and that unvoiced
		// character has a corresponding voiced character.
		return c.getCompatSvsUnichar().compatVs, mmNone
	case vcUnvoiced:
		// TEST_Jt3UaWwr knows that every unvoiced character has a
		// corresponding voiced character.
		return c.compatVs, mmNone
	case vcUndefined:
		switch c.charWidth {
		case cwNarrow:
			return c.codepoint, mmVsmNarrow
		case cwWide:
			return c.codepoint, mmVsmWide
		case cwUndefined:
			// These characters (U+3040, U+3097, U+3098, U+FF00) are not in the UCD.
			return c.codepoint, mmNone
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
func (c *unichar) toTraditionalSemivoiced() (rune, modmark) {
	switch c.voicing {
	case vcSemivoiced:
		return c.codepoint, mmNone
	case vcVoiced:
		unvoiced := c.getCompatVsUnichar()
		if unvoiced.existsCompatSvs() {
			return unvoiced.compatSvs, mmNone
		}
		switch c.charWidth {
		case cwNarrow:
			return c.compatVs, mmSvsmNarrow
		case cwWide:
			return c.compatVs, mmSvsmWide
		default:
			// TEST_T2eKd76G knows that the program never passes here
			panic("unreachable")
		}
	case vcUnvoiced:
		if c.existsCompatSvs() {
			return c.compatSvs, mmNone
		}
		switch c.charWidth {
		case cwNarrow:
			return c.codepoint, mmSvsmNarrow
		case cwWide:
			return c.codepoint, mmSvsmWide
		default:
			// TEST_Mw87qjkF knows that the program never passes here
			panic("unreachable")
		}
	case vcUndefined:
		switch c.charWidth {
		case cwNarrow:
			return c.codepoint, mmSvsmNarrow
		case cwWide:
			return c.codepoint, mmSvsmWide
		case cwUndefined:
			// These characters (U+3040, U+3097, U+3098, U+FF00) are not in the UCD.
			return c.codepoint, mmNone
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
func (c *unichar) toNonspaceVoiced() (rune, modmark) {
	switch c.voicing {
	case vcUnvoiced, vcUndefined:
		return c.codepoint, mmVsmNonspace
	case vcVoiced:
		return c.compatVs, mmVsmNonspace
	case vcSemivoiced:
		return c.compatSvs, mmVsmNonspace
	default:
		// TEST_R8jrnbCz knows that the program never passes here
		panic("unreachable")
	}
}

// for Hiragana-Katakana letters.
// TEST_Pp9gBVj2 knows that this function returns a rune array with
// 1 or 2 elements and no other number of elements.
func (c *unichar) toNonspaceSemivoiced() (rune, modmark) {
	switch c.voicing {
	case vcUnvoiced, vcUndefined:
		return c.codepoint, mmSvsmNonspace
	case vcVoiced:
		return c.compatVs, mmSvsmNonspace
	case vcSemivoiced:
		return c.compatSvs, mmSvsmNonspace
	default:
		// TEST_R8jrnbCz knows that the program never passes here
		panic("unreachable")
	}
}
