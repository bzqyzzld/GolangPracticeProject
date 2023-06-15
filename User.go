package main

import (
	"net"
	"strconv"
	"time"
)

type User struct {
	UserName    string
	UserAddress string
	UserConn    net.Conn
	UserChan    chan string
	Server      *Server
}

func CreateNewUser(conn net.Conn, server *Server) *User {
	userName := strconv.Itoa(int(time.Now().Unix()))
	newUser := &User{
		UserName:    userName,
		UserAddress: conn.RemoteAddr().String(),
		UserConn:    conn,
		UserChan:    make(chan string),
		Server:      server,
	}
	go newUser.UserListenMsg()
	return newUser
}

func (user *User) UserListenMsg() {
	for {
		msg := <-user.UserChan
		user.UserConn.Write([]byte(msg))
	}
}

func (user *User) UserOnLine() {
	// Server.UserGroup 增加该用户
	user.Server.Lock.Lock()
	user.Server.UserGroup[user.UserName] = user
	user.Server.Lock.Unlock()

	// 广播上线消息
	msg := DealMsg("我上线啦", user)
	user.Server.BroadCast(msg)
}

func (user *User) UserOffLine() {
	// Server.UserGroup 增减少该用户
	user.Server.Lock.Lock()
	delete(user.Server.UserGroup, user.UserName)
	user.Server.Lock.Unlock()

	// 广播下线消息
	msg := DealMsg("我下线啦", user)
	user.Server.BroadCast(msg)
}

func (user *User) UserDealMsg(msg string) {
	// 处理用户的消息
	msg = DealMsg(msg, user)
	user.Server.BroadCast(msg)
}
