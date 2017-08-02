package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
	"termp/xterm"

	"github.com/nfnt/resize"
	termbox "github.com/nsf/termbox-go"
)

const SPEED = 5

var shiftx, shifty = 0, 0
var ratex, ratey = 1.0, 1.0
var diax, diay int
var ximg *xterm.XtermImage

func main() {
	img := loadImageFile(os.Args[1])

	println("start to work")
	ximg = xterm.NewXtermImage(img)

	termbox.Init()

	diax, diay = termbox.Size()
	initPicture(diax, diay)

	termbox.SetOutputMode(termbox.Output256)
	defer termbox.Close()

loop:
	for {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		render()

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Ch == 0 {
				switch ev.Key {
				case termbox.KeyArrowUp:
					shifty -= SPEED
				case termbox.KeyArrowDown:
					shifty += SPEED
				case termbox.KeyArrowLeft:
					shiftx -= SPEED
				case termbox.KeyArrowRight:
					shiftx += SPEED
				case termbox.KeySpace:
					break loop
				}
			} else {
				switch ev.Ch {
				case 'w', 'W':
					modifyValue(0, 10)
				case 's', 'S':
					modifyValue(0, -10)
				case 'a', 'A':
					modifyValue(-10, 0)
				case 'd', 'D':
					modifyValue(10, 0)
				}
			}
		}
	}
}

func initPicture(x, y int) {
	w, h := ximg.GetSize()

	if w < x && h < y {
		ximg.Resize(uint(w), uint(h), resize.Lanczos3)
		return
	}

	ximg.Resize(uint(x), uint(y), resize.Lanczos3)

}

func modifyValue(vx, vy int) {
	tvx := vx + diax
	tvy := vy + diay

	if tvx < 0 || tvy < 0 {
		return
	}

	diax = tvx
	diay = tvy

	ximg.Resize(uint(diax), uint(diay), resize.Lanczos3)
}

func render() {
	w, h := ximg.GetSize()

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			xtermColor, _ := ximg.At(j, i)
			termbox.SetCell(j-shiftx, i-shifty, ' ', termbox.ColorDefault, termbox.Attribute(xtermColor+1))
		}
	}

	termbox.Flush()
}

func loadImageFile(filename string) image.Image {

	f, err := os.Open(filename)

	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	if strings.HasSuffix(filename, ".png") {
		image, err := png.Decode(f)
		if err != nil {
			panic(err.Error())
		}

		return image
	} else if strings.HasSuffix(filename, ".jpg") {
		image, err := jpeg.Decode(f)

		if err != nil {
			panic(err.Error())
		}

		return image
	}

	panic("Sorry, your suffix is wrong!")
}
