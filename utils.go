package sgf

import (
	"fmt"
)

type Point struct {
	X				int
	Y				int
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

func SGFFromPoint(x, y int) string {
	const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return fmt.Sprintf("%c%c", alpha[x], alpha[y])
}
