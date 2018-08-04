package main

import (
	dv "github.com/Arafatk/DataViz/trees/btree"
)

func main() {
	tree := dv.NewWithIntComparator(3)
	tree.Put(1, "a")
	tree.Put(1, "b")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Visualizer("heap.png")
}
