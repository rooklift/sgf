package sgf

// Note: boards are created only as needed, and some SGF manipulation
// can be done creating no boards whatsoever.

import (
	"strconv"
)

var TotalBoardsGenerated int			// For debugging.

func (self *Node) Board() *Board {

	// Returns a __COPY__ of the cached board for this node, creating that if needed.
	//
	// The cache relies on the fact that mutating properties B, W, AB, AW, AE, PL
	// cannot be added to a node after creation.

	if self == nil { panic("Node.Board(): called on nil node") }

	if self.__board_cache == nil {
		if self.parent != nil {
			self.__board_cache = self.parent.Board()
		} else {										// We are root
			sz_string, _ := self.GetValue("SZ")
			sz, _ := strconv.Atoi(sz_string)
			if sz < 1  { sz = 19 }
			if sz > 52 { sz = 52 }						// SGF limit
			self.__board_cache = NewBoard(sz)
		}
		self.__board_cache.update_from_node(self)
		TotalBoardsGenerated++
	}

	return self.__board_cache.Copy()
}

func (self *Board) update_from_node(node *Node) {

	for _, p := range node.props["AB"] {
		if len(p) == 5 && p[2] == ':' {
			self.SetStateFromList(p, BLACK)
		} else {
			self.SetState(p, BLACK)
		}
		self.Player = WHITE
	}

	for _, p := range node.props["AW"] {
		if len(p) == 5 && p[2] == ':' {
			self.SetStateFromList(p, WHITE)
		} else {
			self.SetState(p, WHITE)
		}
		self.Player = BLACK			// Prevails in the event of both AB and AW
	}

	for _, p := range node.props["AE"] {
		if len(p) == 5 && p[2] == ':' {
			self.SetStateFromList(p, EMPTY)
		} else {
			self.SetState(p, EMPTY)
		}
	}

	// Play move: B / W. Note that "moves" which are not valid onboard points are passes.

	for _, p := range node.props["B"] {
		self.PlaceStone(p, BLACK)
		self.Player = WHITE
	}

	for _, p := range node.props["W"] {
		self.PlaceStone(p, WHITE)
		self.Player = BLACK
	}

	// Respect PL property

	pl, _ := node.GetValue("PL")
	if pl == "B" || pl == "b" {
		self.Player = BLACK
	}
	if pl == "W" || pl == "w" {
		self.Player = WHITE
	}
}

func (self *Board) PlaceStone(p string, colour Colour) {

	// Other than sanity checks, there is no legality check here.
	// Nor should there be. This only alters a board, and if called
	// by the user program, will have no effect whatsoever on any node.
	//
	// Instead of this, node.PlayMove() is the correct way to make a
	// new node from an existing one.

	if self == nil { panic("Board.PlaceStone(): called on nil board") }

	if colour != BLACK && colour != WHITE {
		panic("Board.PlaceStone(): no colour")
	}

	self.ClearKo()

	if ValidPoint(p, self.Size) == false {		// Consider this a pass
		return
	}

	self.SetState(p, colour)

	opponent := colour.Opposite()
	caps := 0

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.GetState(a) == opponent {
			if self.HasLiberties(a) == false {
				caps += self.DestroyGroup(a)
			}
		}
	}

	self.CapturesBy[colour] += caps

	// Handle suicide...

	if self.HasLiberties(p) == false {
		suicide_caps := self.DestroyGroup(p)
		self.CapturesBy[opponent] += suicide_caps
	}

	// Work out ko square...

	if caps == 1 {
		if self.Singleton(p) {
			if self.Liberties(p) == 1 {					// Yes, the conditions are met, there is a ko
				self.SetKo(self.ko_square_finder(p))
			}
		}
	}

	return
}

func (self *Board) DestroyGroup(p string) int {			// Returns stones removed.

	if self == nil { panic("Board.DestroyGroup(): called on nil board") }

	colour := self.GetState(p)

	if colour != BLACK && colour != WHITE {				// Also happens if p is off board.
		return 0
	}

	self.SetState(p, EMPTY)
	count := 1

	for _, a := range AdjacentPoints(p, self.Size) {

		if self.GetState(a) == colour {
			count += self.DestroyGroup(a)
		}
	}

	return count
}
