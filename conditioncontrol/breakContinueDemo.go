package main

import (
	"fmt"
)

func main() {
	for i := 1; i <= 10; i++ {
		if i > 5 {
			break	// 如果 i 大于5，则跳出循环体
		}
		if i == 2 {
			continue	// 如果 i 等于2，则跳到下一个迭代，相当于过滤掉2
		}
		fmt.Println(i)
	}
	fmt.Println("跳出for循环后的打印")
}