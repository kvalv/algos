package bplus

import "fmt"

type Match struct {
	Node  *Node
	Index int
}

func (m *Match) String() string {
	if m == nil {
		return "nil"
	}
	return fmt.Sprintf("%s/%d", m.Node, m.Index)
}
