package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Put函数源码分析可得以下结论
/*
1. 如果放入的值为空，直接return.
2. 检查当前goroutine的是否设置对象池私有值，如果没有则将x赋值给其私有成员，并将x设置为nil。
3. 如果当前goroutine私有值已经被设置，那么将该值追加到共享列表。
 */

// Get函数源码分析可得以下结论
/*
1. 尝试从本地P对应的那个本地池中获取一个对象值, 并从本地池冲删除该值。
2. 如果获取失败，那么从共享池中获取, 并从共享队列中删除该值。
3. 如果获取失败，那么从其他P的共享池中偷一个过来，并删除共享池中的该值(p.getSlow())。
4. 如果仍然失败，那么直接通过New()分配一个返回值，注意这个分配的值不会被放入池中。New()返回用户注册的New函数的值，如果用户未注册New，那么返回nil。
 */

/*
最后我们来看一下init函数。

func init() {
    runtime_registerPoolCleanup(poolCleanup)
}
可以看到在init的时候注册了一个PoolCleanup函数，他会清除掉sync.Pool中的所有的缓存的对象，这个注册函数会在每次GC的时候运行，
所以sync.Pool中的值只在两次GC中间的时段有效。
 */

// 创建临时对象池，new返回默认值
func GetPool() *sync.Pool {
	return &sync.Pool{New: func() interface{} {
		return 0
	}}
}

// 程序中间无GC效果测试
func TestSyncPoolWithoutGC()  {
	p := GetPool()

	// 获取对象值
	a := p.Get().(int)
	// 添加对象值
	p.Put(2)
	//再次获取
	b := p.Get().(int)

	fmt.Println(a, b)	// 输出0， 2
}

// 程序中间有GC效果测试
func TestSyncPoolWithGC() {
	p := GetPool()

	// 获取对象值
	a := p.Get().(int)
	// 添加对象值
	p.Put(2)
	// 手动GC
	runtime.GC()
	//再次获取
	b := p.Get().(int)

	fmt.Println(a, b)	// 输出0， 0
}

func main() {
	TestSyncPoolWithoutGC()
	TestSyncPoolWithGC()
}
