package rgb

func Luminance(r, g, b byte) byte {
    return byte((54*int(r) + 183*int(g) + 19*int(b)) >> 8)
}

func IsBlack(r, g, b byte) bool {
	return Luminance(r, g, b) < 128
}
