package sgf

// Note: internally, strings are kept in an escaped state e.g. \] and \\
// However, when using the API, your functions will send and receive
// unescaped strings.

var MUTORS = []string{"B", "W", "AB", "AW", "AE", "PL"}

type Node struct {
	props			map[string][]string
	children		[]*Node
	parent			*Node

	board_cache		*Board
}

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

func (self *Node) mutor_check(key string) {

	// If the key changes the board, disallow it if we have children.
	// Otherwise, clear the board_cache.

	for _, s := range MUTORS {
		if key == s {
			if len(self.children) > 0 {
				panic("mutor_check(): node has children; so can't change board altering property " + key)
			}
			self.board_cache = nil
			break
		}
	}
}

func (self *Node) AddValue(key, value string) {			// Handles escaping; no other function should!

	if self == nil { panic("Node.AddValue(): called on nil node") }

	self.mutor_check(key)								// If key is a MUTOR, clear board cache or disallow entirely.

	value = escape_string(value)
	for i := 0; i < len(self.props[key]); i++ {			// Ignore if the value already exists.
		if self.props[key][i] == value {
			return
		}
	}

	self.props[key] = append(self.props[key], value)
}

func (self *Node) SetValue(key, value string) {

	if self == nil { panic("Node.SetValue(): called on nil node") }

	// self.mutor_check(key)							// Not needed because AddValue() will call it.

	self.props[key] = nil
	self.AddValue(key, value)
}

func (self *Node) ValueCount(key string) int {
	return len(self.props[key])
}

func (self *Node) GetValue(key string) (value string, ok bool) {

	// Get the __UNESCAPED__ value for the key, on the assumption that there's only 1 value.

	if self == nil { panic("Node.GetValue(): called on nil node") }

	list := self.props[key]

	if len(list) == 0 {
		return "", false
	}

	return unescape_string(list[0]), true
}

func (self *Node) AllValues(key string) []string {

	// Return all __UNESCAPED__ values for the key, possibly zero

	if self == nil { panic("Node.AllValues(): called on nil node") }

	list := self.props[key]

	var ret []string									// Make a new slice to avoid aliasing.

	for _, s := range list {
		ret = append(ret, unescape_string(s))
	}

	return ret
}

func (self *Node) AllProperties() map[string][]string {

	// Return an __UNESCAPED__ copy of the entire dictionary.

	if self == nil { panic("Node.AllProperties(): called on nil node") }

	ret := make(map[string][]string)

	for key, _ := range self.props {
		ret[key] = self.AllValues(key)					// Will handle the unescaping and copying (anti-aliasing).
	}

	return ret
}

func (self *Node) DeleteValue(key, value string) {

	if self == nil { panic("Node.DeleteValue(): called on nil node") }

	self.mutor_check(key)								// If key is a MUTOR, clear board cache or disallow entirely.

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

func (self *Node) DeleteKey(key string) {

	if self == nil { panic("Node.DeleteKey(): called on nil node") }

	self.mutor_check(key)								// If key is a MUTOR, clear board cache or disallow entirely.

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
