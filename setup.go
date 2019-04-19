package sgf

import (
	"fmt"
	"strconv"
)

func NewTree(size int) *Node {

	// Creates a new root node.

	if size < 1 || size > 52 {
		panic(fmt.Sprintf("NewTree(): invalid size %v", size))
	}

	properties := make(map[string][]string)

	properties["GM"] = []string{"1"}
	properties["FF"] = []string{"4"}
	properties["SZ"] = []string{strconv.Itoa(size)}

	return NewNode(nil, properties)
}

func NewSetup(size int, black, white []string, next_player Colour) *Node {

	// Creates a new root node, with handicap (or other) stones.

	if size < 1 || size > 52 {
		panic(fmt.Sprintf("NewSetup(): invalid size %v", size))
	}

	properties := make(map[string][]string)

	properties["GM"] = []string{"1"}
	properties["FF"] = []string{"4"}
	properties["SZ"] = []string{strconv.Itoa(size)}

	if next_player == WHITE {
		properties["PL"] = []string{"W"}
	} else if next_player == BLACK {
		properties["PL"] = []string{"B"}
	}

	if len(black) > 0 {
		properties["AB"] = []string{}
	}

	if len(white) > 0 {
		properties["AW"] = []string{}
	}

	for _, p := range black {
		properties["AB"] = append(properties["AB"], p)
	}

	for _, p := range white {
		properties["AW"] = append(properties["AW"], p)
	}

	return NewNode(nil, properties)
}
