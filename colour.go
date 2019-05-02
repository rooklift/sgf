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
	if c == BLACK { return WHITE }
	if c == WHITE { return BLACK }
	return EMPTY
}

// Upper returns a single byte string, "B" or "W" or "?", for the colour.
func (c Colour) Upper() string {
	if c == BLACK { return "B" }
	if c == WHITE { return "W" }
	return "?"
}

// Lower returns a single byte string, "b" or "w" or "?", for the colour.
func (c Colour) Lower() string {
	if c == BLACK { return "b" }
	if c == WHITE { return "w" }
	return "?"
}

// Word returns a word, "Black" or "White" or "??", for the colour.
func (c Colour) Word() string {
	if c == BLACK { return "Black" }
	if c == WHITE { return "White" }
	return "??"
}
