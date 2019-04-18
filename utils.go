package sgf

import (
	"fmt"
)

type Point struct {
	X				int
	Y				int
}

func (self Point) String() string {
	return SGFFromPoint(self)
}

func AdjacentPoints(origin Point, size int) []Point {

	var ret []Point

	possibles := []Point{
		Point{origin.X - 1, origin.Y},
		Point{origin.X + 1, origin.Y},
		Point{origin.X, origin.Y - 1},
		Point{origin.X, origin.Y + 1},
	}

	for _, point := range possibles {
		if point.X >= 0 && point.X < size {
			if point.Y >= 0 && point.Y < size {
				ret = append(ret, point)
			}
		}
	}

	return ret
}

func PointFromSGF(s string, size int) (p Point, ok bool) {

	// e.g. "cd" --> 2,3

	// If ok == false, that means the move was a pass.
	// i.e. any non-OK string is a pass in SGF, I guess.

	if len(s) < 2 {
		return Point{-1, -1}, false
	}

	x := int(s[0]) - 97
	y := int(s[1]) - 97

	ok = x >= 0 && x < size && y >= 0 && y < size

	return Point{x, y}, ok
}

func SGFIsPass(s string, size int) bool {		// Our definition of pass
	_, ok := PointFromSGF(s, size)
	return !ok
}

func SGFFromPoint(p Point) string {
	const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return fmt.Sprintf("%c%c", alpha[p.X], alpha[p.Y])
}

func String(x, y int) string {
	return SGFFromPoint(Point{x, y})
}
