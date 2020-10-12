package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

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
		args.in_filename = ""
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
		fmt.Println("page type: f")
	} else {
		fmt.Println("page type: l")
	}
	fmt.Println("dest: ", args.dest)
	fmt.Println("in filename: ", args.in_filename)
}

/*
打印错误信息
*/
func Error(message string) {
	fmt.Fprintf(os.Stderr, "[Error]:"+message)
	os.Exit(1)
}

/*
检查参数格式
*/
func check_args(args *selpg_args) {
	//必须有开始页和结束页并且开始页和结束页不能小于等于0并且结束页不能小于开始页
	if args.start_page <= 0 || args.end_page <= 0 {
		Error("Invalid start page and end page parameter\n")
	}
	//-l和-f参数不能同时输入
	if args.lines_per_page != 72 && args.page_type {
		Error("Cannot input -l and -f at the same time\n")
	}
	//自定义页长不可以小于等于0
	if args.lines_per_page <= 0 {
		Error("the lines per page cannot less than zero\n")
	}

}

/*
执行命令
*/
func Exec(args *selpg_args) {
	//检查输入文件名，如果没有输入文件名从标准输入流中读取
	//如果有输入文件名，则从输入文件名中读取
	//检查输出文件名，如果没有输出文件名则打印到标准输出流中
	//如果有输入文件名，则写入到输出文件中
	var in_file *os.File
	if args.in_filename == "" {
		in_file = os.Stdin
	} else {
		//检查文件是否存在
		_, exist_file := os.Stat(args.in_filename)
		if os.IsNotExist(exist_file) {
			Error("input file not exist\n")
		}

		var err error
		in_file, err = os.Open(args.in_filename)

		//检查文件是否打开
		if err != nil {
			Error("file not open\n")
		}
	}

	if args.dest == "" {
		output(os.Stdout, in_file, args)
	} else {
		cmd := exec.Command("cat", "-d"+args.dest)
		fout, err := cmd.StdinPipe()
		if err != nil {
			Error("StdinPipe")
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		errStart := cmd.Run()
		if errStart != nil {
			Error("CMD Run")
		}

		output(fout, in_file, args)
	}

}

/*
接受dest和in_file的文件句柄，将in_file中的特定内容输出到dest中
*/

func output(fout interface{}, in_file *os.File, args *selpg_args) {
	lines := 0
	pages := 0
	buf := bufio.NewReader(in_file)

	for true {
		var line string
		var err error

		if args.page_type {
			//以换页符定义翻页
			line, err = buf.ReadString('\f')
			pages += 1
		} else {
			//自定义页长
			line, err = buf.ReadString('\n')
			lines += 1

			if lines > args.lines_per_page {
				pages += 1
				lines = 1
			}
		}

		// //读到文件结尾
		// if err == io.EOF {
		// 	break
		// }

		//检查出错
		if err != nil && err != io.EOF {
			Error("file read error")
		}

		//如果在需要输出的范围内，则输出
		if pages >= args.start_page && pages <= args.end_page {
			var output_error error

			stdOutput, ok1 := fout.(*os.File)
			if ok1 {
				_, output_error = fmt.Fprintf(stdOutput, "%s", line)
				if output_error != nil {
					Error("output file error\n")
				}
				continue
			}

			pipeOutput, ok2 := fout.(io.WriteCloser)
			if ok2 {
				_, output_error = pipeOutput.Write([]byte(line))
				if output_error != nil {
					Error("output file error\n")
				}
				continue
			}

			if output_error != nil {
				Error("fout type error\n")
			}
		}

		if err == io.EOF {
			break
		}
	}

	if pages < args.start_page {
		Error("start page too large\n")
	} else if pages < args.end_page {
		Error("end page too large\n")
	}
}

func main() {
	var args selpg_args
	get_arg(&args)
	check_args(&args)
	// print_parameter(&args)
	Exec(&args)
}
