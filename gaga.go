// Package gaga implements Japanese string conversion.
//
// Usage
//
// TODO update package comments
// TODO consider Renaming to doc.go
// Define a normalizer using NewNormalizer() with the normalization flag.
//
// This declares a normalizer, that converts Latin characters to
// half-width and Kana characters to full-width.
// 	import "github.com/y-bash/go-gaga"
// 	n:= gaga.NewNormalizer(gaga.LatinToNarrow | gaga.KanaToWide)
//
// After normalizer is defined, call
// 	s := n.Normalize("ＡＢＣｱｲｳ")
// to normalize the string using the normalization flags.
//
// Then the string is converted to
// 	fmt.Printf("%q", s) // Stdout: "ABCアイウ"
package gaga
