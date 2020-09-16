/*

Norm is a utility to normalize Japanese language text files.

Usage:
	norm [flags] [path ...]

The flags are:
	-v
		Show version
	-h
		Show help
	-f
		Show help of the normalization flags
	-flag string
		Normalization flag (default "Fold")

Examples:

To read standard input:
	$ echo "ＡＢＣｱｲｳ" | norm
	ABCアイウ

	$ echo "ＡＢＣｱｲｳ" | norm -flag KanaToHiragana
	ＡＢＣあいう

	$ echo "ＡＢＣｱｲｳ" | norm -flag "AlphaToNarrow|AlphaToLower|KanaToHiragana"
	abcあいう

If you have the following files,
	$ cat basho_en.txt
	--Keene, Narrow Road 99

To read this file:
	$ norm -flag "AlphaToWide|AlphaToUpper" basho_en.txt
	--ＫＥＥＮＥ, ＮＡＲＲＯＷ ＲＯＡＤ 99

	$ norm -flag "LatinToWide|AlphaToUpper" basho_en.txt
	－－ＫＥＥＮＥ，　ＮＡＲＲＯＷ　ＲＯＡＤ　９９

If you have the following files,
	$ cat basho_jp.txt
	閑さや岩にしみ入る蝉の声

To read this file:
	$ norm -flag HiraganaToKatakana basho_jp.txt
	閑サヤ岩ニシミ入ル蝉ノ声

	$ norm -flag HiraganaToNarrow basho_jp.txt
	閑ｻﾔ岩ﾆｼﾐ入ﾙ蝉ﾉ声
*/
package main
