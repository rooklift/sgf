package sgf

// A Node is the fundamental unit in an SGF tree. Nodes are implemented as maps
// of type map[string][]string. In other words, a key can have multiple values,
// all of which are held as strings. Internally, these strings are kept in an
// escaped state, -- \] and \\ -- however callers must send and receive
// unescaped strings; any required escaping and unescaping is handled
// automatically. A node also contains information about the node's parent (if
// not root) and a list of all child nodes.
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

// -----------------------------------------------------------------------------

// AddValue adds the specified string as a value for the given key. Escaping is
// handled automatically. If the value already exists for the key, nothing
// happens.
func (self *Node) AddValue(key, value string) {			// Handles escaping; no other function should!

	self.mutor_check(key)								// If key is a MUTOR, clear board caches.

	value = escape_string(value)
	for i := 0; i < len(self.props[key]); i++ {			// Ignore if the value already exists.
		if self.props[key][i] == value {
			return
		}
	}

	self.props[key] = append(self.props[key], value)
}

// SetValue sets the specified string as the first and only value for the given
// key. Escaping is handled automatically.
func (self *Node) SetValue(key, value string) {

	// self.mutor_check(key)							// Not needed because AddValue() will call it.

	self.props[key] = nil
	self.AddValue(key, value)
}

// ValueCount returns the number of values a key has.
func (self *Node) ValueCount(key string) int {
	return len(self.props[key])
}

// GetValue returns the first value for the given key, if present, in which case
// ok will be true. Otherwise it returns "" and false. Any string returned will
// have been automatically unescaped.
func (self *Node) GetValue(key string) (value string, ok bool) {

	list := self.props[key]

	if len(list) == 0 {
		return "", false
	}

	return unescape_string(list[0]), true
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
// given key has in this node. These strings are automatically unescaped.
func (self *Node) AllValues(key string) []string {

	list := self.props[key]

	var ret []string									// Make a new slice so that it's safe to modify.

	for _, s := range list {
		ret = append(ret, unescape_string(s))
	}

	return ret
}

// AllProperties returns a copy of the entire dictionary in a node. All values
// contained are automatically unescaped.
func (self *Node) AllProperties() map[string][]string {

	ret := make(map[string][]string)

	for key, _ := range self.props {
		ret[key] = self.AllValues(key)					// Will handle the unescaping and copying.
	}

	return ret
}

// DeleteValue checks if the given key in this node has the given value, and
// removes that value, if it does.
func (self *Node) DeleteValue(key, value string) {

	self.mutor_check(key)								// If key is a MUTOR, clear board caches.

	value = escape_string(value)

	for i := len(self.props[key]) - 1; i >= 0; i-- {
		v := self.props[key][i]
		if v == value {
			self.props[key] = append(self.props[key][:i], self.props[key][i+1:]...)
		}
	}

	if len(self.props[key]) == 0 {
		delete(self.props, key)
	}
}

// DeleteKey deletes the given key and all of its values.
func (self *Node) DeleteKey(key string) {
	self.mutor_check(key)								// If key is a MUTOR, clear board caches.
	delete(self.props, key)
}

// -----------------------------------------------------------------------------

func escape_string(s string) string {

	// Treating the input as a byte sequence, not a sequence of code points. Meh.

	var new_s []byte

	for n := 0; n < len(s); n++ {
		if s[n] == '\\' || s[n] == ']' {
			new_s = append(new_s, '\\')
		}
		new_s = append(new_s, s[n])
	}

	return string(new_s)
}

func unescape_string(s string) string {

	// Treating the input as a byte sequence, not a sequence of code points. Meh.
	// Some issues with unicode.

	var new_s []byte

	forced_accept := false

	for n := 0; n < len(s); n++ {

		if forced_accept {
			new_s = append(new_s, s[n])
			forced_accept = false
			continue
		}

		if s[n] == '\\' {
			forced_accept = true
			continue
		}

		new_s = append(new_s, s[n])
	}

	return string(new_s)
}
