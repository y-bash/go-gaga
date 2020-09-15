package gaga

import (
	"strings"
	"testing"
)

type NormFlag_StringTest struct {
	flag NormFlag
	name string
}

var normflag_stringtests = []NormFlag_StringTest{
	0: {0, "<undefined>"},
	1: {AlphaToNarrow, "AlphaToNarrow"},
	2: {LatinToNarrow, "(AlphaToNarrow | DigitToNarrow | SymbolToNarrow)"},
}

func TestNormFlag_String(t *testing.T) {
	for i, tt := range normflag_stringtests {
		have := tt.flag.String()
		if have != tt.name {
			t.Errorf("#%d have: %s, want: %s", i, have, tt.name)
		}
	}
}

type ParseNormFlagTest struct {
	name string
	flag NormFlag
	errS string
}

var parsenormflagtests = []ParseNormFlagTest{
	0:  {"", 0, "invalid normalization flag"},
	1:  {"<undefined>", 0, "invalid normalization flag"},
	2:  {"AlphaToNarrow", AlphaToNarrow, ""},
	3:  {"LatinToNarrow", LatinToNarrow, ""},
	4:  {"LatinToNarrow", AlphaToNarrow | DigitToNarrow | SymbolToNarrow, ""},
	5:  {"AlphaToNarrow|DigitToNarrow|SymbolToNarrow", LatinToNarrow, ""},
	6:  {"(AlphaToNarrow | DigitToNarrow)", AlphaToNarrow | DigitToNarrow, ""},
	7:  {"((((AlphaToNarrow | DigitToNarrow))))", AlphaToNarrow | DigitToNarrow, ""},
	8:  {"((((AlphaToNarrow))))", AlphaToNarrow, ""},
	9:  {"))))((((AlphaToNarrow", AlphaToNarrow, ""},
	10: {"))))((((|||||", 0, "invalid normalization flag"},
	11: {"(|)(|)(|)(|)(|)", 0, "invalid normalization flag"},
}

func TestParseNormFlag(t *testing.T) {
	for i, tt := range parsenormflagtests {
		have, err := ParseNormFlag(tt.name)
		if err != nil {
			switch {
			case tt.errS == "":
				t.Errorf("#%d have error: %s, want no error", i, err)
			case !strings.Contains(err.Error(), tt.errS):
				t.Errorf("#%d have error: %s, want error: %s", i, err, tt.errS)
			}
			continue
		}
		if tt.errS != "" {
			t.Errorf("#%d have no error, want error: %s", i, tt.errS)
		}
		if have != tt.flag {
			t.Errorf("#%d have: %d, want: %d, %s",
				i, have, tt.flag, tt.flag.String())
		}
	}
}

func TestNormFlagStringAndParse(t *testing.T) {
	for flag, name := range normflagMap {
		haveName := flag.String()
		if haveName != name {
			t.Errorf("NormFlag(%d).String() = %s, want: %s",
				flag, haveName, name)
			continue
		}
		haveFlag, err := ParseNormFlag(name)
		if err != nil {
			t.Errorf("ParseNormFlag(%q) return: %s, want: no error",
				name, err.Error())
			continue
		}
		if haveFlag != flag {
			t.Errorf("ParseNormFlag(%q) = %d, want: %d",
				name, haveFlag, flag)
		}
	}
	for _, combflag := range combflagList {
		var haveFlag NormFlag
		var err error

		haveFlag, err = ParseNormFlag(combflag.name)
		if err != nil {
			t.Errorf("ParseNormFlag(%q) return: %s, want: no error",
				combflag.name, err.Error())
			continue
		}
		if haveFlag != combflag.flag {
			t.Errorf("ParseNormFlag(%q) = %d, want: %d",
				combflag.name, haveFlag, combflag.flag)
			continue
		}

		name := combflag.flag.String()
		haveFlag, err = ParseNormFlag(name)
		if err != nil {
			t.Errorf("ParseNormFlag(%q) return: %s, want: no error",
				name, err.Error())
			continue
		}
		if haveFlag != combflag.flag {
			t.Errorf("ParseNormFlag(%q) = %d, want: %d",
				name, haveFlag, combflag.flag)
		}
	}

}
