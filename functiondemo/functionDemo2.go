package main

import (
	"fmt"
)

// 计算周长和面积的函数，length：长，width：宽
func calcPerimeterArea(length, width float32) (float32, float32) {
	perimeter := (length + width) * 2  // 周长
	area := length * width  // 面积
	return perimeter, area
}

func main() {
	area, perimeter := calcPerimeterArea(3, 4)
	fmt.Printf("Perimeter=%f, Area=%f\n", perimeter, area)
}