package router

import (
	"ac-gateway/pkg/models/ddl"
	"fmt"
	"strings"
)

type methodTree struct {
	method string
	root   *node
}

func (tree methodTree) getValue(path string) (handlers interface{}, tsr bool) {
	return tree.root.getValue(path)
}
func (n *node) addRoute(path string, handlers interface{}) error {
	fullPath := path
	n.priority++
	numParams, err := countWildcards(path)
	if err != nil {
		return err
	}

	// non-empty tree
	if len(n.path) > 0 || len(n.children) > 0 {
	walk:
		for {
			// Update maxParams of the current node

			// Find the longest common prefix.
			// This also implies that the common prefix contains no '*'
			// since the existing key can't contain those chars.
			i := 0
			max := min(len(path), len(n.path))
			for i < max && path[i] == n.path[i] {
				i++
			}

			// Split edge
			if i < len(n.path) {
				child := node{
					path:      n.path[i:],
					wildChild: n.wildChild,
					nType:     static,
					indices:   n.indices,
					children:  n.children,
					handlers:  n.handlers,
					priority:  n.priority - 1,
				}

				n.children = []*node{&child}
				// []byte for proper unicode char conversion, see #65
				n.indices = string([]byte{n.path[i]})
				n.path = path[:i]
				n.handlers = nil
				n.wildChild = false
			}

			// Make new node a child of this node
			if i < len(path) {
				path = path[i:]

				if n.wildChild {
					n = n.children[0]
					n.priority++

					numParams--

					// Check if the wildcard matches
					if len(path) >= len(n.path) && n.path == path[:len(n.path)] &&
						// Check for longer wildcard, e.g. :name and :names
						(len(n.path) >= len(path) || path[len(n.path)] == '/') {
						continue walk
					} else {
						// Wildcard conflict
						pathSeg := strings.SplitN(path, "/", 2)[0]
						prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path
						return fmt.Errorf("'" + pathSeg +
							"' in new path '" + fullPath +
							"' conflicts with existing wildcard '" + n.path +
							"' in existing prefix '" + prefix +
							"'")
					}
				}

				c := path[0]

				// Check if a child with the next path byte exists
				for i := 0; i < len(n.indices); i++ {
					if c == n.indices[i] {
						i = n.incrementChildPrio(i)
						n = n.children[i]
						continue walk
					}
				}

				// Otherwise insert it
				if c != '*' {
					// []byte for proper unicode char conversion, see #65
					n.indices += string([]byte{c})
					child := new(node)
					n.children = append(n.children, child)
					n.incrementChildPrio(len(n.indices) - 1)
					n = child
				}
				if err := n.insertChild(numParams, path, fullPath, handlers); err != nil {
					return err
				}
				return nil

			} else if i == len(path) { // Make node a (in-path) leaf
				if n.handlers != nil {
					return fmt.Errorf("a handle is already registered for path '" + fullPath + "'")
				}
				n.handlers = handlers
			}
			return nil
		}
	} else { // Empty tree
		if err := n.insertChild(numParams, path, fullPath, handlers); err != nil {
			return err
		}
		n.nType = root
	}
	return nil
}

func (n *node) insertChild(numParams uint8, path, fullPath string, handlers interface{}) error {
	// find prefix until first wildcard (beginning with '*')
	for i, max := 0, len(path); numParams > 0; i++ {
		c := path[i]
		if c != '*' {
			continue
		}

		// find wildcard end (either '/' or path end)
		end := i + 1
		for end < max && path[end] != '/' {
			switch path[end] {
			// the wildcard name must not contain '*'
			case '*':
				end++
			default:
				return fmt.Errorf("only one wildcard per path segment is allowed, has: '" +
					path[i:] + "' in path '" + fullPath + "'")
			}
		}

		// check if this Node existing children which would be
		// unreachable if we insert the wildcard here
		if len(n.children) > 0 {
			return fmt.Errorf("wildcard route '" + path[i:end] +
				"' conflicts with existing children in path '" + fullPath + "'")
		}

		// check if the wildcard has a name
		if end-i < 2 {
			return fmt.Errorf("wildcards must be named with a non-empty name in path '" + fullPath + "'")
		}

		//if end != max || numParams > 1 {
		//	return fmt.Errorf("catch-all routes are only allowed at the end of the path in path '" + fullPath + "'")
		//}

		if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
			return fmt.Errorf("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")
		}

		// currently fixed width 1 for '/'
		i--
		if path[i] != '/' {
			return fmt.Errorf("no / before catch-all in path '" + fullPath + "'")
		}

		n.path = path[0:i]

		// first node: catchAll node with empty path
		child := &node{
			wildChild: true,
			nType:     catchAll,
		}
		n.children = []*node{child}
		n.indices = string(path[i])
		n = child
		n.priority++

		// second node: node holding the variable
		child = &node{
			path:     path[i:],
			nType:    catchAll,
			handlers: handlers,
			priority: 1,
		}
		n.children = []*node{child}

		return nil
	}

	// insert remaining path part and handle to the leaf
	n.path = path[0:]
	n.handlers = handlers
	return nil
}

// increments priority of the given child and reorders if necessary
func (n *node) incrementChildPrio(pos int) int {
	n.children[pos].priority++
	prio := n.children[pos].priority

	// adjust position (move to front)
	newPos := pos
	for newPos > 0 && n.children[newPos-1].priority < prio {
		// swap node positions
		tmpN := n.children[newPos-1]
		n.children[newPos-1] = n.children[newPos]
		n.children[newPos] = tmpN

		newPos--
	}

	// build new index char string
	if newPos != pos {
		n.indices = n.indices[:newPos] + // unchanged prefix, might be empty
			n.indices[pos:pos+1] + // the index char we move
			n.indices[newPos:pos] + n.indices[pos+1:] // rest without char at 'pos'
	}

	return newPos
}

func countWildcards(path string) (uint8, error) {
	var n uint
	wildcardIndexes := make([]int, 0, len(path))
	for i := 0; i < len(path); i++ {
		if path[i] != '*' {
			continue
		}
		wildcardIndexes = append(wildcardIndexes, i)
		n++
	}
	if len(wildcardIndexes)%2 != 0 {
		return 0, fmt.Errorf("Error path wildcard format:" + path)
	}
	for i := 0; i < len(wildcardIndexes); i++ {
		if i%2 != 0 {
			if wildcardIndexes[i]-wildcardIndexes[i-1] != 1 {
				return 0, fmt.Errorf("Error wildcard must ** ,path:" + path)
			}
		}

		if i != 0 && i%2 == 0 {
			if wildcardIndexes[i]-wildcardIndexes[i-1] > 0 {
				return 0, fmt.Errorf("Error wildcard must /**/ ,path:" + path)
			}
		}
	}
	return uint8(len(wildcardIndexes) / 2), nil
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func BuildTrieTree(list []ddl.Router) (*TrieTree, error) {
	trieTree := new(TrieTree)
	for _, d := range list {
		for _, m := range d.Methods {
			if err := trieTree.AddRoute(m, d.Uri, nil); err != nil {
				return nil, err
			}
		}
	}
	return trieTree, nil
}

type TrieTree struct {
	paramsIndex int
	methodTrees []methodTree
}

func (tree *TrieTree) AddRoute(method string, path string, handlers interface{}) error {
	r := tree.initRoot(method)
	return r.addRoute(path, handlers)
}

func (tree *TrieTree) GetValue(method, path string) (interface{}, bool) {
	for _, tree := range tree.methodTrees {
		if tree.method == method {
			tree.root.getValue(path)
		}
	}
	return nil, false
}

func (tree *TrieTree) initRoot(method string) *node {
	if tree.methodTrees == nil {
		tree.methodTrees = make([]methodTree, 0)
	}
	for _, item := range tree.methodTrees {
		if item.method == method {
			return item.root
		}
	}
	r := methodTree{
		method: method,
		root:   new(node),
	}
	tree.methodTrees = append(tree.methodTrees, r)
	return r.root
}

type node struct {
	path      string
	wildChild bool
	nType     nodeType
	indices   string
	children  []*node
	handlers  interface{}
	priority  uint32
}

func (n *node) getValue(path string) (handlers interface{}, tsr bool) {
walk: // outer loop for walking the tree
	for {
		if len(path) > len(n.path) {
			if path[:len(n.path)] == n.path {
				path = path[len(n.path):]
				// If this node does not have a wildcard (param or catchAll)
				// child,  we can just look up the next child node and continue
				// to walk down the tree
				if !n.wildChild {
					c := path[0]
					for i := 0; i < len(n.indices); i++ {
						if c == n.indices[i] {
							n = n.children[i]
							continue walk
						}
					}

					// Nothing found.
					// We can recommend to redirect to the same URL without a
					// trailing slash if a leaf exists for that path.
					tsr = (path == "/" && n.handlers != nil)
					return

				}

				// handle wildcard child
				n = n.children[0]
				switch n.nType {
				case catchAll:

					handlers = n.handlers
					return

				default:
					panic("invalid node type")
				}
			}
		} else if path == n.path {
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if handlers = n.handlers; handlers != nil {
				return
			}

			if path == "/" && n.wildChild && n.nType != root {
				tsr = true
				return
			}

			// No handle found. Check if a handle for this path + a
			// trailing slash exists for trailing slash recommendation
			for i := 0; i < len(n.indices); i++ {
				if n.indices[i] == '/' {
					n = n.children[i]
					tsr = (len(n.path) == 1 && n.handlers != nil) ||
						(n.nType == catchAll && n.children[0].handlers != nil)
					return
				}
			}

			return
		}

		// Nothing found. We can recommend to redirect to the same URL with an
		// extra trailing slash if a leaf exists for that path
		tsr = (path == "/") ||
			(len(n.path) == len(path)+1 && n.path[len(path)] == '/' &&
				path == n.path[:len(n.path)-1] && n.handlers != nil)
		return
	}
}

type nodeValue struct {
	handlers interface{}
	tsr      bool
	fullPath string
}
type nodeType uint8

const (
	static nodeType = iota // default
	root
	catchAll
)
