package main

import (
	"log"
	"math"
	"os"
	dc "pngd/decoder"
	"pngd/rgb"

	"github.com/gdamore/tcell/v2"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("usage: pngd <file.png>\n")
	}

	path := args[1]
	source, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read file: %s\n%v\n", path, err.Error())
	}

	decoder := dc.NewDecoder(source)
	if err = decoder.ValidateSignature(); err != nil {
		log.Fatalf("Failed to validate png signature\n%s\n", err.Error())
	}

	buf, err := decoder.Decode()
	if err != nil {
		log.Fatalf("Failed to decode source\n%v\n", err.Error())
	}

	log.Printf("Bits per pixel: %d\n", decoder.IHDR().Bpp)

	// write to screen
	s, err := tcell.NewScreen()
    if err != nil {
        panic(err)
    }

    if err := s.Init(); err != nil {
        panic(err)
    }
    defer s.Fini()

    style := tcell.StyleDefault

	tw, th := s.Size()
	eff_h := th*2
	img_w := int(decoder.IHDR().Width)
	img_h := int(decoder.IHDR().Height)

	scale_x := float64(tw) / float64(img_w)
	scale_y := float64(eff_h) / float64(img_h) 

	scale := math.Min(scale_x, scale_y)

	// offset for centering
	offset_x := (tw - int(float64(img_w) * scale)) / 2
	offset_y := (eff_h - int(float64(img_h) * scale)) / 2

	bpp := int(decoder.IHDR().Bpp)

	white := '█'
	black := '░'

	for y := 0; y < th; y++ {
		for x := 0; x < tw; x++ {
			img_x := int(float64(x) / scale)
			img_y := int(float64(y*2) / scale)

			if img_x >= img_w || img_y >= img_h {
				continue
			}
			
			i := (img_y*img_w + img_x) * bpp
			r := buf[i]
			g := buf[i+1]
			b := buf[i+2]

			ch := white
			if rgb.IsBlack(r, g, b) {
				ch = black
			}

			s.SetContent(x + offset_x, y + offset_y, ch, nil, style)

		}
	}

    s.Show()

    for {
        ev := s.PollEvent()
        if _, ok := ev.(*tcell.EventKey); ok {
            break
        }
    }

}

