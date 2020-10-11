package main

import (
	"fmt"
	"sync"
)

// 创建一个结构体
type Manager struct {
	Name string
	Age int
}

// 声明两个全局变量，一个是Manager结构体指针，一个是用于单例等noce
var m *Manager
var once sync.Once

// 创建Manager单例函数
func GetManage() *Manager {
	once.Do(func() {
		m = &Manager{}
	})
	return m
}

// 用于测试打印Manager信息的函数
func (m *Manager) ManagerInfo() {
	fmt.Printf("Manager info, name=%s, age=%d\n", m.Name, m.Age)
}

func main() {
	m := GetManage()
	m.Name = "Jack"
	m.Age = 23
	m.ManagerInfo()
}
