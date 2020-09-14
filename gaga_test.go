package gaga_test

import (
	"fmt"
	"github.com/y-bash/go-gaga"
)

func ExampleVertFix() {
	s := "閑さや\n岩にしみ入る\n蝉の声"

	ss := gaga.VertFix(s, 6, 6)
	fmt.Println(" 1 2 3 4 5 6")
	fmt.Print(ss[0])

	ss = gaga.VertFix(s, 6, 3)
	fmt.Println("\n 1 2 3 4 5 6")
	fmt.Print(ss[0])

	ss = gaga.VertFix(s, 3, 3)
	fmt.Println("\n-Page1\n 1 2 3")
	fmt.Print(ss[0])

	fmt.Println("\n-Page2\n 1 2 3")
	fmt.Print(ss[1])

	// Output:
	//  1 2 3 4 5 6
	//       蝉岩閑
	//       のにさ
	//       声しや
	//         み
	//         入
	//         る
	//
	//  1 2 3 4 5 6
	//     蝉み岩閑
	//     の入にさ
	//     声るしや
	//
	// -Page1
	//  1 2 3
	// み岩閑
	// 入にさ
	// るしや
	//
	// -Page2
	//  1 2 3
	//     蝉
	//     の
	//     声

}

func ExampleVertShrink() {
	s := "閑さや\n岩にしみ入る\n蝉の声"

	ss := gaga.VertShrink(s, 6, 6)
	fmt.Println(" 1 2 3 4 5 6")
	fmt.Print(ss[0])

	ss = gaga.VertShrink(s, 6, 3)
	fmt.Println("\n 1 2 3 4 5 6")
	fmt.Print(ss[0])

	ss = gaga.VertShrink(s, 3, 3)
	fmt.Println("\n-Page1\n 1 2 3")
	fmt.Print(ss[0])

	fmt.Println("\n-Page2\n 1 2 3")
	fmt.Print(ss[1])

	// Output:
	//  1 2 3 4 5 6
	// 蝉岩閑
	// のにさ
	// 声しや
	//   み
	//   入
	//   る
	//
	//  1 2 3 4 5 6
	// 蝉み岩閑
	// の入にさ
	// 声るしや
	//
	// -Page1
	//  1 2 3
	// み岩閑
	// 入にさ
	// るしや
	//
	// -Page2
	//  1 2 3
	//     蝉
	//     の
	//     声
}

func ExampleVert() {
	in := "閑さや\n岩にしみ入る\n蝉の声"
	out := gaga.Vert(in)
	fmt.Print(out)
	// Output:
	// 蝉岩閑
	// のにさ
	// 声しや
	//   み
	//   入
	//   る
}

func ExampleNormalizer_String() {
	s := "ＧａGa is not がｶﾞガ"
	fmt.Println(0, s)

	n, _ := gaga.Norm(gaga.LatinToNarrow)
	fmt.Println(1, n.String(s))

	n.SetFlag(gaga.KanaToWide)
	fmt.Println(2, n.String(s))

	n.SetFlag(gaga.KanaToHiragana)
	fmt.Println(3, n.String(s))

	n.SetFlag(gaga.KanaToNarrowKatakana)
	fmt.Println(4, n.String(s))

	n.SetFlag(gaga.LatinToNarrow | gaga.AlphaToUpper | gaga.KanaToWideKatakana)
	fmt.Println(5, n.String(s))

	// Output:
	// 0 ＧａGa is not がｶﾞガ
	// 1 GaGa is not がｶﾞガ
	// 2 ＧａGa is not がガガ
	// 3 ＧａGa is not ががが
	// 4 ＧａGa is not ｶﾞｶﾞｶﾞ
	// 5 GAGA IS NOT ガガガ
}
