package main

import (
	avi "avid/decoder"
	"fmt"
	"log"
	"math"
	"os"
	dc "pngd/decoder"
	"pngd/rgb"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("usage: pngd <file.avi> --ascii\n")
	}

	path := args[1]
	src, err := os.ReadFile(path)

	ascii := false
	if len(args) == 3 {
		if args[2] == "--ascii" || args[2] == "-a" {
			ascii = true
		}
	}

	if err != nil {
		log.Fatalf("Failed to read file: %s\n%v\n", path, err.Error())
	}

	avidc := avi.NewDecoder(src)
	avidc.Decode()
	frames := avidc.Frames()
	fps := avidc.FPS()


	png := dc.NewDecoder(frames[0])
	png.ValidateSignature()
	_, err = png.Decode()
	if err != nil {
		fmt.Println(err)
	}

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
	img_w := int(png.IHDR().Width)
	img_h := int(png.IHDR().Height)

	scale_x := float64(tw) / float64(img_w)
	scale_y := float64(eff_h) / float64(img_h) 

	scale := math.Min(scale_x, scale_y)

	// offset for centering
	offset_x := (tw - int(float64(img_w) * scale)) / 2
	offset_y := (eff_h - int(float64(img_h) * scale)) / 2

	bpp := int(png.IHDR().Bpp)
	fmt.Println(bpp)

	white := '█'
	black := '░'
	if ascii {
		white = '@'
		black = ' '
	}

	// 1000s / fps -> s per frame
	ms_per_frame := 1000 / fps
	for j := 1; j < len(frames); j++ {
		png.SetSrc(frames[j])
		err := png.ValidateSignature()

		if err != nil {
			log.Fatalf("Invalid PNG signature\n%v\n", err.Error())
		}

		buf, err := png.Decode()

		if err != nil {
			log.Fatalf("Decoding failed\n%v\n", err.Error())
		}

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
		time.Sleep(time.Duration(ms_per_frame) * time.Millisecond)
	}


	for {
		ev := s.PollEvent()
		if _, ok := ev.(*tcell.EventKey); ok {
			break
		}
	}

}

