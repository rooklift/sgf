package sgf

import (
	"fmt"
	"strconv"
)

type Board struct {				// Contains everything about a go position, except superko stuff
	Size				int
	State				[][]Colour
	Player				Colour
	CapturesBy			map[Colour]int
	Ko					Point				// Uses -1, -1 for no ko
}

func NewBoard(sz int) *Board {

	if sz < 1 || sz > 52 {
		panic(fmt.Sprintf("NewBoard(): bad size %d", sz))
	}

	board := new(Board)
	board.Size = sz
	board.State = make([][]Colour, sz)
	for x := 0; x < sz; x++ {
		board.State[x] = make([]Colour, sz)
	}
	board.Player = BLACK
	board.CapturesBy = make(map[Colour]int)
	board.Ko = Point{-1,-1}
	return board
}

func (self *Node) Board() *Board {

	var sz int

	line := self.GetLine()
	root := line[0]

	sz_string, _ := root.GetValue("SZ")

	if sz_string != "" {
		sz, _ = strconv.Atoi(sz_string)
	}

	if sz < 1 {
		sz = 19		// You can assume
	}

	if sz > 52 {
		sz = 52		// SGF limit
	}

	board := NewBoard(sz)

	// ----------------------------------------------

	for _, node := range line {

		for _, foo := range node.Props["AB"] {
			x, y, ok := PointFromSGF(foo, sz)
			if ok {
				board.State[x][y] = BLACK
				board.Player = WHITE
			}
		}

		for _, foo := range node.Props["AW"] {
			x, y, ok := PointFromSGF(foo, sz)
			if ok {
				board.State[x][y] = WHITE
				board.Player = BLACK			// Prevails in the event of both AB and AW
			}
		}

		for _, foo := range node.Props["AE"] {
			x, y, ok := PointFromSGF(foo, sz)
			if ok {
				board.State[x][y] = EMPTY
			}
		}

		// Play move: B / W

		for _, foo := range node.Props["B"] {
			x, y, ok := PointFromSGF(foo, sz)
			if ok {
				board.modify_with_move(BLACK, x, y)
				board.Player = WHITE
			}
		}

		for _, foo := range node.Props["W"] {
			x, y, ok := PointFromSGF(foo, sz)
			if ok {
				board.modify_with_move(WHITE, x, y)
				board.Player = BLACK
			}
		}

		// Respect PL property

		pl, _ := node.GetValue("PL")
		if pl == "B" || pl == "b" {
			board.Player = BLACK
		}
		if pl == "W" || pl == "w" {
			board.Player = WHITE
		}
	}

	return board
}

func (self *Board) modify_with_move(colour Colour, x, y int) error {

	if colour != BLACK && colour != WHITE {
		return fmt.Errorf("modify_with_move: bad colour")
	}

	if x < 0 || x >= self.Size || y < 0 || y >= self.Size {
		return fmt.Errorf("modify_with_move: bad coordinates %d,%d (size %d)", x, y, self.Size)
	}

	opponent := colour.Opposite()

	self.State[x][y] = colour

	caps := 0

	for _, point := range AdjacentPoints(x, y, self.Size) {
		if self.State[point.X][point.Y] == opponent {
			if self.HasLiberties(point.X, point.Y) == false {
				caps = self.destroy_group(point.X, point.Y)
				self.CapturesBy[colour] += caps
			}
		}
	}

	// Handle suicide...

	if self.HasLiberties(x, y) == false {
		suicide_caps := self.destroy_group(x, y)
		self.CapturesBy[opponent] += suicide_caps
	}

	// Work out ko square...

	if caps == 1 {
		if self.GroupSize(x, y) == 1 {
			if self.Liberties(x, y) == 1 {
				// FIXME
			}
		}
	}

	return nil
}

func (self *Board) destroy_group(x, y int) int {		// Returns stones removed.

	colour := self.State[x][y]

	if colour != BLACK && colour != WHITE {				// Removing this might (conceivably) mess with capture count
		return 0
	}

	self.State[x][y] = EMPTY
	count := 1

	for _, point := range AdjacentPoints(x, y, self.Size) {
		if self.State[point.X][point.Y] == colour {
			count += self.destroy_group(point.X, point.Y)
		}
	}

	return count
}

func (self *Board) Dump() {
	for y := 0; y < self.Size; y++ {
		for x := 0; x < self.Size; x++ {
			c := self.State[x][y]
			if c == BLACK {
				fmt.Printf(" X")
			} else if c == WHITE {
				fmt.Printf(" O")
			} else {
				fmt.Printf(" .")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
	fmt.Printf("Captures by Black: %d\n", self.CapturesBy[BLACK])
	fmt.Printf("Captures by White: %d\n", self.CapturesBy[WHITE])
	fmt.Printf("\n")
	fmt.Printf("Next to play: %v\n", ColourLongNames[self.Player])
}
