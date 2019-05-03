package sgf

// SetState sets the colour at the specified point. The argument should be an
// SGF coordinate, e.g. "dd". This method has no side effects whatsoever: it has
// no effect on ko status, nor the next player, and no captures are performed.
// Illegal positions can be created.
func (self *Board) SetState(p string, colour Colour) {
	x, y, onboard := ParsePoint(p, self.Size)
	if onboard == false {
		return
	}
	self.State[x][y] = colour
}

// AddStone adjusts the board according to the rules of SGF properties AB, AW,
// and AE, setting the board state without making captures. The argument should
// be an SGF coordinate, e.g. "dd".
//
// Any ko square is cleared. If the colour was BLACK or WHITE, the next player
// is set to be the opposite colour.
func (self *Board) AddStone(p string, colour Colour) {
	self.SetState(p, colour)
	self.ClearKo()
	if colour != EMPTY {
		self.Player = colour.Opposite()
	}
}

// AddList is like AddStone, but expects an SGF points list such as "dd:fg".
func (self *Board) AddList(s string, colour Colour) {
	points := ParsePointList(s, self.Size)
	for _, point := range points {
		self.SetState(point, colour)
	}
	self.ClearKo()
	if colour != EMPTY {
		self.Player = colour.Opposite()
	}
}

// ForceStone adjusts the board according to the rules of SGF properties B and
// W. The argument should be an SGF coordinate, e.g. "dd".
//
// A stone of the specified colour is placed at the given location, and makes
// any resulting captures. Aside from the obvious sanity checks, there are no
// legality checks - ko recaptures will succeed, as will playing on an occupied
// point.
//
// The board's Ko and Player fields are updated. Invalid point strings are
// considered passes.
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

// Play attempts to play at point p, with full legality checks. The argument
// should be an SGF coordinate, e.g. "dd". The colour is determined
// intelligently. If successful, the board is changed. If the move is illegal,
// returns an error.
//
// As a reminder, editing a board has no effect on the node in an SGF tree from
// which it was created (if any).
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
