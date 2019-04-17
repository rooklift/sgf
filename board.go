package sgf

import (
	"fmt"
	"strconv"
)

type Board struct {					// Contains everything about a go position, except superko stuff
	Size				int
	State				[][]Colour
	Player				Colour
	CapturesBy			map[Colour]int

	ko					Point
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
	board.CapturesBy[BLACK] = 0					// Not strictly
	board.CapturesBy[WHITE] = 0					// necessary...

	board.ClearKo()
	return board
}

func (self *Board) Copy() *Board {

	ret := new(Board)

	// Easy stuff...

	ret.Size = self.Size
	ret.Player = self.Player
	ret.ko = self.ko

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

func (self *Board) SetKo(p Point) {
	self.ko = p
}

func (self *Board) ClearKo() {
	self.ko = Point{-1, -1}			// Lame way of storing no ko?
}

func (self *Board) HasKo() bool {
	return self.ko.X >= 0 && self.ko.Y >= 0 && self.ko.X < self.Size && self.ko.Y < self.Size
}

func (self *Board) GetKo() Point {
	if self.HasKo() == false {
		return Point{-1, -1}
	}
	return self.ko
}

func (self *Node) BoardFromScratch() *Board {

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
		board.Update(node)
	}

	return board
}

func (self *Board) Update(node *Node) {

	if node == nil {
		panic("board.Update(): called with nil node")
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
		return fmt.Errorf("modify_with_move: bad colour")
	}

	if p.X < 0 || p.X >= self.Size || p.Y < 0 || p.Y >= self.Size {
		return fmt.Errorf("modify_with_move: bad coordinates %d,%d (size %d)", p.X, p.Y, self.Size)
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

func (self *Board) Dump() {

	ko := self.GetKo()		// Usually -1, -1

	for y := 0; y < self.Size; y++ {
		for x := 0; x < self.Size; x++ {
			c := self.State[x][y]
			if c == BLACK {
				fmt.Printf(" X")
			} else if c == WHITE {
				fmt.Printf(" O")
			} else if ko.X == x && ko.Y == y {
				fmt.Printf(" :")
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
