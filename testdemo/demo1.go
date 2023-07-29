package main

type Node struct {
	Key int
	Data int
	Pre *Node
	Next *Node
}

type LRUCache struct {
	Cap int
	DataMap map[int]*Node
	Head *Node
	Tail *Node
}

func NewLRUCache(cap int) *LRUCache {
	cache := &LRUCache{
		Cap:     cap,
		DataMap: make(map[int]*Node, cap),
		Head:    &Node{},
		Tail:    &Node{},
	}
	cache.Head.Next = cache.Tail
	cache.Tail.Pre = cache.Head
	return cache
}

func (l *LRUCache) SetHead(node *Node) {
	node.Next = l.Head.Next
	l.Head.Next.Pre = node
	l.Head.Next = node
	node.Pre = l.Head
}

func (l *LRUCache) Remove(node *Node) {
	node.Pre.Next = node.Next
	node.Next.Pre = node.Pre
}

func (l *LRUCache) Get(key int) int {
	node, ok := l.DataMap[key]
	if !ok {
		return -1
	}
	l.Remove(node)
	l.SetHead(node)
	return node.Data
}

func (l *LRUCache) Put(key, data int) {
	node, ok := l.DataMap[key]
	if ok {
		l.Remove(node)
	} else {
		if len(l.DataMap) == l.Cap {
			delete(l.DataMap, l.Tail.Pre.Key)
			l.Remove(l.Tail.Pre)
		}
		node := &Node{
			Key:  key,
			Data: data,
		}
		l.DataMap[key] = node
	}
	l.SetHead(node)
}

type ListNode struct {
	Val int
	Next *ListNode
}

// 反转链表
func reverseList(head *ListNode) *ListNode {
	if head == nil {
		return head
	}

	cur := head
	var pre *ListNode
	for cur != nil {
		cur, cur.Next, pre = cur.Next, pre, cur
	}
	return pre
}

// 判断环形链表
func hasCycle(head *ListNode) bool {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow, fast = slow.Next, fast.Next.Next
		if slow == fast {
			return true
		}
	}
	return false
}


func main() {

}
