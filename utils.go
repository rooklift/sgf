package sgf

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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
func AdjacentPoints(p string, size int) []string {

	x, y, onboard := ParsePoint(p, size)

	if onboard == false {
		return nil
	}

	var ret []string

	if x > 0 {
		ret = append(ret, byte_to_string(alpha[x - 1]) + byte_to_string(p[1]))		// Left
	}
	if x < size - 1 {
		ret = append(ret, byte_to_string(alpha[x + 1]) + byte_to_string(p[1]))		// Right
	}
	if y > 0 {
		ret = append(ret, byte_to_string(p[0]) + byte_to_string(alpha[y - 1]))		// Up
	}
	if y < size - 1 {
		ret = append(ret, byte_to_string(p[0]) + byte_to_string(alpha[y + 1]))		// Down
	}

	return ret
}

func byte_to_string(b byte) string {	// One cannot do string(b) because string() turns integers into utf-8 strings,
	return string([]byte{b})			// therefore if b > 127 it will necessarily return a string of length >= 2
}

// ParsePoint takes an SGF coordinate (e.g. "dd") and a board size, and returns
// the x and y values (zeroth-indexed) of that point, as well as a boolean value
// indicating whether the coordinates were on the board. If they were not, the
// coordinates returned are always -1, -1.
func ParsePoint(p string, size int) (x, y int, onboard bool) {

	// e.g. "cd" --> 2,3

	// Any string that does not yield an onboard coordinate
	// is considered a pass.

	if len(p) != 2 {
		return -1, -1, false
	}

	x = -1
	y = -1

	if p[0] >= 'a' && p[0] <= 'z' { x = int(p[0]) - 97 }
	if p[1] >= 'a' && p[1] <= 'z' { y = int(p[1]) - 97 }
	if p[0] >= 'A' && p[0] <= 'Z' { x = int(p[0]) - 39 }
	if p[1] >= 'A' && p[1] <= 'Z' { y = int(p[1]) - 39 }

	onboard = x >= 0 && x < size && y >= 0 && y < size

	if onboard == false {
		return -1, -1, false
	} else {
		return x, y, true
	}
}

// ValidPoint takes an SGF coordinate (e.g. "dd") and a board size, and returns
// a boolean indicating whether the coordinate is on the board. Internally, the
// library considers all moves that fail this test to be pass-moves.
func ValidPoint(p string, size int) bool {
	_, _, onboard := ParsePoint(p, size)
	return onboard
}

// Point generates an SGF coordinate (e.g. "dd") from x and y values. The
// arguments are considered zeroth-indexed.
func Point(x, y int) string {
	if x < 0 || x >= 52 || y < 0 || y >= 52 {
		return ""
	}
	return byte_to_string(alpha[x]) + byte_to_string(alpha[y])
}

// HandicapPoints returns a slice of SGF coordinates (e.g. "dd") that are
// Black's handicap stones, for the specified boardsize and handicap (max
// handicap: 9). The tygem argument indicates whether the 3rd stone in an H3
// game should be in the top left. Works poorly for very small board sizes.
func HandicapPoints(boardsize, handicap int, tygem bool) []string {

	if boardsize < 4 || handicap < 2 {
		return nil
	}

	if handicap > 9 {
		handicap = 9
	}

	d := 1; if boardsize >= 7 { d = 2 }; if boardsize >= 13 { d = 3 }

	z := boardsize

	var ret []string

	if handicap >= 2 {
		ret = append(ret, Point(z - d - 1, d))
		ret = append(ret, Point(d, z - d - 1))
	}

	if handicap >= 3 {
		if tygem {
			ret = append(ret, Point(d, d))
		} else {
			ret = append(ret, Point(z - d - 1, z - d - 1))
		}
	}

	if handicap >= 4 {
		if tygem {
			ret = append(ret, Point(z - d - 1, z - d - 1))
		} else {
			ret = append(ret, Point(d, d))
		}
	}

	if boardsize % 2 == 0 {
		return ret
	}

	if handicap == 5 || handicap == 7 || handicap == 9 {
		ret = append(ret, Point(z / 2, z / 2))
	}

	if handicap >= 6 {
		ret = append(ret, Point(d, z / 2))
		ret = append(ret, Point(z - d - 1, z / 2))
	}

	if handicap >= 8 {
		ret = append(ret, Point(z / 2, d))
		ret = append(ret, Point(z / 2, z - d - 1))
	}

	return ret
}

// IsStarPoint takes an SGF coordinate (e.g. "dd") and a board size, and returns
// true if it would be considered a star (hoshi) point.
func IsStarPoint(p string, size int) bool {

	starpoints := HandicapPoints(size, 9, false)

	for _, hoshi := range starpoints {
		if p == hoshi {
			return true
		}
	}

	return false
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

// ParseGTP takes a GTP formatted string (e.g. "D16") and a board size, and
// returns the SGF coordinate (e.g. "dd") or "" if invalid.
func ParseGTP(s string, size int) string {

	if len(s) < 2 || len(s) > 3 {
		return ""
	}

	s = strings.ToUpper(s)

	if s[0] < 'A' || s[0] > 'Z' {
		return ""
	}
	if s[1] < '0' || s[1] > '9' {
		return ""
	}
	if len(s) == 3 && (s[2] < '0' || s[2] > '9') {
		return ""
	}

	x := int(s[0]) - 65
	if x >= 8 {					// Adjust for missing "I"
		x--
	}

	up, _ := strconv.Atoi(s[1:])
	y := size - int(up)

	if x < 0 || x >= size || y < 0 || y >= size {
		return ""
	}

	return Point(x, y)
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
