package main

import (
	"fmt"

	"github.com/Hide-on-bush2/stringutil"
)

func Hello(name string) string {
	if name == "" {
		return "Hello, world"
	}
	return "Hello, " + name
}

func main() {
	fmt.Printf(stringutil.Reverse("!oG ,olleH"))
	fmt.Printf(Hello(""))
}
