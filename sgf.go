package sgf

// Note: internally, strings are kept in an escaped state e.g. \] and \\
// However, when using the API, your functions will send and receive
// unescaped strings.

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const DEFAULT_SIZE = 19
var MUTORS = []string{"B", "W", "AB", "AW", "AE"}

type Colour int

const (
	EMPTY = Colour(iota)
	BLACK
	WHITE
)

// -------------------------------------------------------------------------

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

func NewTree(size int) *Node {

	// Sizes over 25 are not recommended. But 52 is the hard limit for SGF.

	if size < 1 || size > 52 {
		panic(fmt.Sprintf("NewTree(): invalid size %v", size))
	}

	properties := make(map[string][]string)
	size_string := fmt.Sprintf("%d", size)
	properties["SZ"] = []string{size_string}
	properties["GM"] = []string{"1"}
	properties["FF"] = []string{"4"}

	return NewNode(nil, properties)
}

func new_bare_node(parent *Node) *Node {

	// Doesn't accept properties.
	// Used only for file loading.

	node := new(Node)
	node.Parent = parent
	node.Props = make(map[string][]string)

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node
}

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

func (self *Node) RemoveChild(child *Node) {

	if self == nil {
		return
	}

	for i := len(self.Children) - 1; i >= 0; i-- {
		if self.Children[i] == child {
			self.Children = append(self.Children[:i], self.Children[i+1:]...)
		}
	}
}

func (self *Node) GetRoot() *Node {
	node := self
	for {
		if node.Parent != nil {
			node = node.Parent
		} else {
			return node
		}
	}
}

func (self *Node) GetEnd() *Node {
	node := self
	for {
		if len(node.Children) > 0 {
			node = node.Children[0]
		} else {
			return node
		}
	}
}

func (self *Node) Save(filename string) error {

	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outfile.Close()

	w := bufio.NewWriter(outfile)						// bufio for speedier output if file is huge.
	defer w.Flush()

	self.GetRoot().WriteTree(w)

	return nil
}

func (self *Node) WriteTree(outfile io.Writer) {		// Relies on values already being correctly backslash-escaped

	node := self

	fmt.Fprintf(outfile, "(")

	for {

		fmt.Fprintf(outfile, ";")

		for key, _ := range node.Props {

			fmt.Fprintf(outfile, "%s", key)

			for _, value := range node.Props[key] {
				fmt.Fprintf(outfile, "[%s]", value)
			}
		}

		if len(node.Children) > 1 {

			for _, child := range node.Children {
				child.WriteTree(outfile)
			}

			break

		} else if len(node.Children) == 1 {

			node = node.Children[0]
			continue

		} else {

			break

		}

	}

	fmt.Fprintf(outfile, ")\n")
	return
}

// -------------------------------------------------------------------------

func Load(filename string) (*Node, error) {

	sgf_bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	root, err := load_sgf(string(sgf_bytes))

	if err != nil {
		return nil, err
	}

	return root, nil
}

func load_sgf_tree(sgf string, parent_of_local_root *Node) (*Node, int, error) {

	// FIXME: this is not unicode aware. Potential problems exist
	// if a unicode code point contains a meaningful character.

	var root *Node
	var node *Node

	var inside bool
	var value string
	var key string
	var keycomplete bool
	var chars_to_skip int

	var err error

	for i := 0; i < len(sgf); i++ {

		c := sgf[i]

		if chars_to_skip > 0 {
			chars_to_skip--
			continue
		}

		if inside {

			if c == '\\' {
				if len(sgf) <= i + 1 {
					return nil, 0, fmt.Errorf("load_sgf_tree: escape character at end of input")
				}
				value += string('\\')
				value += string(sgf[i + 1])
				chars_to_skip = 1
			} else if c == ']' {
				inside = false
				if node == nil {
					return nil, 0, fmt.Errorf("load_sgf_tree: node == nil after: else if c == ']'")
				}
				node.add_value(key, value)
			} else {
				value += string(c)
			}

		} else {

			if c == '[' {
				value = ""
				inside = true
				keycomplete = true
			} else if c == '(' {
				if node == nil {
					return nil, 0, fmt.Errorf("load_sgf_tree: node == nil after: else if c == '('")
				}
				_, chars_to_skip, err = load_sgf_tree(sgf[i + 1:], node)
				if err != nil {
					return nil, 0, err
				}
			} else if c == ')' {
				if root == nil {
					return nil, 0, fmt.Errorf("load_sgf_tree: root == nil after: else if c == ')'")
				}
				return root, i + 1, nil		// Return characters read.
			} else if c == ';' {
				if node == nil {
					newnode := new_bare_node(parent_of_local_root)
					root = newnode
					node = newnode
				} else {
					newnode := new_bare_node(node)
					node = newnode
				}
			} else {
				if c >= 'A' && c <= 'Z' {
					if keycomplete {
						key = ""
						keycomplete = false
					}
					key += string(c)
				}
			}
		}
	}

	if root == nil {
		return nil, 0, fmt.Errorf("load_sgf_tree: root == nil at function end")
	}

	return root, len(sgf), nil		// Return characters read.
}

func load_sgf(sgf string) (*Node, error) {

	sgf = strings.TrimSpace(sgf)
	if sgf[0] == '(' {				// the load_sgf_tree() function assumes the
		sgf = sgf[1:]				// leading "(" has already been discarded.
	}

	root, _, err := load_sgf_tree(sgf, nil)
	return root, err
}

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
