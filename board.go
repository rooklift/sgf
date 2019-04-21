package sgf

import (
	"fmt"
)

var HoshiString = "."

type Board struct {					// Contains everything about a go position, except superko stuff
	Size				int
	Player				Colour
	Ko					string

	State				[][]Colour
	CapturesBy			map[Colour]int
}

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

func (self *Board) GetState(p string) Colour {
	x, y, onboard := ParsePoint(p, self.Size)
	if onboard == false {
		return EMPTY
	}
	return self.State[x][y]
}

func (self *Board) SetState(p string, c Colour) {
	x, y, onboard := ParsePoint(p, self.Size)
	if onboard == false {
		return
	}
	self.State[x][y] = c
}

func (self *Board) SetStateFromList(s string, c Colour) {
	points := ParsePointList(s, self.Size)
	for _, point := range points {
		self.SetState(point, c)
	}
}

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

func (self *Board) HasKo() bool {
	return self.Ko != ""
}

func (self *Board) SetKo(p string) {
	if ValidPoint(p, self.Size) == false {
		self.Ko = ""
	} else {
		self.Ko = p
	}
}

func (self *Board) ClearKo() {
	self.Ko = ""
}

func (self *Board) Dump() {
	self.DumpBoard()
	fmt.Printf("\n")
	fmt.Printf("Captures by Black: %d\n", self.CapturesBy[BLACK])
	fmt.Printf("Captures by White: %d\n", self.CapturesBy[WHITE])
	fmt.Printf("\n")
	fmt.Printf("Next to play: %v\n", ColourLongNames[self.Player])
}

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

func (self *Board) PlaceStone(p string, colour Colour) {

	// Other than sanity checks, there is no legality check here.
	// Nor should there be. This only alters a board, and if called
	// by the user program, will have no effect whatsoever on any node.
	//
	// Instead of this, node.PlayMove() is the correct way to make a
	// new node from an existing one.

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
