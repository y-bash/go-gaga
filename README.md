# go-gaga (Japanese language utility)

## Installation

```
make install
```

## Usage

### Commands

#### vert


```
$ echo -e "閑さや\n岩にしみ入る\n蝉の声" | vert
蝉岩閑
のにさ
声しや
  み  
  入  
  る
```

### Library

#### Normalizer

```
import "github.com/y-bash/go-gaga"

n:= gaga.Norm(gaga.LatinToNarrow | gaga.KanaToWide)
s := n.String("ＡＢＣｱｲｳ")
fmt.Printf("%q", s) // => "ABCアイウ"

```

## License
MIT

## Author
y-bash

