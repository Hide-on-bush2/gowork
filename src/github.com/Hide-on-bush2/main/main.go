package main

import (
	"fmt"

	"github.com/Hide-on-bush2/MyMarshal"
)

// func main() {
// 	hide_on_bush := read_ini.Listen_methods{read_ini.MyListen}
// 	read_ini.Watch("../read_ini/test.ini", hide_on_bush)

// }

// type tmp struct {
// 	name string
// 	team string
// 	year int
// }

func main() {
	type tmp struct {
		Name     []string          `json:"name"`
		Team     string            `json:"team"`
		Year     int               `json:"old"`
		Test_map map[string]string `json:"-"`
	}
	m := make(map[string]string)
	m["Faker"] = "Mid"
	m["Theshy"] = "top"
	m["UZI"] = "ADCarry"

	t := tmp{Name: []string{"Faker", "bengi", "bang"}, Team: "SKT", Year: 10, Test_map: m}
	json, _ := MyMarshal.JsonMarshal(t)
	fmt.Println(string(json))
	// t := tmp{"Faker", "SKT", 10}
	// obj_type := reflect.TypeOf(t)
	// t1 := reflect.TypeOf(obj_type.Field(0).Name)
	// fmt.Println(t1)
}
