package sgf

import (
	"fmt"
	"strconv"
)

var board_cache = make(map[*Node]*Board)

func ClearBoardCache() { board_cache = make(map[*Node]*Board) }

func (self *Node) Board() *Board {

	// The cache relies on the fact that mutating properties B, W, AB, AW, AE cannot
	// be added to a node after creation.
	//
	// Every return should be returning a copy, never the cached thing itself, so
	// that the caller can safely manipulate its copy.

	cached, ok := board_cache[self]

	if ok {
		return cached.Copy()
	}

	var my_board *Board

	if self.Parent != nil {
		my_board = self.Parent.Board()
	} else {
		// We are root.
		sz_string, _ := self.GetValue("SZ")
		sz, _ := strconv.Atoi(sz_string)
		if sz < 1  { sz = 19 }
		if sz > 52 { sz = 52 }		// SGF limit
		my_board = new_board(sz)
	}

	my_board.update(self)
	board_cache[self] = my_board
	return my_board.Copy()
}

func (self *Board) update(node *Node) {

	for _, p := range node.Props["AB"] {
		self.SetState(p, BLACK)
		self.Player = WHITE
	}

	for _, p := range node.Props["AW"] {
		self.SetState(p, WHITE)
		self.Player = BLACK			// Prevails in the event of both AB and AW
	}

	for _, p := range node.Props["AE"] {
		self.SetState(p, EMPTY)
	}

	// Play move: B / W. Note that "moves" which are not valid onboard points are passes.

	for _, p := range node.Props["B"] {
		self.modify_with_move(p, BLACK)
		self.Player = WHITE
	}

	for _, p := range node.Props["W"] {
		self.modify_with_move(p, WHITE)
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

func (self *Board) modify_with_move(p string, colour Colour) {

	// Other than sanity checks, there is no legality check here.
	// Nor should there be.

	if colour != BLACK && colour != WHITE {
		panic("modify_with_move(): no colour")
	}

	x, y, onboard := XYFromSGF(p, self.Size)

	if onboard == false {		// Consider this a pass
		return
	}

	self.SetState(p, colour)

	opponent := colour.Opposite()
	caps := 0

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.GetState(a) == opponent {
			if self.HasLiberties(a) == false {
				caps += self.destroy_group(a)
			}
		}
	}

	self.CapturesBy[colour] += caps

	// Handle suicide...

	if self.HasLiberties(p) == false {
		suicide_caps := self.destroy_group(p)
		self.CapturesBy[opponent] += suicide_caps
	}

	// Work out ko square...

	self.clear_ko()

	if caps == 1 {
		if self.Singleton(p) {
			if self.Liberties(p) == 1 {					// Yes, the conditions are met, there is a ko
				self.set_ko(self.ko_square_finder(p))
			}
		}
	}

	return nil
}

func (self *Board) destroy_group(p string) int {		// Returns stones removed.

	colour := self.GetState(p)

	if colour != BLACK && colour != WHITE {				// Also happens if p is off board.
		return 0
	}

	self.SetState(p, EMPTY)
	count := 1

	for _, a := range AdjacentPoints(p, self.Size) {

		if self.GetState(a) == colour {
			count += self.destroy_group(a)
		}
	}

	return count
}
