package main

import "net"

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

// 用户处理消息的业务
func (user *User) DoMessage(message string) {
	user.server.BroadCast(user, message)
}

// 监听当前User channel的方法，一旦有消息，就直接发送给客户端
func (user *User) ListenMessage() {
	for {
		msg := <-user.C

		user.conn.Write([]byte(msg + "\n"))
	}
}
