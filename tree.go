package ntree

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	logg "github.com/sirupsen/logrus"
)

// Node struct, represent a name tree
type Node struct {
	component string
	entry     interface{}
	children  []*Node
}

// New return a new tree with given root element at the top of the tree
func New(root string) *Node {
	logg.WithFields(logg.Fields{
		"root": root,
	}).Info("[NameTree] New() new name tree with")
	return &Node{
		component: root,
	}
}

// Size returns the number of name components in the tree
func (n *Node) Size() (size int) {
	if n.component != "" {
		size = 1
	}
	if len(n.children) > 0 {
		for _, child := range n.children {
			size += child.Size()
		}
	}
	return size
}

// FindExactMatch finds the exact match for a given prefix
// returns the entry at the node found and true if an exact match was found, false otherwise
func (n *Node) FindExactMatch(prefix string) (interface{}, bool) {
	_, entry, found, isExact := n.find(prefix)
	if found && isExact {
		return entry, true
	}
	return nil, false
}

// FindLongestMatch finds the longest match for the given prefix
// returns the longest prefix found, the entry at that prefix
// found is true if a longest match was found, false otherwise
func (n *Node) FindLongestMatch(prefix string) (longest string, entry interface{}, found bool) {
	longest, entry, found, _ = n.find(prefix)
	return
}

// find looks up given prefix in the tree, this can do both longest and exact match
// returns longest prefix found, longest is the prefix itself in case of exact match
// entry corresponding to the prefix, found is true if a longest or exact was found, false if none
// isExact is true if exact match is found
func (n *Node) find(prefix string) (longest string, entry interface{}, found bool, isExact bool) {
	components := strings.Split(prefix, "/")[1:]
	someNode, _, i := n.walkTree(components)
	if someNode != nil {
		found = true
		entry = someNode.entry
	}
	if i > 0 {
		// we reached last component, means found exact match
		isExact = true
		longest = prefix
		entry = someNode.entry
		return
	}
	var index int
	for k, c := range components {
		// components does not contain root "/" element
		// so we're matching against components ohter than the root component
		// (check first line in this block)
		if c == someNode.component {
			index = k
		}
	}
	longest = strings.Join(components[:index+1], "/")
	return
}

// Insert adds a prefix to the tree with the given entry.
// this will add multiple nodes to the tree if prefix has has multiple components.
func (n *Node) Insert(prefix string, entry interface{}) {
	components := strings.Split(prefix, "/")[1:]
	count := len(components)
	for count != 0 {
		// log.Println("current count =", count)
		someNode, component, i := n.walkTree(components)
		if i > 0 { // we reached the last component
			log.Printf("added prefix \"%s\"\n", prefix)
			return
		}
		new := &Node{component: component, entry: entry}
		someNode.children = append(someNode.children, new)
		// log.Printf("added new child \"%s\" to parent \"%s\" \n", component, someNode.component)
		count--
		if count == 0 {
			log.Printf("added prefix \"%s\"\n", prefix)
			return
		}
	}
}

// walkTree moves along the tree adding components.
// returns the current Node and component found, and length of current children.
func (n *Node) walkTree(components []string) (*Node, string, int) {
	component := components[0]
	if len(n.children) > 0 { // no children, then return.
		for _, child := range n.children {
			if component == child.component {
				next := components[1:]
				if len(next) > 0 {
					return child.walkTree(next) // tail recursion is it's own reward.
				}
				return child, component, len(n.children)
			}
		}
	}
	return n, component, 0
}

// String implements Stringer interface
func (n *Node) String() string {
	buf := new(bytes.Buffer)
	n.print(buf, 0)
	return string(buf.Bytes())
}

// JSON return tree in json format with given indent level
func (n *Node) JSON(indent int) string {
	obj := newJSONTree(n)
	bytes, err := json.MarshalIndent(obj, "", strings.Repeat(" ", indent))
	if err != nil {
		log.Fatalln(err)
	}
	return string(bytes)
}

// jsonTree struct used to print a tree as JSON object
// this makes lief easier, since json package comes with MarshalIndent
type jsonTree struct {
	Component string      `json:"component"`
	Children  []*jsonTree `json:"children,omitempty"`
}

// newJSONTree return a JSONTree copied from the given name tree
func newJSONTree(tree *Node) *jsonTree {
	jt := new(jsonTree)
	jt.Component = tree.component
	for _, child := range tree.children {
		jt.Children = append(jt.Children, newJSONTree(child))
	}
	return jt
}

// print prints the tree to the given writer with given indent level
func (n *Node) print(writer io.Writer, indent int) {
	fmt.Fprintf(writer, "%s%s \n", strings.Repeat(" ", indent), string(n.component))
	for _, child := range n.children {
		child.print(writer, indent+2)
	}
}
