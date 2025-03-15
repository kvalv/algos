package bplus

import "log/slog"

type PageID = int

type PageCache struct {
	log   *slog.Logger
	nodes []*Node
}

func NewPageCache(log *slog.Logger) *PageCache {
	return &PageCache{
		log:   log,
		nodes: nil,
	}
}

func (pc *PageCache) Read(id PageID) *Node {
	for _, node := range pc.nodes {
		if node.PageID == id {
			return node
		}
	}
	return nil
}

func (pc *PageCache) Write(n *Node) *Node {
	_, med := n.median()
	pc.log.Debug("Disk write", "node", keyString(med))
	var found bool
	for _, node := range pc.nodes {
		if node == n {
			found = true
		}
	}
	if !found {
		pc.nodes = append(pc.nodes, n)
	}
	return n
}
func (pc *PageCache) Allocate() *Node {
	pc.log.Debug("Allocate-Node")
	n := &Node{
		PageID: len(pc.nodes),
	}
	pc.nodes = append(pc.nodes, n)
	return n
}
