package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

//有三个需要并发执行的函数，三个函数会返回同一类型的值，且三个函数执行时间在2ms到10ms之间不等，
//而主程序要求在5ms内返回结果，若5ms内没有执行完毕，则强制返回结果，这个时候某个函数可能还没有返回因此没有值。
func practice1() {
	ch := make(chan int)

	go demo1(ch)
	go demo2(ch)
	go demo3(ch)

	select {
	case val := <-ch:
		fmt.Println("get ch data:", val)
	case <-time.After(time.Millisecond * 5):
		fmt.Println("after 5 ms timeout.")
	}
}

func demo1(ch chan int) {
	time.Sleep(time.Millisecond * 8)

	ch <- 1
}

func demo2(ch chan int) {
	time.Sleep(time.Millisecond * 2)
	ch <- 2
}

func demo3(ch chan int) {
	time.Sleep(time.Millisecond * 7)
	ch <- 3
}

func onceDemo() {
	var once sync.Once
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			//fmt.Println(i)
			once.Do(func() {
				fmt.Println(i)
			})
		}(i)
	}

	wg.Wait()
}

type instance struct {

}

var (
	singleton *instance
	once sync.Once
	initialized uint32
	lock sync.Mutex
)

func Instance() *instance {
	once.Do(func() {
		singleton = &instance{}
	})
	return singleton
}

func Instance2() *instance {
	if atomic.LoadUint32(&initialized) == 1 {
		return singleton
	}

	lock.Lock()
	defer lock.Unlock()

	if singleton == nil {
		atomic.StoreUint32(&initialized, 1)
		singleton = &instance{}
	}
	return singleton
}

func instanceDemo() {
	Instance()
}

func main() {
	//practice1()
	onceDemo()
}