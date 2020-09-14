/*
Package gaga implements simple functions to manipulate UTF-8 encoded
Japanese strings.

Here is a simple example, converting the character type and printing
it vertically.

First, import gaga.
	import "github.com/y-bash/go-gaga"
Define a normalizer using Norm() with the normalization flag.
This declares a normalizer, that converts Latin characters to
half-width and Hiragana-Katakana characters to full-width.
	n, err := gaga.Norm(gaga.LatinToNarrow | gaga.KanaToWide)
	if err != nil {
		log.Fatal(err)
	}
After normalizer is defined, call to normalize the string using the
normalization flags.
	s := n.String("Ｇａｇａはｶﾞｶﾞｶﾞではありません")
	fmt.Println(s)
Output is:
	Gagaはガガガではありません
Using Vert(), make this string vertical.
	v := gaga.Vert(s, 3, 5)
	fmt.Print(v[0])
Output is:
	あガ G
	りガ a
	まガ g
	せで a
	んはは
*/
package gaga
