# imgur

## Download/Install

install is to run `go get -u github.com/cnphpbb/imgur`.

`github.com/cnphpbb/imgur`

## examples

```go
package main

import (
	"github.com/cnphpbb/imgur"
	"golang.org/x/image/colornames"
)

func main() {
	i := imgur.LoadFromUrl("https://go.dev/blog/go-brand/logos.jpg").Resize(600, 400)
	imgur.Canvas(760, 1440, colornames.Deepskyblue).
		Insert(i, 80, 1440/3).
		Save("out.png")
}
```
