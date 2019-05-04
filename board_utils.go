package sgf

import (
	"fmt"
)

// Stones returns all stones in the group at point p, in arbitrary order. The
// argument should be an SGF coordinate, e.g. "dd".
func (self *Board) Stones(p string) []string {

	colour := self.Get(p)
	if colour == EMPTY {					// true also if offboard / invalid
		return nil
	}

	touched := make(map[string]bool)
	self.stones_recurse(p, colour, touched)

	var ret []string
	for key, _ := range touched {
		ret = append(ret, key)
	}

	return ret
}

func (self *Board) stones_recurse(p string, colour Colour, touched map[string]bool) {

	touched[p] = true

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.get_fast(a) == colour {
			if touched[a] == false {
				self.stones_recurse(a, colour, touched)
			}
		}
	}
}

// HasLiberties checks whether the group at point p has any liberties. The
// argument should be an SGF coordinate, e.g. "dd". For groups of stones on
// normal boards, this is always true, but can be false if the calling program
// is manipulating the board directly.
//
// If the point p is empty, returns false.
func (self *Board) HasLiberties(p string) bool {

	colour := self.Get(p)
	if colour == EMPTY {					// true also if offboard / invalid
		return false
	}

	touched := make(map[string]bool)
	return self.has_liberties_recurse(p, colour, touched)
}

func (self *Board) has_liberties_recurse(p string, colour Colour, touched map[string]bool) bool {

	touched[p] = true

	for _, a := range AdjacentPoints(p, self.Size) {
		a_colour := self.get_fast(a)
		if a_colour == EMPTY {
			return true
		} else if a_colour == colour {
			if touched[a] == false {
				if self.has_liberties_recurse(a, colour, touched) {
					return true
				}
			}
		}
	}

	return false
}

// Liberties returns the liberties of the group at point p, in arbitrary order.
// The argument should be an SGF coordinate, e.g. "dd".
func (self *Board) Liberties(p string) []string {

	colour := self.Get(p)
	if colour == EMPTY {					// true also if offboard / invalid
		return nil
	}

	touched := make(map[string]bool)
	touched[p] = true						// Note this
	return self.liberties_recurse(p, colour, touched, nil)
}

func (self *Board) liberties_recurse(p string, colour Colour, touched map[string]bool, ret []string) []string {

	// Note that this function uses the touched map in a different way from others.
	// Also note the constant returning and updating of ret since appends are not visible to caller otherwise.

	for _, a := range AdjacentPoints(p, self.Size) {
		t := touched[a]
		if t == false {
			touched[a] = true
			a_colour := self.get_fast(a)
			if a_colour == EMPTY {
				ret = append(ret, a)
			} else if a_colour == colour {
				ret = self.liberties_recurse(a, colour, touched, ret)
			}
		}
	}

	return ret
}

// Singleton returns true if the specified stone is a group of size 1. The
// argument should be an SGF coordinate, e.g. "dd".
func (self *Board) Singleton(p string) bool {

	colour := self.Get(p)
	if colour == EMPTY {					// true also if offboard / invalid
		return false
	}

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.get_fast(a) == colour {
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
		if self.get_fast(a) == EMPTY {
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
			if self.get_fast(a) == colour.Opposite() {
				if len(self.Liberties(a)) == 1 {
					allowed = true
					break
				}
			} else if self.get_fast(a) == colour {
				if len(self.Liberties(a)) >= 2 {
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
