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

func IsStarPoint(p string, size int) bool {

	if size < 9 {
		return false
	}

	x, y, onboard := XYFromSGF(p, size)

	if onboard == false {
		return false
	}

	var good_x, good_y bool

	if size >= 15 || x == y {
		if x * 2 + 1 == size {
			good_x = true
		}
		if y * 2 + 1 == size {
			good_y = true
		}
	}

	if size >= 12 {
		if x == 3 || x + 4 == size {
			good_x = true
		}
		if y == 3 || y + 4 == size {
			good_y = true
		}
	} else {
		if x == 2 || x + 3 == size {
			good_x = true
		}
		if y == 2 || y + 3 == size {
			good_y = true
		}
	}

	return good_x && good_y
}

func HandicapPoints19(handicap int, tygem bool) []string {

	if handicap > 9 {
		handicap = 9
	}

	var ret []string

	if handicap >= 1 { ret = append(ret, "pd") }
	if handicap >= 2 { ret = append(ret, "dp") }
	if handicap >= 3 { ret = append(ret, "pp") }
	if handicap >= 4 { ret = append(ret, "dd") }

	if handicap >= 6 { ret = append(ret, "dj", "pj") }
	if handicap >= 8 { ret = append(ret, "jd", "jp") }

	if handicap >= 5 && handicap % 2 == 1 { ret = append(ret, "jj") }

	// Tygem seems to put its 3rd handicap stone in the top left...

	if tygem && handicap == 3 {
		ret[2] = "dd"
	}

	return ret
}
