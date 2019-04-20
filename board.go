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

func (self *Board) GetState(p string) Colour {
	if self == nil { panic("Board.GetState(): called on nil board") }
	x, y, onboard := ParsePoint(p, self.Size)
	if onboard == false {
		return EMPTY
	}
	return self.State[x][y]
}

func (self *Board) SetState(p string, c Colour) {
	if self == nil { panic("Board.SetState(): called on nil board") }
	x, y, onboard := ParsePoint(p, self.Size)
	if onboard == false {
		return
	}
	self.State[x][y] = c
}

func (self *Board) SetStateFromList(s string, c Colour) {
	if self == nil { panic("Board.SetStateFromList(): called on nil board") }
	points := ParsePointList(s, self.Size)
	for _, point := range points {
		self.SetState(point, c)
	}
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

func (self *Board) Copy() *Board {

	if self == nil { panic("Board.Copy(): called on nil board") }

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
	if self == nil { panic("Board.HasKo(): called on nil board") }
	return self.Ko != ""
}

func (self *Board) SetKo(p string) {
	if self == nil { panic("Board.SetKo(): called on nil board") }
	if ValidPoint(p, self.Size) == false {
		self.Ko = ""
	} else {
		self.Ko = p
	}
}

func (self *Board) ClearKo() {
	if self == nil { panic("Board.ClearKo(): called on nil board") }
	self.Ko = ""
}

func (self *Board) Dump() {
	if self == nil { panic("Board.Dump(): called on nil board") }
	self.DumpBoard()
	fmt.Printf("\n")
	fmt.Printf("Captures by Black: %d\n", self.CapturesBy[BLACK])
	fmt.Printf("Captures by White: %d\n", self.CapturesBy[WHITE])
	fmt.Printf("\n")
	fmt.Printf("Next to play: %v\n", ColourLongNames[self.Player])
}

func (self *Board) DumpBoard() {

	if self == nil { panic("Board.DumpBoard(): called on nil board") }

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
