package decoder

import (
	"fmt"
	"pngd/errors"
	"pngd/internal"
)

/*

c b
a x  <- processed byte

a is coresponding byte of left pixel
b is byte above
c is coresponding byte of left and up pixel
r is reconstructed byte
*/

type Filter struct {
	Scanlines [][]byte
	recon[][]byte
	bpp int
}


func (f *Filter) none(idx int) {
	scanline := f.Scanlines[idx][1:]
	recon := make([]byte, len(scanline))
	copy(recon, scanline)
	f.recon[idx] = recon
}

func (f *Filter) sub(idx int) {
	scanline := f.Scanlines[idx]
	recon := make([]byte, len(scanline) - 1)
	var a byte
	for i := 0; i < len(recon); i++ {
		if i < f.bpp {
			a = 0
		} else {
			a = recon[i - f.bpp]
		}
		recon[i] = scanline[i+1] + a
	}
	f.recon[idx] = recon
}

func (f *Filter) up(idx int) {
	if idx == 0 {
		f.none(idx)
		return
	}

	scanline := f.Scanlines[idx]
	prev := f.recon[idx - 1]
	recon := make([]byte, len(scanline) - 1)
	for i := 0; i < len(recon); i++ {
		recon[i] = scanline[i + 1] + prev[i]
	}
	f.recon[idx] = recon
}

func (f *Filter) average(idx int) {
	scanline := f.Scanlines[idx]
	recon := make([]byte, len(scanline) - 1)
	var prev []byte
	if idx > 0 {
		prev = f.recon[idx - 1]
	} else {
		prev = make([]byte, len(scanline) - 1)
	}

	var a byte
	for i := 0; i < len(recon); i++ {
		if i < f.bpp {
			a = 0
		} else {
			a = recon[i - f.bpp]
		}
		avg := (int(a) + int(prev[i])) / 2
		recon[i] = scanline[i + 1] + byte(avg)
	}

	f.recon[idx] = recon
}

func (f *Filter) paeth(idx int) {
	scanline := f.Scanlines[idx]
	recon := make([]byte, len(scanline) - 1)

	var prev []byte
	if idx > 0 {
		prev = f.recon[idx - 1]
	} else {
		prev = make([]byte, len(scanline) - 1)
	}

	var a, b, c byte
	for i := 0; i < len(recon); i++ {
		b = prev[i]
		if i < f.bpp {
			c = 0
			a = 0
		} else {
			a = recon[i - f.bpp]
			c = prev[i - f.bpp]
		}
		recon[i] = scanline[i + 1] + internal.PaethPredictor(a, b, c)
	}
	f.recon[idx] = recon

}

type FilterType int
const (
	NONE FilterType = iota
	SUB
	UP
	AVG
	PAETH
)

var filter_type_map = map[byte]FilterType {
	0:NONE,
	1:SUB,
	2:UP,
	3:AVG,
	4:PAETH,
}

func (f *Filter) Reconstruct(bpp byte) ([]byte,error) {
	f.bpp = int(bpp)
	f.recon = make([][]byte, len(f.Scanlines))
	for i, cs := range f.Scanlines {

		filter_type, ok := filter_type_map[cs[0]]
		if !ok {
			return nil, errors.NewInvalidFilterError(
				fmt.Sprintf("Byte %d is an invalid filter", cs[0]))
		}
		switch filter_type {
		case NONE:
			f.none(i)
		case SUB:
			f.sub(i)
		case UP:
			f.up(i)
		case AVG:
			f.average(i)
		case PAETH:
			f.paeth(i)
		default:
			f.none(i)
		}
	}
	return internal.Flatten(f.recon), nil
}



