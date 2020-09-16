# go-gaga (Japanese language utility)

## Installation

For using the library:

```
$ go get github.com/y-bash/go-gaga
```

Next, to install the command (execution binary)

```
$ cd $GOPATH/github.com/y-bash/go-gaga
$ make install
```

## Usage

### Library

#### Norm

```
import "github.com/y-bash/go-gaga"

n := gaga.Norm(gaga.LatinToNarrow | gaga.KanaToWide)
s := n.String("ＡＢＣｱｲｳ")
fmt.Println(s)

```

Output:

```
ABCアイウ
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

## License
MIT

## Author
y-bash

