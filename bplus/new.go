package bplus

import (
	"io"
)

func New(n int, w io.Writer) *BTree {
	log := NewLogger(w)
	b := &BTree{
		n:         n,
		log:       log,
		dbg:       true,
		pageCache: NewPageCache(log),
	}
	x := b.allocate()
	x.Leaf = true
	b.write(x)
	b.Root = x
	return b
}

func FromString(n int, input string, w io.Writer) *BTree {
	T := New(n, w)

	var stack []*Node
	top := func() *Node {
		if len(stack) == 0 {
			return nil
		}
		return stack[len(stack)-1]
	}
	var root *Node
	pop := func() {
		if len(stack) == 0 {
			panic("FromString: invalid input: too many parantheses")
		}
		if len(stack) == 1 {
			root = stack[0]
		}
		stack = stack[:len(stack)-1]
	}

	toKey := func(c rune) int {
		// if 0123...f -> parse as hex {
		if c >= '0' && c <= '9' {
			return int(c - '0')
		}
		return int(c)
	}

	for _, c := range input {
		switch c {
		case '(':
			tmp := T.pageCache.Allocate()
			tmp.Leaf = true
			p := top()
			if p != nil {
				p.Leaf = false
				p.Children = append(p.Children, tmp.PageID)
			}
			stack = append(stack, tmp)
		case ')':
			pop()
		default:
			tmp := top()
			digit := toKey(c)
			tmp.Keys = append(tmp.Keys, digit)
		}
	}
	if len(stack) > 0 {
		panic("FromString: invalid input: unclosed parantheses")
	}

	T.Root = root
	T.validate()

	// Add next child
	var prev *Node
	T.WalkNodes(T.Root, func(n *Node) {
		if n.Leaf {
			if prev != nil {
				prev.RightSibling = n
			}
			prev = n
		}
	})

	return T
}
