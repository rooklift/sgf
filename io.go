package sgf

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func (self *Node) Save(filename string) error {

	if self == nil {
		return fmt.Errorf("Save() called on nil node")
	}

	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outfile.Close()

	w := bufio.NewWriter(outfile)						// bufio for speedier output if file is huge.
	defer w.Flush()

	self.GetRoot().write_tree(w)

	return nil
}

func (self *Node) write_tree(outfile io.Writer) {		// Relies on values already being correctly backslash-escaped

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
				child.write_tree(outfile)
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

func Load(filename string, clear_cache bool) (*Node, error) {

	// Caller must say whether to clear the board cache or not.
	// Any app that's gonna open a ton of files should.

	if clear_cache {
		ClearBoardCache()
	}

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
	// if a unicode code point contains a meaningful character,
	// especially the bytes ] and \

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
				// value += string('\\')		// Do not do this. Discard the escape slash.
				value += string(sgf[i + 1])
				chars_to_skip = 1
			} else if c == ']' {
				inside = false
				if node == nil {
					return nil, 0, fmt.Errorf("load_sgf_tree: node == nil after: else if c == ']'")
				}
				node.add_value(key, value)		// This handles escaping.
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
