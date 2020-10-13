package main

import "fmt"

// 声明定义了一个计算价钱的函数，函数名：calculatePrice，
// 参数：price和num，返回值：totalPrice
func calculatePrice(price int, num int) int {
	// 计算总价钱：单价 x 数量
	totalPrice := price * num
	return totalPrice
}

// main入口函数
func main() {
	price, num := 45, 2		// 定义 price 和 num变量，推断类型为 int
	totalPrice := calculatePrice(price, num)	// 调用calculatePrice函数
	fmt.Println("Total price is", totalPrice)	// 打印返回结果
}