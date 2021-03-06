/*

Vert is a utility to convert text files to vertical printing.

Usage:
	vert [flags] [path ...]

The flags are:
	-v
		Show version
	-h
		Show help
	-width
		Maximum width of output (default: 40)
	-height
		Maximum height of output (default: 25)

Examples:

To read standard input:
	$ echo -e "閑さや\n岩にしみ入る\n蝉の声" | vert
	蝉岩閑
	のにさ
	声しや
	  み
	  入
	  る

If you have the following files,
	$ cat basho.txt
	閑さや
	岩にしみ入る
	蝉の声

	芭蕉

To read this file:
	$ vert basho.txt
	芭  蝉岩閑
	蕉  のにさ
	    声しや
	      み
	      入
	      る

To limit the height:
	$ vert -height 4 basho.txt
	芭  蝉入岩閑
	蕉  のるにさ
	    声  しや
	        み

To limit height and width:
	$ vert -width 3 -height 4 basho.txt
	入岩閑
	るにさ
	  しや
	  み

	芭  蝉
	蕉  の
	    声

*/
package main
