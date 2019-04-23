package sgf

type Colour int8

const (
	EMPTY = Colour(iota)
	BLACK
	WHITE
)

var colour_short_names = map[Colour]string{
	EMPTY: "?",
	BLACK: "B",
	WHITE: "W",
}

var colour_long_names = map[Colour]string {
	EMPTY: "??",
	BLACK: "Black",
	WHITE: "White",
}

// Opposite returns the opposite colour (if called on BLACK or WHITE) otherwise
// it returns EMPTY.
func (c Colour) Opposite() Colour {
	if c == WHITE { return BLACK }
	if c == BLACK { return WHITE }
	return EMPTY
}
