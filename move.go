package sgf

// Play attempts to play the specified move at the node. The argument should be
// an SGF coordinate, e.g. "dd". The colour is determined intelligently.
//
// If successful, a new node is created, and attached as a child. That child is
// then returned and the error is nil. However, if the specified move already
// existed in a child, that child is returned instead and no new node is
// created; the error is still nil. On failure, the original node is returned,
// along with an error. Failure indicates the move was illegal.
//
// Note that passes cannot be played with Play.
func (self *Node) Play(p string) (*Node, error) {
	return self.PlayColour(p, self.Board().Player)							// Uses board info to determine colour.
}

// PlayColour is like Play, except the colour is specified rather than being
// automatically determined.
func (self *Node) PlayColour(p string, colour Colour) (*Node, error) {		// Returns new node on success; self on failure.

	legal, err := self.Board().LegalColour(p, colour)
	if legal == false {
		return self, err
	}

	// Return the already-extant child if there is such a thing...

	key := "B"; if colour == WHITE { key = "W" }

	for _, child := range self.children {
		if child.ValueCount(key) == 1 {											// Ignore any illegal nodes with 2 or more...
			mv, _ := child.GetValue(key)
			if mv == p {
				return child, nil
			}
		}
	}

	new_node := NewNode(self)													// Attaches new_node to self.
	new_node.SetValue(key, p)
	return new_node, nil
}

// Pass passes. The colour is determined intelligently. Normally, a new node is
// created, attached as a child, and returned. However, if the specified pass
// already existed in a child, that child is returned instead and no new node is
// created.
func (self *Node) Pass() *Node {

	if self.__board_cache != nil {
		return self.PassColour(self.Board().Player)
	}

	// Since we don't really need a board, try and not make one...

	keys := self.AllKeys()

	seen_black := false
	seen_white := false

	for _, key := range keys {
		if key == "B" || key == "AB" { seen_black = true }
		if key == "W" || key == "AW" { seen_white = true }
	}

	if seen_black && seen_white == false {
		return self.PassColour(WHITE)
	} else if seen_white && seen_black == false {
		return self.PassColour(BLACK)
	}

	// As a last resort, generate a board...

	return self.PassColour(self.Board().Player)
}

// PassColour is like Pass, except the colour is specified rather than being
// automatically determined.
func (self *Node) PassColour(colour Colour) *Node {

	if colour != BLACK && colour != WHITE {
		panic("Node.PassColour(): no colour specified")							// This is a programming error, so panic, not error.
	}

	board := self.Board()

	key := "B"; if colour == WHITE { key = "W" }

	// Return the already-extant child if there is such a thing...

	for _, child := range self.children {
		if child.ValueCount(key) == 1 {											// Ignore any illegal nodes with 2 or more...
			mv, _ := child.GetValue(key)
			if ValidPoint(mv, board.Size) == false {
				return child
			}
		}
	}

	new_node := NewNode(self)
	new_node.SetValue(key, "")

	return new_node
}
