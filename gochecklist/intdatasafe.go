package main

import "fmt"

// bad: 未限制长度，导致整数溢出
func overflowBad(numControlByUser int32) {
	var numInt int32 = 0
	numInt = numControlByUser + 1
	//对长度限制不当，导致整数溢出
	fmt.Printf("%d\n", numInt)
	//使用numInt，可能导致其他错误
}

// good:
func overflowGood(numControlByUser int32) {
	var numInt int32 = 0
	numInt = numControlByUser + 1
	if numInt < 0 {
	fmt.Println("integer overflow")
	return;
	}
	fmt.Println("integer ok")
	}

func main() {
	overflowBad(2147483647)
	overflowGood(2147483647)
}