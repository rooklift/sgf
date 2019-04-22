package sgf

import (
	"fmt"
)

var HoshiString = "."

// A Board contains information about a Go position. It is possible to generate
// boards from nodes in an SGF tree, but modifying boards created in this way
// has no effect on the SGF nodes themselves.
type Board struct {
	Size				int
	Player				Colour
	Ko					string

	State				[][]Colour
	CapturesBy			map[Colour]int
}

// NewBoard returns an empty board of specified size.
func NewBoard(sz int) *Board {

	if sz < 1 || sz > 52 {
		panic(fmt.Sprintf("NewBoard(): bad size %d", sz))
	}

	board := new(Board)

	board.Size = sz
	board.Player = BLACK
	board.ClearKo()

	board.State = make([][]Colour, sz)
	for x := 0; x < sz; x++ {
		board.State[x] = make([]Colour, sz)
	}

	board.CapturesBy = make(map[Colour]int)
	board.CapturesBy[BLACK] = 0					// Not strictly
	board.CapturesBy[WHITE] = 0					// necessary...

	return board
}

// GetState returns the colour at the specified location. The argument should be
// an SGF coordinate, e.g. "dd".
func (self *Board) GetState(p string) Colour {
	x, y, onboard := ParsePoint(p, self.Size)
	if onboard == false {
		return EMPTY
	}
	return self.State[x][y]
}

// SetState sets the colour at the specified location. The argument should be an
// SGF coordinate, e.g. "dd".
func (self *Board) SetState(p string, c Colour) {
	x, y, onboard := ParsePoint(p, self.Size)
	if onboard == false {
		return
	}
	self.State[x][y] = c
}

// SetStateFromList sets the colour at the specified locations. The argument
// should be an SGF rectangle, e.g. "dd:fg".
func (self *Board) SetStateFromList(s string, c Colour) {
	points := ParsePointList(s, self.Size)
	for _, point := range points {
		self.SetState(point, c)
	}
}

// Copy returns a deep copy of the board.
func (self *Board) Copy() *Board {

	ret := new(Board)

	// Easy stuff...

	ret.Size = self.Size
	ret.Player = self.Player
	ret.Ko = self.Ko

	// State...

	ret.State = make([][]Colour, ret.Size)
	for x := 0; x < ret.Size; x++ {
		ret.State[x] = make([]Colour, ret.Size)
		for y := 0; y < ret.Size; y++ {
			ret.State[x][y] = self.State[x][y]
		}
	}

	// Captures...

	ret.CapturesBy = make(map[Colour]int)
	ret.CapturesBy[BLACK] = self.CapturesBy[BLACK]
	ret.CapturesBy[WHITE] = self.CapturesBy[WHITE]

	return ret
}

// HasKo returns true if the board has a ko square, on which capture by the
// current player to move is prohibited.
func (self *Board) HasKo() bool {
	return self.Ko != ""
}

// SetKo sets the ko square. The argument should be an SGF coordinate, e.g.
// "dd".
func (self *Board) SetKo(p string) {
	if ValidPoint(p, self.Size) == false {
		self.Ko = ""
	} else {
		self.Ko = p
	}
}

// ClearKo removes the ko square, if any.
func (self *Board) ClearKo() {
	self.Ko = ""
}

// Dump prints the board, and some information about captures and next player.
func (self *Board) Dump() {
	self.DumpBoard()
	fmt.Printf("\n")
	fmt.Printf("Captures by Black: %d\n", self.CapturesBy[BLACK])
	fmt.Printf("Captures by White: %d\n", self.CapturesBy[WHITE])
	fmt.Printf("\n")
	fmt.Printf("Next to play: %v\n", ColourLongNames[self.Player])
}

// DumpBoard prints the board.
func (self *Board) DumpBoard() {

	ko_x, ko_y, _ := ParsePoint(self.Ko, self.Size)		// Usually -1, -1

	for y := 0; y < self.Size; y++ {
		for x := 0; x < self.Size; x++ {
			c := self.State[x][y]
			if c == BLACK {
				fmt.Printf(" X")
			} else if c == WHITE {
				fmt.Printf(" O")
			} else if ko_x == x && ko_y == y {
				fmt.Printf(" :")
			} else {
				if IsStarPoint(Point(x, y), self.Size) {
					fmt.Printf(" ")
					fmt.Printf(HoshiString)
				} else {
					fmt.Printf(" .")
				}
			}
		}
		fmt.Printf("\n")
	}
}

// PlaceStone places a stone of the specified colour at the given location. The
// argument should be an SGF coordinate, e.g. "dd". Aside from the obvious
// sanity checks, there are no legality checks. As a reminder, editing a board
// has no effect on the node in an SGF tree from which it was created (if any).
func (self *Board) PlaceStone(p string, colour Colour) {

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

// DestroyGroup deletes the group at the specified location. The argument should
// be an SGF coordinate, e.g. "dd", referring to any stone in the group to be
// destroyed. The number of stones removed is returned.
func (self *Board) DestroyGroup(p string) int {

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
