package sgf

import (
	"fmt"
)

type Board struct {					// Contains everything about a go position, except superko stuff
	size				int
	player				Colour
	ko					string

	state				[][]Colour
	captures_by			map[Colour]int
}

func (self *Board) Size() int {
	return self.size
}

func (self *Board) Player() Colour {
	return self.player
}

func (self *Board) HasKo() bool {
	return self.ko != ""
}

func (self *Board) GetKo() string {
	return self.ko
}

func (self *Board) set_ko(s string) {
	if Onboard(s, self.size) == false {
		self.ko = ""
	} else {
		self.ko = s
	}
}

func (self *Board) clear_ko() {
	self.ko = ""
}

func (self *Board) GetState(p string) Colour {
	x, y, onboard := XYFromSGF(p, self.size)
	if onboard == false {
		return EMPTY
	}
	return self.state[x][y]
}

func (self *Board) set_state(p string, c Colour) {
	x, y, onboard := XYFromSGF(p, self.size)
	if onboard == false {
		return
	}
	self.state[x][y] = c
}

func (self *Board) CapturesBy(c Colour) int {
	return self.captures_by[c]
}

// -----------------------------------------------------------------------------

func new_board(sz int) *Board {

	if sz < 1 || sz > 52 {
		panic(fmt.Sprintf("new_board(): bad size %d", sz))
	}

	board := new(Board)

	board.size = sz
	board.player = BLACK
	board.clear_ko()

	board.state = make([][]Colour, sz)
	for x := 0; x < sz; x++ {
		board.state[x] = make([]Colour, sz)
	}

	board.captures_by = make(map[Colour]int)
	board.captures_by[BLACK] = 0					// Not strictly
	board.captures_by[WHITE] = 0					// necessary...

	return board
}

func (self *Board) Copy() *Board {

	ret := new(Board)

	// Easy stuff...

	ret.size = self.size
	ret.player = self.player
	ret.ko = self.ko

	// State...

	ret.state = make([][]Colour, ret.size)
	for x := 0; x < ret.size; x++ {
		ret.state[x] = make([]Colour, ret.size)
		for y := 0; y < ret.size; y++ {
			ret.state[x][y] = self.state[x][y]
		}
	}

	// Captures...

	ret.captures_by = make(map[Colour]int)
	ret.captures_by[BLACK] = self.captures_by[BLACK]
	ret.captures_by[WHITE] = self.captures_by[WHITE]

	return ret
}

func (self *Board) Dump() {

	ko_x, ko_y, _ := XYFromSGF(self.GetKo(), self.size)		// Usually -1, -1

	for y := 0; y < self.size; y++ {
		for x := 0; x < self.size; x++ {
			c := self.state[x][y]
			if c == BLACK {
				fmt.Printf(" X")
			} else if c == WHITE {
				fmt.Printf(" O")
			} else if ko_x == x && ko_y == y {
				fmt.Printf(" :")
			} else {
				fmt.Printf(" .")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
	fmt.Printf("Captures by Black: %d\n", self.captures_by[BLACK])
	fmt.Printf("Captures by White: %d\n", self.captures_by[WHITE])
	fmt.Printf("\n")
	fmt.Printf("Next to play: %v\n", ColourLongNames[self.player])
}
