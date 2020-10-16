package main

import "fmt"

func main() {
	// 定义局部变量 title，和num
	var num = 6
	var title = "Golang流程控制语句"

	if num == 6  {
		fmt.Println(title)
	}else if num < 6 {
		fmt.Println("num less than 6")
	}else {
		fmt.Println("num greater than 6")
	}
}
