package util

import "fmt"

func GotExpectedFmt(got any, expected any) string {
	return fmt.Sprintf("got:\t\t%v\nexpected:\t%v", got, expected)
}
