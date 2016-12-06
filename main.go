package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
)

var (
	dpi  = flag.Float64("dpi", 102, "screen resolution in Dots Per Inch")
	size = flag.Float64("size", 12, "font size in points")
)

const title = "HelloWorld你好"

func main() {
	flag.Parse()

	f, err := truetype.Parse(gomono.TTF) //gomono.TTF
	if err != nil {
		log.Println(err)
		return
	}

	fg := image.Black

	//原始图片是test.jpeg
	imgb, _ := os.Open("test.jpeg")
	img, imgErr := jpeg.Decode(imgb)
	if imgErr != nil {
		log.Println("0:", imgErr)
		os.Exit(1)
	}
	fmt.Println("imgErr:", imgErr)

	imgW, imgH := img.Bounds().Dx(), img.Bounds().Dy()
	rgba := image.NewRGBA(image.Rect(0, 0, imgW, imgH))

	defer imgb.Close()

	draw.Draw(rgba, rgba.Bounds(), img, image.ZP, draw.Src)

	h := font.HintingNone
	d := &font.Drawer{
		Dst: rgba,
		Src: fg,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    *size,
			DPI:     *dpi,
			Hinting: h,
		}),
	}
	y := imgH - 10
	d.Dot = fixed.Point26_6{
		X: fixed.I(imgW-10) - d.MeasureString(title),
		Y: fixed.I(y),
	}
	d.DrawString(title)

	// Save that RGBA image to disk.
	outFile, err := os.Create("out.png")
	if err != nil {
		log.Println("1:", err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println("2:", err)
		os.Exit(2)
	}
	err = b.Flush()
	if err != nil {
		log.Println("3:", err)
		os.Exit(3)
	}
	fmt.Println("Wrote out.png OK.")
}
