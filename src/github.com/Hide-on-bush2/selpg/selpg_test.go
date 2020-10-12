package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
)

func Test_myselpg(t *testing.T) {
	cmd := exec.Command("selpg", "-s1", "-e4", "-f", "faker.txt")
	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Run()
	if err != nil {
		// t.Error(err)
		log.Fatal(err)
	}
	fmt.Println(output.String())

	file, err := os.Open("faker.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)
	bytesread, err := file.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("bytes read:", bytesread)
	fmt.Println(string(buffer))

	if string(buffer) != output.String() {
		t.Error("Fail")
	}
}
