package sgf

// Note: boards are created only as needed, and some SGF manipulation
// can be done creating no boards whatsoever.

var mutors = []string{"B", "W", "AB", "AW", "AE", "PL", "SZ"}

var TotalBoardsGenerated int			// For debugging.
var TotalBoardsDeleted int				// For debugging.

// -----------------------------------------------------------------------------------------------
// clear_board_cache_recursive() needs to be called whenever a node's board cache becomes invalid.
// This can be due to:
//
//		* Changing a board-altering property.
//		* Changing the identity of its parent.

func (self *Node) clear_board_cache_recursive() {
	if self.__board_cache == nil {						// If nil, all descendent caches are nil also.
		return											// See note in the Node struct about this.
	}
	self.__board_cache = nil
	TotalBoardsDeleted++
	for _, child := range self.children {
		child.clear_board_cache_recursive()
	}
}

func (self *Node) mutor_check(key string) {

	// If the key changes the board, all descendent boards are also invalid.

	for _, s := range mutors {
		if key == s {
			self.clear_board_cache_recursive()
			break
		}
	}
}

// -----------------------------------------------------------------------------------------------

// Board uses the entire history of the tree up to this point to return a board.
// A copy of the result is cached intelligently; the cached board is also purged
// automatically if it becomes invalid (e.g. because a board-altering property
// changed in a relevant part of the SGF tree). Note that modifying a board has
// no effect on the SGF node which created it.
func (self *Node) Board() *Board {

	// Return cache if it exists...

	if self.__board_cache != nil {
		return self.__board_cache.Copy()
	}

	// Otherwise, generate boards for the line, avoiding deep recursion...
	// We do call Board() but the depth is only ever 2.

	line := self.GetLine()

	for _, node := range line {

		if node.__board_cache != nil {
			continue
		}

		// For a node that doesn't have a cache, first get a copy of its parent board...

		if node.parent == nil {							// node is root, so make new.
			sz := node.RootBoardSize()
			node.__board_cache = NewBoard(sz)
		} else {
			node.__board_cache = node.parent.Board()	// fetch a copy.
		}

		// Now update the node's board from its own SGF properties...

		node.__board_cache.update_from_node(node)
	}

	return self.__board_cache.Copy()
}

func (self *Board) update_from_node(node *Node) {

	for _, p := range node.AllValues("AB") {
		if len(p) == 5 && p[2] == ':' {
			self.SetStateFromList(p, BLACK)
		} else {
			self.SetState(p, BLACK)
		}
		self.Player = WHITE
	}

	for _, p := range node.AllValues("AW") {
		if len(p) == 5 && p[2] == ':' {
			self.SetStateFromList(p, WHITE)
		} else {
			self.SetState(p, WHITE)
		}
		self.Player = BLACK			// Prevails in the event of both AB and AW
	}

	for _, p := range node.AllValues("AE") {
		if len(p) == 5 && p[2] == ':' {
			self.SetStateFromList(p, EMPTY)
		} else {
			self.SetState(p, EMPTY)
		}
	}

	// Play move: B / W. Note that "moves" which are not valid onboard points are passes.

	for _, p := range node.AllValues("B") {
		self.PlaceStone(p, BLACK)
		self.Player = WHITE
	}

	for _, p := range node.AllValues("W") {
		self.PlaceStone(p, WHITE)
		self.Player = BLACK
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
