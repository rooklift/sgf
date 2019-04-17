package sgf

import (
	"fmt"
	"strconv"
)

var board_cache = make(map[*Node]*Board)

func (self *Node) Board() *Board {

	// The cache relies on the fact that mutating properties B, W, AB, AW, AE cannot
	// be added to a node after creation.
	//
	// Every return should be returning a copy, never the cached thing itself.

	cached, ok := board_cache[self]

	if ok {
		return cached.Copy()
	}

	var my_board *Board

	if self.Parent != nil {
		my_board = self.Parent.Board().Copy()
	} else {
		// We are root.
		sz_string, _ := self.GetValue("SZ")
		sz, _ := strconv.Atoi(sz_string)
		if sz < 1  { sz = 19 }
		if sz > 52 { sz = 52 }		// SGF limit
		my_board = NewBoard(sz)
	}

	my_board.update(self)
	board_cache[self] = my_board
	return my_board.Copy()
}

func (self *Board) update(node *Node) {

	if node == nil {
		panic("Board.update(): called with nil node")
	}

	for _, foo := range node.Props["AB"] {
		point, ok := PointFromSGF(foo, self.Size)
		if ok {
			self.State[point.X][point.Y] = BLACK
			self.Player = WHITE
		}
	}

	for _, foo := range node.Props["AW"] {
		point, ok := PointFromSGF(foo, self.Size)
		if ok {
			self.State[point.X][point.Y] = WHITE
			self.Player = BLACK			// Prevails in the event of both AB and AW
		}
	}

	for _, foo := range node.Props["AE"] {
		point, ok := PointFromSGF(foo, self.Size)
		if ok {
			self.State[point.X][point.Y] = EMPTY
		}
	}

	// Play move: B / W

	for _, foo := range node.Props["B"] {
		point, ok := PointFromSGF(foo, self.Size)
		if ok {
			self.modify_with_move(BLACK, point)
			self.Player = WHITE
		}
	}

	for _, foo := range node.Props["W"] {
		point, ok := PointFromSGF(foo, self.Size)
		if ok {
			self.modify_with_move(WHITE, point)
			self.Player = BLACK
		}
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

func (self *Board) modify_with_move(colour Colour, p Point) error {

	if colour != BLACK && colour != WHITE {
		return fmt.Errorf("Board.modify_with_move(): bad colour")
	}

	if p.X < 0 || p.X >= self.Size || p.Y < 0 || p.Y >= self.Size {
		return fmt.Errorf("Board.modify_with_move(): bad coordinates %d,%d (size %d)", p.X, p.Y, self.Size)
	}

	opponent := colour.Opposite()

	self.State[p.X][p.Y] = colour

	caps := 0

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.State[a.X][a.Y] == opponent {
			if self.HasLiberties(a) == false {
				caps = self.destroy_group(a)
				self.CapturesBy[colour] += caps
			}
		}
	}

	// Handle suicide...

	if self.HasLiberties(p) == false {
		suicide_caps := self.destroy_group(p)
		self.CapturesBy[opponent] += suicide_caps
	}

	// Work out ko square...

	self.ClearKo()

	if caps == 1 {
		if self.Singleton(p) {
			if self.Liberties(p) == 1 {					// Yes, the conditions are met, there is a ko
				self.SetKo(self.ko_square_finder(p))
			}
		}
	}

	return nil
}

func (self *Board) destroy_group(p Point) int {			// Returns stones removed.

	colour := self.State[p.X][p.Y]

	if colour != BLACK && colour != WHITE {				// Removing this might (conceivably) mess with capture count
		return 0
	}

	self.State[p.X][p.Y] = EMPTY
	count := 1

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.State[a.X][a.Y] == colour {
			count += self.destroy_group(a)
		}
	}

	return count
}
