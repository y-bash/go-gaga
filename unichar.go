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
	case widthFirst <= r && r <= widthLast:
		return &widthTable[r-widthFirst], true
	default:
		return nil, false
	}
}

func neverBeCalled() []rune {
	// TODO Consider whether to include panic() in release version
	// return []rune{}
	panic("Unexpectedly called neverBeCalled()")
}

// This function can only be called if r is known to exist in the tables.
func getUnicharForSure(r rune) *unichar {
	c, ok := getUnichar(r)
	if !ok {
		panic(fmt.Sprintf("Unexpectedly %#U.cmptCase %#U was not found in the table",
			c.codepoint, c.cmptCase))
	}
	return c
}

func (c *unichar) getCmptCaseUnichar() *unichar {
	// TEST_fm8XjZTB ensured that all cmptCases are in the tables.
	return getUnicharForSure(c.cmptCase)
}

func (c *unichar) getCmptWidthUnichar() *unichar {
	// TEST_T3bc4Nh7 ensured that all cmptWidth are in the table.
	return getUnicharForSure(c.cmptWidth)
}

func (c *unichar) getCmptVsUnichar() *unichar {
	// TEST_Cu8iKMxF ensured that all cmptVs are in the tables.
	return getUnicharForSure(c.cmptVs)
}

func (c *unichar) getCmptSvsUnichar() *unichar {
	// TEST_rW4UiNHC ensured that all cmptSvs are in the tables.
	return getUnicharForSure(c.cmptSvs)
}

func (c *unichar) toUpper() rune {
	if c.charCase != ccLower {
		return c.codepoint
	}
	return c.cmptCase
}

func (c *unichar) toLower() rune {
	if c.charCase != ccUpper {
		return c.codepoint
	}
	return c.cmptCase
}

func (c *unichar) toHiraganaUnichar() *unichar {
	if c.charCase != ccKatakana {
		return c
	}
	return c.getCmptCaseUnichar()
}

func (c *unichar) toKatakanaUnichar() *unichar {
	if c.charCase != ccHiragana {
		return c
	}
	return c.getCmptCaseUnichar()
}

func (c *unichar) toWide() rune {
	if c.charWidth != cwNarrow {
		return c.codepoint
	}
	return c.cmptWidth
}

func (c *unichar) toWideUnichar() *unichar {
	if c.charWidth != cwNarrow {
		return c
	}
	return c.getCmptWidthUnichar()
}

func (c *unichar) toNarrow() rune {
	if c.charWidth != cwWide {
		return c.codepoint
	}
	return c.cmptWidth
}

func (c *unichar) toNarrowUnichar() *unichar {
	if c.charWidth != cwWide {
		return c
	}
	return c.getCmptWidthUnichar()
}

func (c *unichar) existsCmptVs() bool {
	return c.codepoint != c.cmptVs
}

func (c *unichar) existsCmptSvs() bool {
	return c.codepoint != c.cmptSvs
}

func (c *unichar) toClassicalVoiced() []rune {
	switch c.voicing {
	case vcVoiced:
		return []rune{c.codepoint}
	case vcSemivoiced:
		// TEST_fW6auXUi knows that every semi-voiced character has
		// a corresponding unvoiced character, and that unvoiced
		// character has a corresponding voiced character.
		return []rune{c.getCmptSvsUnichar().cmptVs}
	case vcUnvoiced:
		// TEST_Jt3UaWwr knows that every unvoiced character has a
		// corresponding voiced character.
		return []rune{c.cmptVs}
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
			return neverBeCalled()
		}
	default:
		// TEST_R8jrnbCz knows that the program never passes here
		return neverBeCalled()
	}
}

func (c *unichar) toClassicalSemivoiced() []rune {
	switch c.voicing {
	case vcSemivoiced:
		return []rune{c.codepoint}
	case vcVoiced:
		unvoiced := c.getCmptVsUnichar()
		if unvoiced.existsCmptSvs() {
			return []rune{unvoiced.cmptSvs}
		}
		switch c.charWidth {
		case cwNarrow:
			return []rune{c.cmptVs, svsmNarrow}
		case cwWide:
			return []rune{c.cmptVs, svsmWide}
		default:
			// TEST_T2eKd76G knows that the program never passes here
			return neverBeCalled()
		}
	case vcUnvoiced:
		if c.existsCmptSvs() {
			return []rune{c.cmptSvs}
		}
		switch c.charWidth {
		case cwNarrow:
			return []rune{c.codepoint, svsmNarrow}
		case cwWide:
			return []rune{c.codepoint, svsmWide}
		default:
			// TEST_Mw87qjkF knows that the program never passes here
			return neverBeCalled()
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
			return neverBeCalled()
		}
	default:
		// TEST_R8jrnbCz knows that the program never passes here
		return neverBeCalled()
	}
}

func (c *unichar) toCombiningVoiced() []rune {
	switch c.voicing {
	case vcUnvoiced, vcUndefined:
		return []rune{c.codepoint, vsmCombining}
	case vcVoiced:
		return []rune{c.cmptVs, vsmCombining}
	case vcSemivoiced:
		return []rune{c.cmptSvs, vsmCombining}
	default:
		// TEST_R8jrnbCz knows that the program never passes here
		return neverBeCalled()
	}
}

func (c *unichar) toCombiningSemivoiced() []rune {
	switch c.voicing {
	case vcUnvoiced, vcUndefined:
		return []rune{c.codepoint, svsmCombining}
	case vcVoiced:
		return []rune{c.cmptVs, svsmCombining}
	case vcSemivoiced:
		return []rune{c.cmptSvs, svsmCombining}
	default:
		// TEST_R8jrnbCz knows that the program never passes here
		return neverBeCalled()
	}
}
