package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var configure map[string]string

type iniReader struct {
	r io.Reader
}

func (rr iniReader) Read(p []byte) (n int, err error) {
	return rr.r.Read(p)
}

func checkError(err error, message string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error]:"+message)
		os.Exit(1)
	}
}

/*
打印错误信息
*/
func Error(message string) {
	fmt.Fprintf(os.Stderr, "[Error]:"+message)
}

func getReader(filename string) (*bufio.Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		Error("Read File\n")
		return nil, err
	}

	// defer file.Close()
	fileReader := bufio.NewReader(file)
	// rr := iniReader{fileReader}
	// io.Copy(os.Stdout, &rr)
	// getConf(fileReader)
	return fileReader, nil
}

/*
处理字符串
*/
func processString(str string) (key string, value string) {
	k := ""
	v := ""

	i := 0
	//找出key
	for i < len(str) && str[i] != ' ' {
		k += string(str[i])
		i++
	}
	//过滤掉=号
	for i < len(str) && (str[i] == ' ' || str[i] == '=') {
		i++
	}

	//找出value
	for i < len(str) {
		v += string(str[i])
		i++
	}

	return k, v
}

func getConf(reader *bufio.Reader) (map[string]string, error) {
	res := make(map[string]string)

	for {
		//读取文件的一行
		linestr, err := reader.ReadString('\n')
		last_line := false
		if err == io.EOF {
			last_line = true
		} else if err != nil {
			Error("Read Line\n")
			return nil, err
		}

		//切掉行的左右两边的空白字符
		linestr = strings.TrimSpace(linestr)

		//忽略空行
		if linestr == "" {
			continue
		}

		//忽略注释
		if linestr[0] == '#' {
			continue
		}

		//检测段名，但不做处理，后面有需求再写
		if linestr[0] == '[' {
			continue
		}

		//处理字符串，把k,v放到map中
		k, v := processString(linestr)
		res[k] = v

		// fmt.Printf("key:%s, value:%s\n", k, v)

		if last_line {
			break
		}
	}

	return res, nil
}

type ListenFunc func(string)

type Listener interface {
	listen(inifile string)
}

type configuration map[string]string

func Watch(filename string, listener Listener) (configuration, error) {
	reader, err := getReader(filename)
	if err != nil {
		return nil, err
	}

	configure, err = getConf(reader)
	if err != nil {
		return nil, err
	}

	printConfigue()
	// for k, v := range configure {
	// 	fmt.Printf("key:%s, value:%s\n", k, v)
	// }
	for {
		listener.listen(filename)
	}
	return configure, nil
}

type listen_methods struct {
	f1 ListenFunc
}

func printConfigue() {
	for k, v := range configure {
		fmt.Printf("key:%s, value:%s\n", k, v)
	}
	fmt.Println()
}

func (methods listen_methods) listen(inifile string) {
	methods.f1(inifile)
}

func myListen(filename string) {
	reader, err := getReader(filename)
	if err != nil {
		Error("Listening Read File\n")
		os.Exit(1)
	}

	tmpConf, err := getConf(reader)
	if err != nil {
		Error("Listening Read Line\n")
		os.Exit(1)
	}

	for k, v := range tmpConf {
		if _, ok := configure[k]; !ok || v != configure[k] {
			configure = tmpConf
			fmt.Println("the configure file has been changed")
			printConfigue()
			break
		}
	}

	return

}

func main() {
	// sr := strings.NewReader("Lbh penpxrq gur pbqr!")
	// rr := iniReader{sr}
	// io.Copy(os.Stdout, &rr)
	// readFile("$GOPATH/src/github.com/Hide-on-bush2/read_ini/test.ini")
	// readFile("./test.ini")

	hide_on_bush := listen_methods{myListen}
	Watch("./test.ini", hide_on_bush)
}
