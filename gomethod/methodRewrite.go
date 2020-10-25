package main

import "fmt"

type Human2 struct {
	name  string
	age   int
	phone string
}
type Student2 struct {
	Human2  //匿名字段
	school string
}
type Employee3 struct {
	Human2   //匿名字段
	company string
}

//Human定义method
func (h *Human2) SayHi() {
	fmt.Printf("Hi, I am %s you can call me on %s\n", h.name, h.phone)
}

//Employee的method重写Human的method
func (e *Employee3) SayHi() {
	fmt.Printf("Hi, I am %s, I work at %s. Call me on %s\n", e.name,
		e.company, e.phone) //Yes you can split into 2 lines here.
}

func main() {
	mark := Student2{Human2{"Mark", 25, "222-222-YYYY"}, "MIT"}
	sam := Employee3{Human2{"Sam", 45, "111-888-XXXX"}, "Golang Inc"}
	mark.SayHi()
	sam.SayHi()
}