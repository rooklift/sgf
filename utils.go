package sgf

type Point struct {
	X				int
	Y				int
}

func AdjacentPoints(x, y, size int) []Point {

	var ret []Point

	possibles := []Point{
		Point{x - 1, y},
		Point{x + 1, y},
		Point{x, y - 1},
		Point{x, y + 1},
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

func PointFromSGF(s string, size int) (x int, y int, ok bool) {		// Doesn't actually return a Point{}, hmm...

	// e.g. "cd" --> 2,3

	// If ok == false, that means the move was a pass.
	// i.e. any non-OK string is a pass in SGF, I guess.

	if len(s) < 2 {
		return -1, -1, false
	}

	x = int(s[0]) - 97
	y = int(s[1]) - 97

	ok = x >= 0 && x < size && y >= 0 && y < size

	return x, y, ok
}
