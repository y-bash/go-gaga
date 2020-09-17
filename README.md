# go-gaga (Japanese language utility)

<img width="300" src="https://raw.githubusercontent.com/y-bash/go-gaga/master/gaga.png">
## Installation

### For using the library:

Linux:

```
$ go get github.com/y-bash/go-gaga
```

Windows:

```
>go get github.com/y-bash/go-gaga
```

### Next, to install the command (If you use binaries):

Linux:

```
$ cd $GOPATH/src/github.com/y-bash/go-gaga
$ make install
```

Windows:

```
>cd %GOPATH%\src\github.com\y-bash\go-gaga
>go install ./...
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

#### Vert

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

Linux:

```
$ echo "ＡＢＣｱｲｳ" | norm
ABCアイウ
```

Windows:
(with wecho comand in the gaga)

```
>wecho ＡＢＣｱｲｳ | norm
ABCアイウ
```


#### Vert

Linux:

```
$ echo -e "閑さや\n岩にしみ入る\n蝉の声" | vert
蝉岩閑
のにさ
声しや
  み  
  入  
  る
```

Windows:
(with wecho comand in the gaga)

```
>wecho 閑さや\n岩にしみ入る\n蝉の声 | vert
蝉岩閑
のにさ
声しや
  み
  入
  る
```

#### Norm & Vert

Linux:

```
$ echo -e "閑さや\n岩にしみ入る\n蝉の声" | norm -flag KanaToWideKatakana | vert
蝉岩閑
ノニサ
声シヤ
  ミ
  入
  ル
```

Windows:
(with wecho comand in the gaga)

```
>wecho 閑さや\n岩にしみ入る\n蝉の声 | norm -flag KanaToWideKatakana | vert
蝉岩閑
ノニサ
声シヤ
  ミ
  入
  ル
```

## License

This software is released under the MIT License, see LICENSE.

## Author

y-bash

