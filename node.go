package sgf

// Note: internally, strings are kept in an escaped state e.g. \] and \\
// However, when using the API, your functions will send and receive
// unescaped strings.

var MUTORS = []string{"B", "W", "AB", "AW", "AE", "PL"}

// -----------------------------------------------------------------------------

type Node struct {
	props			map[string][]string
	children		[]*Node
	parent			*Node

	board_cache		*Board
}

func NewNode(parent *Node, props map[string][]string) *Node {

	node := new(Node)
	node.parent = parent
	node.props = make(map[string][]string)

	for key, _ := range props {
		for _, s := range props[key] {
			node.add_value(key, s)
		}
	}

	if node.parent != nil {
		node.parent.children = append(node.parent.children, node)
	}

	return node
}

// -----------------------------------------------------------------------------

func (self *Node) add_value(key, value string) {			// Handles escaping; no other function should

	value = escape_string(value)

	for i := 0; i < len(self.props[key]); i++ {				// Ignore if the value already exists
		if self.props[key][i] == value {
			return
		}
	}

	self.props[key] = append(self.props[key], value)
}

func (self *Node) AddValue(key, value string) {

	if self == nil { panic("Node.AddValue(): called on nil node") }

	// Disallow keys that change the board...

	for _, s := range MUTORS {
		if key == s {
			panic("Node.AddValue(): Can't change board-altering properties")
		}
	}

	self.add_value(key, value)
}

func (self *Node) SetValue(key, value string) {

	if self == nil { panic("Node.SetValue(): called on nil node") }

	// Disallow keys that change the board...

	for _, s := range MUTORS {
		if key == s {
			panic("Node.SetValue(): Can't change board-altering properties")
		}
	}

	self.props[key] = nil
	self.add_value(key, value)
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

	var ret []string		// Make a new slice to avoid aliasing.

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
		ret[key] = self.AllValues(key)		// Will handle the unescaping and copying (anti-aliasing).
	}

	return ret
}

func (self *Node) DeleteValue(key, value string) {

	if self == nil { panic("Node.DeleteValue(): called on nil node") }

	// Disallow keys that change the board...

	for _, s := range MUTORS {
		if key == s {
			panic("Node.DeleteValue(): Can't change board-altering properties")
		}
	}

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

	// Disallow keys that change the board...

	for _, s := range MUTORS {
		if key == s {
			panic("Node.DeleteKey(): Can't change board-altering properties")
		}
	}

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
