package main

import (
	"fmt"
	"iter"
	"math/rand"
	"strings"
)

type Tree struct {
	Value int
	Left  *Tree
	Right *Tree
}

func Insert(root *Tree, value int) *Tree {
	if root == nil {
		return &Tree{
			Value: value,
			Left:  nil,
			Right: nil,
		}
	}

	if root.Value != value {
		if root.Value > value {
			root.Left = Insert(root.Left, value)
		} else {
			root.Right = Insert(root.Right, value)
		}
	}

	return root
}

func Print(root *Tree, level int) {
	if root == nil {
		return
	}

	Print(root.Right, level+1)

	fmt.Printf("%s%d\n", strings.Repeat("  ", level), root.Value)

	Print(root.Left, level+1)
}

// iterHelper helps to define a recursive iterator because it cannot be an unnamed
// func
func iterHelper(root *Tree, yield func(int) bool) {
	if root == nil {
		return
	}

	iterHelper(root.Left, yield)
	yield(root.Value)
	iterHelper(root.Right, yield)
}

func Iter(root *Tree) iter.Seq[int] {
	return func(yield func(int) bool) {
		iterHelper(root, yield)
	}
}

func main() {
	var root *Tree

	for range 8 {
		root = Insert(root, rand.Intn(100))
	}

	Print(root, 0)

	println("Iterate over the tree")

	for v := range Iter(root) {
		fmt.Println(v)
	}
}
