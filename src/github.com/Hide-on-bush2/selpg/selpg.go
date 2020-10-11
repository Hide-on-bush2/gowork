package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

type selpg_args struct {
	start_page     int
	end_page       int
	lines_per_page int
	page_type      bool // -f: true, -lNumber: false
	in_filename    string
	dest           string
}

func get_arg(args *selpg_args) {
	pflag.IntVarP(&(args.start_page), "start_page", "s", -1, "start page")
	pflag.IntVarP(&(args.end_page), "end_page", "e", -1, "end page")
	pflag.IntVarP(&(args.lines_per_page), "lines_per_page", "l", 72, "lines per page")
	pflag.StringVarP(&(args.dest), "dest", "d", "", "dest")
	pflag.BoolVarP(&(args.page_type), "page_type", "f", false, "page type")
	pflag.Parse() //把用户传递的命令行参数解析为对应变量的值

	file_name := pflag.Args()
	if len(file_name) > 0 {
		args.in_filename = string(file_name[0])
	} else {
		args.in_filename = ""
	}

}

func print_parameter(args *selpg_args) {
	fmt.Println("start page: ", args.start_page)
	fmt.Println("end page: ", args.end_page)
	fmt.Println("lines per page: ", args.lines_per_page)
	if args.page_type {
		fmt.Println("page type: ")
	} else {
		fmt.Println("page type: ")
	}
	fmt.Println("dest: ", args.dest)
	fmt.Println("in filename: ", args.in_filename)
}

func main() {
	var args selpg_args
	get_arg(&args)
	print_parameter(&args)
}
