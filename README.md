# go-gaga (Japanese language utility)

## Installation

For using the library:

```
$ go get github.com/y-bash/go-gaga
```

Next, to install the command (execution binary)

```
$ cd $GOPATH/src/github.com/y-bash/go-gaga
$ make install
```

## Usage

### Library

#### Norm

```
import "github.com/y-bash/go-gaga"

s := "ＧａGa is not がｶﾞガ"
fmt.Println(s)

n := gaga.Norm(gaga.Fold) // gaga.Fold == gaga.LatinToNarrow | gaga.KanaToWide
fmt.Println(n.String(s))

n.setFlag(gaga.LatinToWide | gaga.AlphaToUpper | gaga.KanaToHiragana)
fmt.Println(n.String(s))
```

Output:

```
ＧａGa is not がｶﾞガ
GaGa is not がガガ
ＧＡＧＡ　ＩＳ　ＮＯＴ　ががが
```

### Vert

```
import "github.com/y-bash/go-gaga"

s := gaga.Vert("閑さや\n岩にしみ入る\n蝉の声")
fmt.Print(s)
```

Output:

```
蝉岩閑
のにさ
声しや
  み  
  入  
  る
```

### Commands

#### Norm

```
$ echo "ＡＢＣｱｲｳ" | norm
ABCアイウ
```

#### Vert

```
$ echo -e "閑さや\n岩にしみ入る\n蝉の声" | vert
蝉岩閑
のにさ
声しや
  み  
  入  
  る
```

### Norm & Vert

```
$ echo -e "閑さや\n岩にしみ入る\n蝉の声" | norm -flag KanaToWideKatakana | vert
蝉岩閑
ノニサ
声シヤ
  ミ
  入
  ル
```

## License
MIT

## Author
y-bash

