# Vert

Vert is a utility to convert text to vertical writing.

## Usage:

    vert [flags] [path ...]

## The flags are:

    -v
    	Show version
    -h
    	Show help
    -width
    	Maximum width of output
    -height
    	Maximum height of output


## Examples:

### To read standard input:

    $ echo -e "閑さや\n岩にしみ入る\n蝉の声" | vert
    蝉岩閑
    のにさ
    声しや
      み
      入
      る

### If you have the following files,

    $ cat basho.txt
    閑さや
    岩にしみ入る
    蝉の声

    芭蕉

### To read this file:

    $ vert basho.txt
    芭  蝉岩静
    蕉  のにか
        声染さ
        みや
        入
        る

### To limit the height:

    $ vert -height 4 basho.txt
    芭  蝉入岩静
    蕉  のるにか
        声  染さ
            みや

### To limit height and width:

    $ vert -width 3 -height 4 basho.txt
    入岩静
    るにか
      染さ
      みや

    芭  蝉
    蕉  の
        声

