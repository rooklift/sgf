package sgf

type Colour int8

const (
	EMPTY = Colour(iota)
	BLACK
	WHITE
)

// Opposite returns the opposite colour (if called on BLACK or WHITE) otherwise
// it returns EMPTY.
func (c Colour) Opposite() Colour {
	if c == WHITE { return BLACK }
	if c == BLACK { return WHITE }
	return EMPTY
}

// Upper returns a single byte string, "B" or "W" or "?", for the colour.
func (c Colour) Upper() string {
	if c == WHITE { return "W" }
	if c == BLACK { return "B" }
	return "?"
}

// Lower returns a single byte string, "b" or "w" or "?", for the colour.
func (c Colour) Lower() string {
	if c == WHITE { return "w" }
	if c == BLACK { return "b" }
	return "?"
}

// Word returns a word, "Black" or "White" or "??", for the colour.
func (c Colour) Word() string {
	if c == WHITE { return "White" }
	if c == BLACK { return "Black" }
	return "??"
}
