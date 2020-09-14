package gaga

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"
)

func (c *unichar) String() string {
	return fmt.Sprintf("(%x,%x,%x,%x,%x,%x,%x,%x,%x)",
		c.codepoint, c.category, c.charCase, c.charWidth,
		c.voicing, c.compatCase, c.compatWidth, c.compatVoiced, c.compatSemivoiced)
}

type tableInfo struct {
	table   []unichar
	name    string
	wantN   int
	wantSHA string
}

var tables = []tableInfo{
	{latinTable, "latinTable", 96, "d90c9a10f72b6c029cb6aba58128f534d1935aa760fd5a073c232e16fc4eca22"},
	{kanaTable, "kanaTable", 256, "d7c8dd3e5efb23e2447498e580852728f5663b19ae146de5ad527262bcbdf138"},
	{kanaExtTable, "kanaExtTable", 16, "4c23e39933c8312f42c602a7b602041fed92634c8936cd57786c45ebe8df52b2"},
	{widthTable, "widthTable", 160, "79f5d0526d0696730a8e55526da769628e8280d100e954ea62151ff54699708d"},
}

func TestTableChecksums(t *testing.T) {
	for _, ti := range tables {
		gotN := len(ti.table)
		if gotN <= 0 {
			t.Errorf("table = %s, length is 0", ti.name)
			continue
		}
		var sb strings.Builder
		bufSize := len([]byte(ti.table[0].String())) * gotN * 2
		sb.Grow(bufSize)
		for _, c := range ti.table {
			sb.WriteString(c.String())
		}
		buf := []byte(sb.String())
		gotSHA := fmt.Sprintf("%x", sha256.Sum256(buf))
		if gotN != ti.wantN || gotSHA != ti.wantSHA {
			t.Errorf("table = %s,\n\tn\t\tsha256\n\tgot:  %d\t%s\n\twant: %d\t%s",
				ti.name, gotN, gotSHA, ti.wantN, ti.wantSHA)
		}
	}
}

func testTableSequence(t *testing.T, table []unichar, name string) {
	if len(table) <= 0 {
		t.Errorf("len(%s) is 0 or less", name)
		return
	}
	base := table[0].codepoint
	for i, c := range table {
		if c.codepoint != base+rune(i) {
			t.Errorf("%s[%d] %#U is non-sequential", name, i, c.codepoint)
		}
	}
}

func TestTableSequence(t *testing.T) {
	testTableSequence(t, latinTable, "latinTable")
	testTableSequence(t, kanaTable, "kanaTable")
	testTableSequence(t, kanaExtTable, "kanaExtTable")
	testTableSequence(t, widthTable, "widthTable")
}

func testUnicharTable(t *testing.T, table []unichar, first, last rune, name string) {
	for i := rune(0); i <= last-first; i++ {
		c := &table[i]

		// codepoint
		if _, ok := getUnichar(c.codepoint); !ok {
			t.Errorf("%s[%#U] is not found by getUnichar()", name, c.codepoint)
		}

		// category
		switch c.category {
		case ctUndefined, ctLatinLetter, ctLatinDigit, ctLatinSymbol,
			ctKanaLetter, ctKanaSymbol, ctKanaVom:
		default: // TEST_P8w4qtsm
			t.Errorf("%s[%#U].category == %d, want %d <= category < %d",
				name, c.codepoint, c.category, ctUndefined, ctMax)
		}

		// character case
		if c.charCase < ccUndefined || c.charCase >= ccMax {
			t.Errorf("%s[%#U].charCase == %d, want %d <= charCase < %d",
				name, c.codepoint, c.charCase, ccUndefined, ccMax)
		}
		if c.charCase == ccUndefined && c.compatCase != c.codepoint {
			t.Errorf("%s[%#U].charCase == %d, but compatCase == %#U, want compatCase == %#U",
				name, c.codepoint, c.charCase, c.compatCase, c.codepoint)
		}
		if c.compatCase != c.codepoint {
			if compatCase, ok := getUnichar(c.compatCase); !ok { // TEST_fm8XjZTB
				t.Errorf("%s[%#U].compatCase %#U is not found by getUnichar()",
					name, c.codepoint, c.compatCase)
			} else {
				if compatCase.category != c.category {
					t.Errorf("%s[%#U].category == %d, but compatCase %#U.category == %d, want same value",
						name, c.codepoint, c.category, compatCase.codepoint, compatCase.category)
				}
				if compatCase.charCase == c.charCase {
					t.Errorf("%s[%#U].charCase == %d, but compatCase %#U.charCase == %d, want another one",
						name, c.codepoint, c.charCase, compatCase.codepoint, compatCase.charCase)
				}
			}
		}
		if c.charCase != ccUndefined && c.compatCase == c.codepoint {
			switch c.codepoint {
			case 'ゟ', 'ー', 'ヿ', 'ｰ':
			default:
				t.Errorf("%s[%#U].charCase is %d and compatCase is %#U, want another one",
					name, c.codepoint, c.charCase, c.compatCase)
			}
		}

		// character width
		switch c.charWidth {
		case cwUndefined, cwNarrow, cwWide:
		default: // TEST_U2mt8xTY
			t.Errorf("%s[%#U].charWidth == %d, want %d <= charWidth < %d",
				name, c.codepoint, c.charWidth, cwUndefined, cwMax)
		}
		if c.charWidth == cwUndefined {
			if c.category != ctUndefined {
				t.Errorf("%s[%#U].charWidth == %d, want charWidth != %d",
					name, c.codepoint, c.charWidth, cwUndefined)
			}
			if c.compatWidth != c.codepoint {
				t.Errorf("%s[%#U].charWidth == %d, but compatWidth == %#U, want compatWidth == %#U",
					name, c.codepoint, c.charWidth, c.compatWidth, c.codepoint)
			}
		}
		if c.compatWidth != c.codepoint {
			if compatWidth, ok := getUnichar(c.compatWidth); !ok { // TEST_T3bc4Nh7
				t.Errorf("%s[%#U].compatWidth %#U is not found by getUnichar()",
					name, c.codepoint, c.compatWidth)
			} else {
				if compatWidth.category != c.category {
					t.Errorf("%s[%#U].category == %d, but compatWidth %#U.category == %d, want same value",
						name, c.codepoint, c.category, compatWidth.codepoint, compatWidth.category)
				}
				if compatWidth.charWidth == c.charWidth {
					t.Errorf("%s[%#U].charWidth is %d, but compatWidth %#U.charWidth == %d, want another one",
						name, c.codepoint, c.charWidth, compatWidth.codepoint, compatWidth.charWidth)
				}
			}
		}

		// voicing
		switch c.voicing {
		case vcUndefined, vcUnvoiced, vcVoiced, vcSemivoiced:
		default: // TEST_R8jrnbCz
			t.Errorf("%s[%#U].voicing == %d, want %d <= voicing < %d",
				name, c.codepoint, c.voicing, vcUndefined, vcMax)
		}
		if c.voicing == vcUndefined {
			if c.compatVoiced != c.codepoint {
				t.Errorf("%s[%#U].voicing == %d, but compatVoiced == %#U, want compatVoiced == %#U",
					name, c.codepoint, c.voicing, c.compatVoiced, c.codepoint)
			}
			if c.compatSemivoiced != c.codepoint {
				t.Errorf("%s[%#U].voicing == %d, but compatSemivoiced == %#U, want compatSemivoiced == %#U",
					name, c.codepoint, c.voicing, c.compatSemivoiced, c.codepoint)
			}
		}
		if c.voicing == vcUnvoiced {
			if c.compatVoiced == c.codepoint { // TEST_Jt3UaWwr
				t.Errorf("%s[%#U].voicing == %d, but compatVoiced == %#U, want another one",
					name, c.codepoint, c.voicing, c.compatVoiced)
			}
			if c.charWidth != cwNarrow && c.charWidth != cwWide { // TEST_Mw87qjkF
				t.Errorf("%s[%#U].voicing == %d, but charWidth == %d, want %d or %d",
					name, c.codepoint, c.voicing, c.charWidth, cwNarrow, cwWide)
			}
		}
		if c.voicing == vcVoiced {
			if c.compatVoiced == c.codepoint {
				t.Errorf("%s[%#U].voicing == %d, but compatVoiced == %#U, want another one",
					name, c.codepoint, c.voicing, c.compatVoiced)
			}
			if c.charWidth != cwNarrow && c.charWidth != cwWide { // TEST_T2eKd76G
				t.Errorf("%s[%#U].voicing == %d, but charWidth == %d, want %d or %d",
					name, c.codepoint, c.voicing, c.charWidth, cwNarrow, cwWide)
			}
		}
		if c.voicing == vcSemivoiced {
			if c.compatSemivoiced == c.codepoint {
				t.Errorf("%s[%#U].voicing is %d and compatSemivoiced is %#U, want another one",
					name, c.codepoint, c.voicing, c.compatSemivoiced)
			}
			if c.charWidth != cwNarrow && c.charWidth != cwWide {
				t.Errorf("%s[%#U].voicing == %d, but charWidth == %d, want %d or %d",
					name, c.codepoint, c.voicing, c.charWidth, cwNarrow, cwWide)
			}
			unvoiced := c.getCompatSemivoicedC()
			if !unvoiced.existsCompatVoiced() { // TEST_fW6auXUi
				t.Errorf("%s[%#U].voicing == %d, but getCompatSemivoicedC().existsCompatVoiced() == false, want true",
					name, c.codepoint, c.voicing)
			}
		}

		if c.compatVoiced != c.codepoint {
			if compatVoiced, ok := getUnichar(c.compatVoiced); !ok { // TEST_Cu8iKMxF
				t.Errorf("%s[%#U].compatVoiced %#U is not found by getUnichar()", name, c.codepoint, c.compatVoiced)
			} else {
				if c.category != compatVoiced.category {
					t.Errorf("%s[%#U].category is %d and compatVoiced %#U.category is %d, want same value",
						name, c.codepoint, c.category, compatVoiced.codepoint, compatVoiced.category)
				}
				if c.charWidth != compatVoiced.charWidth {
					t.Errorf("%s[%#U].charWidth is %d and compatVoiced %#U.charWidth is %d, want same value",
						name, c.codepoint, c.charWidth, compatVoiced.codepoint, compatVoiced.charWidth)
				}
				if c.voicing == vcUnvoiced && compatVoiced.voicing != vcVoiced {
					t.Errorf("%s[%#U].voicing is %d and compatVoiced %#U.voicing is %d, want %d",
						name, c.codepoint, c.voicing, compatVoiced.codepoint, compatVoiced.voicing, vcVoiced)
				}
				if c.voicing == vcVoiced && compatVoiced.voicing != vcUnvoiced {
					t.Errorf("%s[%#U].voicing is %d and compatVoiced %#U.voicing is %d, want %d",
						name, c.codepoint, c.voicing, compatVoiced.codepoint, compatVoiced.voicing, vcUnvoiced)
				}
			}
		}
		if c.compatSemivoiced != c.codepoint {
			if compatSemivoiced, ok := getUnichar(c.compatSemivoiced); !ok { // TEST_rW4UiNHC
				t.Errorf("%s[%#U].compatSemivoiced %#U is not found by getUnichar()", name, c.codepoint, c.compatSemivoiced)
			} else {
				if c.category != compatSemivoiced.category {
					t.Errorf("%s[%#U].category is %d and compatSemivoiced %#U.category is %d, want same value",
						name, c.codepoint, c.category, compatSemivoiced.codepoint, compatSemivoiced.category)
				}
				if c.charWidth != compatSemivoiced.charWidth {
					t.Errorf("%s[%#U].charWidth is %d and compatSemivoiced %#U.charWidth is %d, want same value",
						name, c.codepoint, c.charWidth, compatSemivoiced.codepoint, compatSemivoiced.charWidth)
				}
				if c.voicing == vcUnvoiced && compatSemivoiced.voicing != vcSemivoiced {
					t.Errorf("%s[%#U].voicing is %d and compatSemivoiced %#U.voicing is %d, want %d",
						name, c.codepoint, c.voicing, compatSemivoiced.codepoint, compatSemivoiced.voicing, vcSemivoiced)
				}
				if c.voicing == vcSemivoiced && compatSemivoiced.voicing != vcUnvoiced {
					t.Errorf("%s[%#U].voicing is %d and compatSemivoiced %#U.voicing is %d, want %d",
						name, c.codepoint, c.voicing, compatSemivoiced.codepoint, compatSemivoiced.voicing, vcUnvoiced)
				}
			}
		}

		switch c.category {
		// Latin
		case ctLatinLetter, ctLatinDigit, ctLatinSymbol:
			switch c.charCase {
			case ccUndefined, ccUpper, ccLower:
			default:
				t.Errorf("%s[%#U].charCase is %d, want %d or %d or %d",
					name, c.codepoint, c.charCase, ccUndefined, ccUpper, ccLower)
			}
			if c.charWidth != cwNarrow && c.charWidth != cwWide {
				t.Errorf("%s[%#U].charWidth is %d, want 1 or 2", name, c.codepoint, c.charWidth)
			}
			if c.voicing != vcUndefined {
				t.Errorf("%s[%#U].voicing is %d, want 0", name, c.codepoint, c.voicing)
			}
		// Kana
		case ctKanaLetter:
			if c.charCase != ccHiragana && c.charCase != ccKatakana { // TEST_gT8YJdBc
				t.Errorf("%s[%#U].charCase is %d, want 3 or 4", name, c.codepoint, c.charCase)
			}
			if c.charWidth != cwNarrow && c.charWidth != cwWide {
				t.Errorf("%s[%#U].charWidth is %d, want 1 or 2", name, c.codepoint, c.charWidth)
			}
		case ctKanaSymbol:
			if c.charCase != ccUndefined {
				t.Errorf("%s[%#U].charCase is %d, want 0", name, c.codepoint, c.charCase)
			}
			if c.charWidth != cwNarrow && c.charWidth != cwWide {
				t.Errorf("%s[%#U].charWidth is %d, want 1 or 2", name, c.codepoint, c.charWidth)
			}
			if c.voicing != vcUndefined {
				t.Errorf("%s[%#U].voicing is %d, want 0", name, c.codepoint, c.voicing)
			}
		case ctKanaVom:
			if !vom(c.codepoint).isVom() {
				t.Errorf("%s[%#U] is not Vom, want Vom", name, c.codepoint)
			}
			if c.charCase != ccLegacy && c.charCase != ccCombining {
				t.Errorf("%s[%#U].charCase is %d, want 3 or 4", name, c.codepoint, c.charCase)
			}
			if c.charWidth != cwNarrow && c.charWidth != cwWide {
				t.Errorf("%s[%#U].charWidth is %d, want 1 or 2", name, c.codepoint, c.charWidth)
			}
			if c.voicing != vcUndefined {
				t.Errorf("%s[%#U].voicing is %d, want 0", name, c.codepoint, c.voicing)
			}
		// Undefined
		case ctUndefined:
			if c.charCase != ccUndefined {
				t.Errorf("%s[%#U].charCase is %d, want 0", name, c.codepoint, c.charCase)
			}
			if c.charWidth != cwUndefined {
				t.Errorf("%s[%#U].charWidth is %d, want 0", name, c.codepoint, c.charWidth)
			}
			if c.voicing != vcUndefined {
				t.Errorf("%s[%#U].voicing is %d, want 0", name, c.codepoint, c.voicing)
			}
		}
	}
}

func TestUnicharTable(t *testing.T) {
	testUnicharTable(t, latinTable, latinFirst, latinLast, "latinTable")
	testUnicharTable(t, kanaTable, kanaFirst, kanaLast, "kanaTable")
	testUnicharTable(t, widthTable, widthFirst, widthLast, "widthTable")
}

type ToVoicedTest struct {
	in          rune
	outTradChar rune
	outTradMark vom
	outNonsChar rune
	outNonsMark vom
}

var tovoicedtests = []ToVoicedTest{
	// vcUnvoiced
	0: {'か', 'が', vmNone, 'か', '\u3099'}, // compatVoiced is exists, compatSemivoiced is not exists
	1: {'は', 'ば', vmNone, 'は', '\u3099'}, // compatVoiced is exists, compatSemivoiced is exists
	// vcVoiced
	2: {'が', 'が', vmNone, 'か', '\u3099'}, // compatVoiced.compatSemivoiced is not exists
	3: {'ば', 'ば', vmNone, 'は', '\u3099'}, // compatVoiced.compatSemivoiced is exists
	// vcSemivoiced
	4: {'ぱ', 'ば', vmNone, 'は', '\u3099'},
	// vcUndefined, cwWide
	5:  {'あ', 'あ', '゛', 'あ', '\u3099'}, // ctKanaLetter, ccHiragana
	6:  {'ア', 'ア', '゛', 'ア', '\u3099'}, // ctKanaLetter, ccKatakana
	7:  {'・', '・', '゛', '・', '\u3099'}, // ctKanaSymbol, ccUndefined
	8:  {'Ａ', 'Ａ', '゛', 'Ａ', '\u3099'}, // ctLatinLetter, ccUpper
	9:  {'ａ', 'ａ', '゛', 'ａ', '\u3099'}, // ctLatinLetter, ccLower
	10: {'１', '１', '゛', '１', '\u3099'}, // ctLatinDigit, ccUndefined
	11: {'＃', '＃', '゛', '＃', '\u3099'}, // ctLatinSymbol, ccUndefined
	// vcUndefined, cwNarrow
	12: {'ｱ', 'ｱ', 'ﾞ', 'ｱ', '\u3099'}, // ctKanaLetter, ccKatakana
	13: {'ｶ', 'ｶ', 'ﾞ', 'ｶ', '\u3099'}, // ctKanaLetter, ccKatakana, compatWidth.compatVoiced is exists, compatWidth.compatSemivoiced is not exists
	14: {'ﾊ', 'ﾊ', 'ﾞ', 'ﾊ', '\u3099'}, // ctKanaLetter, ccKatakana, compatWidth.compatVoiced is exists, compatWidth.compatSemivoiced is exists
	15: {'･', '･', 'ﾞ', '･', '\u3099'}, // ctKanaSymbol, ccUndefined
	16: {'A', 'A', 'ﾞ', 'A', '\u3099'}, // ctLatinLetter, ccUpper
	17: {'a', 'a', 'ﾞ', 'a', '\u3099'}, // ctLatinLetter, ccLower
	18: {'1', '1', 'ﾞ', '1', '\u3099'}, // ctLatinDigit, ccUndefined
	19: {'#', '#', 'ﾞ', '#', '\u3099'}, // ctLatinDigit, ccUndefined
	// ctUndefined
	20: {'\u3040', '\u3040', vmNone, '\u3040', '\u3099'},
	// VSM
	21: {'゛', '゛', '゛', '゛', '\u3099'},
	22: {'\u3099', '\u3099', '゛', '\u3099', '\u3099'},
	23: {'ﾞ', 'ﾞ', 'ﾞ', 'ﾞ', '\u3099'},
	// SVSM
	24: {'゜', '゜', '゛', '゜', '\u3099'},
	25: {'\u309A', '\u309A', '゛', '\u309A', '\u3099'},
	26: {'ﾟ', 'ﾟ', 'ﾞ', 'ﾟ', '\u3099'},
}

func TestToVoiced(t *testing.T) {
	for n, tt := range tovoicedtests {
		c, ok := getUnichar(tt.in)
		if !ok {
			t.Errorf("%d: %#U is not found by getUnichar()", n, tt.in)
			continue
		}
		var have, want rune
		var haveMm, wantMm vom

		have, haveMm = c.composeVoiced()
		want = tt.outTradChar
		wantMm = tt.outTradMark
		if have != want || haveMm != wantMm {
			t.Errorf("%d: have: (%q, %q), want: (%q, %q)",
				n, have, haveMm, tt.outTradChar, tt.outTradMark)
			break
		}

		have, haveMm = c.decomposeVoiced()
		want = tt.outNonsChar
		wantMm = tt.outNonsMark
		if have != want || haveMm != wantMm {
			t.Errorf("%d: have: (%q, %q), want: (%q, %q)",
				n, have, haveMm, tt.outNonsChar, tt.outNonsMark)
			break
		}

	}
}

type ToSemivoicedTest struct {
	in          rune
	outTradChar rune
	outTradMark vom
	outNonsChar rune
	outNonsMark vom
}

var tosemivoicedtests = []ToSemivoicedTest{
	// vcUnvoiced
	0: {'か', 'か', '゜', 'か', '\u309A'},    // compatVoiced is exists, compatSemivoiced is not exists
	1: {'は', 'ぱ', vmNone, 'は', '\u309A'}, // compatVoiced is exists, compatSemivoiced is exists
	// vcVoice'
	2: {'が', 'か', '゜', 'か', '\u309A'},    // compatVoiced.compatSemivoiced is not exists
	3: {'ば', 'ぱ', vmNone, 'は', '\u309A'}, // compatVoiced.compatSemivoiced is exists
	// vcSemivoiced
	4: {'ぱ', 'ぱ', vmNone, 'は', '\u309A'},
	// vcUndefined, cwWide
	5:  {'あ', 'あ', '゜', 'あ', '\u309A'}, // ctKanaLetter, ccHiragana
	6:  {'ア', 'ア', '゜', 'ア', '\u309A'}, // ctKanaLetter, ccKatakana
	7:  {'・', '・', '゜', '・', '\u309A'}, // ctKanaSymbol, ccUndefined
	8:  {'Ａ', 'Ａ', '゜', 'Ａ', '\u309A'}, // ctLatinLetter, ccUpper
	9:  {'ａ', 'ａ', '゜', 'ａ', '\u309A'}, // ctLatinLetter, ccLower
	10: {'１', '１', '゜', '１', '\u309A'}, // ctLatinDigit, ccUndefined
	11: {'＃', '＃', '゜', '＃', '\u309A'}, // ctLatinSymbol, ccUndefined
	// vcUndefined, cwNarrow
	12: {'ｱ', 'ｱ', 'ﾟ', 'ｱ', '\u309A'}, // ctKanaLetter, ccKatakana
	13: {'ｶ', 'ｶ', 'ﾟ', 'ｶ', '\u309A'}, // ctKanaLetter, ccKatakana, compatWidth.compatVoiced is exists, compatWidth.compatSemivoiced is not exists
	14: {'ﾊ', 'ﾊ', 'ﾟ', 'ﾊ', '\u309A'}, // ctKanaLetter, ccKatakana, compatWidth.compatVoiced is exists, compatWidth.compatSemivoiced is exists
	15: {'･', '･', 'ﾟ', '･', '\u309A'}, // ctKanaSymbol, ccUndefined
	16: {'A', 'A', 'ﾟ', 'A', '\u309A'}, // ctLatinLetter, ccUpper
	17: {'a', 'a', 'ﾟ', 'a', '\u309A'}, // ctLatinLetter, ccLower
	18: {'1', '1', 'ﾟ', '1', '\u309A'}, // ctLatinDigit, ccUndefined
	19: {'#', '#', 'ﾟ', '#', '\u309A'}, // ctLatinDigit, ccUndefined
	// ctUndefined
	20: {'\u3040', '\u3040', vmNone, '\u3040', '\u309A'},
	// VSM
	21: {'゛', '゛', '゜', '゛', '\u309A'},
	22: {'\u3099', '\u3099', '゜', '\u3099', '\u309A'},
	23: {'ﾞ', 'ﾞ', 'ﾟ', 'ﾞ', '\u309A'},
	// SVSM
	24: {'゜', '゜', '゜', '゜', '\u309A'},
	25: {'\u309A', '\u309A', '゜', '\u309A', '\u309A'},
	26: {'ﾟ', 'ﾟ', 'ﾟ', 'ﾟ', '\u309A'},
}

func TestToSemivoiced(t *testing.T) {
	for n, tt := range tosemivoicedtests {
		c, ok := getUnichar(tt.in)
		if !ok {
			t.Errorf("%d: %#U is not found by getUnichar()", n, tt.in)
			continue
		}
		var have, want rune
		var haveMm, wantMm vom

		have, haveMm = c.composeSemivoiced()
		want = tt.outTradChar
		wantMm = tt.outTradMark
		if have != want || haveMm != wantMm {
			t.Errorf("%d: have: (%q, %q), want: (%q, %q)",
				n, have, haveMm, tt.outTradChar, tt.outTradMark)
			break
		}

		have, haveMm = c.decomposeSemivoiced()
		want = tt.outNonsChar
		wantMm = tt.outNonsMark
		if have != want || haveMm != wantMm {
			t.Errorf("%d: have: (%q, %q), want: (%q, %q)",
				n, have, haveMm, tt.outNonsChar, tt.outNonsMark)
			break
		}

	}
}

func BenchmarkGetUnichar(b *testing.B) {
	s := "\t Aa#　Ａａ＃あア。ｱ｡”ﾞ漢字\U0010FFFF"
	rs := []rune(s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, r := range rs {
			getUnichar(r)
		}
	}
	b.StopTimer()
}
