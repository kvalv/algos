package rb

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Graphviz(t Tree, fname string) {
	s := &strings.Builder{}
	fmt.Fprintf(s, "digraph G {\n")
	t.Walk(func(n *Node) {
		if n.Parent != nil {
			src := n.Parent.String()
			dst := n.String()
			fmt.Fprintf(s, "%s -> %s\n", src, dst)
		}
		fmt.Printf("node %s", n.String())
	})
	s.WriteString("}\n")
	fmt.Printf("got %s\n", s.String())

	dotFile := fmt.Sprintf("%s.dot", fname)
	f, err := os.OpenFile(dotFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	// defer f.Close()
	f.WriteString(s.String())
	f.Close()

	if err := exec.Command("dot", "-Tpng", dotFile, "-o", fname).Run(); err != nil {
		panic(err)
	}

}
