package sgf

func (self *Board) GroupSize(x, y int) int {

	// If the point is empty, should this return 0, or the size of the empty "string"? Hmm.

	if self.State[x][y] == EMPTY {
		return 0
	}

	touched := make(map[Point]bool)
	return self.group_size_recurse(x, y, touched)
}

func (self *Board) group_size_recurse(x, y int, touched map[Point]bool) int {

	touched[Point{x, y}] = true
	colour := self.State[x][y]

	count := 1

	for _, point := range AdjacentPoints(x, y, self.Size) {
		if self.State[point.X][point.Y] == colour {
			if touched[point] == false {
				count += self.group_size_recurse(x, y, touched)
			}
		}
	}

	return count
}

func (self *Board) HasLiberties(x, y int) bool {		// Faster than checking if Liberties() == 0
	touched := make(map[Point]bool)
	return self.has_liberties_recurse(x, y, touched)
}

func (self *Board) has_liberties_recurse(x, y int, touched map[Point]bool) bool {

	touched[Point{x, y}] = true
	colour := self.State[x][y]

	for _, point := range AdjacentPoints(x, y, self.Size) {
		if self.State[point.X][point.Y] == EMPTY {
			return true
		} else if self.State[point.X][point.Y] == colour {
			if touched[point] == false {
				if self.has_liberties_recurse(point.X, point.Y, touched) {
					return true
				}
			}
		}
	}

	return false
}

func (self *Board) Liberties(x, y int) int {

	// What on earth is the correct answer to how many liberties an empty square has?

	if self.State[x][y] == EMPTY {
		return -1
	}

	touched := make(map[Point]bool)
	return self.liberties_recurse(x, y, touched)
}

func (self *Board) liberties_recurse(x, y int, touched map[Point]bool) int {

	// Note that this function uses the touched map in a different way from others.
	// Literally every point that's examined is flagged as touched.

	touched[Point{x, y}] = true
	colour := self.State[x][y]

	count := 0

	for _, point := range AdjacentPoints(x, y, self.Size) {
		if touched[point] == false {
			touched[point] = true							// This is fine regardless of what's on the point
			if self.State[point.X][point.Y] == EMPTY {
				count += 1
			} else if self.State[point.X][point.Y] == colour {
				count += self.liberties_recurse(point.X, point.Y, touched)
			}
		}
	}

	return count
}
