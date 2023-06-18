package main

import (
	"fmt"
	"net"
	"regexp"
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
	for msg := range user.UserChan {
		_, err := user.UserConn.Write([]byte(msg))
		if err != nil {
			fmt.Println("监听用户管道错误", err)
			panic(err)
		}
	}

	err := user.UserConn.Close()
	if err != nil {
		fmt.Println("关闭客户端连接错误", err)
		panic(err)
	}
}

func (user *User) UserOnLine() {
	// Server.UserGroup 增加该用户
	user.Server.Lock.Lock()
	user.Server.UserGroup[user.UserName] = user
	user.Server.Lock.Unlock()

	// 广播上线消息
	user.Server.BroadCast(user, "I'm online now")
}

func (user *User) UserOffLine() {
	// Server.UserGroup 增减少该用户
	user.Server.Lock.Lock()
	delete(user.Server.UserGroup, user.UserName)
	user.Server.Lock.Unlock()

	// 广播下线消息
	user.Server.BroadCast(user, "I'm offline now!!")

}

func (user *User) UserDealMsg(msg string) {
	// 处理用户的消息
	switch true {
	case regexp.MustCompile("^who$").MatchString(msg): // 显示我是谁,直接返回当前的用户名
		user.Server.Private(user, user.UserName, user.UserName)

	case regexp.MustCompile("^@(\\w+)\\s+(.*)").MatchString(msg): // 私聊某人
		r := regexp.MustCompile("^@(\\w+) (.*)").FindStringSubmatch(msg)
		toUserName := r[1]
		sendMsg := r[2]
		user.Server.Private(user, toUserName, sendMsg)

	case regexp.MustCompile("^changeName\\s+(\\w+)").MatchString(msg): // 修改用户名字
		r := regexp.MustCompile("^changeName\\s+(\\w+)").FindStringSubmatch(msg)
		newName := r[1]
		user.ChangeUserName(newName)

	default: // 直接广播用户的消息
		user.Server.BroadCast(user, msg)
	}
}

func (user *User) ChangeUserName(newName string) {
	// 修改用户名字
	user.Server.Lock.Lock()
	delete(user.Server.UserGroup, user.UserName)
	user.Server.UserGroup[newName] = user
	user.UserName = newName
	user.Server.Lock.Unlock()
}
