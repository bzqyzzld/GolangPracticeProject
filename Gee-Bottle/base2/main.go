package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.NewEngine()
	r.GET("/", func(context *gee.Context) {
		context.HTML(http.StatusOK, "<h1>Hello World</h1>")
	})

	r.GET("/hello", func(context *gee.Context) {
		// /hello?name=haha
		context.String(http.StatusOK, "hello %s, you're at %s\n", context.Query("name"), context.Path)
	})

	r.POST("/login", func(context *gee.Context) {
		context.Json(http.StatusOK, gee.H{
			"username": context.PostForm("username"),
			"password": context.PostForm("password"),
		})
	})

	r.Run(":9999")
}
