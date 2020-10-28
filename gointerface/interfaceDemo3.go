package main

import "fmt"

type Student2 struct {
}

func main() {​
	var i1 interface{} = new(Student2)
	s := i1.(Student2) //不安全，如果断言失败，会直接panic
	fmt.Println(s)

	var i2 interface{} = new(Student2)
	//安全，断言失败，也不会panic，只是ok的值为false
	s, ok := i2.(Student2)
	if ok {
		fmt.Println(s)
	}
}