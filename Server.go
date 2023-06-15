package main

import (
	"fmt"
	"net"
	"sync"
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

	// 处理用户消息
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
	}
}

func (server *Server) BroadCast(msg string) {
	server.BroadChan <- msg
}
