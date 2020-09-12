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
		c.voicing, c.compatCase, c.compatWidth, c.compatVs, c.compatSvs)
}

type tableInfo struct {
	table   []unichar
	name    string
	wantN   int
	wantSHA string
}

var tables = []tableInfo{
	{latinTable, "latinTable", 96, "d90c9a10f72b6c029cb6aba58128f534d1935aa760fd5a073c232e16fc4eca22"},
	{kanaTable, "kanaTable", 256, "9f30b8ed44761d8667e55a08181f4f9db584e9694db352255291db0102014e54"},
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
			ctKanaLetter, ctKanaSymbol, ctKanaVsm:
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
		// TODO Check for problems in the following cases
		/*
			if c.charCase != ccUndefined && c.compatCase == c.codepoint {
				t.Errorf("%s[%#U].charCase is %d and compatCase is %#U, want another one",
					name, c.codepoint, c.charCase, c.compatCase)
			}
		*/

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
		// TODO Check for problems in the following cases
		/*
			if c.charWidth != cwUndefined && c.compatWidth == c.codepoint {
				t.Errorf("%s[%#U].charWidth is %d and compatWidth is %#U, want another one",
					name, c.codepoint, c.charWidth, c.compatWidth)
			}
		*/

		// voicing
		switch c.voicing {
		case vcUndefined, vcUnvoiced, vcVoiced, vcSemivoiced:
		default: // TEST_R8jrnbCz
			t.Errorf("%s[%#U].voicing == %d, want %d <= voicing < %d",
				name, c.codepoint, c.voicing, vcUndefined, vcMax)
		}
		if c.voicing == vcUndefined {
			if c.compatVs != c.codepoint {
				t.Errorf("%s[%#U].voicing == %d, but compatVs == %#U, want compatVs == %#U",
					name, c.codepoint, c.voicing, c.compatVs, c.codepoint)
			}
			if c.compatSvs != c.codepoint {
				t.Errorf("%s[%#U].voicing == %d, but compatSvs == %#U, want compatSvs == %#U",
					name, c.codepoint, c.voicing, c.compatSvs, c.codepoint)
			}
		}
		if c.voicing == vcUnvoiced {
			if c.compatVs == c.codepoint { // TEST_Jt3UaWwr
				t.Errorf("%s[%#U].voicing == %d, but compatVs == %#U, want another one",
					name, c.codepoint, c.voicing, c.compatVs)
			}
			if c.charWidth != cwNarrow && c.charWidth != cwWide { // TEST_Mw87qjkF
				t.Errorf("%s[%#U].voicing == %d, but charWidth == %d, want %d or %d",
					name, c.codepoint, c.voicing, c.charWidth, cwNarrow, cwWide)
			}
		}
		if c.voicing == vcVoiced {
			if c.compatVs == c.codepoint {
				t.Errorf("%s[%#U].voicing == %d, but compatVs == %#U, want another one",
					name, c.codepoint, c.voicing, c.compatVs)
			}
			if c.charWidth != cwNarrow && c.charWidth != cwWide { // TEST_T2eKd76G
				t.Errorf("%s[%#U].voicing == %d, but charWidth == %d, want %d or %d",
					name, c.codepoint, c.voicing, c.charWidth, cwNarrow, cwWide)
			}
		}
		if c.voicing == vcSemivoiced {
			if c.compatSvs == c.codepoint {
				t.Errorf("%s[%#U].voicing is %d and compatSvs is %#U, want another one",
					name, c.codepoint, c.voicing, c.compatSvs)
			}
			if c.charWidth != cwNarrow && c.charWidth != cwWide {
				t.Errorf("%s[%#U].voicing == %d, but charWidth == %d, want %d or %d",
					name, c.codepoint, c.voicing, c.charWidth, cwNarrow, cwWide)
			}
			unvoiced := c.getCompatSvsUnichar()
			if !unvoiced.existsCompatVs() { // TEST_fW6auXUi
				t.Errorf("%s[%#U].voicing == %d, but getCompatSvsUnichar().existsCompatVs() == false, want true",
					name, c.codepoint, c.voicing)
			}
		}

		if c.compatVs != c.codepoint {
			if compatVs, ok := getUnichar(c.compatVs); !ok { // TEST_Cu8iKMxF
				t.Errorf("%s[%#U].compatVs %#U is not found by getUnichar()", name, c.codepoint, c.compatVs)
			} else {
				if c.category != compatVs.category {
					t.Errorf("%s[%#U].category is %d and compatVs %#U.category is %d, want same value",
						name, c.codepoint, c.category, compatVs.codepoint, compatVs.category)
				}
				if c.charWidth != compatVs.charWidth {
					t.Errorf("%s[%#U].charWidth is %d and compatVs %#U.charWidth is %d, want same value",
						name, c.codepoint, c.charWidth, compatVs.codepoint, compatVs.charWidth)
				}
				if c.voicing == vcUnvoiced && compatVs.voicing != vcVoiced {
					t.Errorf("%s[%#U].voicing is %d and compatVs %#U.voicing is %d, want %d",
						name, c.codepoint, c.voicing, compatVs.codepoint, compatVs.voicing, vcVoiced)
				}
				if c.voicing == vcVoiced && compatVs.voicing != vcUnvoiced {
					t.Errorf("%s[%#U].voicing is %d and compatVs %#U.voicing is %d, want %d",
						name, c.codepoint, c.voicing, compatVs.codepoint, compatVs.voicing, vcUnvoiced)
				}
			}
		}
		if c.compatSvs != c.codepoint {
			if compatSvs, ok := getUnichar(c.compatSvs); !ok { // TEST_rW4UiNHC
				t.Errorf("%s[%#U].compatSvs %#U is not found by getUnichar()", name, c.codepoint, c.compatSvs)
			} else {
				if c.category != compatSvs.category {
					t.Errorf("%s[%#U].category is %d and compatSvs %#U.category is %d, want same value",
						name, c.codepoint, c.category, compatSvs.codepoint, compatSvs.category)
				}
				if c.charWidth != compatSvs.charWidth {
					t.Errorf("%s[%#U].charWidth is %d and compatSvs %#U.charWidth is %d, want same value",
						name, c.codepoint, c.charWidth, compatSvs.codepoint, compatSvs.charWidth)
				}
				if c.voicing == vcUnvoiced && compatSvs.voicing != vcSemivoiced {
					t.Errorf("%s[%#U].voicing is %d and compatSvs %#U.voicing is %d, want %d",
						name, c.codepoint, c.voicing, compatSvs.codepoint, compatSvs.voicing, vcSemivoiced)
				}
				if c.voicing == vcSemivoiced && compatSvs.voicing != vcUnvoiced {
					t.Errorf("%s[%#U].voicing is %d and compatSvs %#U.voicing is %d, want %d",
						name, c.codepoint, c.voicing, compatSvs.codepoint, compatSvs.voicing, vcUnvoiced)
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
		case ctKanaVsm:
			if !isVoicedSoundMark(c.codepoint) && !isSemivoicedSoundMark(c.codepoint) {
				t.Errorf("%s[%#U] is not VSM or SVSM, want VSM or SVSM", name, c.codepoint)
			}
			if c.charCase != ccTraditional && c.charCase != ccCombining {
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

		// Testing return values of methods
		var rs []rune

		// TEST_Vs4Ad89Z
		rs = c.toTraditionalVoiced()
		switch len(rs) {
		case 1, 2:
		default:
			t.Errorf("%s[%#U].toTraditionalVoiced() is %v, want 1 or 2 elements", name, c.codepoint, rs)
		}

		// TEST_s8U59Hzf
		rs = c.toTraditionalSemivoiced()
		switch len(rs) {
		case 1, 2:
		default:
			t.Errorf("%s[%#U].toTraditionalSemivoiced() is %v, want 1 or 2 elements", name, c.codepoint, rs)
		}

		// TEST_R4gNVpGj
		rs = c.toCombiningVoiced()
		switch len(rs) {
		case 1, 2:
		default:
			t.Errorf("%s[%#U].toCombiningVoiced() is %v, want 1 or 2 elements", name, c.codepoint, rs)
		}

		// TEST_Pp9gBVj2
		rs = c.toCombiningSemivoiced()
		switch len(rs) {
		case 1, 2:
		default:
			t.Errorf("%s[%#U].toCombiningSemivoiced() is %v, want 1 or 2 elements", name, c.codepoint, rs)
		}
	}
}

func TestUnicharTable(t *testing.T) {
	testUnicharTable(t, latinTable, latinFirst, latinLast, "latinTable")
	testUnicharTable(t, kanaTable, kanaFirst, kanaLast, "kanaTable")
	testUnicharTable(t, widthTable, widthFirst, widthLast, "widthTable")
}

type ToVoicedTest struct {
	in             rune
	outTraditional string
	outCombining   string
}

var tovoicedtests = []ToVoicedTest{
	// vcUnvoiced
	0: {'か', "が", "か\u3099"}, // compatVs is exists, compatSvs is not exists
	1: {'は', "ば", "は\u3099"}, // compatVs is exists, compatSvs is exists
	// vcVoiced
	2: {'が', "が", "か\u3099"}, // compatVs.compatSvs is not exists
	3: {'ば', "ば", "は\u3099"}, // compatVs.compatSvs is exists
	// vcSemivoiced
	4: {'ぱ', "ば", "は\u3099"},
	// vcUndefined, cwWide
	5:  {'あ', "あ゛", "あ\u3099"}, // ctKanaLetter, ccHiragana
	6:  {'ア', "ア゛", "ア\u3099"}, // ctKanaLetter, ccKatakana
	7:  {'・', "・゛", "・\u3099"}, // ctKanaSymbol, ccUndefined
	8:  {'Ａ', "Ａ゛", "Ａ\u3099"}, // ctLatinLetter, ccUpper
	9:  {'ａ', "ａ゛", "ａ\u3099"}, // ctLatinLetter, ccLower
	10: {'１', "１゛", "１\u3099"}, // ctLatinDigit, ccUndefined
	11: {'＃', "＃゛", "＃\u3099"}, // ctLatinSymbol, ccUndefined
	// vcUndefined, cwNarrow
	12: {'ｱ', "ｱﾞ", "ｱ\u3099"}, // ctKanaLetter, ccKatakana
	13: {'ｶ', "ｶﾞ", "ｶ\u3099"}, // ctKanaLetter, ccKatakana, compatWidth.compatVs is exists, compatWidth.compatSvs is not exists
	14: {'ﾊ', "ﾊﾞ", "ﾊ\u3099"}, // ctKanaLetter, ccKatakana, compatWidth.compatVs is exists, compatWidth.compatSvs is exists
	15: {'･', "･ﾞ", "･\u3099"}, // ctKanaSymbol, ccUndefined
	16: {'A', "Aﾞ", "A\u3099"}, // ctLatinLetter, ccUpper
	17: {'a', "aﾞ", "a\u3099"}, // ctLatinLetter, ccLower
	18: {'1', "1ﾞ", "1\u3099"}, // ctLatinDigit, ccUndefined
	19: {'#', "#ﾞ", "#\u3099"}, // ctLatinDigit, ccUndefined
	// ctUndefined
	20: {'\u3040', "\u3040", "\u3040\u3099"},
	// VSM
	21: {'゛', "゛゛", "゛\u3099"},
	22: {'\u3099', "\u3099゛", "\u3099\u3099"},
	23: {'ﾞ', "ﾞﾞ", "ﾞ\u3099"},
	// SVSM
	24: {'゜', "゜゛", "゜\u3099"},
	25: {'\u309A', "\u309A゛", "\u309A\u3099"},
	26: {'ﾟ', "ﾟﾞ", "ﾟ\u3099"},
}

func TestToVoiced(t *testing.T) {
	for n, tt := range tovoicedtests {
		c, ok := getUnichar(tt.in)
		if !ok {
			t.Errorf("%d: %#U is not found by getUnichar()", n, tt.in)
			continue
		}
		var got, want []rune
		got = c.toTraditionalVoiced()
		want = []rune(tt.outTraditional)
		if len(got) != len(want) {
			t.Errorf("%d: got: %q, want: %q", n, string(got), tt.outTraditional)
			continue
		}
		for i := 0; i < len(got); i++ {
			if got[i] != want[i] {
				t.Errorf("%d: got: %q, want: %q", n, string(got), tt.outTraditional)
				break
			}
		}
		got = c.toCombiningVoiced()
		want = []rune(tt.outCombining)
		if len(got) != len(want) {
			t.Errorf("%d: got: %q, want: %q", n, string(got), tt.outCombining)
			continue
		}
		for i := 0; i < len(got); i++ {
			if got[i] != want[i] {
				t.Errorf("%d: got: %q, want: %q", n, string(got), tt.outCombining)
				break
			}
		}

	}
}

type ToSemivoicedTest struct {
	in             rune
	outTraditional string
	outCombining   string
}

var tosemivoicedtests = []ToSemivoicedTest{
	// vcUnvoiced
	0: {'か', "か゜", "か\u309A"}, // compatVs is exists, compatSvs is not exists
	1: {'は', "ぱ", "は\u309A"},  // compatVs is exists, compatSvs is exists
	// vcVoiced
	2: {'が', "か゜", "か\u309A"}, // compatVs.compatSvs is not exists
	3: {'ば', "ぱ", "は\u309A"},  // compatVs.compatSvs is exists
	// vcSemivoiced
	4: {'ぱ', "ぱ", "は\u309A"},
	// vcUndefined, cwWide
	5:  {'あ', "あ゜", "あ\u309A"}, // ctKanaLetter, ccHiragana
	6:  {'ア', "ア゜", "ア\u309A"}, // ctKanaLetter, ccKatakana
	7:  {'・', "・゜", "・\u309A"}, // ctKanaSymbol, ccUndefined
	8:  {'Ａ', "Ａ゜", "Ａ\u309A"}, // ctLatinLetter, ccUpper
	9:  {'ａ', "ａ゜", "ａ\u309A"}, // ctLatinLetter, ccLower
	10: {'１', "１゜", "１\u309A"}, // ctLatinDigit, ccUndefined
	11: {'＃', "＃゜", "＃\u309A"}, // ctLatinSymbol, ccUndefined
	// vcUndefined, cwNarrow
	12: {'ｱ', "ｱﾟ", "ｱ\u309A"}, // ctKanaLetter, ccKatakana
	13: {'ｶ', "ｶﾟ", "ｶ\u309A"}, // ctKanaLetter, ccKatakana, compatWidth.compatVs is exists, compatWidth.compatSvs is not exists
	14: {'ﾊ', "ﾊﾟ", "ﾊ\u309A"}, // ctKanaLetter, ccKatakana, compatWidth.compatVs is exists, compatWidth.compatSvs is exists
	15: {'･', "･ﾟ", "･\u309A"}, // ctKanaSymbol, ccUndefined
	16: {'A', "Aﾟ", "A\u309A"}, // ctLatinLetter, ccUpper
	17: {'a', "aﾟ", "a\u309A"}, // ctLatinLetter, ccLower
	18: {'1', "1ﾟ", "1\u309A"}, // ctLatinDigit, ccUndefined
	19: {'#', "#ﾟ", "#\u309A"}, // ctLatinDigit, ccUndefined
	// ctUndefined
	20: {'\u3040', "\u3040", "\u3040\u309A"},
	// VSM
	21: {'゛', "゛゜", "゛\u309A"},
	22: {'\u3099', "\u3099゜", "\u3099\u309A"},
	23: {'ﾞ', "ﾞﾟ", "ﾞ\u309A"},
	// SVSM
	24: {'゜', "゜゜", "゜\u309A"},
	25: {'\u309A', "\u309A゜", "\u309A\u309A"},
	26: {'ﾟ', "ﾟﾟ", "ﾟ\u309A"},
}

func TestToSemivoiced(t *testing.T) {
	for n, tt := range tosemivoicedtests {
		c, ok := getUnichar(tt.in)
		if !ok {
			t.Errorf("%d: %#U is not found by getUnichar()", n, tt.in)
			continue
		}
		var got, want []rune
		got = c.toTraditionalSemivoiced()
		want = []rune(tt.outTraditional)
		if len(got) != len(want) {
			t.Errorf("%d: got: %q, want: %q", n, string(got), tt.outTraditional)
			continue
		}
		for i := 0; i < len(got); i++ {
			if got[i] != want[i] {
				t.Errorf("%d: got: %q, want: %q", n, string(got), tt.outTraditional)
				break
			}
		}
		got = c.toCombiningSemivoiced()
		want = []rune(tt.outCombining)
		if len(got) != len(want) {
			t.Errorf("%d: got: %q, want: %q", n, string(got), tt.outCombining)
			continue
		}
		for i := 0; i < len(got); i++ {
			if got[i] != want[i] {
				t.Errorf("%d: got: %q, want: %q", n, string(got), tt.outCombining)
				break
			}
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
