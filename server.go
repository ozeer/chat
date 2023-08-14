package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建一个Server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (server *Server) Handler(conn net.Conn) {
	// 当前连接的业务
	fmt.Println("连接建立成功!")
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
