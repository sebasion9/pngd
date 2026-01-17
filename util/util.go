package util

import "fmt"

func GotExpectedFmt(got any, expected any) string {
	return fmt.Sprintf("got:\t\t%v\nexpected:\t%v", got, expected)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func PaethPredictor(a, b, c byte) byte {
	p := int(a) + int(b) - int(c)
	pa := abs(p - int(a))
	pb := abs(p - int(b))
	pc := abs(p - int(c))

	if pa <= pb && pa <= pc {
		return a
	}
	if pb <= pc {
		return b
	}
	return c
}
