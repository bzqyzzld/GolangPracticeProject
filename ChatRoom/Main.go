package main

func main() {
	server := CreateNewServer("127.0.0.1", 8899)
	server.Start()
}
