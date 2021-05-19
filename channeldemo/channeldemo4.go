package main

import (
	"fmt"
	"sync"
	"time"
)

func chanDemo() {
	// channel 使用前必须初始化，否则是一个nil
	// 无缓冲channel，必须同时要有输入和输出
	ch := make(chan int)

	fmt.Println("****init: ", ch)	// 输出地址：0xc00008c060

	go func() {
		// 输出channel
		ret := <-ch
		fmt.Println(ret)
	}()

	// 输入channel
	ch <- 1

	// 关闭channel
	close(ch)

	// 关闭后不可chan写
	//ch <- 2	// 这里会报错：panic: send on closed channel

	// 关闭后的 chan 是还可以读的
	ret, ok := <-ch
	fmt.Println("****close: ", ret, ok)	// 输出：0 false
}

// 模拟 100 个人抢 10 个鸡蛋
func eggsChanDemo() {
	// 初始化 10 个鸡蛋的 chan
	eggsChan := make(chan int, 10)

	// 输入 10 个值
	for i := 1; i <= 10; i++ {
		eggsChan <- i
	}

	var wg sync.WaitGroup

	// 模拟 100 个人抢 10 个鸡蛋
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func(i int) {
			select {
			case egg := <-eggsChan:
				fmt.Printf("People %d get egg %d\n", i, egg)
			default:
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}

var dataCh = make(chan int, 10)
// 模拟生产者和消费者案例
func producerConsumer()  {
	// ********* 正常情况 ********
	// 10 个生产者
	//for i := 0; i < 10; i++ {
	//	//	go producer(i)
	//	//}

	// 10 个消费者
	//for i := 0; i < 10; i++ {
	//	go consumer(i)
	//}

	// ********* 阻塞情况一 ********
	// 100 个生产者
	//for i := 0; i < 100; i++ {
	//	go producer(i)
	//}

	// 10 个消费者
	//for i := 0; i < 10; i++ {
	//	go consumer(i)
	//}

	// ********* 阻塞情况二 ********
	// 10 个生产者
	for i := 0; i < 10; i++ {
		go producer(i)
	}

	// 100 个消费者
	for i := 0; i < 100; i++ {
		go consumer(i)
	}


	time.Sleep(time.Second * 2)
}

// 生产者
func producer(index int) {
	dataCh <- index
	//fmt.Printf("Producer %d send %d\n", index, index)
}

// 消费者
func consumer(index int) {
	fmt.Printf("Consumer %d recevie %d\n", index, <-dataCh)
}

func main() {
	//chanDemo()
	//eggsChanDemo()
	producerConsumer()

}
