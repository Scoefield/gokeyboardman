package main

import (
	"fmt"
)

func main() {
	// 判断 num 是否为偶数
	if num := 10; num % 2 == 0 {
		fmt.Println(num, "is even")
	}  else {
		fmt.Println(num, "is odd")
	}
}