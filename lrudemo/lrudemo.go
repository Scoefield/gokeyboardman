package main

import "fmt"

// LRU缓存结构
type LRUCache struct {
	Cap int				// 缓存大小
	Map map[int]*Node	// 存储节点数据的map
	Head *Node			// 头节点
	Last *Node			// 尾节点
}

// 双向链表 节点
type Node struct {
	Key int			// key
	Val int			// 数据值
	Pre *Node		// 前一个节点的指针
	Next *Node		// 下一个节点的指针
}

// 创建 LRU 缓存结构，一些初始化操作
func NewLRUCache(cap int) *LRUCache {
	cache := &LRUCache{
		Cap:  cap,
		Map:  make(map[int]*Node, cap),
		Head: &Node{},
		Last: &Node{},
	}

	// 双向链表 初始化
	cache.Head.Next = cache.Last
	cache.Last.Pre = cache.Head
	return cache
}

// 设置头节点
func (l *LRUCache) setHeader(node *Node) {
	l.Head.Next.Pre = node
	node.Next = l.Head.Next
	l.Head.Next = node
	node.Pre = l.Head
}

// 删除节点
func (l *LRUCache) remove(node *Node) {
	//l.Head.Next = node.Next
	//node.Next.Pre = l.Head
	node.Pre.Next = node.Next
	node.Next.Pre = node.Pre
}

// 通过 key 获取数据
// 获取不到直接返回 -1
// 获取到，则先删除获取到的节点，在将该节点放到头节点
// 返回获取的值
func (l *LRUCache) Get(key int) int {
	node, ok := l.Map[key]
	if !ok {
		fmt.Println(-1)
		return -1
	}
	l.remove(node)
	l.setHeader(node)
	fmt.Println(node.Val)
	return node.Val
}

// 加入缓存操作
// 先通过 key 获取节点数据
// 如果获取到，则删除掉该节点，然后添加到头部
// 如果获取不到，则判断缓存是否满了
// 如果缓存满了，则删掉最后一个节点数据
// 最后将节点数据放到头部
func (l *LRUCache) Put(key, value int) {
	node, ok := l.Map[key]
	if ok {
		l.remove(node)
	} else {
		if len(l.Map) == l.Cap {
			delete(l.Map, l.Last.Pre.Key)
			l.remove(l.Last.Pre)
		}
		node = &Node{Key:key, Val:value}
		l.Map[key] = node
	}

	node.Val = value
	l.setHeader(node)
}

// 主函数，LRU算法测试
func main() {
	// 测试一：
	//lruCache := NewLRUCache(3)
	//
	//val := lruCache.Get(2)
	//fmt.Println(val)
	//
	//lruCache.Put(2, 22)
	//lruCache.Put(3, 33)
	//
	//val = lruCache.Get(2)
	//fmt.Println(val)
	//
	//lruCache.Put(4, 44)
	//lruCache.Put(5, 55)
	//
	//val = lruCache.Get(3)
	//fmt.Println(val)

	// 测试二：
	//lruCache := NewLRUCache(2)
	//lruCache.Put(1, 1)
	//lruCache.Put(2, 2)
	//lruCache.Get(1)
	//lruCache.Put(3, 3)
	//lruCache.Get(2)
	//lruCache.Put(4, 4)
	//lruCache.Get(1)
	//lruCache.Get(3)
	//lruCache.Get(4)

	// 测试三：
	lruCache := NewLRUCache(2)
	lruCache.Put(2, 1)
	lruCache.Put(2, 2)
	lruCache.Get(2)
	lruCache.Put(1, 1)
	lruCache.Put(4, 1)
	lruCache.Get(2)

}
