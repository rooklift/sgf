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
// creating a tree of SGF nodes, and returning the root. If the file has more
// than one SGF tree (a rarity) only the first is loaded. The input file is
// closed automatically.
func Load(filename string) (*Node, error) {

	file_bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	data := string(file_bytes)

	// If RAM wastage was ever an issue, one can do the super-spooky:
	// data := *(*string)(unsafe.Pointer(&file_bytes))
	// See https://github.com/golang/go/issues/25484 for details.

	root, err := load_sgf(data)

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

func load_sgf(sgf string) (*Node, error) {

	sgf = strings.TrimSpace(sgf)

	if sgf[0] == '(' {				// the load_sgf_tree() function assumes the
		sgf = sgf[1:]				// leading "(" has already been discarded.
	} else {
		return nil, fmt.Errorf("load_sgf(): Unexpected character before SGF tree")
	}

	root, _, err := load_sgf_tree(sgf, nil)
	return root, err
}

func load_sgf_tree(sgf string, parent_of_local_root *Node) (*Node, int, error) {

	// FIXME: this is not unicode aware. Potential problems exist if
	// a unicode code point contains a meaningful character, especially
	// the bytes ] and \ although this is impossible for utf-8.

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
					return nil, 0, fmt.Errorf("load_sgf_tree(): escape character at end of input")
				}
				// value += string('\\')		// Do not do this. Discard the escape slash.
				value += string(sgf[i + 1])
				chars_to_skip = 1
			} else if c == ']' {
				inside = false
				if node == nil {
					return nil, 0, fmt.Errorf("load_sgf_tree(): node == nil after: else if c == ']'")
				}
				node.AddValue(key, value)
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
					return nil, 0, fmt.Errorf("load_sgf_tree(): node == nil after: else if c == '('")
				}
				_, chars_to_skip, err = load_sgf_tree(sgf[i + 1:], node)		// substrings are memory efficient in Golang
				if err != nil {
					return nil, 0, err
				}
			} else if c == ')' {
				if root == nil {
					return nil, 0, fmt.Errorf("load_sgf_tree(): root == nil after: else if c == ')'")
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
		return nil, 0, fmt.Errorf("load_sgf_tree(): root == nil at function end")
	}

	return root, len(sgf), nil		// Return characters read.
}

// LoadCollection loads an SGF file and returns a slice of all root nodes found
// in it. It is useful for reading the rare SGF files that are in such a format.
// The input file is closed automatically.
func LoadCollection(filename string) ([]*Node, error) {

	var ret []*Node

	file_bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return ret, err
	}

	data := string(file_bytes)

	for {
		data = strings.TrimSpace(data)

		if len(data) == 0 {
			return ret, nil
		}

		if data[0] == '(' {				// the load_sgf_tree() function assumes the
			data = data[1:]				// leading "(" has already been discarded.
		} else {
			return ret, fmt.Errorf("LoadCollection(): Unexpected character outside of SGF trees")
		}

		root, chars_read, err := load_sgf_tree(data, nil)

		if err != nil {
			return ret, nil
		}

		ret = append(ret, root)

		data = data[chars_read:]
	}
}
