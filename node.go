package sgf

import (
	"fmt"
)

// A Node is the fundamental unit in an SGF tree. Nodes are implemented as maps
// of type map[string][]string. In other words, a key can have multiple values,
// all of which are held as strings. These strings are kept in an unescaped
// state; escaping and unescaping is handled during loading and saving of files.
// A node also contains information about the node's parent (if not root) and a
// list of all child nodes.
type Node struct {
	props			[][]string		// e.g. ["B" "dd"]["C" "good move!"]["TE" "1"]
	children		[]*Node
	parent			*Node

	// Note: generating a board_cache always involves generating all the ancestor
	// board_caches first, so if a board_cache is nil, all the node's descendents
	// will have nil caches as well. We actually rely on this fact in the method
	// clear_board_cache_recursive(). Therefore, to ensure this is so, this should
	// never be set directly except by a very few functions, hence its name.

	__board_cache	*Board
}

// NewNode creates a new node with the specified parent.
func NewNode(parent *Node) *Node {

	node := new(Node)
	node.parent = parent

	if node.parent != nil {
		node.parent.children = append(node.parent.children, node)
	}

	return node
}

// Copy provides a deep copy of the node with no attached parent or children.
func (self *Node) Copy() *Node {

	ret := new(Node)
	ret.props = make([][]string, len(self.props))

	for ki := 0; ki < len(self.props); ki++ {
		ret.props[ki] = make([]string, len(self.props[ki]))
		for j := 0; j < len(self.props[ki]); j++ {
			ret.props[ki][j] = self.props[ki][j]
		}
	}

	return ret
}

func (self *Node) key_index(key string) int {

	for i, slice := range self.props {

		if len(slice) < 2 {
			panic(fmt.Sprintf("key_index(): self.props had a slice with length %d", len(slice)))
		}

		if slice[0] == key {
			return i
		}
	}

	return -1
}

// ------------------------------------------------------------------------------------------------------------------
// IMPORTANT...
// AddValue(), DeleteKey(), and DeleteValue() adjust the properties directly and
// so need to call mutor_check() to see if they are affecting any cached boards.
// ------------------------------------------------------------------------------------------------------------------

// AddValue adds the specified string as a value for the given key. If the value
// already exists for the key, nothing happens.
func (self *Node) AddValue(key, val string) {

	self.mutor_check(key)								// If key is a MUTOR, clear board caches.

	ki := self.key_index(key)
	if ki == -1 {
		self.props = append(self.props, []string{key, val})
		return
	}

	for _, old_val := range self.props[ki][1:] {
		if old_val == val {
			return
		}
	}

	self.props[ki] = append(self.props[ki], val)
}

// DeleteKey deletes the given key and all of its values.
func (self *Node) DeleteKey(key string) {

	ki := self.key_index(key)
	if ki == -1 {
		return
	}

	self.mutor_check(key)								// If key is a MUTOR, clear board caches.

	self.props = append(self.props[:ki], self.props[ki + 1:]...)
}

// DeleteValue checks if the given key in this node has the given value, and
// removes that value, if it does.
func (self *Node) DeleteValue(key, val string) {

	ki := self.key_index(key)
	if ki == -1 {
		return
	}

	self.mutor_check(key)								// If key is a MUTOR, clear board caches.

	for i := len(self.props[ki]) - 1; i >= 1; i-- {		// Use i >= 1 so we don't delete the key itself.
		if self.props[ki][i] == val {
			self.props[ki] = append(self.props[ki][:i], self.props[ki][i + 1:]...)
		}
	}

	// Delete key if needed...

	if len(self.props[ki]) < 2 {
		self.props = append(self.props[:ki], self.props[ki + 1:]...)
	}
}

// ------------------------------------------------------------------------------------------------------------------
// IMPORTANT...
// The rest of the functions are either read-only, or built up from the safe
// functions above. None of these must adjust the properties directly.
// ------------------------------------------------------------------------------------------------------------------

// GetValue returns the first value for the given key, if present, in which case
// ok will be true. Otherwise it returns "" and false.
func (self *Node) GetValue(key string) (val string, ok bool) {
	ki := self.key_index(key)
	if ki == -1 {
		return "", false
	}
	return self.props[ki][1], true
}

// SetValue sets the specified string as the first and only value for the given
// key.
func (self *Node) SetValue(key, val string) {
	self.DeleteKey(key)
	self.AddValue(key, val)
}

// SetValues sets the values of the key to the values provided. The original
// slice remains safe to modify.
func (self *Node) SetValues(key string, values []string) {
	self.DeleteKey(key)
	for _, val := range values {
		self.AddValue(key, val)
	}
}

// KeyCount returns the number of keys a node has.
func (self *Node) KeyCount() int {
	return len(self.props)
}

// ValueCount returns the number of values a key has.
func (self *Node) ValueCount(key string) int {
	ki := self.key_index(key)
	if ki == -1 {
		return 0
	}
	return len(self.props[ki]) - 1
}

// AllKeys returns a new slice of strings, containing all the keys that the node
// has.
func (self *Node) AllKeys() []string {
	var ret []string
	for _, slice := range self.props {
		ret = append(ret, slice[0])
	}
	return ret
}

// AllValues returns a new slice of strings, containing all the values that a
// given key has in this node.
func (self *Node) AllValues(key string) []string {

	ki := self.key_index(key)
	if ki == -1 {
		return nil
	}

	var ret []string									// Make a new slice so that it's safe to modify.

	for _, val := range self.props[ki][1:] {
		ret = append(ret, val)
	}

	return ret
}

