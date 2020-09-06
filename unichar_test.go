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
		c.voicing, c.cmptCase, c.cmptWidth, c.cmptVs, c.cmptSvs)
}

type tableInfo struct {
	table   []unichar
	name    string
	wantN   int
	wantSHA string
}

var tables = []tableInfo{
	{latinTable, "latinTable", 96, "d90c9a10f72b6c029cb6aba58128f534d1935aa760fd5a073c232e16fc4eca22"},
	{kanaTable, "kanaTable", 256, "1884391088215e839715d7261c5317b45fca7866ff6ddad4e86536fbffb88f05"},
	{widthTable, "widthTable", 160, "e036650516323d6c60ffa1ca4c51f82983835a86ef69cd0cdc97e0f9d349118c"},
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
		case ctUndefined, ctLatinLetter, ctLatinDigit, ctLatinSymbol, ctKanaLetter, ctKanaSymbol:
		default: // TEST_P8w4qtsm
			t.Errorf("%s[%#U].category == %d, want %d <= category < %d",
				name, c.codepoint, c.category, ctUndefined, ctMax)
		}

		// character case
		if c.charCase < ccUndefined || c.charCase >= ccMax {
			t.Errorf("%s[%#U].charCase == %d, want %d <= charCase < %d",
				name, c.codepoint, c.charCase, ccUndefined, ccMax)
		}
		if c.charCase == ccUndefined && c.cmptCase != c.codepoint {
			t.Errorf("%s[%#U].charCase == %d, but cmptCase == %#U, want cmptCase == %#U",
				name, c.codepoint, c.charCase, c.cmptCase, c.codepoint)
		}
		if c.cmptCase != c.codepoint {
			if cmptCase, ok := getUnichar(c.cmptCase); !ok { // TEST_fm8XjZTB
				t.Errorf("%s[%#U].cmptCase %#U is not found by getUnichar()",
					name, c.codepoint, c.cmptCase)
			} else {
				if cmptCase.category != c.category {
					t.Errorf("%s[%#U].category == %d, but cmptCase %#U.category == %d, want same value",
						name, c.codepoint, c.category, cmptCase.codepoint, cmptCase.category)
				}
				if cmptCase.charCase == c.charCase {
					t.Errorf("%s[%#U].charCase == %d, but cmptCase %#U.charCase == %d, want another one",
						name, c.codepoint, c.charCase, cmptCase.codepoint, cmptCase.charCase)
				}
			}
		}
		// TODO Check for problems in the following cases
		/*
			if c.charCase != ccUndefined && c.cmptCase == c.codepoint {
				t.Errorf("%s[%#U].charCase is %d and cmptCase is %#U, want another one",
					name, c.codepoint, c.charCase, c.cmptCase)
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
			if c.cmptWidth != c.codepoint {
				t.Errorf("%s[%#U].charWidth == %d, but cmptWidth == %#U, want cmptWidth == %#U",
					name, c.codepoint, c.charWidth, c.cmptWidth, c.codepoint)
			}
		}
		if c.cmptWidth != c.codepoint {
			if cmptWidth, ok := getUnichar(c.cmptWidth); !ok { // TEST_T3bc4Nh7
				t.Errorf("%s[%#U].cmptWidth %#U is not found by getUnichar()",
					name, c.codepoint, c.cmptWidth)
			} else {
				if cmptWidth.category != c.category {
					t.Errorf("%s[%#U].category == %d, but cmptWidth %#U.category == %d, want same value",
						name, c.codepoint, c.category, cmptWidth.codepoint, cmptWidth.category)
				}
				if cmptWidth.charWidth == c.charWidth {
					t.Errorf("%s[%#U].charWidth is %d, but cmptWidth %#U.charWidth == %d, want another one",
						name, c.codepoint, c.charWidth, cmptWidth.codepoint, cmptWidth.charWidth)
				}
			}
		}
		// TODO Check for problems in the following cases
		/*
			if c.charWidth != cwUndefined && c.cmptWidth == c.codepoint {
				t.Errorf("%s[%#U].charWidth is %d and cmptWidth is %#U, want another one",
					name, c.codepoint, c.charWidth, c.cmptWidth)
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
			if c.cmptVs != c.codepoint {
				t.Errorf("%s[%#U].voicing == %d, but cmptVs == %#U, want cmptVs == %#U",
					name, c.codepoint, c.voicing, c.cmptVs, c.codepoint)
			}
			if c.cmptSvs != c.codepoint {
				t.Errorf("%s[%#U].voicing == %d, but cmptSvs == %#U, want cmptSvs == %#U",
					name, c.codepoint, c.voicing, c.cmptSvs, c.codepoint)
			}
		}
		if c.voicing == vcUnvoiced {
			if c.cmptVs == c.codepoint { // TEST_Jt3UaWwr
				t.Errorf("%s[%#U].voicing == %d, but cmptVs == %#U, want another one",
					name, c.codepoint, c.voicing, c.cmptVs)
			}
			if c.charWidth != cwNarrow && c.charWidth != cwWide { // TEST_Mw87qjkF
				t.Errorf("%s[%#U].voicing == %d, but charWidth == %d, want %d or %d",
					name, c.codepoint, c.voicing, c.charWidth, cwNarrow, cwWide)
			}
		}
		if c.voicing == vcVoiced {
			if c.cmptVs == c.codepoint {
				t.Errorf("%s[%#U].voicing == %d, but cmptVs == %#U, want another one",
					name, c.codepoint, c.voicing, c.cmptVs)
			}
			if c.charWidth != cwNarrow && c.charWidth != cwWide { // TEST_T2eKd76G
				t.Errorf("%s[%#U].voicing == %d, but charWidth == %d, want %d or %d",
					name, c.codepoint, c.voicing, c.charWidth, cwNarrow, cwWide)
			}
		}
		if c.voicing == vcSemivoiced {
			if c.cmptSvs == c.codepoint {
				t.Errorf("%s[%#U].voicing is %d and cmptSvs is %#U, want another one",
					name, c.codepoint, c.voicing, c.cmptSvs)
			}
			if c.charWidth != cwNarrow && c.charWidth != cwWide {
				t.Errorf("%s[%#U].voicing == %d, but charWidth == %d, want %d or %d",
					name, c.codepoint, c.voicing, c.charWidth, cwNarrow, cwWide)
			}
			unvoiced := c.getCmptSvsUnichar()
			if !unvoiced.existsCmptVs() { // TEST_fW6auXUi
				t.Errorf("%s[%#U].voicing == %d, but getCmptSvsUnichar().existsCmptVs() == false, want true",
					name, c.codepoint, c.voicing)
			}
		}

		if c.cmptVs != c.codepoint {
			if cmptVs, ok := getUnichar(c.cmptVs); !ok { // TEST_Cu8iKMxF
				t.Errorf("%s[%#U].cmptVs %#U is not found by getUnichar()", name, c.codepoint, c.cmptVs)
			} else {
				if c.category != cmptVs.category {
					t.Errorf("%s[%#U].category is %d and cmptVs %#U.category is %d, want same value",
						name, c.codepoint, c.category, cmptVs.codepoint, cmptVs.category)
				}
				if c.charWidth != cmptVs.charWidth {
					t.Errorf("%s[%#U].charWidth is %d and cmptVs %#U.charWidth is %d, want same value",
						name, c.codepoint, c.charWidth, cmptVs.codepoint, cmptVs.charWidth)
				}
				if c.voicing == vcUnvoiced && cmptVs.voicing != vcVoiced {
					t.Errorf("%s[%#U].voicing is %d and cmptVs %#U.voicing is %d, want %d",
						name, c.codepoint, c.voicing, cmptVs.codepoint, cmptVs.voicing, vcVoiced)
				}
				if c.voicing == vcVoiced && cmptVs.voicing != vcUnvoiced {
					t.Errorf("%s[%#U].voicing is %d and cmptVs %#U.voicing is %d, want %d",
						name, c.codepoint, c.voicing, cmptVs.codepoint, cmptVs.voicing, vcUnvoiced)
				}
			}
		}
		if c.cmptSvs != c.codepoint {
			if cmptSvs, ok := getUnichar(c.cmptSvs); !ok { // TEST_rW4UiNHC
				t.Errorf("%s[%#U].cmptSvs %#U is not found by getUnichar()", name, c.codepoint, c.cmptSvs)
			} else {
				if c.category != cmptSvs.category {
					t.Errorf("%s[%#U].category is %d and cmptSvs %#U.category is %d, want same value",
						name, c.codepoint, c.category, cmptSvs.codepoint, cmptSvs.category)
				}
				if c.charWidth != cmptSvs.charWidth {
					t.Errorf("%s[%#U].charWidth is %d and cmptSvs %#U.charWidth is %d, want same value",
						name, c.codepoint, c.charWidth, cmptSvs.codepoint, cmptSvs.charWidth)
				}
				if c.voicing == vcUnvoiced && cmptSvs.voicing != vcSemivoiced {
					t.Errorf("%s[%#U].voicing is %d and cmptSvs %#U.voicing is %d, want %d",
						name, c.codepoint, c.voicing, cmptSvs.codepoint, cmptSvs.voicing, vcSemivoiced)
				}
				if c.voicing == vcSemivoiced && cmptSvs.voicing != vcUnvoiced {
					t.Errorf("%s[%#U].voicing is %d and cmptSvs %#U.voicing is %d, want %d",
						name, c.codepoint, c.voicing, cmptSvs.codepoint, cmptSvs.voicing, vcUnvoiced)
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
			if c.cmptVs != c.codepoint {
				t.Errorf("len(%s[%#U].cmptVs) = %#U, want: %#U", name, c.codepoint, c.cmptVs, c.codepoint)
			}
			if c.cmptSvs != c.codepoint {
				t.Errorf("len(%s[%#U].cmptSvs) = %#U, want: %#U", name, c.codepoint, c.cmptSvs, c.codepoint)
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
			if c.cmptCase != c.codepoint {
				t.Errorf("len(%s[%#U].cmptCase) = %#U, want: %#U",
					name, c.codepoint, c.cmptCase, c.codepoint)
			}
			if c.cmptWidth != c.codepoint {
				t.Errorf("len(%s[%#U].cmptWidth) = %#U, want: %#U",
					name, c.codepoint, c.cmptWidth, c.codepoint)
			}
			if c.cmptVs != c.codepoint {
				t.Errorf("len(%s[%#U].cmptVs) = %#U, want: %#U",
					name, c.codepoint, c.cmptVs, c.codepoint)
			}
			if c.cmptSvs != c.codepoint {
				t.Errorf("len(%s[%#U].cmptSvs) = %#U, want: %#U",
					name, c.codepoint, c.cmptSvs, c.codepoint)
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
	in           rune
	outClassical string
	outCombining string
}

var tovoicedtests = []ToVoicedTest{
	// vcUnvoiced
	0: {'か', "が", "か\u3099"}, // cmptVs is exists, cmptSvs is not exists
	1: {'は', "ば", "は\u3099"}, // cmptVs is exists, cmptSvs is exists
	// vcVoiced
	2: {'が', "が", "か\u3099"}, // cmptVs.cmptSvs is not exists
	3: {'ば', "ば", "は\u3099"}, // cmptVs.cmptSvs is exists
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
	13: {'ｶ', "ｶﾞ", "ｶ\u3099"}, // ctKanaLetter, ccKatakana, cmptWidth.cmptVs is exists, cmptWidth.cmptSvs is not exists
	14: {'ﾊ', "ﾊﾞ", "ﾊ\u3099"}, // ctKanaLetter, ccKatakana, cmptWidth.cmptVs is exists, cmptWidth.cmptSvs is exists
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
		got = c.toClassicalVoiced()
		want = []rune(tt.outClassical)
		if len(got) != len(want) {
			t.Errorf("%d: got: %q, want: %q", n, string(got), tt.outClassical)
			continue
		}
		for i := 0; i < len(got); i++ {
			if got[i] != want[i] {
				t.Errorf("%d: got: %q, want: %q", n, string(got), tt.outClassical)
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
	in           rune
	outClassical string
	outCombining string
}

var tosemivoicedtests = []ToSemivoicedTest{
	// vcUnvoiced
	0: {'か', "か゜", "か\u309A"}, // cmptVs is exists, cmptSvs is not exists
	1: {'は', "ぱ", "は\u309A"},  // cmptVs is exists, cmptSvs is exists
	// vcVoiced
	2: {'が', "か゜", "か\u309A"}, // cmptVs.cmptSvs is not exists
	3: {'ば', "ぱ", "は\u309A"},  // cmptVs.cmptSvs is exists
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
	13: {'ｶ', "ｶﾟ", "ｶ\u309A"}, // ctKanaLetter, ccKatakana, cmptWidth.cmptVs is exists, cmptWidth.cmptSvs is not exists
	14: {'ﾊ', "ﾊﾟ", "ﾊ\u309A"}, // ctKanaLetter, ccKatakana, cmptWidth.cmptVs is exists, cmptWidth.cmptSvs is exists
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
		got = c.toClassicalSemivoiced()
		want = []rune(tt.outClassical)
		if len(got) != len(want) {
			t.Errorf("%d: got: %q, want: %q", n, string(got), tt.outClassical)
			continue
		}
		for i := 0; i < len(got); i++ {
			if got[i] != want[i] {
				t.Errorf("%d: got: %q, want: %q", n, string(got), tt.outClassical)
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
