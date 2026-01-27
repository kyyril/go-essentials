package main

import (
	"fmt"
)

func main () {
	// mapBasic()
	// mapIterating()
	// mapAreReference()
	// mapCounting()
	mapGrouping()
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


// counting with map
func mapCounting() {
	// Count word occurrences
	text := []string{
		"apple", "banana", "apple", "cherry",
		"banana", "apple", "date", "cherry",
	}

	wordCount := make(map[string]int)

	for _, word := range text {
		wordCount[word]++ // zero value (0) + 1 works
	}

	fmt.Println(wordCount)

	fmt.Println("Word counts:")
	for word, count := range wordCount {
		fmt.Printf("%s : %d\n", word, count)
	}
}


func mapGrouping (){
	type Person struct {
		Name string
		Age int
		City string
	}
	people := []Person {
		{Name:"kiki", Age:21, City: "Sibolga"},
		{Name:"sidiq", Age:21, City: "Sibolga"},
		{Name:"jamal", Age:22, City: "Medan"},
		{Name:"wahyu", Age:23, City: "Medan"},
		{Name:"apalah", Age:23, City: "Medan"},
	}
	fmt.Println("origin:",people)

	// group by age
	byAge:= make(map[int] []Person)

	for _, p := range people {
		byAge[p.Age] = append(byAge[p.Age], p)
	}
	fmt.Println(byAge)
	// filter
	for age, group := range byAge {
		fmt.Println("age:",age)
		for _, p := range group { 
		fmt.Println("people:",p)
		}
	}

	// group by city
	byCity := make(map[string] []Person)
	
	for _, p := range people {
		byCity[p.City] = append(byCity[p.City], p)
	}
	for city, name := range byCity {
		fmt.Printf("%s :%v", city, name)
	}
}

