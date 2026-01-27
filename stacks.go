package main

import "fmt"

func strateStack() {
	stack := []int{1, 2, 3, 4, 5, 6}
	stack = append(stack, 7) // push to stack
	fmt.Println(stack)
	
	stack = stack[:len(stack) - 1] // pop last item in stack
	fmt.Println(stack)
}

// When Do We Use Stack?
// valid parentheses
// undo / redo
// reverse data
// DFS