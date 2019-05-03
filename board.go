package sgf

import (
	"bytes"
	"fmt"
	"os"
)

var HoshiString = "."	// Can be changed. Used when printing the board.

// A Board contains information about a Go position. It is possible to generate
// boards from nodes in an SGF tree, but modifying boards created in this way
// has no effect on the SGF nodes themselves. Creating boards from nodes is
// relatively costly, and should probably be avoided if batch processing many
// files.
//
// The state of points on the board can be altered in 3 different ways:
//
// SetState() - changes array.
//
// ForceStone() - changes array, makes captures, updates Ko and Player.
//
// Play() - performs legality checks, changes array, makes captures, updates Ko
// and Player.
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

// Equals returns true if the two boards are the same, including ko status,
// captures, and next player to move.
func (self *Board) Equals(other *Board) bool {
	if self.Size != other.Size || self.Player != other.Player || self.Ko != other.Ko {
		return false
	}
	if self.CapturesBy[BLACK] != other.CapturesBy[BLACK] || self.CapturesBy[WHITE] != other.CapturesBy[WHITE] {
		return false
	}
	for x := 0; x < self.Size; x++ {
		for y := 0; y < self.Size; y++ {
			if self.State[x][y] != other.State[x][y] {
				return false
			}
		}
	}
	return true
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
// SGF coordinate, e.g. "dd". This method has no effect on ko status, nor next
// player, and no captures are performed. Illegal positions can be created.
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

	// In all our printing, try to build up the whole
	// string first to avoid jerky printouts...

	s := self.String()
	s += fmt.Sprintf("Captures: %d by Black - %d by White\n", self.CapturesBy[BLACK], self.CapturesBy[WHITE])
	s += fmt.Sprintf("Next to play: %v\n", self.Player.Word())

	fmt.Printf(s)
	os.Stdout.Sync()		// Same reasoning.
}

// DumpBoard prints the board; it is equivalent to fmt.Printf(board.String()).
func (self *Board) DumpBoard() {
	fmt.Printf(self.String())
}

// String returns an ASCII representation of the board.
func (self *Board) String() string {

	var b bytes.Buffer

	ko_x, ko_y, _ := ParsePoint(self.Ko, self.Size)		// Usually -1, -1

	for y := 0; y < self.Size; y++ {
		for x := 0; x < self.Size; x++ {
			c := self.State[x][y]
			if c == BLACK {
				b.WriteString(" X")
			} else if c == WHITE {
				b.WriteString(" O")
			} else if ko_x == x && ko_y == y {
				b.WriteString(" :")
			} else {
				if IsStarPoint(Point(x, y), self.Size) {
					b.WriteString(" ")
					b.WriteString(HoshiString)
				} else {
					b.WriteString(" .")
				}
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

// ForceStone places a stone of the specified colour at the given location, and
// makes any resulting captures. The argument should be an SGF coordinate, e.g.
// "dd". Aside from the obvious sanity checks, there are no legality checks - ko
// recaptures will succeed, as will playing on an occupied point.
//
// The board's Ko and Player fields are updated.
//
// As a reminder, editing a board has no effect on the node in an SGF tree from
// which it was created (if any).
func (self *Board) ForceStone(p string, colour Colour) {

	if colour != BLACK && colour != WHITE {
		panic("Board.ForceStone(): no colour")
	}

	self.ClearKo()

	if ValidPoint(p, self.Size) == false {		// Consider this a pass
		self.Player = colour.Opposite()
		return
	}

	self.SetState(p, colour)

	caps := 0

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.GetState(a) == colour.Opposite() {
			if self.HasLiberties(a) == false {
				caps += self.DestroyGroup(a)
			}
		}
	}

	self.CapturesBy[colour] += caps

	// Handle suicide...

	if self.HasLiberties(p) == false {
		suicide_caps := self.DestroyGroup(p)
		self.CapturesBy[colour.Opposite()] += suicide_caps
	}

	// Work out ko square...

	if caps == 1 {
		if self.Singleton(p) {
			if self.Liberties(p) == 1 {					// Yes, the conditions are met, there is a ko
				self.SetKo(self.ko_square_finder(p))
			}
		}
	}

	self.Player = colour.Opposite()
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

// Legal returns true if a play at point p would be legal. The argument should
// be an SGF coordinate, e.g. "dd". The colour is determined intelligently. The
// board is not changed. If false, the reason is given in the error.
func (self *Board) Legal(p string) (bool, error) {
	return self.LegalColour(p, self.Player)
}

// LegalColour is like Legal, except the colour is specified rather than being
// automatically determined.
func (self *Board) LegalColour(p string, colour Colour) (bool, error) {

	if colour != BLACK && colour != WHITE {
		return false, fmt.Errorf("Board.LegalColour(): colour not BLACK or WHITE")
	}

	x, y, onboard := ParsePoint(p, self.Size)

	if onboard == false {
		return false, fmt.Errorf("Board.LegalColour(): invalid or off-board string %q", p)
	}

	if self.State[x][y] != EMPTY {
		return false, fmt.Errorf("Board.LegalColour(): point %q (%v,%v) was not empty", p, x, y)
	}

	if self.Ko == p {
		if colour == self.Player {												// i.e. we've not forced a move by the wrong colour.
			return false, fmt.Errorf("Board.LegalColour(): ko recapture forbidden")
		}
	}

	if self.HasLiberties(p) == false {

		// The move we are playing will have no liberties of its own.
		// Therefore, it will be legal iff it has a neighbour which:
		//
		//		- Is an enemy group with 1 liberty, or
		//		- Is a friendly group with 2 or more liberties.

		allowed := false

		for _, a := range AdjacentPoints(p, self.Size) {
			if self.GetState(a) == colour.Opposite() {
				if self.Liberties(a) == 1 {
					allowed = true
					break
				}
			} else if self.GetState(a) == colour {
				if self.Liberties(a) >= 2 {
					allowed = true
					break
				}
			} else {
				panic("wat")
			}
		}

		if allowed == false {
			return false, fmt.Errorf("Board.LegalColour(): suicide forbidden")
		}
	}

	// The move is legal!

	return true, nil
}

// Play attempts to play at point p, with full legality checks. The argument
// should be an SGF coordinate, e.g. "dd". The colour is determined
// intelligently. If successful, the board is changed. If the move is illegal,
// returns an error.
func (self *Board) Play(p string) error {
	return self.PlayColour(p, self.Player)
}

// PlayColour is like Play, except the colour is specified rather than
// being automatically determined.
func (self *Board) PlayColour(p string, colour Colour) error {
	legal, err := self.LegalColour(p, colour)
	if legal == false {
		return err
	}
	self.ForceStone(p, colour)
	return nil
}

// Pass swaps the identity of the next player, and clears any ko.
func (self *Board) Pass() {
	self.ClearKo()
	self.Player = self.Player.Opposite()
}

func (self *Board) ko_square_finder(p string) string {

	// Only called when we know there is indeed a ko.
	// Argument is the location of the capturing stone that caused it.

	var hits []string

	for _, a := range AdjacentPoints(p, self.Size) {
		if self.GetState(a) == EMPTY {
			hits = append(hits, a)
		}
	}

	if len(hits) != 1 {
		panic(fmt.Sprintf("ko_square_finder(): got %d hits", len(hits)))
	}

	return hits[0]
}



/*
func (self *Board) get_state_fast(p string) Colour {
	x := int(p[0]) - 97
	y := int(p[1]) - 97
	if p[0] <= 'Z' { x = int(p[0]) - 39 }
	if p[1] <= 'Z' { y = int(p[1]) - 39 }
	return self.State[x][y]
}

func (self *Board) set_state_fast(p string, c Colour) {
	x := int(p[0]) - 97
	y := int(p[1]) - 97
	if p[0] <= 'Z' { x = int(p[0]) - 39 }
	if p[1] <= 'Z' { y = int(p[1]) - 39 }
	self.State[x][y] = c
}
*/
