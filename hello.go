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
	return greetingPrefix(language) + name
}

func greetingPrefix(lang string) (prefix string) {
	switch lang {
	case french:
		prefix = frenchHelloPrefix
	case spanish:
		prefix = spanishHelloPrefix
	default:
		prefix = helloPrefix
	}
	return
}

func main() {
	fmt.Println(Hello("World", ""))
}
