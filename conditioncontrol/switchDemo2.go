package main

import (
	"fmt"
)

type data [2]int

func main() {
	switch a := 3; a {
		case 3:
			a += 10
			fmt.Println(a)	// 输出13
			fallthrough		// 往下执行
		case 6:
			a += 20
			fmt.Println(a)	// 输出33
		default:
			fmt.Println(a)
	}
}