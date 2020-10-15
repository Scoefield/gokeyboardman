package main

import "fmt"

// 定义单例的结构体
type Singleton struct {

}

// 声明一个Singleton结构体指针的变量
var singleton *Singleton

// 获取单例的函数，返回Singleton结构体指针类型
func GetSingleton() *Singleton {
	// 如果为空，则创建单例
	if singleton == nil {
		singleton = &Singleton{}
	}
	return singleton
}

// 打印测试函数
func (s *Singleton) PrintInfo() {
	fmt.Println("test print info")
}

func main() {
	s := GetSingleton()
	s.PrintInfo()
}
