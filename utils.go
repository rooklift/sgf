package sgf

import (
	"fmt"
)

const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func AdjacentPoints(s string, size int) []string {

	x, y, onboard := XYFromSGF(s, size)

	if onboard == false {
		return nil
	}

	var ret []string

	if x > 0 {
		ret = append(ret, string(alpha[x - 1]) + string(s[1]))		// Left
	}
	if x < size - 1 {
		ret = append(ret, string(alpha[x + 1]) + string(s[1]))		// Right
	}
	if y > 0 {
		ret = append(ret, string(s[0]) + string(alpha[y - 1]))		// Up
	}
	if y < size - 1 {
		ret = append(ret, string(s[0]) + string(alpha[y + 1]))		// Down
	}

	return ret
}

func XYFromSGF(s string, size int) (x, y int, onboard bool) {

	// e.g. "cd" --> 2,3

	// Any string that does not yield an onboard coordinate
	// is considered a pass.

	if len(s) != 2 {
		return -1, -1, false
	}

	x = -1
	y = -1

	if s[0] >= 'a' && s[0] <= 'z' { x = int(s[0]) - 97 }
	if s[1] >= 'a' && s[1] <= 'z' { y = int(s[1]) - 97 }
	if s[0] >= 'A' && s[0] <= 'Z' { x = int(s[0]) - 39 }
	if s[1] >= 'A' && s[1] <= 'Z' { y = int(s[1]) - 39 }

	onboard = x >= 0 && x < size && y >= 0 && y < size

	if onboard == false {
		return -1, -1, false
	} else {
		return x, y, true
	}
}

func Onboard(s string, size int) bool {
	_, _, onboard := XYFromSGF(s, size)
	return onboard
}

func Point(x, y int) string {
	if x < 0 || x >= 52 || y < 0 || y >= 52 {
		return ""
	}
	return fmt.Sprintf("%c%c", alpha[x], alpha[y])
}
