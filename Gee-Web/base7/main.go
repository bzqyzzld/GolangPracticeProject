package main

import (
	"fmt"
	"gee"
)

// 错误处理

func main() {
	r := gee.NewEngine()
	r.Use(gee.Recovery())
	r.GET("/panic", func(context *gee.Context) {
		var name = []string{"age"}
		fmt.Println(name[2])
	})

	r.Run(":9999")
}
