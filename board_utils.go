package sgf

import (
	"fmt"
)

// Stones returns all stones in the group at point p, in arbitrary order. The
// argument should be an SGF coordinate, e.g. "dd".
func (self *Board) Stones(p string) []string {

	if self.GetState(p) == EMPTY {
		return nil
	}

	touched := make(map[string]bool)
	self.stones_recurse(p, touched)

	var ret []string
	for key, _ := range touched {
		ret = append(ret, key)
	}

	return ret
}

func (self *Board) stones_recurse(p string, touched map[string]bool) {

	touched[p] = true
	colour := self.GetState(p)

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.GetState(a) == colour {
			if touched[a] == false {
				self.stones_recurse(a, touched)
			}
		}
	}
}

// HasLiberties checks whether the group at point p has any liberties. The
// argument should be an SGF coordinate, e.g. "dd". For groups of stones on
// normal boards, this is always true, but can be false if the calling program
// is manipulating the board directly.
//
// If the point p is empty, returns true if any of its neighbours are also
// empty, otherwise false.
func (self *Board) HasLiberties(p string) bool {
	touched := make(map[string]bool)
	return self.has_liberties_recurse(p, touched)
}

func (self *Board) has_liberties_recurse(p string, touched map[string]bool) bool {

	// Also works if the point p is EMPTY.
	// Offboard p returns false.

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

// Liberties counts the liberties of the group at point p. The argument should
// be an SGF coordinate, e.g. "dd".
func (self *Board) Liberties(p string) int {

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

// Singleton returns true if the specified stone is a group of size 1. The
// argument should be an SGF coordinate, e.g. "dd".
func (self *Board) Singleton(p string) bool {

	colour := self.GetState(p)

	if colour == EMPTY {
		return false
	}

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.GetState(a) == colour {
			return false
		}
	}

	return true
}

// Legal returns true if a play at point p would be legal. The argument should
// be an SGF coordinate, e.g. "dd". The colour is determined intelligently. The
// board is not changed. If false, the reason is given in the error.
func (self *Board) Legal(p string) (bool, error) {
	return self.LegalColour(p, self.Player)
}

// LegalColour is like Legal, except the colour is specified rather than being
// automatically determined.
func (self *Board) LegalColour(p string, colour Colour) (bool, error) {

	if colour != BLACK && colour != WHITE {
		return false, fmt.Errorf("colour not BLACK or WHITE")
	}

	x, y, onboard := ParsePoint(p, self.Size)

	if onboard == false {
		return false, fmt.Errorf("invalid or off-board string %q", p)
	}

	if self.State[x][y] != EMPTY {
		return false, fmt.Errorf("point %q (%v,%v) was not empty", p, x, y)
	}

	if self.Ko == p {
		if colour == self.Player {						// i.e. we've not forced a move by the wrong colour.
			return false, fmt.Errorf("ko recapture at %q (%v,%v) forbidden", p, x, y)
		}
	}

	has_own_liberties := false
	for _, a := range AdjacentPoints(p, self.Size) {
		if self.GetState(a) == EMPTY {
			has_own_liberties = true
			break
		}
	}

	if has_own_liberties == false {

		// The move we are playing will have no liberties of its own.
		// Therefore, it will be legal iff it has a neighbour which:
		//
		//		- Is an enemy group with 1 liberty, or
		//		- Is a friendly group with 2 or more liberties.

		allowed := false

		for _, a := range AdjacentPoints(p, self.Size) {
			if self.GetState(a) == colour.Opposite() {
				if self.Liberties(a) == 1 {
					allowed = true
					break
				}
			} else if self.GetState(a) == colour {
				if self.Liberties(a) >= 2 {
					allowed = true
					break
				}
			} else {
				panic("wat")
			}
		}

		if allowed == false {
			return false, fmt.Errorf("suicide at %q (%v,%v) forbidden", p, x, y)
		}
	}

	// The move is legal!

	return true, nil
}
