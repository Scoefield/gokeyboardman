package main

import (
	"fmt"
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

func CreateLinkList() *LinkList {
	return &LinkList{
		Header: CreateNode(nil),
	}
}

func (l *LinkList) Append(data interface{}) {
	newNode := CreateNode(data)
	current := l.Header
	for current.Next != nil {
		current = current.Next
	}
	current.Next = newNode
}

func (l *LinkList) turnNode() {
	tmp := l.Header
	var pre *Node
	for tmp.Next != nil {
		tmp, tmp.Next, pre = tmp.Next, pre, tmp
	}
	l.Header = pre
}

func (l *LinkList) scanLink() {
	tmp := l.Header
	for tmp.Next != nil {
		fmt.Println(tmp.Next.Data)
		tmp = tmp.Next
	}
}


func main() {
	l := CreateLinkList()
	l.Append(1)
	l.Append(2)
	l.Append(3)
	l.Append(4)
	l.Append(5)
	l.turnNode()
	l.scanLink()
}