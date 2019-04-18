package sgf

import (
	"fmt"
)

func (self *Board) GroupSize(p string) int {

	if self == nil { panic("Board.GroupSize(): called on nil board") }

	// If the point is empty, should this return 0, or the size of the empty "string"? Hmm.

	if self.GetState(p) == EMPTY {
		return 0
	}

	touched := make(map[string]bool)
	return self.group_size_recurse(p, touched)
}

func (self *Board) group_size_recurse(p string, touched map[string]bool) int {

	touched[p] = true
	colour := self.GetState(p)

	count := 1

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.GetState(a) == colour {
			if touched[a] == false {
				count += self.group_size_recurse(a, touched)
			}
		}
	}

	return count
}

func (self *Board) HasLiberties(p string) bool {		// Faster than checking if Liberties() == 0
	if self == nil { panic("Board.HasLiberties(): called on nil board") }
	touched := make(map[string]bool)
	return self.has_liberties_recurse(p, touched)
}

func (self *Board) has_liberties_recurse(p string, touched map[string]bool) bool {

	touched[p] = true
	colour := self.GetState(p)

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.GetState(a) == EMPTY {
			return true
		} else if self.GetState(a) == colour {
			if touched[a] == false {
				if self.has_liberties_recurse(a, touched) {
					return true
				}
			}
		}
	}

	return false
}

func (self *Board) Liberties(p string) int {

	if self == nil { panic("Board.Liberties(): called on nil board") }

	// What on earth is the correct answer to how many liberties an empty square has?

	if self.GetState(p) == EMPTY {
		return -1
	}

	touched := make(map[string]bool)
	return self.liberties_recurse(p, touched)
}

func (self *Board) liberties_recurse(p string, touched map[string]bool) int {

	// Note that this function uses the touched map in a different way from others.
	// Literally every point that's examined is flagged as touched.

	touched[p] = true
	colour := self.GetState(p)

	count := 0

	for _, a := range AdjacentPoints(p, self.Size) {
		if touched[a] == false {
			touched[a] = true							// This is fine regardless of what's on the point
			if self.GetState(a) == EMPTY {
				count += 1
			} else if self.GetState(a) == colour {
				count += self.liberties_recurse(a, touched)
			}
		}
	}

	return count
}

func (self *Board) Singleton(p string) bool {

	if self == nil { panic("Board.Singleton(): called on nil board") }

	colour := self.GetState(p)

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.GetState(a) == colour {
			return false
		}
	}

	return true
}

func (self *Board) ko_square_finder(p string) string {

	// Only called when we know there is indeed a ko.
	// Argument is the location of the capturing stone that caused it.

	var hits []string

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.GetState(a) == EMPTY {
			hits = append(hits, a)
		}
	}

	if len(hits) != 1 {
		panic(fmt.Sprintf("ko_square_finder(): got %d hits", hits))
	}

	return hits[0]
}
