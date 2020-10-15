package main

import "sync"

// 定义单例的结构体
type Singleton2 struct {
}

// 声明一个Singleton结构体指针的变量
var instance *Singleton2
// 声明普通锁变量，mu
var mu sync.Mutex

// 获取带锁的单例函数
func GetInstance() *Singleton2 {
	// 加锁
	mu.Lock()
	// 函数退出前解锁
	defer mu.Unlock()

	// 如果为空，则创建单例
	if instance == nil {
		instance = &Singleton2{}
	}
	return instance
}

func main() {
	
}
