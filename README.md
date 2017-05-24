# name-tree
tree to store URL like prefixes in key/value fashion (for learning purposes)

### example of usage

```go
package main

import (
	"log"

	"github.com/elmiomar/named-data/name-tree"
)

// example of fib entry
// (this is for demonstration purposes, a fib entry is more than this)
type fibEntry struct {
	nexthops []int
}

// example of pit entry
// (this is for demonstration purposes, a pit entry is more than this)
type pitEntry struct {
	in, out []int
}

// example of cs entry
// (this is for demonstration purposes, a cs entry is more than this)
type csEntry struct {
	data []byte
}

func main() {
	tree := ntree.New("/")

	// adding a prefix with a fib entry as value
	tree.Insert("/gov/nist/antd", fibEntry{
		nexthops: []int{258, 266},
	})

	// adding a prefix with a pit entry as value
	tree.Insert("/hello/omar/ilias", pitEntry{
		in:  []int{258},
		out: []int{260, 266},
	})

	p := "/gov/nist/antd/hello"
	// performing lookup for exact match of prefix p, exact match should not be found at this point
	_, ok := tree.FindExactMatch(p)
	if ok {
		log.Printf("found exact match for prefix \"%s\"\n", p)
	} else {
		log.Printf("no exact match for prefix \"%s\"\n", p)
	}

	// however a longest prefix should exist in the tree
	longest, value, ok := tree.FindLongestMatch(p)
	if ok {
		log.Printf("found longest match \"%s\" with value \"%v\"\n", longest, value)
	} else {
		log.Printf("prefix \"%s\" not found\n", p)
	}

	// adding a prefix with a cs entry as value for the prefix
	tree.Insert(p, csEntry{
		data: []byte("some cached data"),
	})

	// now exact match should be found for prefix p
	value, ok = tree.FindExactMatch(p)
	if ok {
		log.Printf("found exact match for prefix \"%s\" with value \"%v\"\n", p, value)
	} else {
		log.Printf("no exact match for prefix \"%s\"\n", p)
	}

	// number of components int the tree
	log.Printf("number of components in tree is  \"%d\"\n", tree.Size())

	tree.Insert("/mheni/merzouki/pdf.123", csEntry{
		data: []byte("mheni's data"),
	})

	// print tree
	log.Printf("name tree looks like \n%s", tree)
}

```
