package main

import (
	"fmt"

	RxGo "github.com/Hide-on-bush2/rxgo"
)

func main() {
	RxGo.Just("Hello", "World", "!", "!", "!", "!", "Hello").Distinct().Subscribe(func(x string) {
		fmt.Println(x)
	})

	// RxGo.Just("Hello", "World", "!", "!", "!", "!", "Hello").ElementAt("!").Subscribe(func(x string) {
	// 	fmt.Println(x)
	// })
}
