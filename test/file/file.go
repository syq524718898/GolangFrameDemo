package main

import (
	"bufio"
	"fmt"
	"os"
)

func main()  {
	file, e := os.OpenFile("test.txt", os.O_RDWR, 0777)
	if e != nil {
		panic(e)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	b := make([]byte, 1024, 1024)
	n, e := reader.Read(b)
	if e != nil {
		panic(e)
	}
	fmt.Println(n)
	fmt.Println(string(b))
}

func writeTest()  {
	file, _ := os.Create("test.txt")
	defer file.Close()
	writer := bufio.NewWriter(file)
	writer.WriteString("adsfasdf")
	writer.Flush()
}
