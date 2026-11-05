package main

func Hello(name string) string {
	if name == "" {
		name = "World"
	}
	return "Hello, " + name
}
