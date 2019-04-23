package sgf

import (
	"fmt"
	"os"
	"strconv"
)

const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// NewTree returns a root node for a game of the given board size, with various
// sensible properties.
func NewTree(size int) *Node {

	// Creates a new root node with standard properties.

	if size < 1 || size > 52 {
		panic(fmt.Sprintf("NewTree(): invalid size %v", size))
	}

	node := NewNode(nil)

	node.SetValue("GM", "1")
	node.SetValue("FF", "4")
	node.SetValue("SZ", strconv.Itoa(size))

	return node
}

// AdjacentPoints returns a slice of all points (formatted as SGF coordinates, e.g. "dd")
// that are adjacent to the given point, on the given board size.
func AdjacentPoints(s string, size int) []string {

	x, y, onboard := ParsePoint(s, size)

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

// ParsePoint takes an SGF coordinate (e.g. "dd") and a board size, and returns
// the x and y values (zeroth-indexed) of that point, as well as a boolean value
// indicating whether the coordinates were on the board. If they were not, the
// coordinates returned are always -1, -1.
func ParsePoint(s string, size int) (x, y int, onboard bool) {

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

// ValidPoint takes an SGF coordinate (e.g. "dd") and a board size, and returns
// a boolean indicating whether the coordinate is on the board. Internally, the
// library considers all moves that fail this test to be passes.
func ValidPoint(s string, size int) bool {
	_, _, onboard := ParsePoint(s, size)
	return onboard
}

// Point generates an SGF coordinate (e.g. "dd") from x and y values. The
// arguments are considered zeroth-indexed.
func Point(x, y int) string {
	if x < 0 || x >= 52 || y < 0 || y >= 52 {
		return ""
	}
	return string(alpha[x]) + string(alpha[y])
}

// IsStarPoint takes an SGF coordinate (e.g. "dd") and a board size, and returns
// true if it would be considered a star (hoshi) point.
func IsStarPoint(p string, size int) bool {

	if size < 9 {
		return false
	}

	x, y, onboard := ParsePoint(p, size)

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

// HandicapPoints19 returns a slice of SGF coordinates (e.g. "dd") that
// are Black's handicap stones in a 19x19 game of go, for the specified
// handicap. The tygem argument indicates whether the 3rd stone in an
// H3 game should be in the top left.
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

// ParsePointList takes an SGF rectangle (e.g. "dd:fg") and a board size, and
// returns a slice containing all points indicated.
func ParsePointList(s string, size int) []string {

	if len(s) != 5 || s[2] != ':' {
		return nil
	}

	first := s[:2]
	second := s[3:]

	x1, y1, onboard1 := ParsePoint(first, size)
	x2, y2, onboard2 := ParsePoint(second, size)

	if onboard1 == false || onboard2 == false {
		return nil
	}

	if x1 > x2 {
		x1, x2 = x2, x1
	}

	if y1 > y2 {
		y1, y2 = y2, y1
	}

	var ret []string

	for x := x1; x <= x2; x++ {					// <= is correct
		for y := y1; y <= y2; y++ {				// <= is correct
			ret = append(ret, Point(x, y))
		}
	}

	return ret
}

// LoadArgOrQuit loads the filename given in os.Args[n] and returns the root
// node. If this is not possible, the program exits.
func LoadArgOrQuit(n int) *Node {

	if len(os.Args) <= n {
		fmt.Printf("LoadArgOrQuit(): no such arg\n")
		os.Exit(1)
	}

	node, err := Load(os.Args[n])

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	return node
}
