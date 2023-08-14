package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	// 启动监听当前User channel消息的Goroutine
	go user.ListenMessage()

	return user
}

// 用户上线功能
func (user *User) Online() {
	// 用户上线消息
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	// 广播当前用户上线消息
	user.server.BroadCast(user, "已上线")
}

// 用户下线功能
func (user *User) Offline() {
	// 用户下线消息
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	// 广播当前用户下线消息
	user.server.BroadCast(user, "已下线")
}

// 给当前User对应的客户端发送消息
func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (user *User) DoMessage(message string) {
	if message == "who" {
		// 查询当前在线用户有哪些
		user.server.mapLock.Lock()
		for _, v := range user.server.OnlineMap {
			onlineMsg := fmt.Sprintf("[%s]%s:在线...\n", v.Addr, v.Name)
			user.SendMsg(onlineMsg)
		}
		user.server.mapLock.Unlock()
	} else if len(message) > 7 && message[:7] == "rename|" {
		// 消息格式 rename|张三
		newName := strings.Split(message, "|")[1]

		// 判断name是否存在
		_, ok := user.server.OnlineMap[newName]
		if ok {
			user.SendMsg("当前用户名已被使用!\n")
		} else {
			user.server.mapLock.Lock()
			delete(user.server.OnlineMap, user.Name)
			user.server.OnlineMap[newName] = user
			user.server.mapLock.Unlock()

			user.Name = newName
			user.SendMsg("你已经更新用户名:" + newName + "\n")
		}
	} else {
		user.server.BroadCast(user, message)
	}
}

// 监听当前User channel的方法，一旦有消息，就直接发送给客户端
func (user *User) ListenMessage() {
	for {
		msg := <-user.C

		user.conn.Write([]byte(msg + "\n"))
	}
}
