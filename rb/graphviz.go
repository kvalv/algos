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

			var color string
			if n.Color == RED {
				color = "red"
			} else {
				color = "black"
			}
			fmt.Fprintf(s, "%d [color=%q]\n", n.Key, color)
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

	if err := exec.Command("dot", "-Tpng", dotFile, "-o", fname).Run(); err != nil {
		panic(err)
	}

}
