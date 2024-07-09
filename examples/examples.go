package main

import (
	"github.com/cnphpbb/imgur"
	"golang.org/x/image/colornames"
)

func main() {
	i := imgur.LoadFromUrl("https://go.dev/blog/go-brand/logos.jpg").Resize(600, 400)
	imgur.Canvas(760, 1440, colornames.Deepskyblue).
		Insert(i, (760-600)/2, (1440-400)/10).
		Insert(i, (760-600)/2, (1440-400)/2).
		Insert(i, (760-600)/2, 400*2+(1440-400)/8+5).
		Save("out.png")
}
