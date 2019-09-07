package main

import (
	"fmt"
)

func Hello(name string) string {
	return "Hey, " + name
}

func main() {
	fmt.Println(Hello("World"))
}
