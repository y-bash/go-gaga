package gaga

import (
	"fmt"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/width"
	"log"
	"strings"
	"testing"
	"unicode/utf8"
)

const (
	maxr = utf8.MaxRune
	excr = maxr + 1
)

type Normalizer_normalizeRuneTest struct {
	name string
	flag NormFlag
	lo   rune
	hi   rune
	diff rune
}

var normalizer_normalizerunetests = []Normalizer_normalizeRuneTest{
	// latin letter
	0: {"a-z -> A-Z", AlphaToUpper, 'a', 'z', 'A' - 'a'},
	1: {"A-Z -> a-ｚ", AlphaToLower, 'A', 'Z', 'a' - 'A'},
	2: {"A-Z -> Ａ-Ｚ", AlphaToWide, 'A', 'Z', 'Ａ' - 'A'},
	3: {"Ａ-Ｚ -> A-Z", AlphaToNarrow, 'Ａ', 'Ｚ', 'A' - 'Ａ'},
	4: {"a-z -> Ａ-Ｚ", AlphaToUpper | AlphaToWide, 'a', 'z', 'Ａ' - 'a'},
	5: {"A-Z->ａ-ｚ", AlphaToLower | AlphaToWide, 'A', 'Z', 'ａ' - 'A'},
	6: {"ａ-ｚ -> A-Z", AlphaToUpper | AlphaToNarrow, 'ａ', 'ｚ', 'A' - 'ａ'},
	7: {"Ａ-Ｚ -> a-z", AlphaToLower | AlphaToNarrow, 'Ａ', 'Ｚ', 'a' - 'Ａ'},
	// latin digit
	8: {"0-9 -> ０-９", DigitToWide, '0', '9', '０' - '0'},
	9: {"０-９ -> 0-9", DigitToNarrow, '０', '９', '0' - '０'},
	// latin symbol
	10: {"!-/ -> !-/", SymbolToWide, '!', '/', '！' - '!'},
	11: {":-@ -> :-@", SymbolToWide, ':', '@', '：' - ':'},
	12: {"[-` -> [-`", SymbolToWide, '[', '`', '［' - '['},
	13: {"{-~ -> {-~", SymbolToWide, '{', '~', '｛' - '{'},
	14: {"！-／ -> !-/", SymbolToNarrow, '！', '／', '!' - '！'},
	15: {"：-＠ -> :-@", SymbolToNarrow, '：', '＠', ':' - '：'},
	16: {"［-｀ -> [-`", SymbolToNarrow, '［', '｀', '[' - '［'},
	17: {"｛-〜 -> {-~", SymbolToNarrow, '｛', '～', '{' - '｛'},
	// latin all
	18: {"!-~ -> ！-〜", LatinToWide, '!', '~', '！' - '!'},
	19: {"！-〜 -> !-~", LatinToNarrow, '！', '～', '!' - '！'},
	// kana letter
	20: {"ぁ-ゖ -> ァ-ヶ", HiraganaToKatakana, 'ぁ', 'ゖ', 'ァ' - 'ぁ'},
	21: {"ァ-ヶ -> ぁ-ゖ", KatakanaToHiragana, 'ァ', 'ヶ', 'ぁ' - 'ァ'},
	// no effect latin letter
	22: {"a-z -> a-z", DigitToWide | SymbolToWide | KanaToWide, 'a', 'z', 0},
	23: {"A-Z -> A-Z", DigitToWide | SymbolToWide | KanaToWide, 'A', 'Z', 0},
	24: {"ａ-ｚ -> ａ-ｚ", DigitToWide | SymbolToWide | KanaToWide, 'ａ', 'ｚ', 0},
	25: {"Ａ-Ｚ -> Ａ-Ｚ", DigitToWide | SymbolToWide | KanaToWide, 'Ａ', 'Ｚ', 0},
	// no effect latin number
	26: {"0-9 -> 0-9", AlphaToUpper | AlphaToWide | SymbolToWide | KanaToWide, '0', '9', 0},
	27: {"0-9 -> 0-9", AlphaToUpper | AlphaToNarrow | SymbolToNarrow | KanaToNarrow, '0', '9', 0},
	28: {"0-9 -> 0-9", AlphaToLower | AlphaToWide | SymbolToWide | KanaToWide, '0', '9', 0},
	29: {"0-9 -> 0-9", AlphaToLower | AlphaToNarrow | SymbolToNarrow | KanaToNarrow, '0', '9', 0},
	30: {"０-９ -> ０-９", AlphaToUpper | AlphaToWide | SymbolToWide | KanaToWide, '０', '９', 0},
	31: {"０-９ -> ０-９", AlphaToUpper | AlphaToNarrow | SymbolToNarrow | KanaToNarrow, '０', '９', 0},
	32: {"０-９ -> ０-９", AlphaToLower | AlphaToWide | SymbolToWide | KanaToWide, '０', '９', 0},
	33: {"０-９ -> ０-９", AlphaToLower | AlphaToNarrow | SymbolToNarrow | KanaToNarrow, '０', '９', 0},
	// no effect latin symbol
	34: {"!-/ -> !-/", AlphaToUpper | AlphaToWide | DigitToWide | KanaToWide, '!', '/', 0},
	35: {":-@ -> :-@", AlphaToUpper | AlphaToWide | DigitToWide | KanaToWide, ':', '@', 0},
	36: {"[-` -> [-`", AlphaToUpper | AlphaToWide | DigitToWide | KanaToWide, '[', '`', 0},
	37: {"{-~ -> {-~", AlphaToUpper | AlphaToWide | DigitToWide | KanaToWide, '{', '~', 0},
	38: {"!-/ -> !-/", AlphaToLower | AlphaToWide | DigitToWide | KanaToWide, '!', '/', 0},
	39: {":-@ -> :-@", AlphaToLower | AlphaToWide | DigitToWide | KanaToWide, ':', '@', 0},
	40: {"[-` -> [-`", AlphaToLower | AlphaToWide | DigitToWide | KanaToWide, '[', '`', 0},
	41: {"{-~ -> {-~", AlphaToLower | AlphaToWide | DigitToWide | KanaToWide, '{', '~', 0},
	42: {"!-/ -> !-/", AlphaToUpper | AlphaToNarrow | DigitToNarrow | KanaToNarrow, '!', '/', 0},
	43: {":-@ -> :-@", AlphaToUpper | AlphaToNarrow | DigitToNarrow | KanaToNarrow, ':', '@', 0},
	44: {"[-` -> [-`", AlphaToUpper | AlphaToNarrow | DigitToNarrow | KanaToNarrow, '[', '`', 0},
	45: {"{-~ -> {-~", AlphaToUpper | AlphaToNarrow | DigitToNarrow | KanaToNarrow, '{', '~', 0},
	46: {"!-/ -> !-/", AlphaToLower | AlphaToNarrow | DigitToNarrow | KanaToNarrow, '!', '/', 0},
	47: {":-@ -> :-@", AlphaToLower | AlphaToNarrow | DigitToNarrow | KanaToNarrow, ':', '@', 0},
	48: {"[-` -> [-`", AlphaToLower | AlphaToNarrow | DigitToNarrow | KanaToNarrow, '[', '`', 0},
	49: {"{-~ -> {-~", AlphaToLower | AlphaToNarrow | DigitToNarrow | KanaToNarrow, '{', '~', 0},
	50: {"！-／ -> ！-／ ", AlphaToUpper | AlphaToWide | DigitToWide | KanaToWide, '！', '／', 0},
	51: {"：-＠ -> ：-＠ ", AlphaToUpper | AlphaToWide | DigitToWide | KanaToWide, '：', '＠', 0},
	52: {"［-｀ -> ［-｀ ", AlphaToUpper | AlphaToWide | DigitToWide | KanaToWide, '［', '｀', 0},
	53: {"｛-〜 -> ｛-〜 ", AlphaToUpper | AlphaToWide | DigitToWide | KanaToWide, '｛', '〜', 0},
	54: {"！-／ -> ！-／ ", AlphaToLower | AlphaToWide | DigitToWide | KanaToWide, '！', '／', 0},
	55: {"：-＠ -> ：-＠ ", AlphaToLower | AlphaToWide | DigitToWide | KanaToWide, '：', '＠', 0},
	56: {"［-｀ -> ［-｀ ", AlphaToLower | AlphaToWide | DigitToWide | KanaToWide, '［', '｀', 0},
	57: {"｛-〜 -> ｛-〜 ", AlphaToLower | AlphaToWide | DigitToWide | KanaToWide, '｛', '〜', 0},
	58: {"！-／ -> ！-／ ", AlphaToUpper | AlphaToNarrow | DigitToNarrow | KanaToNarrow, '！', '／', 0},
	59: {"：-＠ -> ：-＠ ", AlphaToUpper | AlphaToNarrow | DigitToNarrow | KanaToNarrow, '：', '＠', 0},
	60: {"［-｀ -> ［-｀ ", AlphaToUpper | AlphaToNarrow | DigitToNarrow | KanaToNarrow, '［', '｀', 0},
	61: {"｛-〜 -> ｛-〜 ", AlphaToUpper | AlphaToNarrow | DigitToNarrow | KanaToNarrow, '｛', '〜', 0},
	62: {"！-／ -> ！-／ ", AlphaToLower | AlphaToNarrow | DigitToNarrow | KanaToNarrow, '！', '／', 0},
	63: {"：-＠ -> ：-＠ ", AlphaToLower | AlphaToNarrow | DigitToNarrow | KanaToNarrow, '：', '＠', 0},
	64: {"［-｀ -> ［-｀ ", AlphaToLower | AlphaToNarrow | DigitToNarrow | KanaToNarrow, '［', '｀', 0},
	65: {"｛-〜 -> ｛-〜 ", AlphaToLower | AlphaToNarrow | DigitToNarrow | KanaToNarrow, '｛', '〜', 0},
	// no effect kana letter
	66: {"ぁ-ゖ -> ぁ-ゖ", LatinToWide | AlphaToUpper, 'ぁ', 'ゖ', 0},
	67: {"ぁ-ゖ -> ぁ-ゖ", LatinToWide | AlphaToLower, 'ぁ', 'ゖ', 0},
	68: {"ぁ-ゖ -> ぁ-ゖ", LatinToNarrow | AlphaToUpper, 'ぁ', 'ゖ', 0},
	69: {"ぁ-ゖ -> ぁ-ゖ", LatinToNarrow | AlphaToLower, 'ぁ', 'ゖ', 0},
	70: {"ァ-ヺ -> ァ-ヺ", LatinToWide | AlphaToUpper, 'ァ', 'ヺ', 0},
	71: {"ァ-ヺ -> ァ-ヺ", LatinToWide | AlphaToLower, 'ァ', 'ヺ', 0},
	72: {"ァ-ヺ -> ァ-ヺ", LatinToNarrow | AlphaToUpper, 'ァ', 'ヺ', 0},
	73: {"ァ-ヺ -> ァ-ヺ", LatinToNarrow | AlphaToLower, 'ァ', 'ヺ', 0},
	// no effect kana symbol
	74: {"、-〠 -> 、-〠", LatinToWide | AlphaToUpper, '、', '〠', 0},
	75: {"、-〠 -> 、-〠", LatinToWide | AlphaToLower, '、', '〠', 0},
	76: {"、-〠 -> 、-〠", LatinToNarrow | AlphaToUpper, '、', '〠', 0},
	77: {"、-〠 -> 、-〠", LatinToNarrow | AlphaToLower, '、', '〠', 0},
	78: {"｡-･ -> ｡-･", LatinToWide | AlphaToUpper, '｡', '･', 0},
	79: {"｡-･ -> ｡-･", LatinToWide | AlphaToLower, '｡', '･', 0},
	80: {"｡-･ -> ｡-･", LatinToNarrow | AlphaToUpper, '｡', '･', 0},
	81: {"｡-･ -> ｡-･", LatinToNarrow | AlphaToLower, '｡', '･', 0},
}

func TestNormalizer_normalizeRune(t *testing.T) {
	for i, tt := range normalizer_normalizerunetests {
		n, err := Norm(tt.flag)
		if err != nil {
			t.Errorf("#%d: %s", i, err.Error())
			continue
		}
		for src := tt.lo; src <= tt.hi; src++ {
			want, wantVm := src+tt.diff, vmNone
			have, haveVm := n.normalizeRune(src)
			if have != want || haveVm != vmNone {
				t.Errorf("#%d %s, %s\nnormalizeRune(%#U)\nhave:(%#U, %#U), \nwant:(%#U, %#U)",
					i, tt.name, tt.flag, src, have, haveVm, want, wantVm)
			}
		}
	}
}

type Normalizer_RuneTest struct {
	flag NormFlag
	in   rune
	out  string
}

var normalizer_runetests = []Normalizer_RuneTest{
	0: {AlphaToUpper, 'a', "A"},
	1: {HiraganaToNarrow, 'が', "ｶﾞ"},
}

func TestNormalizer_Rune(t *testing.T) {
	for i, tt := range normalizer_runetests {
		n, err := Norm(tt.flag)
		if err != nil {
			t.Errorf("#%d: %s", i, err.Error())
			continue
		}
		have := n.Rune(tt.in)
		if have != tt.out {
			t.Errorf("#%d %s, Rune(%#U) = %q, want: %q", i, tt.flag, tt.in, have, tt.out)
		}
	}
}

type Normalizer_StringTest struct {
	flag NormFlag
	in   string
	out  string
}

var normalizer_stringtests = []Normalizer_StringTest{
	// simple latin conversion <- zero length string
	0: {AlphaToUpper, "", ""},
	1: {AlphaToLower, "", ""},
	2: {AlphaToWide, "", ""},
	3: {AlphaToNarrow, "", ""},
	4: {DigitToWide, "", ""},
	5: {DigitToNarrow, "", ""},
	6: {SymbolToWide, "", ""},
	7: {SymbolToNarrow, "", ""},
	8: {LatinToWide, "", ""},
	9: {LatinToNarrow, "", ""},

	// simple latin conversion <- Latin (Letter/Digit/Symbol)
	10: {AlphaToUpper,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZAZ09!~　ＡＺＡＺ０９！～"},
	11: {AlphaToLower,
		" AZaz09!~　ＡＺａｚ０９！～",
		" azaz09!~　ａｚａｚ０９！～"},
	12: {AlphaToWide,
		" AZaz09!~　ＡＺａｚ０９！～",
		" ＡＺａｚ09!~　ＡＺａｚ０９！～"},
	13: {AlphaToNarrow,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　AZaz０９！～"},
	14: {DigitToWide,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz０９!~　ＡＺａｚ０９！～"},
	15: {DigitToNarrow,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ09！～"},
	16: {SymbolToWide,
		" AZaz09!~　ＡＺａｚ０９！～",
		"　AZaz09！～　ＡＺａｚ０９！～"},
	17: {SymbolToNarrow,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~ ＡＺａｚ０９!~"},
	18: {LatinToWide,
		" AZaz09!~　ＡＺａｚ０９！～",
		"　ＡＺａｚ０９！～　ＡＺａｚ０９！～"},
	19: {LatinToNarrow,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~ AZaz09!~"},

	// simple latin conversion <- CJK (Hiragana/Katakana/Symbol/Han), Emoji
	20: {AlphaToUpper,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	21: {AlphaToLower,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	22: {AlphaToWide,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	23: {AlphaToNarrow,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	24: {DigitToWide,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	25: {DigitToNarrow,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	26: {SymbolToWide,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	27: {SymbolToNarrow,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	28: {LatinToWide,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	29: {LatinToNarrow,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},

	// multiple latin conversion <- Latin (Letter/Digit/Symbol)
	30: {AlphaToUpper | LatinToWide,
		" AZaz09!~　ＡＺａｚ０９！～",
		"　ＡＺＡＺ０９！～　ＡＺＡＺ０９！～"},
	31: {AlphaToUpper | LatinToNarrow,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZAZ09!~ AZAZ09!~"},
	32: {AlphaToLower | LatinToWide,
		" AZaz09!~　ＡＺａｚ０９！～",
		"　ａｚａｚ０９！～　ａｚａｚ０９！～"},
	33: {AlphaToLower | LatinToNarrow,
		" AZaz09!~　ＡＺａｚ０９！～",
		" azaz09!~ azaz09!~"},
	34: {AlphaToUpper | AlphaToWide,
		" AZaz09!~　ＡＺａｚ０９！～",
		" ＡＺＡＺ09!~　ＡＺＡＺ０９！～"},
	35: {AlphaToUpper | AlphaToNarrow,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZAZ09!~　AZAZ０９！～"},
	36: {AlphaToLower | AlphaToWide,
		" AZaz09!~　ＡＺａｚ０９！～",
		" ａｚａｚ09!~　ａｚａｚ０９！～"},
	37: {AlphaToLower | AlphaToNarrow,
		" AZaz09!~　ＡＺａｚ０９！～",
		" azaz09!~　azaz０９！～"},
	38: {AlphaToUpper | DigitToWide,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZAZ０９!~　ＡＺＡＺ０９！～"},
	39: {AlphaToLower | SymbolToNarrow,
		" AZaz09!~　ＡＺａｚ０９！～",
		" azaz09!~ ａｚａｚ０９!~"},

	// multiple latin conversion <- CJK (Hiragana/Katakana/Symbol/Han), Emoji
	40: {AlphaToUpper | LatinToWide,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	41: {AlphaToUpper | LatinToNarrow,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	42: {AlphaToLower | LatinToWide,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	43: {AlphaToLower | LatinToNarrow,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	44: {AlphaToUpper | AlphaToWide,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	45: {AlphaToUpper | AlphaToNarrow,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	46: {AlphaToLower | AlphaToWide,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	47: {AlphaToLower | AlphaToNarrow,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	48: {AlphaToUpper | DigitToWide,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	49: {AlphaToLower | SymbolToNarrow,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},

	// simple Japanese conversion <- zero length string
	50: {ComposeVom, "", ""},
	51: {HiraganaToNarrow | ComposeVom, "", ""},
	52: {HiraganaToKatakana | ComposeVom, "", ""},
	53: {KatakanaToWide | ComposeVom, "", ""},
	54: {KatakanaToNarrow | ComposeVom, "", ""},
	55: {KatakanaToHiragana | ComposeVom, "", ""},
	56: {KanaSymbolToWide | ComposeVom, "", ""},
	57: {KanaSymbolToNarrow | ComposeVom, "", ""},

	// simple Japanese conversion <- CJK (Hiragana/Katakana/Symbol/Han), Emoji
	58: {ComposeVom,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	59: {HiraganaToNarrow | ComposeVom,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」ｱｳﾞｧｹアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	60: {HiraganaToKatakana | ComposeVom,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」アヴァヶアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	61: {KatakanaToWide | ComposeVom,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣アヴァケ漢👻"},
	62: {KatakanaToNarrow | ComposeVom,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖｱｦﾞｧｹ､｣ｱｳﾞｧｹ漢👻"},
	63: {KatakanaToHiragana | ComposeVom,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖあを゛ぁゖ､｣あゔぁけ漢👻"},
	64: {KanaSymbolToWide | ComposeVom,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ、」ｱｳﾞｧｹ漢👻"},
	65: {KanaSymbolToNarrow | ComposeVom,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"､｣あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},

	// simple Japanese conversion <- Latin (Letter/Digit/Symbol)
	66: {ComposeVom,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ０９！～"},
	67: {HiraganaToNarrow | ComposeVom,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ０９！～"},
	68: {HiraganaToKatakana | ComposeVom,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ０９！～"},
	69: {KatakanaToWide | ComposeVom,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ０９！～"},
	70: {KatakanaToNarrow | ComposeVom,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ０９！～"},
	71: {KatakanaToHiragana | ComposeVom,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ０９！～"},
	72: {KanaSymbolToWide | ComposeVom,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ０９！～"},
	73: {KanaSymbolToNarrow | ComposeVom,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ０９！～"},

	// multiple Japanese conversion <- CJK (Hiragana/Katakana/Symbol/Han), Emoji
	74: {KatakanaToHiragana | KanaSymbolToWide | ComposeVom,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖあを゛ぁゖ、」あゔぁけ漢👻"},
	75: {KanaToNarrow | ComposeVom,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"､｣ｱｳﾞｧｹｱｦﾞｧｹ､｣ｱｳﾞｧｹ漢👻"},
	76: {KatakanaToHiragana | ComposeVom,
		"ｦｧｨｮｯｱｲﾛﾝﾜｲｴｶｹ",
		"をぁぃょっあいろんわいえかけ"},
	77: {KatakanaToHiragana | ComposeVom,
		"ァアィイレロヮワヰヱヲンヵカヶケヷヸヹヺ",
		"ぁあぃいれろゎわゐゑをんゕかゖけわ゛ゐ゛ゑ゛を゛"},

	// simple Japanese conversion (no vsm directive) <- Voiced sound character/Voiced sound mark
	78: {KatakanaToHiragana,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"がか゛かﾞか\u3099がか゛かﾞか\u3099か゛かﾞか\u3099"},
	79: {KatakanaToHiragana,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"あ゛あﾞあ\u3099あ゛あﾞあ\u3099あ゛あﾞあ\u3099"},
	80: {KatakanaToHiragana,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099"},

	// multiple Japanese conversion (no vsm directive) <- Voiced sound character/Voiced sound mark
	81: {KatakanaToNarrow | HiraganaToNarrow,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ｶﾞｶ゛ｶﾞｶ\u3099ｶﾞｶ゛ｶﾞｶ\u3099ｶ゛ｶﾞｶ\u3099"},
	82: {KatakanaToNarrow | HiraganaToNarrow,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"ｱ゛ｱﾞｱ\u3099ｱ゛ｱﾞｱ\u3099ｱ゛ｱﾞｱ\u3099"},
	83: {KatakanaToNarrow | HiraganaToNarrow,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099"},

	// simple Japanese conversion (voiced kana traditional directive) <- Voiced sound character/Voiced sound mark
	84: {KatakanaToHiragana | ComposeVom,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ががががががががががが"},
	85: {KatakanaToHiragana | ComposeVom,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"あ゛あ゛あ゛あ゛あ゛あ゛あ゛あ゛あ゛"},
	86: {KatakanaToHiragana | ComposeVom,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099"},

	// multiple Japanese conversion (voiced kana traditional directive) <- Voiced sound character/Voiced sound mark
	87: {KatakanaToNarrow | HiraganaToNarrow | ComposeVom,
		"か゛かﾞか\u3099カ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞ"},
	88: {KatakanaToNarrow | HiraganaToNarrow | ComposeVom,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"ｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞ"},
	89: {KatakanaToNarrow | HiraganaToNarrow | ComposeVom,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099"},

	// simple Japanese conversion (voiced kana traditional directive) <- Voiced sound character/Voiced sound mark
	90: {KatakanaToHiragana | DecomposeVom,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099"},
	91: {KatakanaToHiragana | DecomposeVom,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099"},
	92: {KatakanaToHiragana | DecomposeVom,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099"},

	// multiple Japanese conversion (voiced kana traditional directive) <- Voiced sound character/Voiced sound mark
	93: {KatakanaToNarrow | HiraganaToNarrow | DecomposeVom,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099"},
	94: {KatakanaToNarrow | HiraganaToNarrow | DecomposeVom,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099"},
	95: {KatakanaToNarrow | HiraganaToNarrow | DecomposeVom,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099"},

	// simple Japanese conversion (isolated vsm directive) <- Voiced sound character/Voiced sound mark
	96: {KatakanaToHiragana | IsolatedKanaVomToNarrow,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"がかﾞかﾞかﾞがかﾞかﾞかﾞかﾞかﾞかﾞ"},
	97: {KatakanaToHiragana | IsolatedKanaVomToNarrow,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"あﾞあﾞあﾞあﾞあﾞあﾞあﾞあﾞあﾞ"},
	98: {KatakanaToHiragana | IsolatedKanaVomToNarrow,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"AﾞAﾞAﾞ日ﾞ日ﾞ日ﾞäﾞäﾞäﾞ"},
	99: {KatakanaToHiragana | IsolatedKanaVomToWide,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"がか゛か゛か゛がか゛か゛か゛か゛か゛か゛"},
	100: {KatakanaToHiragana | IsolatedKanaVomToWide,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"あ゛あ゛あ゛あ゛あ゛あ゛あ゛あ゛あ゛"},
	101: {KatakanaToHiragana | IsolatedKanaVomToWide,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A゛A゛A゛日゛日゛日゛ä゛ä゛ä゛"},
	102: {KatakanaToHiragana | IsolatedKanaVomToNonspace,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"がか\u3099か\u3099か\u3099がか\u3099か\u3099か\u3099か\u3099か\u3099か\u3099"},
	103: {KatakanaToHiragana | IsolatedKanaVomToNonspace,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099"},
	104: {KatakanaToHiragana | IsolatedKanaVomToNonspace,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A\u3099A\u3099A\u3099日\u3099日\u3099日\u3099ä\u3099ä\u3099ä\u3099"},

	// multiple Japanese conversion (isolated vsm directive) <- Voiced sound character/Voiced sound mark
	105: {KatakanaToNarrow | HiraganaToNarrow | IsolatedKanaVomToNarrow,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞ"},
	106: {KatakanaToNarrow | HiraganaToNarrow | IsolatedKanaVomToNarrow,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"ｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞ"},
	107: {KatakanaToNarrow | HiraganaToNarrow | IsolatedKanaVomToNarrow,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"AﾞAﾞAﾞ日ﾞ日ﾞ日ﾞäﾞäﾞäﾞ"},
	108: {KatakanaToNarrow | HiraganaToNarrow | IsolatedKanaVomToWide,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ｶﾞｶ゛ｶ゛ｶ゛ｶﾞｶ゛ｶ゛ｶ゛ｶ゛ｶ゛ｶ゛"},
	109: {KatakanaToNarrow | HiraganaToNarrow | IsolatedKanaVomToWide,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"ｱ゛ｱ゛ｱ゛ｱ゛ｱ゛ｱ゛ｱ゛ｱ゛ｱ゛"},
	110: {KatakanaToNarrow | HiraganaToNarrow | IsolatedKanaVomToWide,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A゛A゛A゛日゛日゛日゛ä゛ä゛ä゛"},
	111: {KatakanaToNarrow | HiraganaToNarrow | IsolatedKanaVomToNonspace,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ｶﾞｶ\u3099ｶ\u3099ｶ\u3099ｶﾞｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099"},
	112: {KatakanaToNarrow | HiraganaToNarrow | IsolatedKanaVomToNonspace,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099"},
	113: {KatakanaToNarrow | HiraganaToNarrow | IsolatedKanaVomToNonspace,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A\u3099A\u3099A\u3099日\u3099日\u3099日\u3099ä\u3099ä\u3099ä\u3099"},

	// simple Japanese conversion (voiced kana traditional directive/isolated vsm directive) <- Voiced sound character/Voiced sound mark
	114: {KatakanaToHiragana | ComposeVom | IsolatedKanaVomToNarrow,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ががががががががががが"},
	115: {KatakanaToHiragana | ComposeVom | IsolatedKanaVomToNarrow,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"あ゛あ゛あ゛あ゛あ゛あ゛あ゛あ゛あ゛"},
	116: {KatakanaToHiragana | ComposeVom | IsolatedKanaVomToNarrow,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"AﾞAﾞAﾞ日ﾞ日ﾞ日ﾞäﾞäﾞäﾞ"},
	117: {KatakanaToHiragana | ComposeVom | IsolatedKanaVomToWide,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ががががががががががが"},
	118: {KatakanaToHiragana | ComposeVom | IsolatedKanaVomToWide,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"あ゛あ゛あ゛あ゛あ゛あ゛あ゛あ゛あ゛"},
	119: {KatakanaToHiragana | ComposeVom | IsolatedKanaVomToWide,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A゛A゛A゛日゛日゛日゛ä゛ä゛ä゛"},
	120: {KatakanaToHiragana | ComposeVom | IsolatedKanaVomToNonspace,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ががががががががががが"},
	121: {KatakanaToHiragana | ComposeVom | IsolatedKanaVomToNonspace,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"あ゛あ゛あ゛あ゛あ゛あ゛あ゛あ゛あ゛"},
	122: {KatakanaToHiragana | ComposeVom | IsolatedKanaVomToNonspace,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A\u3099A\u3099A\u3099日\u3099日\u3099日\u3099ä\u3099ä\u3099ä\u3099"},

	// multiple Japanese conversion (voiced kana traditional directive/isolated vsm directive) <- Voiced sound character/Voiced sound mark
	123: {KatakanaToNarrow | HiraganaToNarrow | ComposeVom | IsolatedKanaVomToNarrow,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞ"},
	124: {KatakanaToNarrow | HiraganaToNarrow | ComposeVom | IsolatedKanaVomToNarrow,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"ｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞ"},
	125: {KatakanaToNarrow | HiraganaToNarrow | ComposeVom | IsolatedKanaVomToNarrow,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"AﾞAﾞAﾞ日ﾞ日ﾞ日ﾞäﾞäﾞäﾞ"},
	126: {KatakanaToNarrow | HiraganaToNarrow | ComposeVom | IsolatedKanaVomToWide,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞ"},
	127: {KatakanaToNarrow | HiraganaToNarrow | ComposeVom | IsolatedKanaVomToWide,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"ｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞ"},
	128: {KatakanaToNarrow | HiraganaToNarrow | ComposeVom | IsolatedKanaVomToWide,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A゛A゛A゛日゛日゛日゛ä゛ä゛ä゛"},
	129: {KatakanaToNarrow | HiraganaToNarrow | ComposeVom | IsolatedKanaVomToNonspace,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞ"},
	130: {KatakanaToNarrow | HiraganaToNarrow | ComposeVom | IsolatedKanaVomToNonspace,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"ｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞ"},
	131: {KatakanaToNarrow | HiraganaToNarrow | ComposeVom | IsolatedKanaVomToNonspace,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A\u3099A\u3099A\u3099日\u3099日\u3099日\u3099ä\u3099ä\u3099ä\u3099"},

	// simple Japanese conversion (voiced kana combining directive/isolated vsm directive) <- Voiced sound character/Voiced sound mark
	132: {KatakanaToHiragana | DecomposeVom | IsolatedKanaVomToNarrow,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099"},
	133: {KatakanaToHiragana | DecomposeVom | IsolatedKanaVomToNarrow,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099"},
	134: {KatakanaToHiragana | DecomposeVom | IsolatedKanaVomToNarrow,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"AﾞAﾞAﾞ日ﾞ日ﾞ日ﾞäﾞäﾞäﾞ"},
	135: {KatakanaToHiragana | DecomposeVom | IsolatedKanaVomToWide,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099"},
	136: {KatakanaToHiragana | DecomposeVom | IsolatedKanaVomToWide,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099"},
	137: {KatakanaToHiragana | DecomposeVom | IsolatedKanaVomToWide,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A゛A゛A゛日゛日゛日゛ä゛ä゛ä゛"},
	138: {KatakanaToHiragana | DecomposeVom | IsolatedKanaVomToNonspace,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099"},
	139: {KatakanaToHiragana | DecomposeVom | IsolatedKanaVomToNonspace,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099"},
	140: {KatakanaToHiragana | DecomposeVom | IsolatedKanaVomToNonspace,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A\u3099A\u3099A\u3099日\u3099日\u3099日\u3099ä\u3099ä\u3099ä\u3099"},

	// multiple Japanese conversion (voiced kana combining directive/isolated vsm directive) <- Voiced sound character/Voiced sound mark
	141: {KatakanaToNarrow | HiraganaToNarrow | DecomposeVom | IsolatedKanaVomToNarrow,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099"},
	142: {KatakanaToNarrow | HiraganaToNarrow | DecomposeVom | IsolatedKanaVomToNarrow,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099"},
	143: {KatakanaToNarrow | HiraganaToNarrow | DecomposeVom | IsolatedKanaVomToNarrow,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"AﾞAﾞAﾞ日ﾞ日ﾞ日ﾞäﾞäﾞäﾞ"},
	144: {KatakanaToNarrow | HiraganaToNarrow | DecomposeVom | IsolatedKanaVomToWide,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099"},
	145: {KatakanaToNarrow | HiraganaToNarrow | DecomposeVom | IsolatedKanaVomToWide,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099"},
	146: {KatakanaToNarrow | HiraganaToNarrow | DecomposeVom | IsolatedKanaVomToWide,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A゛A゛A゛日゛日゛日゛ä゛ä゛ä゛"},
	147: {KatakanaToNarrow | HiraganaToNarrow | DecomposeVom | IsolatedKanaVomToNonspace,
		"がか゛かﾞか\u3099ガカ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099"},
	148: {KatakanaToNarrow | HiraganaToNarrow | DecomposeVom | IsolatedKanaVomToNonspace,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099"},
	149: {KatakanaToNarrow | HiraganaToNarrow | DecomposeVom | IsolatedKanaVomToNonspace,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A\u3099A\u3099A\u3099日\u3099日\u3099日\u3099ä\u3099ä\u3099ä\u3099"},

	// VSM/SVSM testing that cannot be combined with the previous character
	150: {ComposeVom,
		"は゛ば゛ぱ゛は゜ば゜ぱ゜ハ゛バ゛パ゛ハ゜バ゜パ゜ﾊ゛ﾊ゜",
		"ばば゛ぱ゛ぱば゜ぱ゜ババ゛パ゛パバ゜パ゜ﾊﾞﾊﾟ"},
	151: {ComposeVom | IsolatedKanaVomToNarrow,
		"は゛ば゛ぱ゛は゜ば゜ぱ゜ハ゛バ゛パ゛ハ゜バ゜パ゜ﾊ゛ﾊ゜",
		"ばばﾞぱﾞぱばﾟぱﾟババﾞパﾞパバﾟパﾟﾊﾞﾊﾟ"},
	152: {ComposeVom | IsolatedKanaVomToWide,
		"は゛ば゛ぱ゛は゜ば゜ぱ゜ハ゛バ゛パ゛ハ゜バ゜パ゜ﾊ゛ﾊ゜",
		"ばば゛ぱ゛ぱば゜ぱ゜ババ゛パ゛パバ゜パ゜ﾊﾞﾊﾟ"},
	153: {ComposeVom | IsolatedKanaVomToNonspace,
		"は゛ば゛ぱ゛は゜ば゜ぱ゜ハ゛バ゛パ゛ハ゜バ゜パ゜ﾊ゛ﾊ゜",
		"ばば\u3099ぱ\u3099ぱば\u309Aぱ\u309Aババ\u3099パ\u3099パバ\u309Aパ\u309Aﾊﾞﾊﾟ"},
	154: {DecomposeVom,
		"は゛ば゛ぱ゛は゜ば゜ぱ゜ハ゛バ゛パ゛ハ゜バ゜パ゜ﾊ゛ﾊ゜",
		"は\u3099は\u3099゛は\u309A゛は\u309Aは\u3099゜は\u309A゜ハ\u3099ハ\u3099゛ハ\u309A゛ハ\u309Aハ\u3099゜ハ\u309A゜ﾊ\u3099ﾊ\u309A"},
	155: {DecomposeVom | IsolatedKanaVomToNarrow,
		"は゛ば゛ぱ゛は゜ば゜ぱ゜ハ゛バ゛パ゛ハ゜バ゜パ゜ﾊ゛ﾊ゜",
		"は\u3099は\u3099ﾞは\u309Aﾞは\u309Aは\u3099ﾟは\u309Aﾟハ\u3099ハ\u3099ﾞハ\u309Aﾞハ\u309Aハ\u3099ﾟハ\u309Aﾟﾊ\u3099ﾊ\u309A"},
	156: {DecomposeVom | IsolatedKanaVomToWide,
		"は゛ば゛ぱ゛は゜ば゜ぱ゜ハ゛バ゛パ゛ハ゜バ゜パ゜ﾊ゛ﾊ゜",
		"は\u3099は\u3099゛は\u309A゛は\u309Aは\u3099゜は\u309A゜ハ\u3099ハ\u3099゛ハ\u309A゛ハ\u309Aハ\u3099゜ハ\u309A゜ﾊ\u3099ﾊ\u309A"},
	157: {DecomposeVom | IsolatedKanaVomToNonspace,
		"は゛ば゛ぱ゛は゜ば゜ぱ゜ハ゛バ゛パ゛ハ゜バ゜パ゜ﾊ゛ﾊ゜",
		"は\u3099は\u3099\u3099は\u309A\u3099は\u309Aは\u3099\u309Aは\u309A\u309Aハ\u3099ハ\u3099\u3099ハ\u309A\u3099ハ\u309Aハ\u3099\u309Aハ\u309A\u309Aﾊ\u3099ﾊ\u309A"},
	158: {ComposeVom,
		"はﾞばﾞぱﾞはﾟばﾟぱﾟハﾞバﾞパﾞハﾟバﾟパﾟﾊﾞﾊﾟ",
		"ばばﾞぱﾞぱばﾟぱﾟババﾞパﾞパバﾟパﾟﾊﾞﾊﾟ"},
	159: {ComposeVom | IsolatedKanaVomToNarrow,
		"はﾞばﾞぱﾞはﾟばﾟぱﾟハﾞバﾞパﾞハﾟバﾟパﾟﾊﾞﾊﾟ",
		"ばばﾞぱﾞぱばﾟぱﾟババﾞパﾞパバﾟパﾟﾊﾞﾊﾟ"},
	160: {ComposeVom | IsolatedKanaVomToWide,
		"はﾞばﾞぱﾞはﾟばﾟぱﾟハﾞバﾞパﾞハﾟバﾟパﾟﾊﾞﾊﾟ",
		"ばば゛ぱ゛ぱば゜ぱ゜ババ゛パ゛パバ゜パ゜ﾊﾞﾊﾟ"},
	161: {ComposeVom | IsolatedKanaVomToNonspace,
		"はﾞばﾞぱﾞはﾟばﾟぱﾟハﾞバﾞパﾞハﾟバﾟパﾟﾊﾞﾊﾟ",
		"ばば\u3099ぱ\u3099ぱば\u309Aぱ\u309Aババ\u3099パ\u3099パバ\u309Aパ\u309Aﾊﾞﾊﾟ"},
	162: {DecomposeVom,
		"はﾞばﾞぱﾞはﾟばﾟぱﾟハﾞバﾞパﾞハﾟバﾟパﾟﾊﾞﾊﾟ",
		"は\u3099は\u3099ﾞは\u309Aﾞは\u309Aは\u3099ﾟは\u309Aﾟハ\u3099ハ\u3099ﾞハ\u309Aﾞハ\u309Aハ\u3099ﾟハ\u309Aﾟﾊ\u3099ﾊ\u309A"},
	163: {DecomposeVom | IsolatedKanaVomToNarrow,
		"はﾞばﾞぱﾞはﾟばﾟぱﾟハﾞバﾞパﾞハﾟバﾟパﾟﾊﾞﾊﾟ",
		"は\u3099は\u3099ﾞは\u309Aﾞは\u309Aは\u3099ﾟは\u309Aﾟハ\u3099ハ\u3099ﾞハ\u309Aﾞハ\u309Aハ\u3099ﾟハ\u309Aﾟﾊ\u3099ﾊ\u309A"},
	164: {DecomposeVom | IsolatedKanaVomToWide,
		"はﾞばﾞぱﾞはﾟばﾟぱﾟハﾞバﾞパﾞハﾟバﾟパﾟﾊﾞﾊﾟ",
		"は\u3099は\u3099゛は\u309A゛は\u309Aは\u3099゜は\u309A゜ハ\u3099ハ\u3099゛ハ\u309A゛ハ\u309Aハ\u3099゜ハ\u309A゜ﾊ\u3099ﾊ\u309A"},
	165: {DecomposeVom | IsolatedKanaVomToNonspace,
		"はﾞばﾞぱﾞはﾟばﾟぱﾟハﾞバﾞパﾞハﾟバﾟパﾟﾊﾞﾊﾟ",
		"は\u3099は\u3099\u3099は\u309A\u3099は\u309Aは\u3099\u309Aは\u309A\u309Aハ\u3099ハ\u3099\u3099ハ\u309A\u3099ハ\u309Aハ\u3099\u309Aハ\u309A\u309Aﾊ\u3099ﾊ\u309A"},

	// VSM/SVSM testing that cannot be combined with the previous character (out of range in unichars table)
	166: {ComposeVom,
		"日゛本゜語ﾞ平ﾟ仮\u3099名\u309A",
		"日゛本゜語ﾞ平ﾟ仮\u3099名\u309A"},
	167: {ComposeVom | IsolatedKanaVomToNarrow,
		"日゛本゜語ﾞ平ﾟ仮\u3099名\u309A",
		"日ﾞ本ﾟ語ﾞ平ﾟ仮ﾞ名ﾟ"}, // TEST_N9x6dneg
	168: {ComposeVom | IsolatedKanaVomToWide,
		"日゛本゜語ﾞ平ﾟ仮\u3099名\u309A",
		"日゛本゜語゛平゜仮゛名゜"},
	169: {ComposeVom | IsolatedKanaVomToNonspace,
		"日゛本゜語ﾞ平ﾟ仮\u3099名\u309A",
		"日\u3099本\u309A語\u3099平\u309A仮\u3099名\u309A"},
	170: {DecomposeVom,
		"日゛本゜語ﾞ平ﾟ仮\u3099名\u309A",
		"日゛本゜語ﾞ平ﾟ仮\u3099名\u309A"},
	171: {DecomposeVom | IsolatedKanaVomToNarrow,
		"日゛本゜語ﾞ平ﾟ仮\u3099名\u309A",
		"日ﾞ本ﾟ語ﾞ平ﾟ仮ﾞ名ﾟ"},
	172: {DecomposeVom | IsolatedKanaVomToWide,
		"日゛本゜語ﾞ平ﾟ仮\u3099名\u309A",
		"日゛本゜語゛平゜仮゛名゜"},
	173: {DecomposeVom | IsolatedKanaVomToNonspace,
		"日゛本゜語ﾞ平ﾟ仮\u3099名\u309A",
		"日\u3099本\u309A語\u3099平\u309A仮\u3099名\u309A"},

	// VSM testing with or without ComposeVom/IsolatedKanaVomToNarrow flags
	174: {ComposeVom,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ｺﾞｺﾞｺﾞゴゴゴゴ"},
	175: {IsolatedKanaVomToNarrow,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ｺﾞｺﾞｺﾞコﾞコﾞコﾞゴ"},
	176: {ComposeVom | IsolatedKanaVomToNarrow,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ｺﾞｺﾞｺﾞゴゴゴゴ"},
	177: {KatakanaToHiragana,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"こﾞこ゛こ\u3099こﾞこ゛こ\u3099ご"},
	178: {KatakanaToHiragana | ComposeVom,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ごごごごごごご"},
	179: {KatakanaToHiragana | IsolatedKanaVomToNarrow,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"こﾞこﾞこﾞこﾞこﾞこﾞご"},
	180: {KatakanaToHiragana | ComposeVom | IsolatedKanaVomToNarrow,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ごごごごごごご"},
	181: {KatakanaToWide,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"コﾞコ゛コ\u3099コﾞコ゛コ\u3099ゴ"},
	182: {KatakanaToWide | ComposeVom,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ゴゴゴゴゴゴゴ"},
	183: {KatakanaToWide | IsolatedKanaVomToNarrow,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"コﾞコﾞコﾞコﾞコﾞコﾞゴ"},
	184: {KatakanaToWide | ComposeVom | IsolatedKanaVomToNarrow,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ゴゴゴゴゴゴゴ"},
	185: {KatakanaToNarrow,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ｺﾞｺ゛ｺ\u3099ｺﾞｺ゛ｺ\u3099ｺﾞ"}, // TEST_L7tADs2z
	186: {KatakanaToNarrow | ComposeVom,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞ"},
	187: {KatakanaToNarrow | IsolatedKanaVomToNarrow,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞ"},
	188: {KatakanaToNarrow | ComposeVom | IsolatedKanaVomToNarrow,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞ"},

	// SVSM testing with or without ComposeVom/IsolatedKanaVomToNarrow flags
	189: {ComposeVom,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ﾎﾟﾎﾟﾎﾟポポポポ"},
	190: {IsolatedKanaVomToNarrow,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ﾎﾟﾎﾟﾎﾟホﾟホﾟホﾟポ"},
	191: {ComposeVom | IsolatedKanaVomToNarrow,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ﾎﾟﾎﾟﾎﾟポポポポ"},
	192: {KatakanaToHiragana,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ほﾟほ゜ほ\u309Aほﾟほ゜ほ\u309Aぽ"},
	193: {KatakanaToHiragana | ComposeVom,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ぽぽぽぽぽぽぽ"},
	194: {KatakanaToHiragana | IsolatedKanaVomToNarrow,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ほﾟほﾟほﾟほﾟほﾟほﾟぽ"},
	195: {KatakanaToHiragana | ComposeVom | IsolatedKanaVomToNarrow,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ぽぽぽぽぽぽぽ"},
	196: {KatakanaToWide,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ホﾟホ゜ホ\u309Aホﾟホ゜ホ\u309Aポ"},
	197: {KatakanaToWide | ComposeVom,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ポポポポポポポ"},
	198: {KatakanaToWide | IsolatedKanaVomToNarrow,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ホﾟホﾟホﾟホﾟホﾟホﾟポ"},
	199: {KatakanaToWide | ComposeVom | IsolatedKanaVomToNarrow,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ポポポポポポポ"},
	200: {KatakanaToNarrow,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ﾎﾟﾎ゜ﾎ\u309Aﾎﾟﾎ゜ﾎ\u309Aﾎﾟ"}, // TEST_K6t8hQYp
	201: {KatakanaToNarrow | ComposeVom,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ﾎﾟﾎﾟﾎﾟﾎﾟﾎﾟﾎﾟﾎﾟ"},
	202: {KatakanaToNarrow | IsolatedKanaVomToNarrow,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ﾎﾟﾎﾟﾎﾟﾎﾟﾎﾟﾎﾟﾎﾟ"},
	203: {KatakanaToNarrow | ComposeVom | IsolatedKanaVomToNarrow,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ﾎﾟﾎﾟﾎﾟﾎﾟﾎﾟﾎﾟﾎﾟ"},

	// VSM/SVSM testing illegal rune value
	204: {ComposeVom,
		string([]rune{excr, '゛', excr, 'ﾞ', excr, '\u3099', excr, '゜', excr, 'ﾟ', excr, '\u309A'}),
		string([]rune{excr, '゛', excr, 'ﾞ', excr, '\u3099', excr, '゜', excr, 'ﾟ', excr, '\u309A'})},
	205: {ComposeVom | IsolatedKanaVomToNarrow,
		string([]rune{excr, '゛', excr, 'ﾞ', excr, '\u3099', excr, '゜', excr, 'ﾟ', excr, '\u309A'}),
		string([]rune{excr, 'ﾞ', excr, 'ﾞ', excr, 'ﾞ', excr, 'ﾟ', excr, 'ﾟ', excr, 'ﾟ'})},
	206: {ComposeVom | IsolatedKanaVomToWide,
		string([]rune{excr, '゛', excr, 'ﾞ', excr, '\u3099', excr, '゜', excr, 'ﾟ', excr, '\u309A'}),
		string([]rune{excr, '゛', excr, '゛', excr, '゛', excr, '゜', excr, '゜', excr, '゜'})},
	207: {ComposeVom | IsolatedKanaVomToNonspace,
		string([]rune{excr, '゛', excr, 'ﾞ', excr, '\u3099', excr, '゜', excr, 'ﾟ', excr, '\u309A'}),
		string([]rune{excr, '\u3099', excr, '\u3099', excr, '\u3099', excr, '\u309A', excr, '\u309A', excr, '\u309A'})},
	208: {DecomposeVom,
		string([]rune{excr, '゛', excr, 'ﾞ', excr, '\u3099', excr, '゜', excr, 'ﾟ', excr, '\u309A'}),
		string([]rune{excr, '゛', excr, 'ﾞ', excr, '\u3099', excr, '゜', excr, 'ﾟ', excr, '\u309A'})},
	209: {DecomposeVom | IsolatedKanaVomToNarrow,
		string([]rune{excr, '゛', excr, 'ﾞ', excr, '\u3099', excr, '゜', excr, 'ﾟ', excr, '\u309A'}),
		string([]rune{excr, 'ﾞ', excr, 'ﾞ', excr, 'ﾞ', excr, 'ﾟ', excr, 'ﾟ', excr, 'ﾟ'})},
	210: {DecomposeVom | IsolatedKanaVomToWide,
		string([]rune{excr, '゛', excr, 'ﾞ', excr, '\u3099', excr, '゜', excr, 'ﾟ', excr, '\u309A'}),
		string([]rune{excr, '゛', excr, '゛', excr, '゛', excr, '゜', excr, '゜', excr, '゜'})},
	211: {DecomposeVom | IsolatedKanaVomToNonspace,
		string([]rune{excr, '゛', excr, 'ﾞ', excr, '\u3099', excr, '゜', excr, 'ﾟ', excr, '\u309A'}),
		string([]rune{excr, '\u3099', excr, '\u3099', excr, '\u3099', excr, '\u309A', excr, '\u309A', excr, '\u309A'})},

	// special Katakana-Hiragana letters
	212: {KatakanaToHiragana | ComposeVom,
		"アイウエオヤユヨツワカケァィゥェォャュョッヮヵヶヷヸヹヺ",
		"あいうえおやゆよつわかけぁぃぅぇぉゃゅょっゎゕゖわ゛ゐ゛ゑ゛を゛"},
	213: {KatakanaToNarrow | ComposeVom,
		"アイウエオヤユヨツワカケァィゥェォャュョッヮヵヶヷヸヹヺ",
		"ｱｲｳｴｵﾔﾕﾖﾂﾜｶｹｧｨｩｪｫｬｭｮｯﾜｶｹﾜﾞｲﾞｴﾞｦﾞ"},
	214: {HiraganaToKatakana | ComposeVom,
		"あいうえおやゆよつわかけぁぃぅぇぉゃゅょっゎゕゖわ゛ゐ゛ゑ゛を゛",
		"アイウエオヤユヨツワカケァィゥェォャュョッヮヵヶヷヸヹヺ"},
	215: {HiraganaToNarrow | ComposeVom,
		"あいうえおやゆよつわかけぁぃぅぇぉゃゅょっゎゕゖわ゛ゐ゛ゑ゛を゛",
		"ｱｲｳｴｵﾔﾕﾖﾂﾜｶｹｧｨｩｪｫｬｭｮｯﾜｶｹﾜﾞｲﾞｴﾞｦﾞ"},
	216: {KatakanaToHiragana | ComposeVom,
		"ｱｲｳｴｵﾔﾕﾖﾂﾜｶｹｧｨｩｪｫｬｭｮｯﾜｶｹﾜﾞｲﾞｴﾞｦﾞ",
		"あいうえおやゆよつわかけぁぃぅぇぉゃゅょっわかけわ゛い゛え゛を゛"},
	217: {KatakanaToWide | ComposeVom,
		"ｱｲｳｴｵﾔﾕﾖﾂﾜｶｹｧｨｩｪｫｬｭｮｯﾜｶｹﾜﾞｲﾞｴﾞｦﾞ",
		"アイウエオヤユヨツワカケァィゥェォャュョッワカケヷイ゛エ゛ヺ"},

	// Hiragana letter YORI, Katakana letter KOTO
	218: {KatakanaToHiragana | ComposeVom, "ゟヿ", "ゟヿ"},
	219: {KatakanaToNarrow | ComposeVom, "ゟヿ", "ゟヿ"},
	220: {KatakanaToWide | ComposeVom, "ゟヿ", "ゟヿ"},
	221: {HiraganaToKatakana | ComposeVom, "ゟヿ", "ゟヿ"},
	222: {HiraganaToNarrow | ComposeVom, "ゟヿ", "ゟヿ"},

	// Katakana Phonetic Extensions
	223: {KatakanaToHiragana,
		"ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿㇰ゛ㇱ゛ㇲ゛ㇳ゛ㇴ゛ㇵ゛ㇶ゛ㇷ゛ㇸ゛ㇹ゛ㇺ゛ㇻ゛ㇼ゛ㇽ゛ㇾ゛ㇿ゛",
		"くしすとぬはひふへほむらりるれろく゛し゛す゛と゛ぬ゛は゛ひ゛ふ゛へ゛ほ゛む゛ら゛り゛る゛れ゛ろ゛"},
	224: {KatakanaToHiragana | ComposeVom,
		"ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿㇰ゛ㇱ゛ㇲ゛ㇳ゛ㇴ゛ㇵ゛ㇶ゛ㇷ゛ㇸ゛ㇹ゛ㇺ゛ㇻ゛ㇼ゛ㇽ゛ㇾ゛ㇿ゛",
		"くしすとぬはひふへほむらりるれろぐじずどぬ゛ばびぶべぼむ゛ら゛り゛る゛れ゛ろ゛"},
	225: {KatakanaToNarrow,
		"ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿㇰ゛ㇱ゛ㇲ゛ㇳ゛ㇴ゛ㇵ゛ㇶ゛ㇷ゛ㇸ゛ㇹ゛ㇺ゛ㇻ゛ㇼ゛ㇽ゛ㇾ゛ㇿ゛",
		"ｸｼｽﾄﾇﾊﾋﾌﾍﾎﾑﾗﾘﾙﾚﾛｸ゛ｼ゛ｽ゛ﾄ゛ﾇ゛ﾊ゛ﾋ゛ﾌ゛ﾍ゛ﾎ゛ﾑ゛ﾗ゛ﾘ゛ﾙ゛ﾚ゛ﾛ゛"},
	226: {KatakanaToNarrow | ComposeVom,
		"ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿㇰ゛ㇱ゛ㇲ゛ㇳ゛ㇴ゛ㇵ゛ㇶ゛ㇷ゛ㇸ゛ㇹ゛ㇺ゛ㇻ゛ㇼ゛ㇽ゛ㇾ゛ㇿ゛",
		"ｸｼｽﾄﾇﾊﾋﾌﾍﾎﾑﾗﾘﾙﾚﾛｸﾞｼﾞｽﾞﾄﾞﾇﾞﾊﾞﾋﾞﾌﾞﾍﾞﾎﾞﾑﾞﾗﾞﾘﾞﾙﾞﾚﾞﾛﾞ"},
	227: {KatakanaToHiragana,
		"ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿㇰ゜ㇱ゜ㇲ゜ㇳ゜ㇴ゜ㇵ゜ㇶ゜ㇷ゜ㇸ゜ㇹ゜ㇺ゜ㇻ゜ㇼ゜ㇽ゜ㇾ゜ㇿ゜",
		"くしすとぬはひふへほむらりるれろく゜し゜す゜と゜ぬ゜は゜ひ゜ふ゜へ゜ほ゜む゜ら゜り゜る゜れ゜ろ゜"},
	228: {KatakanaToHiragana | ComposeVom,
		"ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿㇰ゜ㇱ゜ㇲ゜ㇳ゜ㇴ゜ㇵ゜ㇶ゜ㇷ゜ㇸ゜ㇹ゜ㇺ゜ㇻ゜ㇼ゜ㇽ゜ㇾ゜ㇿ゜",
		"くしすとぬはひふへほむらりるれろく゜し゜す゜と゜ぬ゜ぱぴぷぺぽむ゜ら゜り゜る゜れ゜ろ゜"},
	229: {KatakanaToNarrow,
		"ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿㇰ゜ㇱ゜ㇲ゜ㇳ゜ㇴ゜ㇵ゜ㇶ゜ㇷ゜ㇸ゜ㇹ゜ㇺ゜ㇻ゜ㇼ゜ㇽ゜ㇾ゜ㇿ゜",
		"ｸｼｽﾄﾇﾊﾋﾌﾍﾎﾑﾗﾘﾙﾚﾛｸ゜ｼ゜ｽ゜ﾄ゜ﾇ゜ﾊ゜ﾋ゜ﾌ゜ﾍ゜ﾎ゜ﾑ゜ﾗ゜ﾘ゜ﾙ゜ﾚ゜ﾛ゜"},
	230: {KatakanaToNarrow | ComposeVom,
		"ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿㇰ゜ㇱ゜ㇲ゜ㇳ゜ㇴ゜ㇵ゜ㇶ゜ㇷ゜ㇸ゜ㇹ゜ㇺ゜ㇻ゜ㇼ゜ㇽ゜ㇾ゜ㇿ゜",
		"ｸｼｽﾄﾇﾊﾋﾌﾍﾎﾑﾗﾘﾙﾚﾛｸﾟｼﾟｽﾟﾄﾟﾇﾟﾊﾟﾋﾟﾌﾟﾍﾟﾎﾟﾑﾟﾗﾟﾘﾟﾙﾟﾚﾟﾛﾟ"},

	// the Yen mark
	231: {SymbolToWide, "\\￥", "＼￥"},
	232: {SymbolToNarrow, "＼￥", "\\￥"},

	// overflow
	233: {ComposeVom,
		string([]rune{-1, '\u0000', maxr, excr}),
		string([]rune{-1, '\u0000', maxr, excr})},
	234: {DecomposeVom,
		string([]rune{-1, '\u0000', maxr, excr}),
		string([]rune{-1, '\u0000', maxr, excr})},
	235: {LatinToNarrow | KanaToWide,
		string([]rune{-1, '\u0000', maxr, excr}),
		string([]rune{-1, '\u0000', maxr, excr})},

	// the whitespace
	236: {SymbolToWide, "\u0020\u3000", "\u3000\u3000"},
	237: {SymbolToNarrow, "\u0020\u3000", "\u0020\u0020"},

	// the wave dash and the characters similar to it
	238: {SymbolToWide, "~\uFF5E\u301C\U0001301C", "\uFF5E\uFF5E\u301C\U0001301C"},
	239: {SymbolToNarrow, "~\uFF5E\u301C\U0001301C", "~~\u301C\U0001301C"},

	// the hyphen-minus and the characters similar to it
	240: {SymbolToWide, "-\uFF0D\uFF70\u30FC", "\uFF0D\uFF0D\uFF70\u30FC"},
	241: {SymbolToNarrow, "-\uFF0D\uFF70\u30FC", "--\uFF70\u30FC"},
	242: {KatakanaToWide, "-\uFF0D\uFF70\u30FC", "-\uFF0D\u30FC\u30FC"},
	243: {KatakanaToNarrow, "-\uFF0D\uFF70\u30FC", "-\uFF0D\uFF70\uFF70"},

	// ctUndefined
	244: {Fold, "\u3040\u3097\u3098\uFF00", "\u3040\u3097\u3098\uFF00"},

	// user perspective testing
	245: {Fold, "Ｇo言語のﾊﾟｯケｰｼﾞ (Ｐａｃｋａｇｅ）", "Go言語のパッケージ (Package)"},
	246: {KanaToHiragana, "ふりがな | ｽｽﾞｷ イチロウ", "ふりがな | すずき いちろう"},
	247: {KanaToWideKatakana, "フリガナ | すす゛き ｲﾁﾛｰ", "フリガナ | スズキ イチロー"},
	248: {KanaToNarrowKatakana, "ﾌﾘｶﾞﾅ | スズキ いちろう", "ﾌﾘｶﾞﾅ | ｽｽﾞｷ ｲﾁﾛｳ"},
	249: {LatinToNarrow|AlphaToUpper, "ローマ字(半角) | Ｓｕｚｕｋｉ, ichiro", "ローマ字(半角) | SUZUKI, ICHIRO"},
}

func hexs(s string) string {
	var ss []string
	for _, r := range []rune(s) {
		e := fmt.Sprintf("%04X", r)
		ss = append(ss, e)
	}
	switch len(ss) {
	case 0:
		return "<empty>"
	case 1:
		return "[" + ss[0] + "]"
	default:
		return "[" + strings.Join(ss, " ") + "]"
	}
}

func TestNormalizer_String(t *testing.T) {
	for i, tt := range normalizer_stringtests {
		n, err := Norm(tt.flag)
		if err != nil {
			t.Errorf("#%d: %s", i, err.Error())
			continue
		}
		out := n.String(tt.in)
		if out != tt.out {
			t.Errorf("TestNormalize #%d\n\tflag:"+
				"\t%s, \n\targs:\t%q\n\thave:\t%q\n\twant:\t%q\n"+
				"\targs16:\t%s\n\thave16:\t%s\n\twant16:\t%s",
				i, tt.flag, tt.in, out, tt.out,
				hexs(tt.in), hexs(out), hexs(tt.out))
		}
	}
}

type NormTest struct {
	flag NormFlag
	errS string
}

var normtests = []NormTest{
	{KatakanaToWide | KatakanaToNarrow, "invalid normalization flag"},
	{KatakanaToWide, ""},
}

func TestNorm(t *testing.T) {
	for i, tt := range normtests {
		_, err := Norm(tt.flag)
		if tt.errS == "" && err != nil {
			t.Errorf("#%d have error: %s, want no error", i, err.Error())
			continue
		}
		if tt.errS != "" {
			if err == nil {
				t.Errorf("#%d have no error, want error:%s", i, tt.errS)
				continue
			}
			if !strings.Contains(err.Error(), tt.errS) {
				t.Errorf("#%d have error: %s, want error: %s", i, err.Error(), tt.errS)
			}
		}
	}
}

func TestNormalizer_SetFlag(t *testing.T) {
	for i, tt := range normtests {
		n, err := Norm(AlphaToNarrow)
		if err != nil {
			log.Fatalf("unexpectedly error: %s", err.Error())
			continue
		}

		err = n.SetFlag(tt.flag)
		if tt.errS == "" && err != nil {
			t.Errorf("#%d have error: %s, want no error", i, err.Error())
			continue
		}
		if tt.errS != "" {
			if err == nil {
				t.Errorf("#%d have no error, want error:%s", i, tt.errS)
				continue
			}
			if !strings.Contains(err.Error(), tt.errS) {
				t.Errorf("#%d have error: %s, want error: %s", i, err.Error(), tt.errS)
			}
		}
	}
}

func normflags() []int {
	flags := make([]int, 0, len(normflagMap))
	for key := range normflagMap {
		flags = append(flags, int(key))
	}
	return flags
}

// nCm combinations algorithm
func comb(arr []int, n int) (result [][]int) {
	if n <= 0 || n > len(arr) {
		return result
	}
	if n == 1 {
		for _, e := range arr {
			result = append(result, []int{e})
		}
		return result
	}
	if n == len(arr) {
		return append(result, arr)
	}
	for _, a := range comb(arr[1:], n-1) {
		c := append([]int{arr[0]}, a...)
		result = append(result, c)
	}
	return append(result, comb(arr[1:], n)...)
}

func comball(arr []int) (result [][]int) {
	for i := 1; i <= len(arr); i++ {
		c := comb(arr, i)
		log.Printf("  %02dC%02d = %7d\n", len(arr), i, len(c))
		result = append(result, c...)
	}
	return result
}

func normflagcombs() []NormFlag {
	flags := normflags()
	flagcombs := comball(flags)
	result := make([]NormFlag, len(flagcombs))
	for i, comb := range flagcombs {
		for _, flag := range comb {
			result[i] |= NormFlag(flag)
		}
	}
	return result
}

func parenormflagcombs() (valid, invalid []NormFlag) {
	flagcombs := normflagcombs()
	log.Printf("  total : %7d\n", len(flagcombs))
	valid = make([]NormFlag, 0, len(flagcombs))
	invalid = make([]NormFlag, 0, len(flagcombs))
outer:
	for _, flagcomb := range flagcombs {
		for _, invalidFlags := range invalidFlagsList {
			if flagcomb&invalidFlags == invalidFlags {
				invalid = append(invalid, flagcomb)
				continue outer
			}
		}
		valid = append(valid, flagcomb)
	}
	return valid, invalid
}

func TestHeavyNormFlags(t *testing.T) {
	valid, invalid := parenormflagcombs()
	log.Printf("  valid : %7d\n", len(valid))
	log.Printf("invalid : %7d\n", len(invalid))
	// testing invalid flags
	for _, flag := range invalid {
		_, err := Norm(flag)
		if err == nil {
			t.Errorf("TestInvalidFlags: %s is valid, want: invalid", flag)
		}
	}
	// testing valid flags
	// testing all runes and all flag combinations
outer:
	for i, flag := range valid {
		n, err := Norm(flag)
		if err != nil {
			t.Errorf("TestInvalidFlags: %s is invalid, want: valid", flag)
			continue
		}
		for r := rune(0); r < maxr; r++ {
			c, rOK := findUnichar(r)
			nr, vm := n.normalizeRune(r)
			_, nrOK := findUnichar(nr)
			if rOK != nrOK {
				// TEST_G9amUMTr
				if rOK {
					t.Errorf("normalizeRune(%#U), flags: %s,\n\thave: %#U is not exists in unichars"+
						"\n\twant: %#U is exists in unichars", r, flag, nr, nr)
				} else {
					t.Errorf("normalizeRune(%#U), flags: %s,\n\thave: %#U is exists in unichars"+
						"\n\twant: %#U is not exists in unichars", r, flag, nr, nr)
				}
				break outer
			}
			if !nrOK {
				if vm != vmNone {
					t.Errorf("normalizeRune(%#U), flags: %s,\n\thave: nr=%#U, vm=%#U\n\twant: vm=%#U",
						r, flag, nr, vm, vmNone)
					break outer
				}
				if r != nr {
					t.Errorf("normalizeRune(%#U), flags: %s,\n\thave: nr=%#U, vm=%#U\n\twant: nr=%#U",
						r, flag, nr, vm, r)
					break outer
				}
				continue
			}
			if c.voicing == vcUndefined || c.voicing == vcUnvoiced {
				if vm != vmNone {
					// TEST_nD7FwQUW
					t.Errorf("normalizeRune(%#U), flags: %s,\n\thave: nr=%#U, vm=%#U\n\twant: vm=%#U",
						r, flag, nr, vm, vmNone)
					break outer
				}
				continue
			}
			if (c.charCase == ccHiragana && flag.has(HiraganaToNarrow) && !flag.has(DecomposeVom)) ||
				(c.charCase == ccKatakana && flag.has(KatakanaToNarrow) && !flag.has(DecomposeVom)) {
				if c.voicing == vcVoiced && vm != vmVsmNarrow {
					t.Errorf("normalizeRune(%#U), flags: %s,\n\thave: nr=%#U, vm=%#U\n\twant: vm=%#U",
						r, flag, nr, vm, vmVsmNarrow)
					break outer
				}
				if c.voicing == vcSemivoiced && vm != vmSsmNarrow {
					t.Errorf("normalizeRune(%#U), flags: %s,\n\thave: nr=%#U, vm=%#U\n\twant: vm=%#U",
						r, flag, nr, vm, vmSsmNarrow)
					break outer
				}
			}
			if c.charCase == ccHiragana && flag.has(HiraganaToKatakana) && !flag.has(DecomposeVom) {
				if vm != vmNone {
					t.Errorf("normalizeRune(%#U), flags: %s,\n\thave: nr=%#U, vm=%#U\n\twant: vm=%#U",
						r, flag, nr, vm, vmNone)
					break outer
				}
			}
			if c.charCase == ccKatakana && flag.has(KatakanaToHiragana) && !flag.has(DecomposeVom) {
				switch r {
				case 'ヷ', 'ヸ', 'ヹ', 'ヺ':
					if vm != vmVsmWide {
						t.Errorf("normalizeRune(%#U), flags: %s,\n\thave: nr=%#U, vm=%#U\n\twant: vm=%#U",
							r, flag, nr, vm, vmVsmWide)
						break outer
					}
				default:
					if vm != vmNone {
						t.Errorf("normalizeRune(%#U), flags: %s,\n\thave: nr=%#U, vm=%#U\n\twant: vm=%#U",
							r, flag, nr, vm, vmNone)
						break outer
					}
				}
			}
		}
		if i%500 == 0 {
			log.Printf("%5d/%5d (%3d%% done)", i, len(valid), i*100/len(valid))
		}
	}
}

const normSTR = "\t Aa#　Ａａ＃あア。ｱ｡”ﾞ漢字ｶﾞｷﾞｸﾞｹﾞｺﾞﾊﾟﾋﾟﾌﾟﾍﾟﾎﾟ\U0010FFFF"

func BenchmarkString(b *testing.B) {
	n, _ := Norm(LatinToNarrow | KanaToWide)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := strings.Repeat(normSTR, 1)
		n.String(s)
	}
	b.StopTimer()
}

func BenchmarkNormNFKD(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := strings.Repeat(normSTR, 1)
		norm.NFKD.String(s)
	}
	b.StopTimer()
}

func BenchmarkNormNFKC(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := strings.Repeat(normSTR, 1)
		norm.NFKC.String(s)
	}
	b.StopTimer()
}

func BenchmarkWidthNarrow(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := strings.Repeat(normSTR, 1)
		width.Narrow.String(s)
	}
	b.StopTimer()
}

func BenchmarkWidthWiden(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := strings.Repeat(normSTR, 1)
		width.Widen.String(s)
	}
	b.StopTimer()
}

func BenchmarkWidthFold(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := strings.Repeat(normSTR, 1)
		width.Fold.String(s)
	}
	b.StopTimer()
}
