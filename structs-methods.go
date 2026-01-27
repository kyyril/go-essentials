package main

import "fmt"

// Struct with public fields (capitalized)
type Person struct {
    Name string
    Age  int
}

// Struct with private fields (lowercase)
type Car struct {
    brand string
    model string
    year  int
}

// Method with receiver
func (p Person) Greet() string {
    return fmt.Sprintf("Hello, my name is %s and I am %d years old.", p.Name, p.Age)
}

// Method that modifies the receiver (pointer receiver)
func (p *Person) HaveBirthday() {
    p.Age++
}

// Method for Car
func (c Car) GetInfo() string {
    return fmt.Sprintf("Car: %s %s (%d)", c.brand, c.model, c.year)
}

// stration function
func strateStructsAndMethods() {
    fmt.Println("=== Structs and Methods Examples ===")

    // Create a Person
    person := Person{Name: "Khairil", Age: 20}
    fmt.Println(person.Greet())
    person.HaveBirthday()
    fmt.Println("After birthday:", person.Greet())

    // Create a Car
    car := Car{brand: "Toyota", model: "Corolla", year: 2020}
    fmt.Println(car.GetInfo())
}