package sgf

// A Node is the fundamental unit in an SGF tree. Nodes are implemented as maps
// of type map[string][]string. In other words, a key can have multiple values,
// all of which are held as strings. These strings are kept in an unescaped
// state; escaping and unescaping is handled during loading and saving of files.
// A node also contains information about the node's parent (if not root) and a
// list of all child nodes.
type Node struct {
	props			map[string][]string
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
	node.props = make(map[string][]string)

	if node.parent != nil {
		node.parent.children = append(node.parent.children, node)
	}

	return node
}

// Copy provides a deep copy of the node with no attached parent or children.
func (self *Node) Copy() *Node {
	ret := new(Node)
	ret.props = self.AllProperties()					// This is a deep copy of the map, so safe to use.
	return ret
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

	for i := 0; i < len(self.props[key]); i++ {			// Ignore if the value already exists.
		if self.props[key][i] == val {
			return
		}
	}

	self.props[key] = append(self.props[key], val)
}

// DeleteKey deletes the given key and all of its values.
func (self *Node) DeleteKey(key string) {
	self.mutor_check(key)								// If key is a MUTOR, clear board caches.
	delete(self.props, key)
}

// DeleteValue checks if the given key in this node has the given value, and
// removes that value, if it does.
func (self *Node) DeleteValue(key, val string) {

	self.mutor_check(key)								// If key is a MUTOR, clear board caches.

	for i := len(self.props[key]) - 1; i >= 0; i-- {
		v := self.props[key][i]
		if v == val {
			self.props[key] = append(self.props[key][:i], self.props[key][i+1:]...)
		}
	}

	if len(self.props[key]) == 0 {
		delete(self.props, key)
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

	list := self.props[key]

	if len(list) == 0 {
		return "", false
	}

	return list[0], true
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
	return len(self.props[key])
}

// AllKeys returns a new slice of strings, containing all the keys that the node
// has.
func (self *Node) AllKeys() []string {

	var ret []string

	for key, _ := range self.props {
		ret = append(ret, key)
	}

	return ret
}

// AllValues returns a new slice of strings, containing all the values that a
// given key has in this node.
func (self *Node) AllValues(key string) []string {

	var ret []string									// Make a new slice so that it's safe to modify.

	for _, val := range self.props[key] {
		ret = append(ret, val)
	}

	return ret
}

// AllProperties returns a deep copy of the entire dictionary in a node.
func (self *Node) AllProperties() map[string][]string {

	ret := make(map[string][]string)

	for key, _ := range self.props {
		ret[key] = self.AllValues(key)					// Will handle the copying.
	}

	return ret
}
