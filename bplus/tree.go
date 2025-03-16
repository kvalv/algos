package bplus

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
)

type BTree struct {
	// t = n - 1
	n         int // n pointers, n-1 keys
	log       *slog.Logger
	dbg       bool
	Root      *Node
	pageCache *PageCache
}

func (T *BTree) isValid() error {
	var err error
	T.WalkNodes(T.Root, func(n *Node) {
		if len(n.Keys) > 2*T.n-1 {
			err = (fmt.Errorf("node %s has %d keys", n, len(n.Keys)))
		}
		if !n.Leaf && len(n.Keys)+1 != len(n.Children) {
			err = fmt.Errorf("BTree violation: node %s has %d keys and %d children - expected %d children",
				n,
				len(n.Keys),
				len(n.Children),
				len(n.Keys)+1)
		}
		if len(n.Children) > 0 && n.Leaf {
			var childRepr []string
			for _, id := range n.Children {
				child := T.pageCache.Read(id)
				childRepr = append(childRepr, child.String())
			}
			err = (fmt.Errorf("Node %q is a leaf node with %d children; children=%q", n, len(n.Children), strings.Join(childRepr, ", ")))
		}
		if len(n.Children) == 0 && !n.Leaf {
			err = (fmt.Errorf("Node %q is not a leaf, but has children", n))
		}
	})
	return err
}
func (T *BTree) validate() {
	if !T.dbg {
		return
	}
	if err := T.isValid(); err != nil {
		panic(err)
	}
}
func (T *BTree) read(n *Node, i int) *Node {
	if len(n.Children) < i-1 || i < 0 {
		return nil
	}
	id := n.Children[i]
	return T.pageCache.Read(id)
}
func (T *BTree) write(n *Node) *Node {
	return T.pageCache.Write(n)
}
func (T *BTree) allocate() *Node {
	return T.pageCache.Allocate()
}

func (T *BTree) String(n *Node) string {
	if n == nil {
		panic("BTree.String(): n is nil")
	}
	var s strings.Builder
	s.WriteString("(")
	for _, k := range n.Keys {
		s.WriteString(keyString(k))
	}
	// if !n.Leaf { // I think we're using .Children for pointers, too...
	for i := range n.Children {
		c := T.read(n, i)
		fmt.Fprintf(&s, "%s", T.String(c))
	}
	// }
	s.WriteString(")")

	return s.String()
}

func (T *BTree) WalkNodes(n *Node, f func(n *Node)) {
	if n == nil {
		return
	}
	f(n)
	for _, id := range n.Children {
		c := T.pageCache.Read(id)
		T.WalkNodes(c, f)
	}
}
func (T *BTree) Walk(n *Node, f func(key int)) {
	if n == nil {
		return
	}
	for i, key := range n.Keys {
		if !n.Leaf {
			c := T.read(n, i)
			T.Walk(c, f)
		}
		f(key)
	}
	if !n.Leaf {
		c := T.read(n, len(n.Keys))
		T.Walk(c, f)
	}
}
func (b *BTree) Keys() []int {
	var res []int
	b.Walk(b.Root, func(key int) {
		res = append(res, key)
	})
	return res
}

// returns nil if node does not have any keys, or the key is greater
// than all keys in this set, in which case - consider calling T.lastChild
func (T *BTree) insertionIndex(key int, node *Node) *int {
	if len(node.Keys) == 0 || key > node.Keys[len(node.Keys)-1] {
		return nil
	}
	// otherwise we knwo for sure there's at least one key that is greater
	for i, k := range node.Keys {
		if k >= key {
			return &i
		}
	}
	panic("unreachable")
}

func (T *BTree) lastChild(N *Node) *Node {
	length := len(N.Children)
	if length == 0 {
		return nil
	}
	pageID := N.Children[length-1]
	return T.pageCache.Read(pageID)
}

func (T *BTree) Find(key int) *Match {
	C := T.Root
	var count int
	for !C.Leaf {
		count++
		if count > 100 {
			panic("infinite loop")
		}
		i := T.insertionIndex(key, C)
		if i == nil {
			C = T.lastChild(C)
			continue
		}
		if C.Keys[*i] == key {
			C = T.read(C, *i+1)
		} else {
			C = T.read(C, *i)
		}
	}
	i := T.insertionIndex(key, C)
	if i == nil {
		return nil
	}
	return &Match{C, *i}
}

func (T *BTree) Range(key, upper int) RangeIterator {
	C := T.Root
	for !C.Leaf {
		i := T.insertionIndex(key, C)
		if i == nil {
			C = T.lastChild(C)
			continue
		}
		if C.Keys[*i] == key {
			C = T.read(C, *i+1)
		} else {
			C = T.read(C, *i)
		}
	}
	i := T.insertionIndex(key, C)
	if i == nil {
		return NewEmptyIterator[Match]()
	}

	j := *i
	return NewIterator(func() *Match {
		m := Match{Index: j, Node: C}
		if j >= len(C.Keys) {
			if C.RightSibling == nil {
				return nil
			}
			// otherwise visit next sibling
			C = T.pageCache.Read(*C.RightSibling)
			j = 0
			m = Match{Index: j, Node: C}
		}
		if C.Keys[j] >= upper {
			return nil
		}
		j++
		return &m
	})
}

func (T *BTree) Insert(key int, value PageID) {
	node := T.Root
	defer T.validate()

	var stack []*Node
	for !node.Leaf {
		stack = append(stack, node)
		if i := T.insertionIndex(key, node); i != nil {
			node = T.read(node, *i)
		} else {
			node = T.lastChild(node)
		}
	}
	stack = append(stack, node) // append the leaf, too

	// we'll loop over the nodes, bottom-up - starting with the leaf node
	slices.Reverse(stack)

	pageID := value
	minKey := key
	var par *Node
	fmt.Printf("-- stack --\n")
	for _, node := range stack {
		fmt.Printf("%s\n", node.String())
	}
	fmt.Printf("-- stack end --\n")
	for i, node := range stack {
		// insert a given key and pageID into the parent node.
		// We keep doing this while splitting is necessary
		fmt.Printf("--- iteration start -- node=%s\n", node.String())
		fmt.Printf("will insert %s in node %s\n", keyString(minKey), T.String(node))
		T.insertInNode(node, minKey, pageID)

		if !T.NeedsSplit(node) {
			break
		}
		fmt.Printf("node needs split\n")
		// otherwise we need to split. Split and add new key to parent
		j := (T.n + 1) / 2 // ceil[n/2]
		pre := T.String(node)
		right, mk := T.Split(node, j)
		if i+1 < len(stack) {
			par = stack[i+1]
		} else {
			par = nil
		}
		var parStr string = "<nil>"
		if par != nil {
			parStr = par.String()
		}
		fmt.Printf("split done\n  start=%s\n  left=%s\n  right=%s\n  par=  %s\n", pre, T.String(node), T.String(right), parStr)
		if par == nil {
			fmt.Printf("parent is nil, so we create a new root node\n")
			par = T.allocate()
			par.Keys = []int{mk} // not sure of this
			par.Children = []PageID{node.PageID, right.PageID}
			par.Leaf = false
			T.Root = par
			fmt.Printf("root is now %s\n", T.String(par))
			return // no need to continue down. we know we we're done
		} else {
			// T.insertInNode(par, right.MinKey(), right.PageID)
			fmt.Printf("parent exists, will add %s to %s in next iteration\n", keyString(mk), par.String())
		}

		// otherwise, we have split and we need to register the new
		// key to the parent.
		// fmt.Printf("parent is %s and will insert %s\n", par, keyString(right.MinKey()))
		minKey = mk // what if we remove extra keys?? then we're fucked
		pageID = right.PageID

	}
}

// inserts the key at the appropriate location. value is either a leaf value,
// or a pointer to a child page, which is put to the RIGHT
func (T *BTree) insertInNode(node *Node, key int, value PageID) {
	i := T.insertionIndex(key, node)
	if i == nil {
		node.Keys = append(node.Keys, key)
		if node.Leaf {
			node.Values = append(node.Values, value)
		} else {
			node.Children = append(node.Children, value) // seems about right
		}
	} else {
		node.Keys = slices.Insert(node.Keys, *i, key)
		if node.Leaf {
			node.Values = slices.Insert(node.Values, *i+1, value)
		} else {
			node.Children = slices.Insert(node.Children, *i+1, value)
		}
	}
}

func (T *BTree) NeedsSplit(n *Node) bool {
	return T.n == len(n.Keys)-1
}

// Splits current node at index i, returning the new node, along with the key that should
// be used as the separation key for parent nodes
func (T *BTree) Split(node *Node, i int) (*Node, int) {
	right := T.pageCache.Allocate()
	right.Leaf = node.Leaf

	right.Keys = node.Keys[i:]
	node.Keys = node.Keys[:i]

	right.RightSibling = node.RightSibling
	tmp := right.PageID
	node.RightSibling = &tmp

	if node.Leaf {
		right.Values = node.Values[i:]
		node.Values = node.Values[:i]
		return right, right.MinKey()
	} else {
		right.Children = node.Children[i+1:]
		node.Children = node.Children[:i+1]
		if len(right.Keys) == len(right.Children) {
			fmt.Printf("removing extraneous key %s -> now keys=%v\n", keyString(right.Keys[0]), right.Keys[1:])
			separationKey := right.Keys[0]
			right.Keys = right.Keys[1:]
			return right, separationKey
		}
		return right, right.Keys[0]
	}
}
