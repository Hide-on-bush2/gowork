package main

import (
	"fmt"
	"os"
	"time"
	"github.com/Hide-on-bush2/rxgo"
)

func main() {
	rxgo.Just("Hello", "World", "!", "!", "!", "!", "Hello").First(func(x interface{}) bool { return x == "World" }).Subscribe(func(x string) {
		fmt.Println(x)
	})
	fmt.Println("----------")
	rxgo.Just("Hello", "World", "!", "!", "!", "!", "Hello").Skip(3).Subscribe(func(x string) {
		fmt.Println(x)
	})
	fmt.Println("----------")
	rxgo.Just("Hello", "World", "!", "!", "!", "!", "Hello").Take(3).Subscribe(func(x string) {
		fmt.Println(x)
	})
	fmt.Println("----------")
	rxgo.Just("Hello1", "World", "!1", "!2", "!3", "!4", "Hello2").Debounce(time.Second).Subscribe(func(x string) {
		fmt.Println(x)
	})
	fmt.Println("----------")
	rxgo.Just("Hello1", "World", "!1", "!2", "!3", "!4", "Hello2").Last().Subscribe(func(x string) {
		fmt.Println(x)
	})
	fmt.Println("----------")
	rxgo.Just("Hello1", "World", "!1", "!2", "!3", "!4", "Hello2").Skiplast(3).Subscribe(func(x string) {
		fmt.Println(x)
	})
	fmt.Println("----------")
	rxgo.Just("Hello1", "World", "!1", "!2", "!3", "!4", "Hello2").Takelast(3).Subscribe(func(x string) {
		fmt.Println(x)
	})
	fmt.Println("----------")
	rxgo.Just(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Map(func(x int) int {
		if x != 0 {
			time.Sleep(1 * time.Millisecond)
		}
		return x
	}).Debounce(2 * time.Millisecond).Subscribe(func(x int) {
		if x != 9 {
			fmt.Printf("error Debounce with %d\n", x)
			os.Exit(-1)
		}
		fmt.Printf("Debunce %d\n", x)
	})
	fmt.Println("----------")
	rxgo.Just("Hello", "World", "!", "!", "!", "!", "Hello").ElementAt("!").Subscribe(func(x string) {
		fmt.Println(x)
	})
}
