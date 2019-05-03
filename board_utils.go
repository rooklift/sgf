package sgf

// GroupSize counts the size of the group at the given location. The argument
// should be an SGF coordinate, e.g. "dd".
func (self *Board) GroupSize(p string) int {

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

// HasLiberties checks whether the group at the given location has any
// liberties. The argument should be an SGF coordinate, e.g. "dd". For groups of
// stones on normal boards, this is always true, but can be false if the calling
// program is manipulating the board directly.
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

// GroupSize counts the liberties of the group at the given location. The
// argument should be an SGF coordinate, e.g. "dd".
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
