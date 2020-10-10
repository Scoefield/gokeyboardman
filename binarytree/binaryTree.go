package main

import (
	"encoding/json"
	"fmt"
)

/*
实现二叉树翻转
 */

// 定义节点结构体
type Node struct {
	Key interface{}
	Left *Node
	Right *Node
}

// 二叉树翻转方法
func (n *Node) TurnNode()  {
	n.Left, n.Right = n.Right, n.Left

	if n.Left != nil {
		n.Left.TurnNode()
	}
	if n.Right != nil {
		n.Right.TurnNode()
	}
}

// 测试函数
func TestDemo() {
	root := &Node{
		Key:   4,
	}

	// 创建左子树
	root.Left = &Node{
		Key:   2,
		Left:  &Node{
			Key:   1,
			Left:  nil,
			Right: nil,
		},
		Right: &Node{
			Key:   3,
			Left:  nil,
			Right: nil,
		},
	}

	// 创建左子树
	root.Right = &Node{
		Key:   7,
		Left:  &Node{
			Key:   6,
			Left:  nil,
			Right: nil,
		},
		Right: &Node{
			Key:   9,
			Left:  nil,
			Right: nil,
		},
	}

	root.TurnNode()
	// json序列化 并以友好的方式打印出来
	out, _ := json.MarshalIndent(root, "", "   ")
	fmt.Println(string(out))
}


func main() {
	TestDemo()
}
