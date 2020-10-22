package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

// 使用 File 来 Read
func fileRead1(filename string) int {
	// open 文件，返回文件句柄
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	// 定义一个缓冲变量
	buf := make([]byte, 4096)
	var nbytes int
	for {
		n, err := fi.Read(buf)	// 每次 Read到放 buf 里
		if err != nil && err != io.EOF {
			panic(err)
		}
		// 读不到后退出循环
		if n == 0 {
			break
		}
		// 统计字节数
		nbytes += n
	}
	// 返回字节数
	return nbytes
}

// 使用 bufio 方式读取
func bufioRead2(filename string) int {
	// 同样是打开文件，返回文件句柄
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	// 定义一个缓冲变量
	buf := make([]byte, 4096)
	var nbytes int
	// bufio 方式读取
	rd := bufio.NewReader(fi)
	for {
		n, err := rd.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		nbytes += n
	}
	return nbytes
}

// 使用ioutil方式读取
func ioutilRead3(filename string) int {
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	// ioutil 方式读取
	fd, err := ioutil.ReadAll(fi)
	nbytes := len(fd)
	// 返回读取到的字节数
	return nbytes
}

// 测试函数
func testReadFile(filename string) {
	fmt.Printf("============test fileRead1, filename= %s ===========\n", filename)
	start := time.Now()
	size1 := fileRead1(filename)
	t1 := time.Now()
	fmt.Printf("Read 1 cost: %v, size: %d\n", t1.Sub(start), size1)

	fmt.Printf("============test bufioRead2, filename= %s ===========\n", filename)
	size2 := bufioRead2(filename)
	t2 := time.Now()
	fmt.Printf("Read 2 cost: %v, size: %d\n", t2.Sub(t1), size2)

	fmt.Printf("============test ioutilRead3, filename= %s ===========\n", filename)
	size3 := ioutilRead3(filename)
	t3 := time.Now()
	fmt.Printf("Read 3 cost: %v, size: %d\n", t3.Sub(t2), size3)
}

func main() {
	//fileName := "./test1.txt"
	//	//bytesNum := fileRead1(fileName)
	//	//fmt.Println(bytesNum)
	//	//
	//	//bytesNum2 := bufioRead2(fileName)
	//	//fmt.Println(bytesNum2)

	testReadFile("small.txt")
	testReadFile("middle.txt")
	testReadFile("large.txt")
}
