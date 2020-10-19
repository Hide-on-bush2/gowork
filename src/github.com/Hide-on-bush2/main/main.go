package main

import "github.com/Hide-on-bush2/read_ini"

func main() {
	hide_on_bush := read_ini.Listen_methods{read_ini.MyListen}
	read_ini.Watch("../read_ini/test1.ini", hide_on_bush)

}
