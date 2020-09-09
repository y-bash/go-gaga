package gaga_test

import (
	"fmt"
	"github.com/y-bash/go-gaga"
)

func ExampleVert() {
	s := "閑さや\n岩にしみ入る\n蝉の声"
	ss := gaga.Vert(s, 3, 6)
	fmt.Print(ss[0])
	// Output:
	// 蝉岩閑
	// のにさ
	// 声しや
	//   み
	//   入
	//   る
}

func ExampleNormalizer_Normalize() {
	s := "ＧａGa is not がｶﾞガ"
	fmt.Println(0, s)

	n, _ := gaga.NewNormalizer(gaga.LatinToNarrow)
	fmt.Println(1, n.Normalize(s))

	n.SetFlag(gaga.KanaToWide)
	fmt.Println(2, n.Normalize(s))

	n.SetFlag(gaga.KanaToHiragana)
	fmt.Println(3, n.Normalize(s))

	n.SetFlag(gaga.KanaToNarrowKatakana)
	fmt.Println(4, n.Normalize(s))

	n.SetFlag(gaga.LatinToNarrow | gaga.AlphaToUpper | gaga.KanaToWideKatakana)
	fmt.Println(5, n.Normalize(s))

	// Output:
	// 0 ＧａGa is not がｶﾞガ
	// 1 GaGa is not がｶﾞガ
	// 2 ＧａGa is not がガガ
	// 3 ＧａGa is not ががが
	// 4 ＧａGa is not ｶﾞｶﾞｶﾞ
	// 5 GAGA IS NOT ガガガ
}
