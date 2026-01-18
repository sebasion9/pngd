package internal

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

func Flatten(rows [][]byte) []byte {
    var total int
    for _, row := range rows {
        total += len(row)
    }

    flat := make([]byte, 0, total)
    for _, row := range rows {
        flat = append(flat, row...)
    }
    return flat
}
