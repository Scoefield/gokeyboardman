package main

import (
	"fmt"
)

// 电话结构体
type Phone interface {
	call()
}

// 诺基亚电话
type NokiaPhone struct {
}

// 实现诺基亚电话的 call
func (nokiaPhone NokiaPhone) call() {
	fmt.Println("I am Nokia, I can call you!")
}

// 苹果电话
type IPhone struct {
}

// 实现苹果电话的 call
func (iPhone IPhone) call() {
	fmt.Println("I am iPhone, I can call you!")
}

func main() {
	var phone Phone
	phone = new(NokiaPhone)
	phone.call()

	phone = new(IPhone)
	phone.call()
}