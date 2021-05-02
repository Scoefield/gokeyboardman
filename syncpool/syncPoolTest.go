package main
import (
	"fmt"
	"sync"
)
// 定义一个 Person 结构体，有Name和Age变量
type Person struct {
	Name string
	Age int
}
// 初始化sync.Pool，new函数就是创建Person结构体
func initPool() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			fmt.Println("创建一个 person.")
			return &Person{}
		},
	}
}
// 主函数，入口函数
func main() {
	pool := initPool()
	person := pool.Get().(*Person)
	fmt.Println("首次从sync.Pool中获取person：", person)
	person.Name = "Jack"
	person.Age = 23
	pool.Put(person)
	fmt.Println("设置的对象Name: ", person.Name)
	fmt.Println("设置的对象Age: ", person.Age)
	fmt.Println("Pool 中有一个对象，调用Get方法获取：", pool.Get().(*Person))
	fmt.Println("Pool 中没有对象了，再次调用Get方法：", pool.Get().(*Person))
}
