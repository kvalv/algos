package btree

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func keystr(n *Node) string {
	var parts []string
	for _, k := range n.Keys {
		parts = append(parts, fmt.Sprintf("%d", k))
	}
	return strings.Join(parts, "")
}

func Graphviz(t *BTree, fname string) {
	s := &strings.Builder{}
	fmt.Fprintf(s, "digraph G {\n")
	t.WalkNodes(t.Root, func(n *Node) {

		for _, c := range n.Children {
			src := strings.Replace(strings.Replace(n.String(), "(", "", -1), ")", "", -1)
			dst := strings.Replace(strings.Replace(c.String(), "(", "", -1), ")", "", -1)
			fmt.Fprintf(s, "%s -> %s\n", src, dst)
		}

	})
	s.WriteString("}\n")

	dotFile := fmt.Sprintf("%s.dot", fname)
	f, err := os.OpenFile(dotFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	f.WriteString(s.String())

	if err := exec.Command("dot", "-Tpng", dotFile, "-o", fname).Run(); err != nil {
		panic(err)
	}

}
