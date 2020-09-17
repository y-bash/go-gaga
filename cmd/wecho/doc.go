/*

Wecho is an echo command that writes utf-8 text to standard output.

The Windows echo command usually writes cp932 text to standard output
in the Japanese locale, which causes garbled characters when piped to
a program written in Golang. Wecho is provided to solve this problem.

Wecho also converts the sequence of \n in command line argument to
the newline character.

Usage:
	wecho [flags] [text]

The flags are:
	-v
		Show version
	-h
		Show help

Examples:
	> wecho 閑さや\n岩にしみ入る\n蝉の声
	閑さや
	岩にしみ入る
	蝉の声

*/
package main
