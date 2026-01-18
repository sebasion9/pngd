package bmp

import (
	"os"
	"image"
	"image/color"

	"golang.org/x/image/bmp"
)

func RGBAToBMP(path string, width, height int, rgba []byte) error {
    img := image.NewRGBA(image.Rect(0, 0, width, height))

    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            i := (y*width + x) * 4
            r := rgba[i]
            g := rgba[i+1]
            b := rgba[i+2]
            a := rgba[i+3]
            img.Set(x, y, color.NRGBA{R: r, G: g, B: b, A: a})
        }
    }

    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()

    return bmp.Encode(f, img)
}

func RGBToBMP(path string, width, height int, rgb []byte) error {
    img := image.NewNRGBA(image.Rect(0, 0, width, height))

    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            i := (y*width + x) * 3
            r := rgb[i]
            g := rgb[i+1]
            b := rgb[i+2]
            img.Set(x, y, color.NRGBA{R: r, G: g, B: b, A: 255}) // alpha = 255
        }
    }

    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()

    return bmp.Encode(f, img)
}
