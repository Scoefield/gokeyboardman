/*
	file：chatRoomServer.go
	runCmd：go run chatRoomServer.go
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"runtime"
	"strconv"
	"strings"
)

// 保存连接用户（客户端）的全局变量
var (
	userMgr *UserMgr
)

const (
	ClientJoin = 1	// 客户端加入
	ClientLeave = 2	// 客户端离开
)

// 用户（客户端）结构体
type UserMgr struct {
	// 在线用户（客户端）变量，map 类型
	// key为客户端端口号，用于标识别哪个用户，
	onlineUsers map[string]*UserProcess
}

// 用户（客户端）连接进程信息
type UserProcess struct {
	// 连接句柄
	Conn net.Conn
	// RemoteAddr 字段，表示该Conn是哪个用户
	UserAddr string
}

// 消息结构体
type Message struct {
	DataSource string	`json:"data_source"`	// 消息（数据）来源
	Data string	`json:"data"`					// 消息（数据）内容
}

// 读取消息（数据）函数，入参为 连接句柄，返回消息体 Message
func ReadData(conn net.Conn) (msg Message, err error) {
	// 定义一个缓存变量 buf，可以存放 4096 个字节
	buf := make([]byte, 4096)
	// 读取客户端发送的消息
	n, err := conn.Read(buf[:])
	if err != nil {
		return
	}

	// 将客户端地址和读取到的消息赋值到消息体结构
	msg.DataSource = conn.RemoteAddr().String()
	msg.Data = string(buf[:n])

	return
}

// 写消息（数据）函数，入参为 连接句柄 和 需要发送到消息（数据）data
func WriteData(conn net.Conn, data []byte) (err error) {
	// 写（发送）数据
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write fail", err)
		return
	}
	return
}

// 广播消息给每一个在线用户（客户端）
func SendMesToEachOnlineUser(userKey string, data []byte) {
	// 遍历在线用户（客户端）全局队列，给每一个在线用户（客户端）发消息
	for id, up := range userMgr.onlineUsers {
		// 排除掉本身（当前客户端）
		if id == userKey {
			continue
		}
		// 发送消息（数据）
		err := WriteData(up.Conn, data)
		// 发送出错时，打错误信息
		if err != nil {
			fmt.Printf("转发消息给[%s]客户端失败，失败信息为：%v\n", up.UserAddr, err)
		}
	}
}

// 服务端处理通讯消息函数
func ServerProcessMsg(userKey string, msg *Message) (err error) {
	// 打印提示信息
	fmt.Printf("客户端[%s]发送的消息: %s\n", msg.DataSource, msg.Data)
	// 序列化消息
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 广播消息给每一个在线的用户（客户端），当前客户端除外
	SendMesToEachOnlineUser(userKey, data)
	return
}

// 封装正式处理通讯的函数
func SubProcess(conn net.Conn, userKey string) (err error) {
	// 循环读取客户端发送过来的消息（数据）
	for {
		// 读取消息
		msg, err := ReadData(conn)
		// 报错或者客户端连接断开时，打印提示信息，并将该用户（客户端）在全局队列里剔除掉
		if err != nil {
			delete(userMgr.onlineUsers, userKey)
			if err == io.EOF {
				fmt.Printf("客户端[%s]退出，与服务器端的连接断开.\n", conn.RemoteAddr())

				// 将客户端离开通讯室的消息广播出去
				clientAddr := conn.RemoteAddr().String()
				JoinOrLeaveMsg(userKey, clientAddr, ClientLeave)

				return err
			} else {
				fmt.Println("ReadData err=", err)
				return err
			}
		}

		// 服务端处理消息
		err = ServerProcessMsg(userKey, &msg)
		if err != nil {
			return err
		}
	}
}

// 客户端加入或离开 ChatRoom，需要发送消息通知每个在线的客户端
func JoinOrLeaveMsg(userKey string, clientAddr string, msgType int) {
	var msgStr string
	switch msgType {
	case 1:
		msgStr = "加入"
	case 2:
		msgStr = "离开"
	default:
		fmt.Println("JoinOrLeaveMsg type wrong!")
	}

	// 将客户端加入聊天室的消息广播出去
	var msg Message
	joinStr := fmt.Sprintf("客户端[%s]已%s【ChatRoom】", clientAddr, msgStr)
	msg.DataSource = "服务端"
	msg.Data = joinStr
	joinMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
	}
	SendMesToEachOnlineUser(userKey, joinMsg)
}

// 处理和客户端的通讯，入参 连接句柄和 userKey
// userKey 其实为当前连接客户端的端口号，用于标识是哪个客户端
func Process(conn net.Conn, userKey string) {
	//这里需要延时关闭conn
	defer conn.Close()

	// 处理通讯
	err := SubProcess(conn, userKey)
	// 通讯出错时，打印错误信息
	if err != io.EOF {
		fmt.Printf("客户端[%v]和服务器协程[%v]通讯错误，错误信息为：%s\n",
			conn.RemoteAddr().String(), GetGID(), err)
	}
}

// 获取协程 id，用于标识哪个协程
func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

//完成对userMgr初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers : make(map[string]*UserProcess, 1024),
	}
}

// 服务端主函数（入口函数）
func main() {
	// 提示信息
	fmt.Println("服务器在 8080 端口监听....")
	// 服务端监听端口 8080
	listen, err := net.Listen("tcp", "0.0.0.0:8080")
	defer listen.Close()
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}

	// 一旦监听成功，就等待客户端连接的到来
	fmt.Println("服务器等待客户端的连接.....")
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=" ,err)
		}

		// 连接进来的客户端加入  用户(客户端)全局队列
		// 后面用于广播消息到该队列所有的客户端
		clientAddr := conn.RemoteAddr().String()
		var userProcess UserProcess
		userProcess.Conn = conn		// 当前连接
		userProcess.UserAddr = clientAddr	// 当前连接的客户端地址
		userKey := strings.Split(clientAddr, ":")[1]	// 当前连接的客户端端口号，用与标识是哪个客户端
		userMgr.onlineUsers[userKey] = &userProcess

		// 打印建立连接的信息
		fmt.Printf("客户端[%s]与服务端已建立连接.\n", clientAddr)

		// 客户端加入通讯室时，广播该客户端加入的消息
		JoinOrLeaveMsg(userKey, clientAddr, ClientJoin)

		// 一旦连接建立成功，则单独启动一个协程和客户端保持通讯
		go Process(conn, userKey)
	}
}
