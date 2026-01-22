package main

import "fmt"

// Function strating pointer usage
func modifyValue(ptr *int) {
    *ptr = 20
}

// Function strating pointer to struct
func modifyPerson(p *Person) {
    p.Name = "Bob"
    p.Age = 40
}

// stration function
func stratePointers() {
    fmt.Println("=== Pointers Examples ===")

    // Example 1: Basic pointer
    var a int = 10
    var ptr *int = &a
    fmt.Printf("Value of a: %d\n", a)
    fmt.Printf("Address of a: %v\n", ptr)
    fmt.Printf("Value via pointer: %d\n", *ptr)

    // Example 2: Modifying value through pointer
    modifyValue(ptr)
    fmt.Printf("After modifyValue: %d\n", a)

    // Example 3: Pointer to struct
    person := Person{Name: "Khairil", Age: 30}
    fmt.Printf("Before modifyPerson: %s, %d\n", person.Name, person.Age)
    modifyPerson(&person)
    fmt.Printf("After modifyPerson: %s, %d\n", person.Name, person.Age)

    // Example 4: New function
    ptr2 := new(int)
    *ptr2 = 100
    fmt.Printf("Value from new: %d\n", *ptr2)
}