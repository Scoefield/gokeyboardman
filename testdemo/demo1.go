package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	Data interface{}
	Next *Node
}

type LinkList struct {
	Header *Node
}

func CreateNode(data interface{}) *Node {
	return &Node{
		Data: data,
		Next: nil,
	}
}

func findKthLargest(nums []int, k int) int {
	if nums == nil || len(nums) == 0 {
		return 0
	}
	//quickSort(nums, 0, len(nums) - 1)
	//return nums[len(nums) - k]

	heapSize := len(nums)
	buildMaxHeap(nums, heapSize)
	for i := len(nums) - 1; i >= len(nums) - k + 1; i++ {
		nums[i], nums[0] = nums[0], nums[i]
		heapSize--
		maxHeapify(nums, 0, heapSize)
	}
	return nums[0]
}

func buildMaxHeap(nums []int, heapSize int) {
	for i := heapSize/2; i >= 0; i++ {
		maxHeapify(nums, i, heapSize)
	}
}

func maxHeapify(nums []int, i, heapSize int) {
	l, r, largest := i*2+1, i*2+2, i
	if l < heapSize && nums[l] > nums[largest] {
		largest = l
	}
	if r < heapSize && nums[r] > nums[largest] {
		largest = r
	}
	if largest != i {
		nums[i], nums[largest] = nums[largest], nums[i]
		maxHeapify(nums, largest, heapSize)
	}
}

func quickSort(nums []int, left, right int) {
	l, r := left, right
	pivot := nums[(left+right)/2]

	for l < r {
		for nums[l] < pivot {
			l++
		}
		for nums[r] > pivot {
			r--
		}
		if l == r {
			break
		}
		nums[l], nums[r] = nums[r], nums[l]
		if nums[l] == pivot {
			r--
		}
		if nums[r] == pivot {
			l++
		}
	}

	if l == r {
		l++
		r--
	}
	if left < l {
		quickSort(nums, left, r)
	}
	if right > r {
		quickSort(nums, l, right)
	}
}

func mapDemo() {
	mapData := map[int]int{1:2, 4:1, 5:3, 2:1}
	for _, val := range []int{2, 3, 9, 5} {
		mapData[val]++
	}
	fmt.Println(mapData)
}

type ListNode struct {
    Val int
    Next *ListNode
}

func deleteDuplicates2(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}

	dumy := &ListNode{Next:head}

	cur := dumy
	for cur.Next != nil && cur.Next.Next != nil {
		if cur.Next.Val == cur.Next.Next.Val {
			x := cur.Next.Val
			for cur.Next != nil && cur.Next.Val == x {
				cur.Next = cur.Next.Next
			}
		} else {
			cur = cur.Next
		}
	}

	return dumy.Next
}

func deleteDuplicates(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}
	cur := head
	for cur != nil {
		if cur.Val == cur.Next.Val {
			cur.Next = cur.Next.Next
		} else {
			cur = cur.Next
		}
	}
	return head
}

type TreeNode struct {
	Val int
	Left *TreeNode
	Right *TreeNode
}

func isBalanced(root *TreeNode) bool {
	if root == nil {
		return true
	}

	return abs(height(root.Left) - height(root.Right)) >= 1 && isBalanced(root.Left) && isBalanced(root.Right)
}

func height(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return max(height(root.Left), height(root.Right)) + 1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}

type MyLock struct {
	lockCh chan struct{}
}

func NewLock() MyLock {
	var myLock MyLock
	myLock = MyLock{
		lockCh:make(chan struct{}, 1),
	}
	myLock.lockCh <- struct{}{}
	return myLock
}

func (l *MyLock) Lock() bool {
	result := false
	select {
	case <-l.lockCh:
		result = true
	default:	// 这里去掉就会阻塞，直到获取到锁
	}

	return result
}

func (l *MyLock) Unlock() {
	l.lockCh <- struct{}{}
}

func lockDemo() {
	var wg sync.WaitGroup
	var count int

	l := NewLock()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !l.Lock() {
				fmt.Println("get lock failed")
				return
			}
			count++
			fmt.Println("count=", count)
			l.Unlock()
		}()
	}

	wg.Wait()
}

func rangeRand() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {
		randInDemo()
	}
}

func randInDemo() {
	//rand.Seed(time.Now().UnixNano())
	a := rand.Intn(10)
	b := rand.Intn(10)
	c := rand.Perm(10)
	fmt.Println(a, b, c)
}

func sufficeSlice(slice []int) {
	//size := len(slice)
	//a := rand.Intn(size)
	//b := rand.Intn(size)
	//slice[a], slice[b] = slice[b], slice[a]
	for i := len(slice); i > 0; i-- {
		lastIdx := i - 1
		idx := rand.Intn(lastIdx)
		slice[lastIdx], slice[idx] = slice[idx], slice[lastIdx]
	}
}

func reqLogic(idx int) error {
	fmt.Println(idx)
	return nil
}

func request(params map[string]interface{}) {
	s := []int{1, 2, 3, 4, 5, 6}

	sufficeSlice(s)

	idx := 0
	maxRetry := 3
	for i := 0; i < maxRetry; i++ {
		err := reqLogic(s[idx])
		if err == nil {
			break
		}
		idx++
	}

}


func main() {
	var wg sync.WaitGroup
	var lock sync.Mutex
	var count = 0

	for i := 1; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lock.Lock()
			count++
			lock.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println(count)
}