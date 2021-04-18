/*
	file：chatRoomClient.go
	runCmd：go run chatRoomClient.go
 */

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

// 消息结构体
type Message struct {
	DataSource string	`json:"data_source"`	// 消息（数据）来源
	Data string	`json:"data"`					// 消息（数据）内容
}

// 与服务端建立连接，成功后返回连接句柄 conn
func GetConnect() (conn net.Conn, err error) {
	// 链接到服务器，这里用本地(localhost)举例，端口号为：8080，tcp 协议
	conn, err = net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	return
}

// 发送消息到服务端，入参为 连接句柄conn 和 发送消息内容content
func SendMsgToServer(conn net.Conn, content string) (err error) {
	// 去掉多余到换行符，并打印需要发送的消息内容
	readContent := strings.Trim(content, "\n")
	fmt.Println("发送的消息:", readContent)

	// 发送消息的内容转换成byte的格式，并发送给服务端
	data := []byte(readContent)
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("SendGroupMes err=", err.Error())
		return
	}
	return
}

// 客户端发送消息函数
func ClientSendMsg(conn net.Conn) {
	fmt.Println("请输入你想要发送的消息：")

	// 这里循环获取终端输入的消息
	for {
		// bufio 缓冲的方式获取终端输入的消息
		//fmt.Printf("发送消息: ")
		inputReader := bufio.NewReader(os.Stdin)
		// 读取到换行符就结束此次读取消息，并将消息存到 content 变量
		content, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("ReadString err=", err)
		}

		// 发送终端读取的消息到服务端
		err = SendMsgToServer(conn, content)
		if err != nil {
			fmt.Println("SendGroupMes err=", err)
		}
	}
}

func ClientReceiveMsg(conn net.Conn) {
	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("conn.Read err: ", err)
			return
		}
		var msg Message
		err = json.Unmarshal(buf[:n], &msg)
		if err != nil {
			fmt.Println("json.Unmarshal err: ", err)
		}

		fmt.Printf("收到[%s]发来的消息: %s\n",  msg.DataSource, msg.Data)
	}
}

// 客户端处理函数，入参为：与服务端建立连接的句柄 conn
func ClientProcess(conn net.Conn) {
	// 单独开启一个协程用于接收服务端发送过来的数据
	go ClientReceiveMsg(conn)
	// 主协程处理发送数据
	ClientSendMsg(conn)
}

// 客户端主函数（入口函数）
func main() {
	// 与服务端建立连接，返回连接句柄 conn
	conn, err := GetConnect()
	if err != nil {
		fmt.Println("GetConnect err=", err)
	}
	fmt.Println("与服务端建立连接成功，欢迎来到【ChatRoom】，现在可以进行通讯了.")
	// 客户端处理进程
	ClientProcess(conn)
}
