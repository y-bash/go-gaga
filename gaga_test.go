package gaga_test

import (
	"fmt"
	"github.com/y-bash/go-gaga"
)

func ExampleVertFix() {
	s := "閑さや\n岩にしみ入る\n蝉の声"

	vs := gaga.VertFix(s, 6, 6)
	fmt.Println(" 1 2 3 4 5 6")
	fmt.Print(vs)

	vs = gaga.VertFix(s, 6, 3)
	fmt.Println("\n 1 2 3 4 5 6")
	fmt.Print(vs)

	vs = gaga.VertFix(s, 3, 3)
	fmt.Println("\n 1 2 3")
	fmt.Print(vs)

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
	//  1 2 3
	// み岩閑
	// 入にさ
	// るしや
	//
	//     蝉
	//     の
	//     声
}

func ExampleVertFixStrings() {
	s := "閑さや\n岩にしみ入る\n蝉の声"

	ss := gaga.VertFixStrings(s, 6, 6)
	fmt.Println(" 1 2 3 4 5 6")
	fmt.Print(ss[0])

	ss = gaga.VertFixStrings(s, 6, 3)
	fmt.Println("\n 1 2 3 4 5 6")
	fmt.Print(ss[0])

	ss = gaga.VertFixStrings(s, 3, 3)
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

func ExampleVertToFit() {
	s := "閑さや\n岩にしみ入る\n蝉の声"

	vs := gaga.VertToFit(s, 6, 6)
	fmt.Println(" 1 2 3 4 5 6")
	fmt.Print(vs)

	vs = gaga.VertToFit(s, 6, 3)
	fmt.Println("\n 1 2 3 4 5 6")
	fmt.Print(vs)

	vs = gaga.VertToFit(s, 3, 3)
	fmt.Println("\n 1 2 3")
	fmt.Print(vs)

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
	//  1 2 3
	// み岩閑
	// 入にさ
	// るしや
	//
	//     蝉
	//     の
	//     声
}

func ExampleVertToFitStrings() {
	s := "閑さや\n岩にしみ入る\n蝉の声"

	ss := gaga.VertToFitStrings(s, 6, 6)
	fmt.Println(" 1 2 3 4 5 6")
	fmt.Print(ss[0])

	ss = gaga.VertToFitStrings(s, 6, 3)
	fmt.Println("\n 1 2 3 4 5 6")
	fmt.Print(ss[0])

	ss = gaga.VertToFitStrings(s, 3, 3)
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

func ExampleVertNowrap() {
	in := "閑さや岩にしみ入る蝉の声"
	out := gaga.VertNowrap(in)
	fmt.Print(out)
	// Output:
	// 閑
	// さ
	// や
	// 岩
	// に
	// し
	// み
	// 入
	// る
	// 蝉
	// の
	// 声

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
