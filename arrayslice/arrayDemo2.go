package main

import "fmt"

func main() {
	a := [...]float64{67.7, 89.8, 21, 78}
	for i := 0; i < len(a); i++ { // 循环遍历从0到数组长度，获取数组元素
		fmt.Printf("%d th element of a is %.2f\n", i, a[i])
	}
}