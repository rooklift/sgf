package sgf

// Note: internally, strings are kept in an escaped state e.g. \] and \\
// However, when using the API, your functions will send and receive
// unescaped strings.

var MUTORS = []string{"B", "W", "AB", "AW", "AE"}

// -----------------------------------------------------------------------------

type Node struct {
	Props			map[string][]string
	Children		[]*Node
	Parent			*Node
}

func NewNode(parent *Node, props map[string][]string) *Node {

	node := new(Node)
	node.Parent = parent
	node.Props = make(map[string][]string)

	for key, _ := range props {
		for _, s := range props[key] {
			node.add_value(key, s)
		}
	}

	if node.Parent != nil {
		node.Parent.Children = append(node.Parent.Children, node)
	}

	return node
}

// -----------------------------------------------------------------------------

func (self *Node) AddValue(key, value string) {

	// Disallow keys that change the board...

	for _, s := range MUTORS {
		if key == s {
			panic("AddValue(): Can't change board-altering properties")
		}
	}

	self.add_value(key, value)
}

func (self *Node) add_value(key, value string) {			// Handles escaping; no other function should

	value = escape_string(value)

	for i := 0; i < len(self.Props[key]); i++ {				// Ignore if the value already exists
		if self.Props[key][i] == value {
			return
		}
	}

	self.Props[key] = append(self.Props[key], value)
}

func (self *Node) SetValue(key, value string) {

	// Disallow keys that change the board...

	for _, s := range MUTORS {
		if key == s {
			panic("SetValue(): Can't change board-altering properties")
		}
	}

	self.Props[key] = nil
	self.add_value(key, value)
}

func (self *Node) GetValue(key string) (value string, ok bool) {

	// Get the value for the key, on the assumption that there's only 1 value.

	list := self.Props[key]

	if len(list) == 0 {
		return "", false
	}

	return unescape_string(list[0]), true
}

func (self *Node) AllValues(key string) []string {

	// Return all values for the key, possibly zero

	list := self.Props[key]

	if len(list) == 0 {
		return nil
	}

	var ret []string		// Make a new slice to avoid aliasing.

	for _, s := range list {
		ret = append(ret, unescape_string(s))
	}

	return ret
}

func (self *Node) DeleteValue(key, value string) {

	// Disallow keys that change the board...

	for _, s := range MUTORS {
		if key == s {
			panic("DeleteValue(): Can't change board-altering properties")
		}
	}

	for i := len(self.Props[key]) - 1; i >= 0; i-- {
		v := self.Props[key][i]
		if v == value {
			self.Props[key] = append(self.Props[key][:i], self.Props[key][i+1:]...)
		}
	}

	if len(self.Props[key]) == 0 {
		delete(self.Props, key)
	}
}

func (self *Node) DeleteKey(key string) {

	// Disallow keys that change the board...

	for _, s := range MUTORS {
		if key == s {
			panic("DeleteKey(): Can't change board-altering properties")
		}
	}

	delete(self.Props, key)
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
