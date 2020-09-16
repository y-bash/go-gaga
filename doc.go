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
	n, err := gaga.Norm(gaga.LatinToNarrow | gaga.KanaToHiragana)
	if err != nil {
		log.Fatal(err)
	}

After normalizer is defined, call to normalize the string using the
normalization flags.
	s := n.String("ＧａGaはがｶﾞガではありません")
	fmt.Println(s)

Output is:
	GaGaはがががではありません

Using Vert(), make this string vertical.
	vs := gaga.Vert(s, 3, 5)
	fmt.Print(vs)

Output is:
	あが G
	りが a
	まが G
	せで a
	んはは
*/
package gaga
