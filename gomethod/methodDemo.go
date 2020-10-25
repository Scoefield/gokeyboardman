package main

import (
	"fmt"
)

/*
	golang 方法示例
 */

// 定义一个 Employee 结构体
type Employee struct {
	name     string
	salary   int
	currency string
}

// displaySalary()方法
func (e Employee) displaySalary() {
	fmt.Printf("Salary of %s is %s%d \n", e.name, e.currency, e.salary)
}

func main() {
	emp1 := Employee {
		name:     "Sam Adolf",
		salary:   5000,
		currency: "$",
	}
	// 方法调用
	emp1.displaySalary()
}