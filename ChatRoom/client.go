package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func CreateNewClient(ip string, port int) *Client {
	client := &Client{
		ServerIp:   ip,
		ServerPort: port,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn

	return client
}

var ip string
var port int

func init() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "服务器ip(默认是127.0.0，1)")
	flag.IntVar(&port, "port", 8899, "服务器端口(默认是8899)")
}

func main() {
	flag.Parse()
	client := CreateNewClient(ip, port)
	if client == nil {
		fmt.Print("连接错误")
		return
	}
	fmt.Println("连接成功")
	select {}
}
