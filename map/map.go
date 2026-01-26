package main

import (
	"fmt"
)

func main () {
	// mapBasic()
	// mapIterating()
	mapAreReference()
}

func mapBasic() {
	inventory := map[string]int{
		"bananas": 3,
		"apples":   22,
	}

	inventory["bananas"] = 8
	fmt.Println(inventory, "change bananas value")

	inventory["lemon"] = 10
	fmt.Println(inventory, "add lemon")

	if count, ok := inventory["apples"]; ok{
		fmt.Println(count, "check count apples")
	}


	delete(inventory, "bananas")
	delete(inventory, "apples")
	fmt.Println("after delete bananas and apples >")

	if _, ok := inventory["bananas"]; !ok { //check banana has removed
		fmt.Println("bananas has ben removed")
	}

	fmt.Println(inventory, "curr data")
}

// Map Iterating
func mapIterating () {
	scores:= map[string] int {
		"kiki": 9,
		"wahyu": 8,
		"jamal":7,
	}

	keys:= make([]string, 0, len(scores)) // create slice to store key
	values:= make([]int, 0, len(scores)) // create slice to store value
	for k, v := range scores {
		keys = append(keys, k)
		values = append(values, v)
	}
	fmt.Println(keys)
	fmt.Println(values)
}

// maps in Go are reference types.
// So you are not copying the map’s data, you are copying a reference (pointer-like value) to the same map.
func modifyValue (m map[string] int) {
	m["a"] = 999 //modif a
	m["new"] = 111 //new key
}
func mapAreReference () {
	original:= map [string] int {
		"a": 1,
		"b": 2,
	}

	reference := original
	reference["c"] = 3
	
	fmt.Println(original,"original")
	fmt.Println(reference, "reference") // output will same to original
	
	modifyValue(original) // how about function? is same to
	fmt.Println(original,"original")
	fmt.Println(reference, "reference") // output will same to original
}


// counting map
