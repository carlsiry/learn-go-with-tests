package main

import (
	"fmt"
)

const spanish = "Spanish"
const french = "French"

const helloPrefix = "Hey, "
const spanishHelloPrefix = "Hola, "
const frenchHelloPrefix = "Bonjour, "

const emptyName = ""
const defaultName = "world"

func Hello(name, language string) string {
	if name == emptyName {
		name = defaultName
	}

	var prefix string
	switch language {
	case french:
		prefix = frenchHelloPrefix
	case spanish:
		prefix = spanishHelloPrefix
	default:
		prefix = helloPrefix
	}

	return prefix + name
}

func main() {
	fmt.Println(Hello("World", ""))
}
