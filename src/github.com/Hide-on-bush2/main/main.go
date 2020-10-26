package main

import (
	"fmt"

	"github.com/Hide-on-bush2/MyMarshal"
)

// func main() {
// 	hide_on_bush := read_ini.Listen_methods{read_ini.MyListen}
// 	read_ini.Watch("../read_ini/test.ini", hide_on_bush)

// }

type tmp struct {
	name string
	team string
	year int
}

func main() {
	json, _ := MyMarshal.JsonMarshal(tmp{"Faker", "SKT", 10})
	fmt.Println(string(json))
	// t := tmp{"Faker", "SKT", 10}
	// obj_type := reflect.TypeOf(t)
	// t1 := reflect.TypeOf(obj_type.Field(0).Name)
	// fmt.Println(t1)
}
