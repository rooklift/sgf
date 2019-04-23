package sgf

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// SaveCollection creates a new file, and saves each tree given into that file.
// It is useful for saving the rarely-used SGF collection format. Note that the
// location of the nodes in their trees is irrelevant: in each case, the whole
// tree is always saved.
func SaveCollection(nodes []*Node, filename string) error {

	var roots []*Node

	for _, node := range nodes {
		if node != nil {
			roots = append(roots, node.GetRoot())		// Note that node.Save() relies on GetRoot() here.
		}
	}

	if len(roots) == 0 {
		return fmt.Errorf("SaveCollection(): No non-nil roots supplied")
	}

	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(outfile)		// bufio for speedier output if file is huge.
	for _, root := range roots {
		root.write_tree(w)
	}
	w.Flush()							// "After all data has been written, the client should call the Flush method"

	// We didn't defer outfile.Close() like normal people so we can check its error here, just in case...

	err = outfile.Close()
	if err != nil {
		return err
	}

	return nil
}

// Save saves the entire game tree to the specified file. It does not need to be
// called from the root node, but can be called from any node in an SGF tree -
// the whole tree is always saved.
func (self *Node) Save(filename string) error {
	return SaveCollection([]*Node{self}, filename)
}

func (self *Node) write_tree(outfile io.Writer) {

	node := self

	fmt.Fprintf(outfile, "(")

	for {

		fmt.Fprintf(outfile, ";")

		for key, _ := range node.props {

			fmt.Fprintf(outfile, "%s", key)

			for _, value := range node.props[key] {
				fmt.Fprintf(outfile, "[%s]", escape_string(value))
			}
		}

		if len(node.children) > 1 {

			for _, child := range node.children {
				child.write_tree(outfile)
			}

			break

		} else if len(node.children) == 1 {

			node = node.children[0]
			continue

		} else {

			break

		}

	}

	fmt.Fprintf(outfile, ")\n")
	return
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

// Load reads an SGF file (or GIB file, if the extension .gib is present)
// creating a tree of SGF nodes, and returning the root. The input file is
// closed automatically. If the file has more than one SGF tree (a rarity) only
// the first is loaded - use LoadCollection() for such files instead.
func Load(filename string) (*Node, error) {

	file_bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	data := string(file_bytes)

	// If RAM wastage was ever an issue, one can do the super-spooky:
	// data := *(*string)(unsafe.Pointer(&file_bytes))
	// See https://github.com/golang/go/issues/25484 for details.

	root, _, err := load_sgf_tree(data, nil)

	if err != nil {
		if strings.HasSuffix(filename, ".gib") {
			root, err = load_gib(data)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return root, nil
}

func load_sgf_tree(sgf string, parent_of_local_root *Node) (*Node, int, error) {

	// A tree is whatever is between ( and ).
	//
	// FIXME: this is not unicode aware. Potential problems exist if
	// a unicode code point contains a meaningful character, especially
	// the bytes ] and \ although this is impossible for utf-8.

	var root *Node
	var node *Node
	var tree_started bool
	var inside_value bool
	var value string
	var key string
	var keycomplete bool

	for i := 0; i < len(sgf); i++ {

		c := sgf[i]

		if tree_started == false {
			if c <= ' ' {				// Reasonable definition of whitespace, where ' ' is byte 32.
				continue
			} else if c == '(' {
				tree_started = true
				continue
			} else {
				return nil, 0, fmt.Errorf("load_sgf_tree(): unexpected byte Ox%x before (", c)
			}
		}

		if inside_value {

			if c == '\\' {
				if len(sgf) <= i + 1 {
					return nil, 0, fmt.Errorf("load_sgf_tree(): escape character at end of input")
				}
				// value += string('\\')		// Do not do this. Discard the escape slash.
				value += string(sgf[i + 1])
				i++								// Skip 1 character.
			} else if c == ']' {
				inside_value = false
				if node == nil {
					return nil, 0, fmt.Errorf("load_sgf_tree(): value ended by ] but node was nil")
				}
				node.AddValue(key, value)
			} else {
				value += string(c)
			}

		} else {

			if c == '[' {
				if node == nil {
					return nil, 0, fmt.Errorf("load_sgf_tree(): value started by [ but node was nil")
				}
				value = ""
				inside_value = true
				keycomplete = true
			} else if c == '(' {
				if node == nil {
					return nil, 0, fmt.Errorf("load_sgf_tree(): new subtree started but node was nil")
				}
				_, chars_to_skip, err := load_sgf_tree(sgf[i:], node)	// Substrings are memory efficient in Golang.
				if err != nil {
					return nil, 0, err
				}
				i += chars_to_skip - 1		// Subtract 1: the ( character we have read is also counted by the recurse.
			} else if c == ')' {
				if root == nil {
					return nil, 0, fmt.Errorf("load_sgf_tree(): subtree ended but local root was nil")
				}
				return root, i + 1, nil		// Return characters read.
			} else if c == ';' {
				if node == nil {
					newnode := NewNode(parent_of_local_root)
					root = newnode
					node = newnode
				} else {
					newnode := NewNode(node)
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

	// Just being here must mean we reached the actual end of the file without
	// reading a final ')' character. Still, we can return what we have.

	if root == nil {
		return nil, 0, fmt.Errorf("load_sgf_tree(): local root was nil at function end")
	}

	return root, len(sgf), nil		// Return characters read.
}

// LoadCollection loads an SGF file and returns a slice of all root nodes found
// in it. It is useful for reading the rare SGF files that are in the
// "collection" format. The input file is closed automatically. Note that it is
// OK to use this function on normal, single-tree SGF files, in which case a
// slice of length 1 will be returned.
func LoadCollection(filename string) ([]*Node, error) {

	var ret []*Node

	file_bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return ret, err
	}

	data := string(file_bytes)

	for {
		if len(data) == 0 {
			return ret, nil
		}

		root, chars_read, err := load_sgf_tree(data, nil)

		if err != nil {
			return ret, err
		}

		ret = append(ret, root)

		data = data[chars_read:]
	}
}
