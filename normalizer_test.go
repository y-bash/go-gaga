package gaga

import (
	"fmt"
	"strings"
	"testing"
)

type NormalizeRuneTest struct {
	name string
	flag NormFlag
	lo   rune
	hi   rune
	diff rune
}

var normalizerunetests = []NormalizeRuneTest{
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

func TestNormalizeRune(t *testing.T) {
	for i, tt := range normalizerunetests {
		n, err := NewNormalizer(tt.flag)
		if err != nil {
			t.Errorf("#%d: %s", i, err.Error())
			continue
		}
		for src := tt.lo; src <= tt.hi; src++ {
			want := src + tt.diff
			got := n.NormalizeRune(src)
			if len(got) != 1 {
				t.Errorf("#%d %s NormalizeRune(%#U) = %v, want len(%v) is 1", i, tt.name, src, got, got)
			}
			if got[0] != want {
				t.Errorf("#%d %s NormalizeRune(%#U) = %#U, want %#U", i, tt.name, src, got[0], want)
			}
		}
	}
}

type NormalizeTest struct {
	flag NormFlag
	//	mods string
	in  string
	out string
}

// TODO testing of invalid normalization flag
// TODO testing of whitespace
// TODO Consider whether to test the following characters
//   U+301C  '〜' 1.1 WAVE DASH
//   U+FF5E  '～' 1.1 FULLWIDTH TILDE
//   U+1301C '〜' 5.2 EGYPTIAN HIEROGLYPH A024
var normalizetests = []NormalizeTest{
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

	// simple japanese conversion <- zero length string
	50: {KanaVsmToClassical, "", ""},
	51: {HiraganaToNarrow | KanaVsmToClassical, "", ""},
	52: {HiraganaToKatakana | KanaVsmToClassical, "", ""},
	53: {KatakanaToWide | KanaVsmToClassical, "", ""},
	54: {KatakanaToNarrow | KanaVsmToClassical, "", ""},
	55: {KatakanaToHiragana | KanaVsmToClassical, "", ""},
	56: {KanaSymToWide | KanaVsmToClassical, "", ""},
	57: {KanaSymToNarrow | KanaVsmToClassical, "", ""},

	// simple japanese conversion <- CJK (Hiragana/Katakana/Symbol/Han), Emoji
	58: {KanaVsmToClassical,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	59: {HiraganaToNarrow | KanaVsmToClassical,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」ｱｳﾞｧｹアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	60: {HiraganaToKatakana | KanaVsmToClassical,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」アヴァヶアヺァヶ､｣ｱｳﾞｧｹ漢👻"},
	61: {KatakanaToWide | KanaVsmToClassical,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ､｣アヴァケ漢👻"},
	62: {KatakanaToNarrow | KanaVsmToClassical,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖｱｦﾞｧｹ､｣ｱｳﾞｧｹ漢👻"},
	63: {KatakanaToHiragana | KanaVsmToClassical,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖあを゛ぁゖ､｣あゔぁけ漢👻"},
	64: {KanaSymToWide | KanaVsmToClassical,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖアヺァヶ、」ｱｳﾞｧｹ漢👻"},
	65: {KanaSymToNarrow | KanaVsmToClassical,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"､｣あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻"},

	// simple japanese conversion <- Latin (Letter/Digit/Symbol)
	66: {KanaVsmToClassical,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ０９！～"},
	67: {HiraganaToNarrow | KanaVsmToClassical,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ０９！～"},
	68: {HiraganaToKatakana | KanaVsmToClassical,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ０９！～"},
	69: {KatakanaToWide | KanaVsmToClassical,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ０９！～"},
	70: {KatakanaToNarrow | KanaVsmToClassical,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ０９！～"},
	71: {KatakanaToHiragana | KanaVsmToClassical,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ０９！～"},
	72: {KanaSymToWide | KanaVsmToClassical,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ０９！～"},
	73: {KanaSymToNarrow | KanaVsmToClassical,
		" AZaz09!~　ＡＺａｚ０９！～",
		" AZaz09!~　ＡＺａｚ０９！～"},

	// multiple japanese conversion <- CJK (Hiragana/Katakana/Symbol/Han), Emoji
	74: {KatakanaToHiragana | KanaSymToWide | KanaVsmToClassical,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"、」あゔぁゖあを゛ぁゖ、」あゔぁけ漢👻"},
	75: {KanaToNarrow | KanaVsmToClassical,
		"、」あゔぁゖアヺァヶ､｣ｱｳﾞｧｹ漢👻",
		"､｣ｱｳﾞｧｹｱｦﾞｧｹ､｣ｱｳﾞｧｹ漢👻"},
	76: {KatakanaToHiragana | KanaVsmToClassical,
		"ｦｧｨｮｯｱｲﾛﾝﾜｲｴｶｹ",
		"をぁぃょっあいろんわいえかけ"},
	77: {KatakanaToHiragana | KanaVsmToClassical,
		"ァアィイレロヮワヰヱヲンヵカヶケヷヸヹヺ",
		"ぁあぃいれろゎわゐゑをんゕかゖけわ゛ゐ゛ゑ゛を゛"},

	// simple japanese conversion (no vsm directive) <- Voiced sound character/Voiced sound mark
	78: {KatakanaToHiragana,
		"か゛かﾞか\u3099カ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"か゛かﾞか\u3099か゛かﾞか\u3099か゛かﾞか\u3099"},
	79: {KatakanaToHiragana,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"あ゛あﾞあ\u3099あ゛あﾞあ\u3099あ゛あﾞあ\u3099"},
	80: {KatakanaToHiragana,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099"},

	// multiple japanese conversion (no vsm directive) <- Voiced sound character/Voiced sound mark
	81: {KatakanaToNarrow | HiraganaToNarrow,
		"か゛かﾞか\u3099カ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ｶ゛ｶﾞｶ\u3099ｶ゛ｶﾞｶ\u3099ｶ゛ｶﾞｶ\u3099"},
	82: {KatakanaToNarrow | HiraganaToNarrow,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"ｱ゛ｱﾞｱ\u3099ｱ゛ｱﾞｱ\u3099ｱ゛ｱﾞｱ\u3099"},
	83: {KatakanaToNarrow | HiraganaToNarrow,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099"},

	// simple japanese conversion (vsm classic directive) <- Voiced sound character/Voiced sound mark
	84: {KatakanaToHiragana | KanaVsmToClassical,
		"か゛かﾞか\u3099カ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ががががががががが"},
	85: {KatakanaToHiragana | KanaVsmToClassical,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"あ゛あ゛あ゛あ゛あ゛あ゛あ゛あ゛あ゛"},
	86: {KatakanaToHiragana | KanaVsmToClassical,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"AﾞAﾞAﾞ日゛日ﾞ日\u3099ä゛äﾞä\u3099"}, // TODO Consider whether this specification (return 3099) is good

	// multiple japanese conversion (vsm classic directive) <- Voiced sound character/Voiced sound mark
	87: {KatakanaToNarrow | HiraganaToNarrow | KanaVsmToClassical,
		"か゛かﾞか\u3099カ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞｶﾞ"},
	88: {KatakanaToNarrow | HiraganaToNarrow | KanaVsmToClassical,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"ｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞｱﾞ"},
	89: {KatakanaToNarrow | HiraganaToNarrow | KanaVsmToClassical,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"AﾞAﾞAﾞ日゛日ﾞ日\u3099ä゛äﾞä\u3099"},

	// simple japanese conversion (vsm combining directive) <- Voiced sound character/Voiced sound mark
	90: {KatakanaToHiragana | KanaVsmToCombining,
		"か゛かﾞか\u3099カ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099か\u3099"},
	91: {KatakanaToHiragana | KanaVsmToCombining,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099あ\u3099"},
	92: {KatakanaToHiragana | KanaVsmToCombining,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A\u3099A\u3099A\u3099日\u3099日\u3099日\u3099ä\u3099ä\u3099ä\u3099"},

	// multiple japanese conversion (vsm combining directive) <- Voiced sound character/Voiced sound mark
	93: {KatakanaToNarrow | HiraganaToNarrow | KanaVsmToCombining,
		"か゛かﾞか\u3099カ゛カﾞカ\u3099ｶ゛ｶﾞｶ\u3099",
		"ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099ｶ\u3099"},
	94: {KatakanaToNarrow | HiraganaToNarrow | KanaVsmToCombining,
		"あ゛あﾞあ\u3099ア゛アﾞア\u3099ｱ゛ｱﾞｱ\u3099",
		"ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099ｱ\u3099"},
	95: {KatakanaToNarrow | HiraganaToNarrow | KanaVsmToCombining,
		"A゛AﾞA\u3099日゛日ﾞ日\u3099ä゛äﾞä\u3099",
		"A\u3099A\u3099A\u3099日\u3099日\u3099日\u3099ä\u3099ä\u3099ä\u3099"},

	// VSM/SVSM testing that cannot be combined with the previous character
	96: {KanaVsmToClassical,
		"は゛ば゛ぱ゛は゜ば゜ぱ゜ハ゛バ゛パ゛ハ゜バ゜パ゜ﾊ゛ﾊ゜",
		"ばば゛ぱ゛ぱば゜ぱ゜ババ゛パ゛パバ゜パ゜ﾊﾞﾊﾟ"},
	97: {KanaVsmToClassical,
		"はﾞばﾞぱﾞはﾟばﾟぱﾟハﾞバﾞパﾞハﾟバﾟパﾟﾊﾞﾊﾟ",
		"ばばﾞぱﾞぱばﾟぱﾟババﾞパﾞパバﾟパﾟﾊﾞﾊﾟ"},
	98: {KanaVsmToClassical | KanaSymToNarrow,
		"は゛ば゛ぱ゛は゜ば゜ぱ゜ハ゛バ゛パ゛ハ゜バ゜パ゜ﾊ゛ﾊ゜",
		"ばばﾞぱﾞぱばﾟぱﾟババﾞパﾞパバﾟパﾟﾊﾞﾊﾟ"},
	99: {KanaVsmToClassical | KanaSymToWide,
		"はﾞばﾞぱﾞはﾟばﾟぱﾟハﾞバﾞパﾞハﾟバﾟパﾟﾊﾞﾊﾟ",
		"ばば゛ぱ゛ぱば゜ぱ゜ババ゛パ゛パバ゜パ゜ﾊﾞﾊﾟ"},

	// VSM/SVSM testing that cannot be combined with the previous character (out of range in unichars table)
	100: {KanaVsmToClassical,
		"日゛本゜語ﾞ平ﾟ仮\u3099名\u309A",
		"日゛本゜語ﾞ平ﾟ仮\u3099名\u309A"}, // TODO Consider whether this specification (return 309A) is good
	101: {KanaVsmToClassical | KanaSymToNarrow,
		"日゛本゜語ﾞ平ﾟ仮\u3099名\u309A",
		"日ﾞ本ﾟ語ﾞ平ﾟ仮ﾞ名ﾟ"}, // TEST_N9x6dneg
	102: {KanaVsmToClassical | KanaSymToWide,
		"日゛本゜語ﾞ平ﾟ仮\u3099名\u309A",
		"日゛本゜語゛平゜仮゛名゜"}, // TEST_A9fCxUi6

	// VSM testing with and without KanaSymToXxx/KavaVsmToXxx flags
	103: {KanaVsmToClassical,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ｺﾞｺﾞｺﾞゴゴゴゴ"},
	104: {KatakanaToHiragana | KanaVsmToClassical,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ごごごごごごご"},
	105: {KatakanaToWide,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"コﾞコ゛コ\u3099コﾞコ゛コ\u3099ゴ"},
	106: {KatakanaToWide | KanaSymToWide,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"コ゛コ゛コ゛コ゛コ゛コ゛ゴ"},
	107: {KatakanaToWide | KanaSymToWide | KanaVsmToClassical,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ゴゴゴゴゴゴゴ"},
	108: {KatakanaToNarrow,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ｺﾞｺ゛ｺ\u3099ｺﾞｺ゛ｺ\u3099ｺﾞ"}, // TEST_L7tADs2z
	109: {KatakanaToNarrow | KanaSymToNarrow,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞ"},
	110: {KanaToNarrow | KanaSymToNarrow | KanaVsmToClassical,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞｺﾞ"},
	111: {KatakanaToHiragana,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"こﾞこ゛こ\u3099こﾞこ゛こ\u3099ご"},
	112: {KatakanaToHiragana | KanaSymToNarrow,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"こﾞこﾞこﾞこﾞこﾞこﾞご"},
	113: {KatakanaToHiragana | KanaSymToNarrow | KanaVsmToClassical,
		"ｺﾞｺ゛ｺ\u3099コﾞコ゛コ\u3099ゴ",
		"ごごごごごごご"},

	// SVSM testing with and without KanaSymToXxx/KavaVsmToXxx flags
	114: {KanaVsmToClassical,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ﾎﾟﾎﾟﾎﾟポポポポ"},
	115: {KatakanaToHiragana | KanaVsmToClassical,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ぽぽぽぽぽぽぽ"},
	116: {KatakanaToWide,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ホﾟホ゜ホ\u309Aホﾟホ゜ホ\u309Aポ"},
	117: {KatakanaToWide | KanaSymToWide,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ホ゜ホ゜ホ゜ホ゜ホ゜ホ゜ポ"},
	118: {KatakanaToWide | KanaSymToWide | KanaVsmToClassical,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ポポポポポポポ"},
	119: {KatakanaToNarrow,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ﾎﾟﾎ゜ﾎ\u309Aﾎﾟﾎ゜ﾎ\u309Aﾎﾟ"}, // TEST_K6t8hQYp
	120: {KatakanaToNarrow | KanaSymToNarrow,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ﾎﾟﾎﾟﾎﾟﾎﾟﾎﾟﾎﾟﾎﾟ"},
	121: {KanaToNarrow | KanaSymToNarrow | KanaVsmToClassical,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ﾎﾟﾎﾟﾎﾟﾎﾟﾎﾟﾎﾟﾎﾟ"},
	122: {KatakanaToHiragana,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ほﾟほ゜ほ\u309Aほﾟほ゜ほ\u309Aぽ"},
	123: {KatakanaToHiragana | KanaSymToNarrow,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ほﾟほﾟほﾟほﾟほﾟほﾟぽ"},
	124: {KatakanaToHiragana | KanaSymToNarrow | KanaVsmToClassical,
		"ﾎﾟﾎ゜ﾎ\u309Aホﾟホ゜ホ\u309Aポ",
		"ぽぽぽぽぽぽぽ"},

	// VSM/SVSM testing illegal rune value
	125: {KanaVsmToClassical, // TODO Consider whether this specification (return 3099, 309A) is good
		string([]rune{0x10FFFF + 1, '゛', 0x10FFFF + 1, 'ﾞ', 0x10FFFF + 1, '\u3099', 0x10FFFF + 1, '゜', 0x10FFFF + 1, 'ﾟ', 0x10FFFF + 1, '\u309A'}),
		string([]rune{0x10FFFF + 1, '゛', 0x10FFFF + 1, 'ﾞ', 0x10FFFF + 1, '\u3099', 0x10FFFF + 1, '゜', 0x10FFFF + 1, 'ﾟ', 0x10FFFF + 1, '\u309A'})},
	126: {KanaVsmToCombining,
		string([]rune{0x10FFFF + 1, '゛', 0x10FFFF + 1, 'ﾞ', 0x10FFFF + 1, '\u3099', 0x10FFFF + 1, '゜', 0x10FFFF + 1, 'ﾟ', 0x10FFFF + 1, '\u309A'}),
		string([]rune{0x10FFFF + 1, '\u3099', 0x10FFFF + 1, '\u3099', 0x10FFFF + 1, '\u3099', 0x10FFFF + 1, '\u309A', 0x10FFFF + 1, '\u309A', 0x10FFFF + 1, '\u309A'})},
	127: {KanaSymToWide,
		string([]rune{0x10FFFF + 1, '゛', 0x10FFFF + 1, 'ﾞ', 0x10FFFF + 1, '\u3099', 0x10FFFF + 1, '゜', 0x10FFFF + 1, 'ﾟ', 0x10FFFF + 1, '\u309A'}),
		string([]rune{0x10FFFF + 1, '゛', 0x10FFFF + 1, '゛', 0x10FFFF + 1, '゛', 0x10FFFF + 1, '゜', 0x10FFFF + 1, '゜', 0x10FFFF + 1, '゜'})},
	128: {KanaSymToNarrow,
		string([]rune{0x10FFFF + 1, '゛', 0x10FFFF + 1, 'ﾞ', 0x10FFFF + 1, '\u3099', 0x10FFFF + 1, '゜', 0x10FFFF + 1, 'ﾟ', 0x10FFFF + 1, '\u309A'}),
		string([]rune{0x10FFFF + 1, 'ﾞ', 0x10FFFF + 1, 'ﾞ', 0x10FFFF + 1, 'ﾞ', 0x10FFFF + 1, 'ﾟ', 0x10FFFF + 1, 'ﾟ', 0x10FFFF + 1, 'ﾟ'})},

	// special Katakana-Hiragana letters
	129: {KatakanaToHiragana | KanaVsmToClassical,
		"アイウエオヤユヨツワカケァィゥェォャュョッヮヵヶヷヸヹヺ",
		"あいうえおやゆよつわかけぁぃぅぇぉゃゅょっゎゕゖわ゛ゐ゛ゑ゛を゛"},
	130: {KatakanaToNarrow | KanaVsmToClassical,
		"アイウエオヤユヨツワカケァィゥェォャュョッヮヵヶヷヸヹヺ",
		"ｱｲｳｴｵﾔﾕﾖﾂﾜｶｹｧｨｩｪｫｬｭｮｯﾜｶｹﾜﾞｲﾞｴﾞｦﾞ"},
	131: {HiraganaToKatakana | KanaVsmToClassical,
		"あいうえおやゆよつわかけぁぃぅぇぉゃゅょっゎゕゖわ゛ゐ゛ゑ゛を゛",
		"アイウエオヤユヨツワカケァィゥェォャュョッヮヵヶヷヸヹヺ"},
	132: {HiraganaToNarrow | KanaVsmToClassical,
		"あいうえおやゆよつわかけぁぃぅぇぉゃゅょっゎゕゖわ゛ゐ゛ゑ゛を゛",
		"ｱｲｳｴｵﾔﾕﾖﾂﾜｶｹｧｨｩｪｫｬｭｮｯﾜｶｹﾜﾞｲﾞｴﾞｦﾞ"},
	133: {KatakanaToHiragana | KanaVsmToClassical,
		"ｱｲｳｴｵﾔﾕﾖﾂﾜｶｹｧｨｩｪｫｬｭｮｯﾜｶｹﾜﾞｲﾞｴﾞｦﾞ",
		"あいうえおやゆよつわかけぁぃぅぇぉゃゅょっわかけわ゛い゛え゛を゛"},
	134: {KatakanaToWide | KanaVsmToClassical,
		"ｱｲｳｴｵﾔﾕﾖﾂﾜｶｹｧｨｩｪｫｬｭｮｯﾜｶｹﾜﾞｲﾞｴﾞｦﾞ",
		"アイウエオヤユヨツワカケァィゥェォャュョッワカケヷイ゛エ゛ヺ"},

	// Hiragana letter YORI, Katakana letter KOTO
	135: {KatakanaToHiragana | KanaVsmToClassical, "ゟヿ", "ゟヿ"},
	136: {KatakanaToNarrow | KanaVsmToClassical, "ゟヿ", "ゟヿ"},
	137: {KatakanaToWide | KanaVsmToClassical, "ゟヿ", "ゟヿ"},
	138: {HiraganaToKatakana | KanaVsmToClassical, "ゟヿ", "ゟヿ"},
	139: {HiraganaToNarrow | KanaVsmToClassical, "ゟヿ", "ゟヿ"},

	// Yen mark
	140: {SymbolToWide, "\\￥", "＼￥"},
	141: {SymbolToNarrow, "＼￥", "\\￥"},

	// overflow
	/*
	145: {0,
		string([]rune{-1, '\u0000', '\U0010FFFF', 0x10FFFF + 1}),
		string([]rune{-1, '\u0000', '\U0010ffff', 0x10FFFF + 1})},
		*/
	142: {KanaVsmToClassical,
		string([]rune{-1, '\u0000', '\U0010FFFF', 0x10FFFF + 1}),
		string([]rune{-1, '\u0000', '\U0010ffff', 0x10FFFF + 1})},
	143: {KanaVsmToCombining,
		string([]rune{-1, '\u0000', '\U0010FFFF', 0x10FFFF + 1}),
		string([]rune{-1, '\u0000', '\U0010ffff', 0x10FFFF + 1})},
	144: {LatinToNarrow | KanaToWide,
		string([]rune{-1, '\u0000', '\U0010FFFF', 0x10FFFF + 1}),
		string([]rune{-1, '\u0000', '\U0010ffff', 0x10FFFF + 1})},

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

func TestNormalize(t *testing.T) {
	for i, tt := range normalizetests {
		n, err := NewNormalizer(tt.flag)
		if err != nil {
			t.Errorf("#%d: %s", i, err.Error())
			continue
		}
		out := n.Normalize(tt.in)
		if out != tt.out {
			t.Errorf("TestNormalize #%d\n\tflag:"+
				"\t%s, \n\targs:\t%q\n\thave:\t%q\n\twant:\t%q\n"+
				"\targs16:\t%s\n\thave16:\t%s\n\twant16:\t%s",
				i, tt.flag, tt.in, out, tt.out,
				hexs(tt.in), hexs(out), hexs(tt.out))
		}
	}
}
