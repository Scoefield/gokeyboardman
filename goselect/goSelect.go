package main

import "fmt"

/*
题目：
	已知，当select 存在多个 case时会随机选择一个满足条件的case执行。
	现在我们有一个需求：我们有一个函数会持续不间断地从ch1和ch2中分别接收任务1和任务2，
	如何确保当ch1和ch2同时达到就绪状态时，优先执行任务1，在没有任务1的时候再去执行任务2呢？
*/

// select解法一
func worker(ch1, ch2 <-chan int, stopCh chan struct{}) {

	for {
		select {
		case <-stopCh:
			return
		case job1 := <-ch1:
			fmt.Println(job1)
		default:
			select {
			case job2 := <-ch2:
				fmt.Println(job2)
			default:
			}
		}
	}
}

// select 解法二
func worker2(ch1, ch2 <-chan int, stopCh chan struct{}) {
	for {
		select {
		case <-stopCh:
			return
		case job1 := <-ch1:
			fmt.Println(job1)
		case job2 := <-ch2:
		priority:
			for {
				select {
				case job1 := <-ch1:
					fmt.Println(job1)
				default:
					break priority
				}
			}
			fmt.Println(job2)
		}
	}
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	stop := make(chan struct{})
	go func() {
		for i := 1; i < 11; i++ {
			ch1 <- i
		}
	}()
	go func() {
		for i := 1; i < 11; i++ {
			ch2 <- i
		}
		stop <- struct{}{}
	}()

	//worker(ch1, ch2, stop)
	worker2(ch1, ch2, stop)
}
