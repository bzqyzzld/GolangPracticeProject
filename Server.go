package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	IP        string
	PORT      int
	UserGroup map[string]*User
	BroadChan chan string
	Lock      sync.RWMutex
}

func CreateNewServer(ip string, port int) *Server {
	newServer := &Server{
		IP:        ip,
		PORT:      port,
		UserGroup: make(map[string]*User),
		BroadChan: make(chan string),
		Lock:      sync.RWMutex{},
	}
	go newServer.ListenMsg()
	return newServer
}

func (server *Server) ListenMsg() {
	for {
		msg := <-server.BroadChan
		server.Lock.Lock()
		for _, client := range server.UserGroup {
			client.UserChan <- msg
		}
		server.Lock.Unlock()
	}

}

func (server *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.IP, server.PORT))
	if err != nil {
		fmt.Println("启动监听端口错误", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Print("连接过来错误了", err)
			continue
		}

		go server.Handle(conn)
	}
}

func (server *Server) Handle(conn net.Conn) {
	newUser := CreateNewUser(conn, server)

	// 用户上线
	fmt.Printf("用户%s上线了\n", newUser.UserName)
	newUser.UserOnLine()

	isLive := make(chan bool)

	// 处理用户消息
	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil || n == 0 {
				// 用户连接断开了
				fmt.Printf("用户%s连接已经断开\n", newUser.UserName)
				newUser.UserOffLine()
				break
			}

			// 处理用户的消息
			newUser.UserDealMsg(string(buf[:n]))
			isLive <- true
		}
	}()

	// 用户超时强制踢出,主要配合select和time.After使用
	for {
		select {
		case <-isLive:
			// 重置定时器,可以什么都不写
		case <-time.After(10 * time.Second):
			// 用户已经超时了,可以强制踢出
			server.Private(newUser, newUser.UserName, "You had been kicked out!!!")
			time.Sleep(time.Second) // 暂停1s钟，等待踢人消息完成发送
			close(newUser.UserChan)
			conn.Close()
			return
		}
	}

}

func (server *Server) BroadCast(fromUser *User, msg string) {
	sendMsg := DealMsg(msg, fromUser)
	server.BroadChan <- sendMsg
}

func (server *Server) Private(fromUser *User, toUserName, msg string) {
	isFound := false
	server.Lock.Lock()
	for name, user := range server.UserGroup {
		if name == toUserName {
			isFound = true
			sendMsg := DealMsg(msg, fromUser)
			user.UserChan <- sendMsg
			break
		}
	}
	server.Lock.Unlock()

	// 没有找到toUserName, 返回系统消息
	if isFound == false {
		sendMsg := DealMsg(fmt.Sprintf("没有找到%s用户", toUserName), nil)
		fromUser.UserChan <- sendMsg
	}
}
