package bplus

type Node struct {
	Keys []int
	Leaf bool

	RightSibling *Node

	// leaf: has N-1 keys and N pointers
	// For leaf, the last pointer points to sibling node (next) - not back
	Children []*Node // ??
	Pointers []PageID
}

func (N *Node) lastChild() *Node {
	length := len(N.Children)
	if length == 0 {
		return nil
	}
	return N.Children[length-1]
}

func (n *Node) median() (index int, key int) {
	if len(n.Keys) == 0 {
		return 0, 0
	}

	index = len(n.Keys) / 2
	key = n.Keys[index]
	return
}
