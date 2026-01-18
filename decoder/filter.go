package decoder

import (
	"fmt"
	"pngd/errors"
	"pngd/util"
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
	CompressedScanlines [][]byte
	Recon[][]byte
	Bpp int
}

func (f *Filter) None(idx int) {
	scanline := f.CompressedScanlines[idx][1:]
	recon := make([]byte, len(scanline))
	copy(recon, scanline)
	f.Recon[idx] = recon
}

func (f *Filter) Sub(idx int) {
	scanline := f.CompressedScanlines[idx]
	recon := make([]byte, len(scanline) - 1)
	var a byte
	for i := 0; i < len(recon); i++ {
		if i < f.Bpp {
			a = 0
		} else {
			a = recon[i - f.Bpp]
		}
		recon[i] = scanline[i+1] + a
	}
	f.Recon[idx] = recon
}

func (f *Filter) Up(idx int) {
	if idx == 0 {
		f.None(idx)
		return
	}

	scanline := f.CompressedScanlines[idx]
	prev := f.Recon[idx - 1]
	recon := make([]byte, len(scanline) - 1)
	for i := 0; i < len(recon); i++ {
		recon[i] = scanline[i + 1] + prev[i]
	}
	f.Recon[idx] = recon
}

func (f *Filter) Average(idx int) {
	scanline := f.CompressedScanlines[idx]
	recon := make([]byte, len(scanline) - 1)
	var prev []byte
	if idx > 0 {
		prev = f.Recon[idx - 1]
	} else {
		prev = make([]byte, len(scanline) - 1)
	}

	var a byte
	for i := 0; i < len(recon); i++ {
		if i < f.Bpp {
			a = 0
		} else {
			a = recon[i - f.Bpp]
		}
		avg := (int(a) + int(prev[i])) / 2
		recon[i] = scanline[i + 1] + byte(avg)
	}

	f.Recon[idx] = recon
}

func (f *Filter) Paeth(idx int) {
	scanline := f.CompressedScanlines[idx]
	recon := make([]byte, len(scanline) - 1)

	var prev []byte
	if idx > 0 {
		prev = f.Recon[idx - 1]
	} else {
		prev = make([]byte, len(scanline) - 1)
	}

	var a, b, c byte
	for i := 0; i < len(recon); i++ {
		b = prev[i]
		if i < f.Bpp {
			c = 0
			a = 0
		} else {
			a = recon[i - f.Bpp]
			c = prev[i - f.Bpp]
		}
		recon[i] = scanline[i + 1] + util.PaethPredictor(a, b, c)
	}
	f.Recon[idx] = recon

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

func (f *Filter) Reconstruct(bpp byte) error {
	f.Bpp = int(bpp)
	f.Recon = make([][]byte, len(f.CompressedScanlines))
	for i, cs := range f.CompressedScanlines {

		filter_type, ok := filter_type_map[cs[0]]
		if !ok {
			return errors.NewInvalidFilterError(
				fmt.Sprintf("Byte %d is an invalid filter", cs[0]))
		}
		switch filter_type {
		case NONE:
			f.None(i)
		case SUB:
			f.Sub(i)
		case UP:
			f.Up(i)
		case AVG:
			f.Average(i)
		case PAETH:
			f.Paeth(i)
		default:
			f.None(i)
		}
	}
	return nil
}


