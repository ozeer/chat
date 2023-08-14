package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// 创建一个Server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听Message广播消息的channel的goroutine，一旦有消息就发送给全部在线的User
func (server *Server) ListeningMessage() {
	for {
		message := <-server.Message

		// 将消息发送给全部的User
		server.mapLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.C <- message
		}
		server.mapLock.Unlock()
	}
}

// 广播消息方法
func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]%s:%s", user.Addr, user.Name, msg)

	server.Message <- sendMsg
}

func (server *Server) Handler(conn net.Conn) {
	// 当前连接的业务
	// fmt.Println("连接建立成功!")

	user := NewUser(conn, server)
	user.Online()

	// 接收客户端用户发送的消息
	go func() {
		buf := make([]byte, 4096)

		for {
			n, err := conn.Read(buf)

			// 用户下线
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read error: ", err)
				return
			}

			// 提取用户的消息
			msg := string(buf[:n-1])

			// 将得到的消息进行广播
			user.DoMessage(msg)
		}
	}()

	// 当前handler阻塞
	select {}
}

// 启动服务器的接口
func (server *Server) Start() {
	// socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.Listen error: ", err)
		return
	}

	// 关闭socket
	defer listen.Close()

	// 启动监听
	go server.ListeningMessage()

	for {
		// accept
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("net.Accept error: ", err)
			continue
		}

		// do handler
		go server.Handler(conn)
	}
}
