package main

import (
	"fmt"
	"os"
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

/*
从输入流中获取参数保存到sel_args结构体中
*/
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
		args.in_filename = ""sa
	}

}

/*
打印参数列表
*/
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

/*
打印错误信息
*/
func args_error(message string){
	fmt.Fprintf(os.Stderr, "[Error]:" + message)
	os.Exit(1)
}


/*
检查参数格式
*/
func check_args(args *selpg_args){
	//必须有开始页和结束页并且开始页和结束页不能小于等于0并且结束页不能小于开始页
	if args.start_page <= 0 || args.end_page <= 0 {
		args_error("Invalid start page and end page parameter\n")
	}
	//-l和-f参数不能同时输入
	if args.lines_per_page != 72 && args.page_type {
		args_error("Cannot input -l and -f at the same time\n")
	}
	//自定义页长不可以小于等于0
	if args.lines_per_page <= 0 {
		args_error("the lines per page cannot less than zero\n")
	}
	
}

/*
执行命令
*/
func exec(args *selpg_args){
	//检查输入文件名，如果没有输入文件名从标准输入流中读取
	//如果有输入文件名，则从输入文件名中读取
	//检查输出文件名，如果没有输出文件名则打印到标准输出流中
	//如果有输入文件名，则写入到输出文件中
	var in_file *os.File
	if args.in_filename == "" {
		in_file = os.Stdin
	}else{
		
	}
}

func main() {
	var args selpg_args
	get_arg(&args)
	check_args(&args)
	print_parameter(&args)
}
