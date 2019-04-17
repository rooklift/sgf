package sgf

import (
	"fmt"
)

func (self *Board) GroupSize(p Point) int {

	// If the point is empty, should this return 0, or the size of the empty "string"? Hmm.

	if self.State[p.X][p.Y] == EMPTY {
		return 0
	}

	touched := make(map[Point]bool)
	return self.group_size_recurse(p, touched)
}

func (self *Board) group_size_recurse(p Point, touched map[Point]bool) int {

	touched[p] = true
	colour := self.State[p.X][p.Y]

	count := 1

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.State[a.X][a.Y] == colour {
			if touched[a] == false {
				count += self.group_size_recurse(a, touched)
			}
		}
	}

	return count
}

func (self *Board) HasLiberties(p Point) bool {		// Faster than checking if Liberties() == 0
	touched := make(map[Point]bool)
	return self.has_liberties_recurse(p, touched)
}

func (self *Board) has_liberties_recurse(p Point, touched map[Point]bool) bool {

	touched[p] = true
	colour := self.State[p.X][p.Y]

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.State[a.X][a.Y] == EMPTY {
			return true
		} else if self.State[a.X][a.Y] == colour {
			if touched[a] == false {
				if self.has_liberties_recurse(a, touched) {
					return true
				}
			}
		}
	}

	return false
}

func (self *Board) Liberties(p Point) int {

	// What on earth is the correct answer to how many liberties an empty square has?

	if self.State[p.X][p.Y] == EMPTY {
		return -1
	}

	touched := make(map[Point]bool)
	return self.liberties_recurse(p, touched)
}

func (self *Board) liberties_recurse(p Point, touched map[Point]bool) int {

	// Note that this function uses the touched map in a different way from others.
	// Literally every point that's examined is flagged as touched.

	touched[p] = true
	colour := self.State[p.X][p.Y]

	count := 0

	for _, a := range AdjacentPoints(p, self.Size) {
		if touched[a] == false {
			touched[a] = true							// This is fine regardless of what's on the point
			if self.State[a.X][a.Y] == EMPTY {
				count += 1
			} else if self.State[a.X][a.Y] == colour {
				count += self.liberties_recurse(a, touched)
			}
		}
	}

	return count
}

func (self *Board) Singleton(p Point) bool {

	colour := self.State[p.X][p.Y]

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.State[a.X][a.Y] == colour {
			return false
		}
	}

	return true
}

func (self *Board) ko_square_finder(p Point) Point {

	// Only called when we know there is indeed a ko.
	// Argument is the location of the capturing stone that caused it.

	var hits []Point

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.State[a.X][a.Y] == EMPTY {
			hits = append(hits, a)
		}
	}

	if len(hits) != 1 {
		panic(fmt.Sprintf("ko_square_finder(): got %d hits", hits))
	}

	return hits[0]
}