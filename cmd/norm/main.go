package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/y-bash/go-gaga"
	"io"
	"log"
	"os"
	"strings"
)

var version = "v0.0.0" // set value by go build -ldflags

func read(f io.Reader) string {
	var sb strings.Builder
	sb.Grow(1024)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		sb.WriteString(sc.Text())
		sb.WriteString("\n")
	}
	return sb.String()
}

func readfiles(paths []string) (out []string, err error) {
	if len(paths) == 0 {
		out = []string{read(os.Stdin)}
		return
	}
	for _, path := range paths {
		var f *os.File
		f, err = os.Open(path)
		if err != nil {
			return
		}
		defer f.Close()
		out = append(out, read(f))
	}
	return
}

func normstrs(f io.Writer, in []string, flag gaga.NormFlag) error {
	if len(in) <= 0 {
		return nil
	}
	n, err := gaga.Norm(flag)
	if err != nil {
		return err
	}
	fmt.Fprint(f, n.String(in[0]))
	for i := 1; i < len(in); i++ {
		fmt.Fprintln(f)
		fmt.Fprint(f, n.String(in[i]))
	}
	return nil
}

func main() {
	var v, h, f bool
	var normflag string
	flag.BoolVar(&v, "v", false, "show version")
	flag.BoolVar(&h, "h", false, "show help")
	flag.BoolVar(&f, "f", false, "show help of the normalization flags")
	flag.StringVar(&normflag, "flag", "Fold", "normalization flag")
	flag.Parse()
	if v {
		fmt.Println("version:", version)
		return
	}
	if h {
		flag.Usage()
		return
	}
	if f {
		fmt.Print(flaghelp)
		return
	}
	nf, err := gaga.ParseNormFlag(normflag)
	if err != nil {
		log.Fatal(err)
	}
	var ss []string
	ss, err = readfiles(flag.Args())
	if err != nil {
		log.Fatal(err)
	}
	err = normstrs(os.Stdout, ss, nf)
	if err != nil {
		log.Fatal(err)
	}
}

const flaghelp = `The normalization flags of the Norm command.

Usage:
    norm -flag <flags>


Example:
    norm -flag "AlphaToNarrow | KatakanaToHiragana"


The flags are:


LatinToNarrow
    Equivalent Combination:
        AlphaToNarrow | DigitToNarrow | SymbolToNarrow

    Description:
        LatinToNarrow is a combination of normalization flags for converting
        all the full-width Latin characters to their half-width.

                 | CHARACTER     | CONVERT TO
        ---------+---------------+----------------
                 | Wide Alphabet | Narrow Alphabet
        Category | Wide Digit    | Narrow Digit
                 | Wide Symbol   | Narrow Symbol
        ---------+---------------+----------------
        Example  | "Ａ１？"      | "A1?"

LatinToWide
    Equivalent Combination:
        AlphaToWide | DigitToWide | SymbolToWide

    Description:
        LatinToWide is a combination of normalization flags for converting
        all the half-width Latin characters to their full-width.
        
                 | CHARACTER       | CONVERT TO
        ---------+-----------------+--------------
                 | Narrow Alphabet | Wide Alphabet
        Category | Narrow Digit    | Wide Digit
                 | Narrow Symbol   | Wide Symbol
        ---------+-----------------+--------------
        Example  | "A1?"           | "Ａ１？"
    

KanaToNarrow
    Equivalent Combination:
        HiraganaToNarrow | KatakanaToNarrow | KanaSymbolToNarrow |
        IsolatedVomToNarrow | ComposeVom

    Description:
        KanaToNarrow is a combination of normalization flags for converting
        the full-width Hiragana-Katakana characters to their half-width as
        much as possible.
        
                 | CHARACTER                       | CONVERT TO
        ---------+---------------------------------+-------------------
                 | Hiaragana                       | Narrow Katakana
        Category | Wide Katakana                   | Narrow Katakana
                 | Wide Kana Symbol                | Narrow Kana Symbol
                 | Voiced/Semi-voiced Kana Letter  | Legacy composed
                 | Isolated Voicing Modifier (VOM) | Narrow VOM
        ---------+---------------------------------+-------------------
        Example  | "あイ、が゛"                    | "ｱｲ､ｶﾞﾞ"
    

KanaToWide
    Equivalent Combination:
        KatakanaToWide | KanaSymbolToWide | IsolatedVomToWide |
		ComposeVom

    Description:
        KanaToWide is a combination of normalization flags for converting
        all the half-width Katakana characters to their full-width.
        
                 | CHARACTER                       | CONVERT TO
        ---------+---------------------------------+-----------------
                 | Narrow Katakana                 | Wide Katakana
        Category | Narrow Kana Symbol              | Wide Kana Symbol
                 | Voiced/Semi-voiced Kana Letter  | Legacy composed
                 | Isolated Voicing Modifier (VOM) | Wide VOM
        ---------+---------------------------------+-----------------
        Example  | "ｱ､ｶﾞﾞ"                         | "ア、ガ゛"
    

KanaToWideKatakana
    Equivalent Combination:
        KatakanaToWide | HiraganaToKatakana | KanaSymbolToWide |
        IsolatedVomToWide | ComposeVom

    Description:
        KanaToWideKatakana is a combination of normalization flags for
        converting all the half-width Katakana characters to their full-width,
        and the Hiragana characters to their full-width Katakana as much as
        possible..
        
                 | CHARACTER                       | CONVERT TO
        ---------+---------------------------------+-----------------
                 | Hiragana                        | Wide Katakana
        Category | Narrow Katakana                 | Wide Katakana
                 | Narrow Kana Symbol              | Wide Kana Symbol
                 | Voiced/Semi-voiced Kana Letter  | Legacy composed
                 | Isolated Voicing Modifier (VOM) | Wide VOM
        ---------+---------------------------------+-----------------
        Example  | "あｲ､ｶﾞﾞ"                       | "アイ、ガ゛"
    

KanaToNarrowKatakana
    Equivalent Combination:
        KatakanaToNarrow | HiraganaToNarrow | KanaSymbolToNarrow |
		IsolatedVomToNarrow | ComposeVom

    Description:
        KanaToNarrowKatakana is a combination of normalization flags for
        converting the full-width Katakana characters to their half-width,
        and the Hiragana characters to their half-width Katakana as much as
        possible.
        
                 | CHARACTER                       | CONVERT TO
        ---------+---------------------------------+-------------------
                 | Hiragana                        | Narrow Katakana
        Category | Wide Katakana                   | Narrow Katakana
                 | Wide Kana Symbol                | Narrow Kana Symbol
                 | Voiced/Semi-voiced Kana Letter  | Legacy composed
                 | Isolated Voicing Modifier (VOM) | Narrow VOM
        ---------+---------------------------------+-------------------
        Example  | "あイ、が゛"                    | "ｱｲ､ｶﾞﾞ"
    

KanaToHiragana
    Equivalent Combination:
        KatakanaToHiragana | KanaSymbolToWide | IsolatedVomToWide |
        ComposeVom

    Description:
        KanaToHiragana is a combination of normalization flags for
        converting the full-width Katakana characters to their Hiragana
        as much as possible, and all the half-width Katakana characters
        to their Hiragana.
        
                 | CHARACTER                       | CONVERT TO
        ---------+---------------------------------+----------------------
                 | Wide Katakana                   | Hiragana
        Category | Narrow Katakana                 | Hiragana
                 | Narrow Kana Symbol              | Wide Kana Symbol
                 | Voiced/Semi-voiced Kana Letter  | Legacy composed
                 | Isolated Voicing Modifier (VOM) | Wide VOM
        ---------+---------------------------------+----------------------
        Example  | "アｲ､ガ゛"                      | "あい、が゛"
    

Fold
    Equivalent Combination:
        LatinToNarrow | KanaToWide

    Description:
        Fold is a combination of normalization flags for converting
        the Latin characters and the Hiragana-Katakana characters to
        their canonical width.
        
                 | CHARACTER                       | CONVERT TO
        ---------+---------------------------------+-----------------
                 | Wide Alphabet                   | Narrow Alphabet
                 | Wide Digit                      | Narrow Digit
                 | Wide Symbol                     | Narrow Symbol
        Category | Narrow Katakana                 | Wide Katakana
                 | Narrow Kana Symbol              | Wide Kana Symbol
                 | Voiced/Semi-voiced Kana Letter  | Legacy composed
                 | Isolated Voicing Modifier (VOM) | Wide VOM
        ---------+---------------------------------+-----------------
        Example  | "Ａ１？ｱ､ｶﾞﾞ"                   | "A1?ア、ガ゛"
    

AlphaToNarrow
    Description:
        AlphaToNarrow converts all the full-width Latin letters to
        their half-width.

    Example:
        [Ａ] =>[A]

AlphaToWide
    Description:
        AlphaToWide converts all the half-width Latin letters to
        their full-width.

    Example:
        [A] => [Ａ]

AlphaToUpper
    Description:
        AlphaToUpper converts all the lower case Latin letters to
        their upper case.

    Examples:
        [a] => [A],  [ａ] => [Ａ]

AlphaToLower
    Description:
        AlphaToLower converts all the upper case Latin letters to
        their lower case.

    Examples:
        [A] => [a],  [Ａ] => [ａ]

DigitToNarrow
    Description:
        DigitToNarrow converts all the full-width Latin digits to
        their half-width.

    Example:
        [１] => [1]

DigitToWide
    Description:
        DigitToWide converts all the half-width Latin digits to
        their full-width.

    Example:
        [1] => [１]

SymbolToNarrow
    Description:
        SymbolToNarrow converts all the full-width Latin symbols to
        their half-width.

    Example:
        [？] => [?]

SymbolToWide
    Description:
        SymbolToWide converts all the half-width Latin symbols to
        their full-width.

    Example:
        [?] => [？]

HiraganaToNarrow
    Description:
        HiraganaToNarrow converts the full-width Hiragana letters to
        their half-width Katakana as much as possible.

    Example:
        [あ] => [ｱ]

HiraganaToKatakana
    Description:
        HiraganaToKatakana converts the full-width Hiragana letters to
        their full-width Katakana as much as possible.

    Example:
        [あ] => [ア]

KatakanaToNarrow
    Description:
        KatakanaToNarrow converts the full-width Katakana letters to
        their half-width Katakana as much as possible.

    Example:
        [ア] => [ｱ]

KatakanaToWide
    Description:
        KatakanaToWide converts all the half-width Katakana letters to
        their full-width Katakana.

    Example:
        [ｱ] => [ア]

KatakanaToHiragana
    Description:
        KatakanaToHiragana converts the half-width or full-width Katakana
        letters to their full-width Hiragana as much as possible.

    Examples:
        [ア] => [あ],  [ｱ] => [あ]

KanaSymbolToNarrow
    Description:
        KanaSymbolToNarrow converts the full-width Hiragana-Katakana
        symbols to their half-width as much as possible.

    Example:
        [、] => [､]

KanaSymbolToWide
    Description:
        KanaSymbolToWide converts all the half-width Katakana symbols
        to their full-width.

    Example:
        [､] => [、]

ComposeVom
    Description:
        ComposeVom composes the voiced or semi-voiced sound letters in
        the most conventional way.

    Examples:
        [が]     => [が],  [か][゛] => [が],    [か][U+3099] => [が],
        [か][ﾞ]  => [が],  [ｶ][゛]  => [ｶ][ﾞ],  [ｶ][ﾞ]       => [ｶ][ﾞ],
        [は][゜] => [ぱ],  [ヰ][゛] => [ヸ],    [ゐ][゛]     => [ゐ][゛]

DecomposeVom
    Description:
        DecomposeVom decomposes the voiced or semi-voiced sound letters
        in a way similar to the Unicode canonical decomposition mappings.

    Examples:
        [が]         => [か][U+3099],  [か][゛] => [か][U+3099],
        [か][U+3099] => [か][U+3099],  [か][ﾞ]  => [か][U+3099],
        [ｶ][゛]      => [ｶ][U+3099],   [ｶ][ﾞ]   => [ｶ][U+3099],
        [ぱ]         => [は][U+309A],  [ヰ][゛] => [ヰ][U+3099],
        [ゐ][゛]     => [ゐ][U+3099]

IsolatedVomToNarrow
    Description:
        IsolatedVomToNarrow converts an isolated voicing modifier
        which was not combined into a base letter into a half-width
        voiced or semi-voiced sound letter.

    Examples:
        [゛] => [ﾞ],  [U+3099] => [ﾞ],  [゜] => [ﾟ],  [U+309A] => [ﾟ]

IsolatedVomToWide
    Description:
        IsolatedVomToWide converts an isolated voicing modifier
        which was not combined into a base letter into a full-width
        voiced or semi-voiced sound letter.

    Examples:
        [U+3099] => [゛],  [ﾞ] => [゛],  [U+309A] => [゜],  [ﾟ] => [゜]

IsolatedVomToNonspace
    Description:
        IsolatedVomToCombining converts an isolated voicing
        modifier which was not combined into a base letter into a
        combining voiced or semi-voiced sound letter.

    Examples:
        [゛] => [U+3099],  [ﾞ] => [U+3099],  [゜] = [U+309A],  [ﾟ] => [U+309A]
`
